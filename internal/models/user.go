package models

type User struct {
  ID          ID     `json:"id"`
  Name        string `json:"name"`
  PhoneNumber string `json:"phone_number"`
  // value 0 means it's not disabled and value 1 means it's disabled.
  IsDisabled  uint8  `json:"is_disabled"`
  CreatedBy   string `json:"created_by"`
  CreatedByID ID     `json:"created_by_id"`
}
