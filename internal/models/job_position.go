package models

// List of some permission the job position could have
type Permission struct {
	// Does the current job position is allowed to create a job position as child of himself?
	IsAllowCreateJP bool `json:"is_allow_create_jp"`
}

type JobPosotion struct {
	ID ID `json:"id"`
	// ID of the user the JP is for that.
	UserID      ID         `json:"user_id"`
	Title       string     `json:"title"`
	RegionID    ID         `json:"region_id"`
	RegionName  string     `json:"region_name"`
	ParentID    ID         `json:"parent_id"`
	ParentName  string     `json:"parent_name"`
	Permissions Permission `json:"permissions"`
}
