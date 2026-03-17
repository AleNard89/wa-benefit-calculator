package processes

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"benefit-calculator-api/chat"
	"benefit-calculator-api/core/httpx"
	"benefit-calculator-api/orgs"

	d "benefit-calculator-api/core/decorators"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var allowedExtensions = map[string]bool{".pptx": true, ".pdf": true, ".docx": true}

func uploadDocument(c *gin.Context) {
	companyID := httpx.HeaderCompanyID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Invalid process ID"})
		return
	}

	service := ProcessService{}
	process, err := service.GetByID(id, companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Process not found"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "File richiesto"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Formato non supportato. Usa PPTX, DOCX o PDF."})
		return
	}

	companyService := orgs.CompanyService{}
	company, err := companyService.GetByID(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: "Errore recupero azienda"})
		return
	}

	// Delete old file if exists
	if process.DocumentPath != nil && *process.DocumentPath != "" {
		oldPath := *process.DocumentPath
		_ = os.Remove(oldPath)
		// Clean up old RAG chunks
		chunkService := chat.ChunkService{}
		_ = chunkService.DeleteByFile(companyID, oldPath)
	}

	// Save new file: /data/companies/<slug>/process_<id>_<filename>
	safeFileName := fmt.Sprintf("process_%d_%s", process.ID, sanitizeFilename(file.Filename))
	destPath := filepath.Join(company.StoragePath, safeFileName)

	if err := c.SaveUploadedFile(file, destPath); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: "Errore salvataggio file"})
		return
	}

	// Update process with document info
	process.DocumentPath = &destPath
	process.DocumentName = &file.Filename
	process.UpdatedAt = time.Now()
	if err := service.Update(process); err != nil {
		_ = os.Remove(destPath)
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: "Errore aggiornamento processo"})
		return
	}

	zap.S().Infow("Documento caricato", "process_id", id, "file", file.Filename)
	c.JSON(http.StatusOK, process)

	// Index for RAG asynchronously (serialized via worker queue)
	chat.EnqueueIndexing(companyID, destPath, file.Filename, id)
}

var UploadDocument = d.CompanyRequiredAndPermissionRequired([]string{ProcessUpdate}, uploadDocument)

func downloadDocument(c *gin.Context) {
	companyID := httpx.HeaderCompanyID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Invalid process ID"})
		return
	}

	service := ProcessService{}
	process, err := service.GetByID(id, companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Process not found"})
		return
	}

	if process.DocumentPath == nil || *process.DocumentPath == "" {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Nessun documento associato"})
		return
	}

	filePath := *process.DocumentPath
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "File non trovato su disco"})
		return
	}

	fileName := "documento"
	if process.DocumentName != nil {
		fileName = *process.DocumentName
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.File(filePath)
}

var DownloadDocument = d.CompanyRequiredAndPermissionRequired([]string{ProcessRead}, downloadDocument)

func deleteDocument(c *gin.Context) {
	companyID := httpx.HeaderCompanyID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Invalid process ID"})
		return
	}

	service := ProcessService{}
	process, err := service.GetByID(id, companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Process not found"})
		return
	}

	if process.DocumentPath == nil || *process.DocumentPath == "" {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Nessun documento associato"})
		return
	}

	// Delete file from disk
	_ = os.Remove(*process.DocumentPath)

	// Delete RAG chunks
	chunkService := chat.ChunkService{}
	_ = chunkService.DeleteByFile(companyID, *process.DocumentPath)

	// Clear document fields
	process.DocumentPath = nil
	process.DocumentName = nil
	process.UpdatedAt = time.Now()
	if err := service.Update(process); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: "Errore aggiornamento processo"})
		return
	}

	zap.S().Infow("Documento eliminato", "process_id", id)
	c.JSON(http.StatusOK, process)
}

var DeleteDocument = d.CompanyRequiredAndPermissionRequired([]string{ProcessUpdate}, deleteDocument)

func sanitizeFilename(name string) string {
	name = filepath.Base(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' || r == '.' {
			return r
		}
		return '_'
	}, name)
	return name
}
