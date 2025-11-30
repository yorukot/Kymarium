package notification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/yorukot/knocker/models"
)

func TestSendDiscord(t *testing.T) {
	var payload struct {
		Username string `json:"username"`
		Embeds   []struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		} `json:"embeds"`
	}

	client := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodPost {
				return nil, errors.New("expected POST")
			}

			if ct := r.Header.Get("Content-Type"); ct != "application/json" {
				return nil, fmt.Errorf("unexpected content type: %s", ct)
			}

			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				return nil, fmt.Errorf("failed to decode payload: %w", err)
			}

			return &http.Response{
				StatusCode: http.StatusNoContent,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     make(http.Header),
			}, nil
		}),
	}

	cfg, err := json.Marshal(map[string]string{"webhook_url": "http://example.com/hook"})
	if err != nil {
		t.Fatalf("failed to marshal config: %v", err)
	}

	notification := models.Notification{
		Type:   models.NotificationTypeDiscord,
		Name:   "Knocker",
		Config: cfg,
	}

	if err := SendWithClient(context.Background(), client, notification, "Title", "Description"); err != nil {
		t.Fatalf("SendWithClient returned error: %v", err)
	}

	if payload.Username != "Knocker" {
		t.Fatalf("expected username to be %q, got %q", "Knocker", payload.Username)
	}

	if len(payload.Embeds) != 1 {
		t.Fatalf("expected one embed, got %d", len(payload.Embeds))
	}

	if payload.Embeds[0].Title != "Title" || payload.Embeds[0].Description != "Description" {
		t.Fatalf("unexpected embed payload: %+v", payload.Embeds[0])
	}
}

func TestSendTelegram(t *testing.T) {
	var text string

	client := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Path != "/bottoken/sendMessage" {
				return nil, fmt.Errorf("unexpected path: %s", r.URL.Path)
			}

			var payload struct {
				ChatID string `json:"chat_id"`
				Text   string `json:"text"`
			}

			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				return nil, fmt.Errorf("failed to decode payload: %w", err)
			}

			text = payload.Text
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     make(http.Header),
			}, nil
		}),
	}

	oldBase := telegramAPIBase
	defer func() { telegramAPIBase = oldBase }()
	telegramAPIBase = "http://example.com"

	cfg, err := json.Marshal(map[string]string{
		"bot_token": "token",
		"chat_id":   "12345",
	})
	if err != nil {
		t.Fatalf("failed to marshal config: %v", err)
	}

	notification := models.Notification{
		Type:   models.NotificationTypeTelegram,
		Config: cfg,
	}

	if err := SendWithClient(context.Background(), client, notification, "Title", "Description"); err != nil {
		t.Fatalf("SendWithClient returned error: %v", err)
	}

	if !strings.Contains(text, "Title") || !strings.Contains(text, "Description") {
		t.Fatalf("unexpected telegram text: %q", text)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
