package model

import "time"

type WhatsappSession struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	ClosedAt  time.Time
}
