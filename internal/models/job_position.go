package models

import "encoding/json"

type CommonJobPosition struct {
	ID ID `json:"id" example:"6a79030f-0685-49d1-bbdd-31ab1b4c1613"`
	// ID of the user the JP is for that.
	UserID ID     `json:"user_id" example:"6a79030f-0685-49d1-bbdd-31ab1b4c1613"`
	Title  string `json:"title" example:"معاون مدرسه" validate:"required"`
	// The region the JP belongs to
	RegionID ID `json:"region_id" example:"b11c9be1-b619-4ef5-be1b-a1cd9ef265b7"`
	// The time the JP is created with UTC timezone and unix timestamp in seconds.
	CreatedAt int64 `json:"created_at" example:"1641011200"`
}
type UserJobPosition struct {
	CommonJobPosition
	ParentID ID `json:"parent_id,omitempty" validate:"required" example:"5abcdeff-0685-49d1-bbdd-31ab1b4c1613"`
}

func (s UserJobPosition) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID        string `json:"id"`
		UserID    string `json:"user_id"`
		RegionID  string `json:"region_id"`
		ParentID  string `json:"parent_id"`
		Title     string `json:"title"`
		CreatedAt int64  `json:"created_at"`
	}{
		ParentID:  s.ParentID.ToString(),
		ID:        s.ID.ToString(),
		UserID:    s.UserID.ToString(),
		RegionID:  s.RegionID.ToString(),
		Title:     s.Title,
		CreatedAt: s.CreatedAt,
	})
}

type UserJPWithPermission struct {
	JobPosition UserJobPosition `json:"job_position" validate:"required"`
	Permission  Permission      `json:"permission" validate:"required"`
}

// Admin job position parent id is not required. Means admin jps has no parent.
type AdminJobPosition struct {
	CommonJobPosition
}

type AdminJPWithPermission struct {
	JobPosition AdminJobPosition `json:"job_position" validate:"required"`
	Permission  Permission       `json:"permission" validate:"required"`
}
