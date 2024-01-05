package models

import "time"

type Event struct {
	ID                string
	Header            string
	EventTime         time.Time
	EventDuration     time.Duration
	Description       string
	UserID            string
	InAdvanceDuration time.Duration
}
