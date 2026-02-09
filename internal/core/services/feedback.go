package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shanth1/gotrace/internal/core/domain"
)

type FeedbackServiceImpl struct {
	notifyURL  string
	notifyKey  string
	httpClient *http.Client
}

func NewFeedbackService(notifyURL, notifyKey string) *FeedbackServiceImpl {
	return &FeedbackServiceImpl{
		notifyURL: notifyURL,
		notifyKey: notifyKey,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type notificationPayload struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func (s *FeedbackServiceImpl) SendFeedback(ctx context.Context, cmd domain.SendFeedbackCmd) error {
	msgText := fmt.Sprintf(
		"From: %s\nEmail: %s\nUser ID: %s\n\nContent:\n%s",
		cmd.Name,
		cmd.Email,
		cmd.UserID,
		cmd.Message,
	)

	payload := notificationPayload{
		Title:   "New Feedback Request",
		Message: msgText,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal feedback payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.notifyURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", s.notifyKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send feedback to bot service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("bot service returned error: status %d", resp.StatusCode)
	}

	return nil
}
