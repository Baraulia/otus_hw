package models

import "time"

type Notification struct {
	ID          string
	EventHeader string
	EventTime   time.Time
	UserID      string
}
