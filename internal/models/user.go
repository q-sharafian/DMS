package models

const (
	UserDisabled = 0
	UserEnabled  = 1
)

type User struct {
	// ID of the user
	ID          ID     `json:"id"`
	Name        string `json:"name" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	// value 0 means it's not disabled and value 1 means it's disabled.
	IsDisabled uint8 `json:"is_disabled"`
	// The id of job position created this user
	CreatedBy ID `json:"created_by"`
}
