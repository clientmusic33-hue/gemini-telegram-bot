package gemini

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Message struct {
	Role    string
	Content string
}

type Client struct {
	client     *genai.Client
	model      string
	timeout    time.Duration
	maxRetries int
}

func New(apiKey, model string, timeoutSeconds, maxRetries int) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Gemini API key is required")
	}
	if model == "" {
		model = "gemini-2.0-flash"
	}
	if timeoutSeconds <= 0 {
		timeoutSeconds = 30
	}
	if maxRetries <= 0 {
		maxRetries = 3
	}

	gclient, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("create Gemini client: %w", err)
	}

	return &Client{
		client:     gclient,
		model:      model,
		timeout:    time.Duration(timeoutSeconds) * time.Second,
		maxRetries: maxRetries,
	}, nil
}

func (c *Client) SendMessage(parent context.Context, messages []Message) (string, error) {
	if len(messages) == 0 {
		return "", fmt.Errorf("at least one message is required")
	}

	ctx, cancel := context.WithTimeout(parent, c.timeout)
	defer cancel()

	history := make([]*genai.Content, 0, len(messages)-1)
	for _, message := range messages[:len(messages)-1] {
		role := "user"
		if message.Role == "model" || message.Role == "assistant" {
			role = "model"
		}
		history = append(history, &genai.Content{Role: role, Parts: []genai.Part{genai.Text(message.Content)}})
	}

	last := messages[len(messages)-1].Content
	var lastErr error
	for attempt := 0; attempt < c.maxRetries; attempt++ {
		model := c.client.GenerativeModel(c.model)
		model.History = history
		model.GenerationConfig = &genai.GenerationConfig{
			Temperature:     genai.Float32Ptr(0.7),
			TopP:            genai.Float32Ptr(0.95),
			MaxOutputTokens: genai.Int32Ptr(2048),
		}

		response, err := model.GenerateContent(ctx, genai.Text(last))
		if err == nil {
			text := responseText(response)
			if text != "" {
				return text, nil
			}
			return "", fmt.Errorf("Gemini returned an empty response")
		}
		lastErr = err

		if attempt+1 < c.maxRetries {
			select {
			case <-time.After(time.Duration(1<<attempt) * time.Second):
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
	}
	return "", fmt.Errorf("Gemini request failed after retries: %w", lastErr)
}

func responseText(response *genai.GenerateContentResponse) string {
	if response == nil {
		return ""
	}
	var text strings.Builder
	for _, candidate := range response.Candidates {
		if candidate == nil || candidate.Content == nil {
			continue
		}
		for _, part := range candidate.Content.Parts {
			text.WriteString(fmt.Sprint(part))
		}
	}
	return strings.TrimSpace(text.String())
}

func (c *Client) Close() error {
	if c.client == nil {
		return nil
	}
	return c.client.Close()
}
