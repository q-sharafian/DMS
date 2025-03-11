package models

import (
	"time"
)

type Event struct {
	ID ID `json:"id"`
	// event name
	Name string `json:"name"`
	// ID of job position wants to create event
	CreatedBy   ID        `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Description string    `json:"description"`
}

type ApprovedEvent struct {
	EventID ID `json:"event_id"`
	// Name of Person approved the event
	ApprovedBy string    `json:"approved_by"`
	CreatedAt  time.Time `json:"at"`
}
