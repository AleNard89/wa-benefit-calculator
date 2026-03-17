package chat

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gopdf "github.com/ledongthuc/pdf"
	"go.uber.org/zap"
)

// PPTX streaming text extraction
// Uses xml.Decoder to read token-by-token, extracting only <a:t> text elements.
// Never loads the full slide XML into memory — safe for slides with large embedded images.

func ExtractTextFromPPTX(filePath string) ([]string, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open pptx: %w", err)
	}
	defer r.Close()

	var slideTexts []string

	for i := 1; i <= 200; i++ {
		slideName := fmt.Sprintf("ppt/slides/slide%d.xml", i)
		text, found := streamSlideText(r, slideName)
		if !found {
			break
		}
		if text != "" {
			slideTexts = append(slideTexts, text)
		}
	}

	return slideTexts, nil
}

func streamSlideText(r *zip.ReadCloser, slideName string) (string, bool) {
	for _, f := range r.File {
		if f.Name != slideName {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return "", true
		}
		defer rc.Close()

		decoder := xml.NewDecoder(rc)
		var parts []string
		var inTextElement bool

		for {
			tok, err := decoder.Token()
			if err != nil {
				break
			}
			switch t := tok.(type) {
			case xml.StartElement:
				// <a:t> is the text run element in OOXML (both PPTX and charts)
				if t.Name.Local == "t" && (t.Name.Space == "http://schemas.openxmlformats.org/drawingml/2006/main" || t.Name.Space == "a" || t.Name.Space == "") {
					inTextElement = true
				}
			case xml.EndElement:
				if t.Name.Local == "t" {
					inTextElement = false
				}
			case xml.CharData:
				if inTextElement {
					text := strings.TrimSpace(string(t))
					if text != "" {
						parts = append(parts, text)
					}
				}
			}
		}
		return strings.Join(parts, " "), true
	}
	return "", false
}

// DOCX streaming text extraction
// Uses xml.Decoder to read token-by-token, extracting only <w:t> text elements.

func ExtractTextFromDOCX(filePath string) (string, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot open docx: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name != "word/document.xml" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return "", fmt.Errorf("cannot open document.xml: %w", err)
		}
		defer rc.Close()

		decoder := xml.NewDecoder(rc)
		var paragraphs []string
		var currentPara []string
		var inTextElement bool

		for {
			tok, err := decoder.Token()
			if err != nil {
				break
			}
			switch t := tok.(type) {
			case xml.StartElement:
				if t.Name.Local == "t" && (t.Name.Space == "http://schemas.openxmlformats.org/wordprocessingml/2006/main" || t.Name.Space == "w" || t.Name.Space == "") {
					inTextElement = true
				}
			case xml.EndElement:
				if t.Name.Local == "t" {
					inTextElement = false
				}
				// End of paragraph <w:p>
				if t.Name.Local == "p" {
					if len(currentPara) > 0 {
						paragraphs = append(paragraphs, strings.Join(currentPara, ""))
						currentPara = nil
					}
				}
			case xml.CharData:
				if inTextElement {
					text := string(t)
					if text != "" {
						currentPara = append(currentPara, text)
					}
				}
			}
		}
		if len(currentPara) > 0 {
			paragraphs = append(paragraphs, strings.Join(currentPara, ""))
		}
		return strings.Join(paragraphs, "\n"), nil
	}
	return "", fmt.Errorf("word/document.xml not found in docx")
}

// PDF text extraction

func ExtractTextFromPDF(filePath string) (string, error) {
	f, r, err := gopdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot open pdf: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	for i := 1; i <= r.NumPage(); i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}
		text, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}
		trimmed := strings.TrimSpace(text)
		if trimmed != "" {
			if buf.Len() > 0 {
				buf.WriteString("\n\n")
			}
			buf.WriteString(trimmed)
		}
	}
	return buf.String(), nil
}

// Chunking

func ChunkText(text string, maxChars int, overlap int) []string {
	if len(text) <= maxChars {
		return []string{text}
	}

	var chunks []string
	start := 0
	for start < len(text) {
		end := start + maxChars
		if end >= len(text) {
			chunks = append(chunks, strings.TrimSpace(text[start:]))
			break
		}
		for i := end; i > start+maxChars/2; i-- {
			if text[i] == '\n' || text[i] == ' ' {
				end = i
				break
			}
		}
		chunks = append(chunks, strings.TrimSpace(text[start:end]))
		start = end - overlap
		if start < 0 {
			start = 0
		}
	}
	return chunks
}

// Index a single file

func IndexFile(azure *AzureClient, companyID int, filePath string) error {
	fileName := filepath.Base(filePath)
	ext := strings.ToLower(filepath.Ext(filePath))

	var allText string

	switch ext {
	case ".pptx":
		slides, err := ExtractTextFromPPTX(filePath)
		if err != nil {
			return fmt.Errorf("error extracting pptx: %w", err)
		}
		allText = strings.Join(slides, "\n\n---\n\n")
	case ".docx":
		text, err := ExtractTextFromDOCX(filePath)
		if err != nil {
			return fmt.Errorf("error extracting docx: %w", err)
		}
		allText = text
	case ".pdf":
		text, err := ExtractTextFromPDF(filePath)
		if err != nil {
			return fmt.Errorf("error extracting pdf: %w", err)
		}
		allText = text
	default:
		return fmt.Errorf("unsupported file type: %s", ext)
	}

	if strings.TrimSpace(allText) == "" {
		zap.S().Warnw("No text extracted from file", "file", filePath)
		return nil
	}

	chunks := ChunkText(allText, 2000, 200)

	chunkService := ChunkService{}
	if err := chunkService.DeleteByFile(companyID, filePath); err != nil {
		zap.S().Warnw("Error deleting old chunks", "error", err)
	}

	zap.S().Infow("Indexing file", "file", fileName, "chunks", len(chunks))

	for i, chunk := range chunks {
		embedding, err := azure.CreateEmbedding(chunk)
		if err != nil {
			zap.S().Warnw("Error creating embedding", "file", fileName, "chunk", i, "error", err)
			continue
		}

		doc := &DocumentChunk{
			CompanyID:  companyID,
			FilePath:   filePath,
			FileName:   fileName,
			ChunkIndex: i,
			Content:    chunk,
		}
		if err := chunkService.Insert(doc, embedding); err != nil {
			zap.S().Warnw("Error inserting chunk", "file", fileName, "chunk", i, "error", err)
		}
	}

	zap.S().Infow("File indexed successfully", "file", fileName, "chunks", len(chunks))
	return nil
}

// Index all files in a company folder

func IndexCompanyFolder(azure *AzureClient, companyID int, folderPath string) error {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("cannot read folder: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".pptx" && ext != ".pdf" && ext != ".docx" {
			continue
		}
		fullPath := filepath.Join(folderPath, entry.Name())
		if err := IndexFile(azure, companyID, fullPath); err != nil {
			zap.S().Warnw("Error indexing file", "file", entry.Name(), "error", err)
		}
	}
	return nil
}
