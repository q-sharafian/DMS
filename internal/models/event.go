package models

import (
	"time"
)

type Event struct {
	ID ID `json:"id"`
	// event name
	Name         string    `json:"name"`
	CreatedBy    string    `json:"created_by"`
	CreatedByID  ID        `json:"created_by_id"`
	At           time.Time `json:"at"`
	LastChangeAt time.Time `json:"last_change_at"`
	Description  string    `json:"description"`
}

type ApprovedEvent struct {
	EventID ID `json:"event_id"`
	// Name of Person approved the event
	ApprovedBy string    `json:"approved_by"`
	At         time.Time `json:"at"`
}
