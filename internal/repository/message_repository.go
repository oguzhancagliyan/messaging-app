package repository

import (
	"context"
	"database/sql"
	"errors"
	"messaging-app/internal/model"
	"time"

	"go.uber.org/zap"
)

type MessageRepository interface {
	GetUnsentMessages(ctx context.Context, limit int) ([]model.Message, error)
	MarkAsSent(ctx context.Context, id int64) error
	GetSentMessages(ctx context.Context, limit int) ([]model.Message, error)
}

type messageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *messageRepository {
	return &messageRepository{
		db: db,
	}
}

func (r *messageRepository) GetSentMessages(ctx context.Context, limit int) ([]model.Message, error) {
	query := `
        SELECT id, "to", content, sent, sent_at, created_at
        FROM messages
        WHERE sent = TRUE
        ORDER BY sent_at DESC
        LIMIT $1
    `

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		zap.L().Error("Failed to fetch sent messages from DB", zap.Error(err), zap.Int("limit", limit))
		return nil, err
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var msg model.Message
		if err := rows.Scan(&msg.ID, &msg.To, &msg.Content, &msg.Sent, &msg.SentAt, &msg.CreatedAt); err != nil {
			zap.L().Error("Failed to scan sent message row", zap.Error(err))
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		zap.L().Error("Error during row iteration", zap.Error(err))
		return nil, err
	}

	return messages, nil
}

func (r *messageRepository) GetUnsentMessages(ctx context.Context, limit int) ([]model.Message, error) {

	query := `
    SELECT id, "to", content, sent, sent_at, created_at
    FROM messages
    WHERE sent = FALSE
    ORDER BY sent_at DESC
    LIMIT $1
`
	rows, err := r.db.QueryContext(ctx, query, limit)

	if err != nil {
		zap.L().Error("MessageRepository-GetUnsentMessages got error", zap.Error(err), zap.Int("limit", limit))
		return nil, err
	}

	defer rows.Close()

	var messages []model.Message

	for rows.Next() {
		var msg model.Message
		err := rows.Scan(&msg.ID, &msg.To, &msg.Content, &msg.Sent, &msg.SentAt, &msg.CreatedAt)
		if err != nil {
			zap.L().Error("MessageRepository-GetUnsentMessages Scanning row got error", zap.Error(err), zap.Int("limit", limit))
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		zap.L().Error("MessageRepository-GetUnsentMessages Error for recieving from database", zap.Error(err))
		return nil, err
	}

	return messages, nil
}

func (r *messageRepository) MarkAsSent(ctx context.Context, id int64) error {
	query := `UPDATE messages SET sent = TRUE, sent_at = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		zap.L().Error("MessageRepository-MarkAsSent Message could not be updated", zap.Error(err), zap.Int64("id", id))
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		zap.L().Error("MessageRepository-MarkAsSent RowsAffected alınamadı", zap.Error(err), zap.Int64("id", id))
		return err
	}

	if affected == 0 {
		//TODO: add this error to custom error list
		err = errors.New("hiçbir kayıt güncellenmedi")
		zap.L().Warn("MessageRepository-MarkAsSent message could not be found and updated", zap.Int64("id", id))
		return err
	}

	zap.L().Info("Message marked as sent", zap.Int64("id", id))
	return nil
}
