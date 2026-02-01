package services

import (
	"context"
	"fmt"

	"github.com/shanth1/gotools/notify"
	"github.com/shanth1/gotrace/internal/core/domain"
)

type FeedbackServiceImpl struct {
	notifier    notify.Notifier
	adminChatID string
}

func NewFeedbackService(notifier notify.Notifier, adminChatID string) *FeedbackServiceImpl {
	return &FeedbackServiceImpl{
		notifier:    notifier,
		adminChatID: adminChatID,
	}
}

func (s *FeedbackServiceImpl) SendFeedback(ctx context.Context, cmd domain.SendFeedbackCmd) error {
	text := fmt.Sprintf(
		"📩 *New Feedback Request*\n\n"+
			"*From:* %s\n"+
			"*Email:* %s\n"+
			"*User ID:* %s\n\n"+
			"*Message:*\n%s",
		cmd.Name,
		cmd.Email,
		cmd.UserID,
		cmd.Message,
	)

	msg := notify.Message{
		Subject: "New Feedback",
		Text:    text,
	}

	if err := s.notifier.Send(ctx, s.adminChatID, msg); err != nil {
		return fmt.Errorf("failed to send feedback notification: %w", err)
	}

	return nil
}
