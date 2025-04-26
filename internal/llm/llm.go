package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Client struct {
	baseURL   string
	engineID  string
	modelName string
	client    *http.Client
}

func NewClient(baseURL, engineID, modelName string) *Client {
	return &Client{
		baseURL:   baseURL,
		engineID:  engineID,
		modelName: modelName,
		client:    &http.Client{Timeout: 15 * time.Second},
	}
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

func (c *Client) Summarize(prompt string) (string, error) {
	reqBody := ChatRequest{
		Model: c.modelName,
		Messages: []ChatMessage{
			{
				Role:    "system",
				Content: "You are a helpful assistant. Summarize the following logs briefly highlighting errors or important events.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	bodyBytes, _ := json.Marshal(reqBody)
	url := c.baseURL + "/engines/" + c.engineID + "/v1/chat/completions"

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("LLM call failed")
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", errors.New("no choices returned")
	}

	return chatResp.Choices[0].Message.Content, nil
}
