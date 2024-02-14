package models

import "time"

type Notification struct {
	ID          string    `json:"eventId"`
	EventHeader string    `json:"eventHeader"`
	EventTime   time.Time `json:"eventTime"`
	UserID      string    `json:"userId"`
}
