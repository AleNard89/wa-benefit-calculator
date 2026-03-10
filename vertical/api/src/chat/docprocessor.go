package chat

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

// PPTX text extraction

type pptxSlideBody struct {
	XMLName xml.Name `xml:"sld"`
	CSld    struct {
		SpTree struct {
			Shapes []pptxShape `xml:"sp"`
		} `xml:"spTree"`
	} `xml:"cSld"`
}

type pptxShape struct {
	TxBody struct {
		Paragraphs []struct {
			Runs []struct {
				Text string `xml:"t"`
			} `xml:"r"`
		} `xml:"p"`
	} `xml:"txBody"`
}

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

		data, err := io.ReadAll(rc)
		if err != nil {
			return "", fmt.Errorf("cannot read document.xml: %w", err)
		}

		type docxRun struct {
			Text string `xml:"t"`
		}
		type docxParagraph struct {
			Runs []docxRun `xml:"r"`
		}
		type docxBody struct {
			Paragraphs []docxParagraph `xml:"p"`
		}
		type docxDocument struct {
			XMLName xml.Name `xml:"document"`
			Body    docxBody `xml:"body"`
		}

		var doc docxDocument
		if err := xml.Unmarshal(data, &doc); err != nil {
			return "", fmt.Errorf("cannot parse document.xml: %w", err)
		}

		var paragraphs []string
		for _, p := range doc.Body.Paragraphs {
			var line []string
			for _, run := range p.Runs {
				if run.Text != "" {
					line = append(line, run.Text)
				}
			}
			if len(line) > 0 {
				paragraphs = append(paragraphs, strings.Join(line, ""))
			}
		}
		return strings.Join(paragraphs, "\n"), nil
	}
	return "", fmt.Errorf("word/document.xml not found in docx")
}

func ExtractTextFromPPTX(filePath string) ([]string, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open pptx: %w", err)
	}
	defer r.Close()

	var slideTexts []string

	for i := 1; i <= 100; i++ {
		slideName := fmt.Sprintf("ppt/slides/slide%d.xml", i)
		text := extractSlideText(r, slideName)
		if text == "" {
			break
		}
		slideTexts = append(slideTexts, text)
	}

	return slideTexts, nil
}

func extractSlideText(r *zip.ReadCloser, slideName string) string {
	for _, f := range r.File {
		if f.Name != slideName {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return ""
		}
		defer rc.Close()

		data, err := io.ReadAll(rc)
		if err != nil {
			return ""
		}

		var slide pptxSlideBody
		if err := xml.Unmarshal(data, &slide); err != nil {
			return ""
		}

		var parts []string
		for _, shape := range slide.CSld.SpTree.Shapes {
			for _, para := range shape.TxBody.Paragraphs {
				var line []string
				for _, run := range para.Runs {
					if run.Text != "" {
						line = append(line, run.Text)
					}
				}
				if len(line) > 0 {
					parts = append(parts, strings.Join(line, ""))
				}
			}
		}
		return strings.Join(parts, "\n")
	}
	return ""
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
		if end > len(text) {
			end = len(text)
		}
		// Try to break at a newline or space
		if end < len(text) {
			for i := end; i > start+maxChars/2; i-- {
				if text[i] == '\n' || text[i] == ' ' {
					end = i
					break
				}
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
	default:
		return fmt.Errorf("unsupported file type: %s", ext)
	}

	if strings.TrimSpace(allText) == "" {
		zap.S().Warnw("No text extracted from file", "file", filePath)
		return nil
	}

	// Split into chunks (~2000 chars, ~500 tokens)
	chunks := ChunkText(allText, 2000, 200)

	// Delete existing chunks for this file
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
