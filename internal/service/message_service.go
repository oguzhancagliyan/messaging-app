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
	"strconv"
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
	webhookURL = "https://webhook.site/a37c80f8-11d3-443b-bc32-23feb9359461"
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

		if len(msg.Content) > 160 {
			zap.L().Warn("Message content exceeds 160 characters",
				zap.Int64("id", msg.ID),
				zap.String("content", msg.Content),
				zap.Int("length", len(msg.Content)),
			)
			continue
		}

		zap.L().Info("Attempting to send message",
			zap.Int64("id", msg.ID),
			zap.String("to", msg.To),
			zap.String("content", msg.Content),
		)

		_, err := s.sendMessageToWebhook(msg)
		if err != nil {
			zap.L().Error("Message sending via webhook failed",
				zap.Int64("id", msg.ID),
				zap.Error(err),
				zap.String("to", msg.To),
				zap.String("content", msg.Content),
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

		if err := s.cache.SaveMessageID(strconv.FormatInt(msg.ID, 10), time.Now()); err != nil {
			zap.L().Warn("Failed to cache messageId in Redis",
				zap.String("messageId", strconv.FormatInt(msg.ID, 10)),
				zap.Error(err),
			)
		} else {
			zap.L().Info("MessageId cached in Redis",
				zap.String("messageId", strconv.FormatInt(msg.ID, 10)),
			)
		}
	}

	return nil
}

func (s *messageService) sendMessageToWebhook(msg model.Message) (bool, error) {
	payload := map[string]string{
		"to":      msg.To,
		"content": msg.Content,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		zap.L().Error("Failed to marshal payload",
			zap.String("to", msg.To),
			zap.String("content", msg.Content),
			zap.Error(err),
		)
		return false, errors.ErrMarshalPayload
	}

	// TODO: Use a more robust HTTP client with retry logic
	// and circuit breaker pattern for production code.
	// https://github.com/eapache/go-resiliency
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(data))
	if err != nil {
		zap.L().Error("Failed to create HTTP request",
			zap.String("url", webhookURL),
			zap.Error(err),
			zap.String("payload", string(data)),
		)
		return false, errors.ErrCreateHTTPRequest
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-ins-auth-key", authKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Error("Failed to send HTTP request",
			zap.String("url", webhookURL),
			zap.Error(err),
			zap.String("payload", string(data)),
		)
		return false, errors.ErrHTTPCall
	}
	defer resp.Body.Close()

	zap.L().Info("Webhook response received",
		zap.Int("status", resp.StatusCode),
		zap.String("url", webhookURL),
		zap.String("payload", string(data)),
		zap.String("response", resp.Status),
	)

	if resp.StatusCode != http.StatusOK {
		zap.L().Warn("Unexpected status code from webhook",
			zap.Int("status", resp.StatusCode),
			zap.String("url", webhookURL),
			zap.String("payload", string(data)),
			zap.String("response", resp.Status),
		)
		return false, errors.ErrWebhookFailed
	}

	return true, nil
}
