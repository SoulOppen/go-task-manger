package taskllm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultOpenAIBase = "https://api.openai.com/v1"

type openAIProvider struct {
	cfg    Config
	client *http.Client
}

func newOpenAIProvider(cfg Config) *openAIProvider {
	base := strings.TrimRight(cfg.BaseURL, "/")
	if base == "" {
		base = defaultOpenAIBase
	}
	cfg.BaseURL = base
	return &openAIProvider{
		cfg: cfg,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (o *openAIProvider) CompleteJSON(ctx context.Context, system, user string) (string, error) {
	url := o.cfg.BaseURL + "/chat/completions"
	body := map[string]interface{}{
		"model":       o.cfg.Model,
		"temperature": 0.1,
		"response_format": map[string]string{
			"type": "json_object",
		},
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": user},
		},
	}
	rawBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(rawBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.cfg.APIKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("openai-compatible HTTP %d: %s", resp.StatusCode, truncateForErr(b))
	}
	return extractOpenAIContent(b)
}

func extractOpenAIContent(respBody []byte) (string, error) {
	var root struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(respBody, &root); err != nil {
		return "", fmt.Errorf("respuesta openai-compatible no JSON: %w", err)
	}
	if root.Error != nil {
		return "", fmt.Errorf("API: %s", root.Error.Message)
	}
	if len(root.Choices) == 0 {
		return "", fmt.Errorf("openai-compatible: choices vacio")
	}
	return root.Choices[0].Message.Content, nil
}
