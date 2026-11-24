// Package machine calls the separate translation service for labelled machine
// translation (TRN-3, P3-03). The service is optional: with no URL there is
// none, and the API returns verified phrases only. Its contract is in
// docs/translation-service.md.
package machine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ErrUnsupported means the service does not translate that language pair.
var ErrUnsupported = errors.New("language pair not supported")

// Result is one machine translation.
type Result struct {
	Translation string  `json:"translation"`
	Confidence  float64 `json:"confidence"` // 0 to 1
	Model       string  `json:"model"`
}

// Translator translates text between two language codes.
type Translator interface {
	Translate(ctx context.Context, text, from, to string) (Result, error)
}

// Client calls the translation service over HTTP.
type Client struct {
	url  string
	http *http.Client
}

// New returns a client for the service at baseURL, or nil when baseURL is
// empty: no machine translation.
func New(baseURL string) Translator {
	if baseURL == "" {
		return nil
	}
	return &Client{url: strings.TrimSuffix(baseURL, "/") + "/translate", http: &http.Client{Timeout: 5 * time.Second}}
}

// Translate implements Translator.
func (c *Client) Translate(ctx context.Context, text, from, to string) (Result, error) {
	body, err := json.Marshal(map[string]string{"text": text, "from": from, "to": to})
	if err != nil {
		return Result{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.http.Do(req)
	if err != nil {
		return Result{}, fmt.Errorf("translation service: %w", err)
	}
	defer func() { _ = res.Body.Close() }() // a read-only body; nothing to report
	switch {
	case res.StatusCode == http.StatusUnprocessableEntity:
		return Result{}, ErrUnsupported
	case res.StatusCode != http.StatusOK:
		return Result{}, fmt.Errorf("translation service: %s", res.Status)
	}
	var r Result
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return Result{}, fmt.Errorf("translation service: %w", err)
	}
	if r.Translation == "" || r.Model == "" || r.Confidence < 0 || r.Confidence > 1 {
		return Result{}, errors.New("translation service: incomplete answer")
	}
	return r, nil
}
