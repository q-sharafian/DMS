package models

const (
	UserDisabled = 0
	UserEnabled  = 1
)

type User struct {
	// ID of the user
	ID          ID          `json:"id"`
	Name        string      `json:"name" validate:"required"`
	PhoneNumber PhoneNumber `json:"phone_number" validate:"required"`
	IsDisabled  Disability  `json:"is_disabled"`
	// The id of job position created this user
	CreatedBy ID `json:"created_by"`
}

type PhoneNumber string

func (p PhoneNumber) ToString() string {
	return string(p)
}

// If the input phone was nil, return nil error and set phone to be nil.
func (p *PhoneNumber) FromString(s *string) error {
	if s == nil {
		p = nil
		return nil
	}
	*p = PhoneNumber(*s)
	return nil
}
