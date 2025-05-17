package service

import (
	"bytes"
	"context"
	"encoding/json"
	"messaging-app/internal/cache"
	"messaging-app/internal/errors"
	"messaging-app/internal/model"
	"messaging-app/internal/repository"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type MessageService interface {
	SendUnsentMessages(ctx context.Context) error
}

type messageService struct {
	repo  repository.MessageRepository
	cache cache.Cache
}

func NewMessageService(repo repository.MessageRepository, cache cache.Cache) *messageService {
	return &messageService{repo: repo, cache: cache}
}

const (
	webhookURL = "https://webhook.site/c3f13233-1ed4-429e-9649-8133b3b9c9cd"
	authKey    = "INS.me1x9uMcyYGlhKKQVPoc.bO3j9aZwRTOcA2Ywo"
)

func (s *messageService) SendUnsentMessages(ctx context.Context) error {
	messages, err := s.repo.GetUnsentMessages(ctx, 2)
	if err != nil {
		zap.L().Error("Failed to fetch unsent messages", zap.Error(err))
		return errors.ErrFetchMessages
	}

	if len(messages) == 0 {
		zap.L().Info("No messages found to send")
		return errors.ErrNoUnsentMessages
	}

	for _, msg := range messages {
		zap.L().Info("Attempting to send message",
			zap.Int64("id", msg.ID),
			zap.String("to", msg.To),
			zap.String("content", msg.Content),
		)

		messageId, err := s.sendMessageToWebhook(msg)
		if err != nil {
			zap.L().Error("Message sending via webhook failed",
				zap.Int64("id", msg.ID),
				zap.Error(err),
			)
			continue
		}

		if err := s.repo.MarkAsSent(ctx, msg.ID); err != nil {
			zap.L().Error("Failed to mark message as sent in DB",
				zap.Int64("id", msg.ID),
				zap.Error(err),
			)
			continue
		}

		if err := s.cache.SaveMessageID(messageId, time.Now()); err != nil {
			zap.L().Warn("Failed to cache messageId in Redis",
				zap.String("messageId", messageId),
				zap.Error(err),
			)
		} else {
			zap.L().Info("MessageId cached in Redis",
				zap.String("messageId", messageId),
			)
		}
	}

	return nil
}

func (s *messageService) sendMessageToWebhook(msg model.Message) (string, error) {
	payload := map[string]string{
		"to":      msg.To,
		"content": msg.Content,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", errors.ErrMarshalPayload
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(data))
	if err != nil {
		return "", errors.ErrCreateHTTPRequest
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-ins-auth-key", authKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.ErrHTTPCall
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		zap.L().Warn("Unexpected status code from webhook",
			zap.Int("status", resp.StatusCode),
		)
		return "", errors.ErrWebhookFailed
	}

	var result struct {
		Message   string `json:"message"`
		MessageID string `json:"messageId"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", errors.ErrInvalidWebhookResponse
	}

	return result.MessageID, nil
}
