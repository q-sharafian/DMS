package models

const (
	UserDisabled = 0
	UserEnabled  = 1
)

type User struct {
	// ID of the user
	ID          ID          `json:"id" validate:"uuidv4" example:"6a79030f-0685-49d1-bbdd-31ab1b4c1613"`
	Name        string      `json:"name" validate:"required" example:"John Doe"`
	PhoneNumber PhoneNumber `json:"phone_number" validate:"required" example:"9171234567"`
	// * 0: enabled user
	// * 1: disabled user
	IsDisabled Disability `json:"is_disabled" example:"0" enum:"0,1"`
	// The id of job position created this user
	CreatedBy *ID `json:"created_by" validate:"uuidv4" example:"6a79030f-0685-49d1-bbdd-31ab1b4c1613"`
}

type AdminUser struct {
	// ID of the admin user
	ID          ID          `json:"id" validate:"uuidv4" example:"6a79030f-0685-49d1-bbdd-31ab1b4c1613"`
	Name        string      `json:"name" validate:"required" example:"John Doe"`
	PhoneNumber PhoneNumber `json:"phone_number" validate:"required" example:"9171234567"`
	// * 0: enabled user
	// * 1: disabled user
	IsDisabled Disability `json:"is_disabled" example:"0" enums:"0,1"`
}

type PhoneNumber string

var NilPhone PhoneNumber = ""

// Return true if the phone is nil. (means phone is zero value)
func (p PhoneNumber) IsNil() bool {
	return p == NilPhone
}

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
