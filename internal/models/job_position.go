package models

import "time"

type JobPosotion struct {
	ID ID `json:"id"`
	// ID of the user the JP is for that.
	UserID   ID     `json:"user_id"`
	Title    string `json:"title"`
	RegionID ID     `json:"region_id"`
	// ID of parent job position the current job position is child of that
	ParentID *ID `json:"parent_id"`
	// The time the JP is created
	CreatedAt time.Time `json:"at"`
}
