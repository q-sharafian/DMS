package db

import (
  "time"
)

type ID uint64

type User struct {
  ID          ID
  Name        string
  PhoneNumber string
  // value 0 means it's not disabled and value 1 means it's disabled.
  IsDisabled uint8
  // The user created the current user
  CreatedByID ID `json:"created_by_id"`
}

type Event struct {
  ID ID
  // event name
  Name         string
  CreatedByID  ID
  At           time.Time
  UpdatedAt time.Time
  Description  string
}
