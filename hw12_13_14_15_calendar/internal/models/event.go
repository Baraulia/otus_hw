package models

import "time"

type Event struct {
	ID               string     `json:"id"`
	Header           string     `json:"header"`
	Description      string     `json:"description"`
	UserID           string     `json:"user_id"`
	EventTime        time.Time  `json:"event_time"`
	FinishEventTime  *time.Time `json:"finish_event_time,omitempty"`
	NotificationTime *time.Time `json:"notification_time,omitempty"`
}
