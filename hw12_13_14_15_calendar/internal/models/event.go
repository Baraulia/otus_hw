package models

import "time"

type Event struct {
	ID               string
	Header           string
	Description      string
	UserID           string
	EventTime        time.Time
	FinishEventTime  time.Time
	NotificationTime time.Time
}
