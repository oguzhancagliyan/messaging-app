package model

import "time"

type Message struct {
	ID        int64      `db:"id"`
	To        string     `db:"to"`
	Content   string     `db:"content"`
	Sent      bool       `db:"sent"`
	SentAt    *time.Time `db:"sent_at"`
	CreatedAt time.Time  `db:"created_at"`
}
