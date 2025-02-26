package user

type User struct {
  PhoneNumber string
  Name        string
}

// Creates a new user as the user's child in the hierarchy tree. If be allowed.
func CreateUser(user User) (User, error) {
  return user, nil
}
