package models

import "time"

type Event struct {
	ID               string     `json:"id"`
	Header           string     `json:"header"`
	Description      string     `json:"description"`
	UserID           string     `json:"userId"`
	EventTime        time.Time  `json:"eventTime"`
	FinishEventTime  *time.Time `json:"finishEventTime,omitempty"`
	NotificationTime *time.Time `json:"notificationTime,omitempty"`
}
