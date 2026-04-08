package taskllm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultGeminiBase = "https://generativelanguage.googleapis.com"

type geminiProvider struct {
	cfg    Config
	client *http.Client
}

func newGeminiProvider(cfg Config) *geminiProvider {
	base := strings.TrimRight(cfg.BaseURL, "/")
	if base == "" {
		base = defaultGeminiBase
	}
	cfg.BaseURL = base
	return &geminiProvider{
		cfg: cfg,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (g *geminiProvider) CompleteJSON(ctx context.Context, system, user string) (string, error) {
	model := strings.TrimPrefix(g.cfg.Model, "models/")
	u, err := url.Parse(g.cfg.BaseURL + "/v1beta/models/" + url.PathEscape(model) + ":generateContent")
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("key", g.cfg.APIKey)
	u.RawQuery = q.Encode()
	urlStr := u.String()

	body := map[string]interface{}{
		"systemInstruction": map[string]interface{}{
			"parts": []map[string]string{{"text": system}},
		},
		"contents": []map[string]interface{}{
			{
				"role":  "user",
				"parts": []map[string]string{{"text": user}},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":      0.1,
			"responseMimeType": "application/json",
		},
	}
	rawBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlStr, bytes.NewReader(rawBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("gemini HTTP %d: %s", resp.StatusCode, truncateForErr(b))
	}
	return extractGeminiText(b)
}

func extractGeminiText(respBody []byte) (string, error) {
	var root struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
		Error *struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(respBody, &root); err != nil {
		return "", fmt.Errorf("respuesta gemini no JSON: %w", err)
	}
	if root.Error != nil {
		return "", fmt.Errorf("gemini API: %s", root.Error.Message)
	}
	if len(root.Candidates) == 0 || len(root.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini: candidatos vacios")
	}
	return root.Candidates[0].Content.Parts[0].Text, nil
}

func truncateForErr(b []byte) string {
	s := string(b)
	if len(s) > 400 {
		return s[:400] + "..."
	}
	return s
}
