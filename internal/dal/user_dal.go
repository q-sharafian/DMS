package dal

// Representation of user entity
type user struct {
  ID          ID     `db:"id"`
  Name        string `db:"name"`
  PhoneNumber string `db:"phone_number"`
  // value 0 means it's not disabled and value 1 means it's disabled.
  IsDisabled  uint8 `db:"is_disabled"`
  CreatedByID ID    `db:"created_by_id"`
}

type UserDAL interface {
  CreateUser(name, phoneNumber string) error
  GetUserByID(id int) (*user, error)
}

// It would be a Postgres implementation of userDAL
type psqlUserDAL struct{}

func NewPsqlUserDAL() *psqlUserDAL {
  return &psqlUserDAL{}
}

func (d *psqlUserDAL) CreateUser(name, phoneNumber string) error {
  return nil
}

func (d *psqlUserDAL) GetUserByID(id int) (*user, error) {
  return nil, nil
}
