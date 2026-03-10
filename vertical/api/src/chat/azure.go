package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"go.uber.org/zap"
)

type AzureClient struct {
	endpoint            string
	apiKey              string
	chatDeployment      string
	embeddingDeployment string
	apiVersion          string
}

func NewAzureClient() *AzureClient {
	return &AzureClient{
		endpoint:            os.Getenv("AZURE_OPENAI_ENDPOINT"),
		apiKey:              os.Getenv("AZURE_OPENAI_API_KEY"),
		chatDeployment:      os.Getenv("AZURE_OPENAI_CHAT_DEPLOYMENT"),
		embeddingDeployment: os.Getenv("AZURE_OPENAI_EMBEDDING_DEPLOYMENT"),
		apiVersion:          os.Getenv("AZURE_OPENAI_API_VERSION"),
	}
}

func (c *AzureClient) IsConfigured() bool {
	return c.endpoint != "" && c.apiKey != "" && c.chatDeployment != ""
}

// Embedding

type embeddingRequest struct {
	Input string `json:"input"`
}

type embeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

func (c *AzureClient) CreateEmbedding(text string) ([]float32, error) {
	url := fmt.Sprintf("%s/openai/deployments/%s/embeddings?api-version=%s",
		c.endpoint, c.embeddingDeployment, c.apiVersion)

	body, _ := json.Marshal(embeddingRequest{Input: text})
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("embedding request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		zap.S().Errorw("Embedding API error", "status", resp.StatusCode, "body", string(respBody))
		return nil, fmt.Errorf("embedding API returned %d", resp.StatusCode)
	}

	var result embeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}
	return result.Data[0].Embedding, nil
}

// Chat completion (streaming)

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	Stream      bool          `json:"stream"`
}

func (c *AzureClient) ChatCompletionStream(messages []chatMessage) (*http.Response, error) {
	url := fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=%s",
		c.endpoint, c.chatDeployment, c.apiVersion)

	body, _ := json.Marshal(chatRequest{
		Messages:    messages,
		Temperature: 0.3,
		Stream:      true,
	})

	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", c.apiKey)

	return http.DefaultClient.Do(req)
}

// Chat completion (non-streaming, for title generation)

type chatResponseChoice struct {
	Message chatMessage `json:"message"`
}

type chatResponse struct {
	Choices []chatResponseChoice `json:"choices"`
}

func (c *AzureClient) ChatCompletion(messages []chatMessage) (string, error) {
	url := fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=%s",
		c.endpoint, c.chatDeployment, c.apiVersion)

	body, _ := json.Marshal(chatRequest{
		Messages:    messages,
		Temperature: 0.3,
		Stream:      false,
	})

	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("chat API returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}
	return result.Choices[0].Message.Content, nil
}
