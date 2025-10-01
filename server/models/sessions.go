package models

import "time"

type Session struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	ClosedAt  time.Time
}
