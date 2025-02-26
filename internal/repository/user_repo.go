package repository

import (
  "DMS/internal/dal"
  "DMS/internal/models"
)

type userRepo struct {
  userDAL dal.UserDAL
}

func newUserRepo(userDAL dal.UserDAL) userRepo {
  return userRepo{
    userDAL: userDAL,
  }
}

func (r *userRepo) CreateUser(name, phoneNumber string) error {
  return r.userDAL.CreateUser(name, phoneNumber)
}

func (r *userRepo) GetUserByID(id int) (*models.User, error) {
  var user, _ = r.userDAL.GetUserByID(id)
  return &models.User{
    ID:          toModelID(user.ID),
    Name:        user.Name,
    PhoneNumber: user.PhoneNumber,
    IsDisabled:  user.IsDisabled,
    // TODO: Edit it
    CreatedBy: user.CreatedByID.ToString(),
  }, nil
}
