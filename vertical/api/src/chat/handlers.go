package chat

import (
	"bufio"
	"encoding/json"
	"orbita-api/auth"
	"orbita-api/core/httpx"
	d "orbita-api/core/decorators"
	"orbita-api/orgs"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var azure *AzureClient

func init() {
	azure = NewAzureClient()
	if azure.IsConfigured() {
		zap.S().Info("Azure OpenAI client configured")
	} else {
		zap.S().Warn("Azure OpenAI not configured, chat will return errors")
	}
}

func GetAzureClient() *AzureClient {
	return azure
}

// GET /chat/conversations
func getConversations(c *gin.Context) {
	user := auth.MustCurrentUser(c)
	companyID := httpx.HeaderCompanyID(c)
	service := ConversationService{}
	conversations, err := service.GetByUser(user.ID, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, conversations)
}

var GetConversations = d.AuthenticationRequired(getConversations)

// GET /chat/conversations/:id/messages
func getMessages(c *gin.Context) {
	user := auth.MustCurrentUser(c)
	companyID := httpx.HeaderCompanyID(c)
	convID, _ := strconv.Atoi(c.Param("id"))

	convService := ConversationService{}
	conv, err := convService.GetByID(convID)
	if err != nil || conv.UserID != user.ID || conv.CompanyID != companyID {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Conversazione non trovata"})
		return
	}

	msgService := MessageService{}
	messages, err := msgService.GetByConversation(convID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, messages)
}

var GetMessages = d.AuthenticationRequired(getMessages)

// POST /chat/conversations
type CreateConversationBody struct {
	Title string `json:"title"`
}

func createConversation(c *gin.Context) {
	user := auth.MustCurrentUser(c)
	companyID := httpx.HeaderCompanyID(c)

	body := CreateConversationBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		body.Title = "Nuova conversazione"
	}
	if body.Title == "" {
		body.Title = "Nuova conversazione"
	}

	service := ConversationService{}
	conv, err := service.Create(user.ID, companyID, body.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusCreated, conv)
}

var CreateConversation = d.AuthenticationRequired(createConversation)

// DELETE /chat/conversations/:id
func deleteConversation(c *gin.Context) {
	user := auth.MustCurrentUser(c)
	companyID := httpx.HeaderCompanyID(c)
	convID, _ := strconv.Atoi(c.Param("id"))

	convService := ConversationService{}
	conv, err := convService.GetByID(convID)
	if err != nil || conv.UserID != user.ID || conv.CompanyID != companyID {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Conversazione non trovata"})
		return
	}

	if err := convService.Delete(convID); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

var DeleteConversation = d.AuthenticationRequired(deleteConversation)

// POST /chat/conversations/:id/messages - SSE streaming
type SendMessageBody struct {
	Content string `json:"content" binding:"required"`
}

func sendMessage(c *gin.Context) {
	user := auth.MustCurrentUser(c)
	companyID := httpx.HeaderCompanyID(c)
	convID, _ := strconv.Atoi(c.Param("id"))

	if !azure.IsConfigured() {
		c.JSON(http.StatusServiceUnavailable, httpx.ErrorResponse{Message: "Azure OpenAI non configurato"})
		return
	}

	convService := ConversationService{}
	conv, err := convService.GetByID(convID)
	if err != nil || conv.UserID != user.ID || conv.CompanyID != companyID {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Conversazione non trovata"})
		return
	}

	body := SendMessageBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpx.ErrorResponse{Message: "Messaggio vuoto"})
		return
	}

	// Save user message
	msgService := MessageService{}
	_, err = msgService.Create(convID, "user", body.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: "Errore salvataggio messaggio"})
		return
	}

	// Get conversation history
	history, _ := msgService.GetLastN(convID, 20)

	// Build RAG-augmented messages
	messages := BuildMessages(azure, companyID, history, body.Content)

	// Generate title for new conversations
	if conv.Title == "Nuova conversazione" || conv.Title == "" {
		go func() {
			title := GenerateTitle(azure, body.Content)
			_ = convService.UpdateTitle(convID, title)
		}()
	}

	// Stream response via SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	resp, err := azure.ChatCompletionStream(messages)
	if err != nil {
		fmt.Fprintf(c.Writer, "data: {\"error\": \"Errore connessione Azure OpenAI\"}\n\n")
		c.Writer.Flush()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		zap.S().Errorw("Azure chat error", "status", resp.StatusCode, "body", string(respBody))
		fmt.Fprintf(c.Writer, "data: {\"error\": \"Errore Azure OpenAI (%d)\"}\n\n", resp.StatusCode)
		c.Writer.Flush()
		return
	}

	var fullResponse strings.Builder
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		// Parse the SSE chunk to extract content delta
		type deltaChunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}

		var chunk deltaChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		for _, choice := range chunk.Choices {
			if choice.Delta.Content != "" {
				fullResponse.WriteString(choice.Delta.Content)
				fmt.Fprintf(c.Writer, "data: %s\n\n", data)
				c.Writer.Flush()
			}
		}
	}

	// Save assistant response
	if fullResponse.Len() > 0 {
		_, _ = msgService.Create(convID, "assistant", fullResponse.String())
		_ = convService.Touch(convID)
	}

	fmt.Fprintf(c.Writer, "data: [DONE]\n\n")
	c.Writer.Flush()
}

var SendMessage = d.AuthenticationRequired(sendMessage)

// POST /chat/index - trigger document indexing for current company
func indexDocuments(c *gin.Context) {
	if !azure.IsConfigured() {
		c.JSON(http.StatusServiceUnavailable, httpx.ErrorResponse{Message: "Azure OpenAI non configurato"})
		return
	}

	companyID := httpx.HeaderCompanyID(c)
	if companyID == 0 {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Company ID richiesto"})
		return
	}

	// Get company to find storage path
	companyService := orgs.CompanyService{}
	company, err := companyService.GetByID(companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Azienda non trovata"})
		return
	}

	go func() {
		zap.S().Infow("Starting document indexing", "company", company.Name, "path", company.StoragePath)
		if err := IndexCompanyFolder(azure, companyID, company.StoragePath); err != nil {
			zap.S().Errorw("Indexing failed", "company", company.Name, "error", err)
		} else {
			zap.S().Infow("Indexing completed", "company", company.Name)
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{"message": "Indicizzazione avviata"})
}

var IndexDocuments = d.SuperuserRequired(indexDocuments)
