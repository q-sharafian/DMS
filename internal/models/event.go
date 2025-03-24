package models

import (
	"time"
)

type Event struct {
	ID ID `json:"id" example:"76a79030f-0685-49d1-bbdd-31ab1b4c1613"`
	// event name
	Name string `json:"name" validate:"required"`
	// ID of job position wants to create event
	CreatedBy ID `json:"created_by" example:"6a79030f-0685-49d1-bbdd-31ab1b4c1613"`
	// Date when the event is created. Based on UTC time zone and Unix timestamp. (In seconds)
	CreatedAt int64 `json:"created_at"`
	// Date when the event is updated. Based on UTC time zone and Unix timestamp. (In seconds)
	// If it is nil, means the event is not updated.
	UpdatedAt   *int64 `json:"updated_at"`
	Description string `json:"description"`
}

type ApprovedEvent struct {
	EventID ID `json:"event_id"`
	// Name of Person approved the event
	ApprovedBy string    `json:"approved_by"`
	CreatedAt  time.Time `json:"at"`
}
