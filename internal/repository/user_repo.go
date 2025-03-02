package repository

import (
	"DMS/internal/dal"
	"DMS/internal/models"
)

type UserRepo struct {
	userDAL dal.UserDAL
}

func newUserRepo(userDAL dal.UserDAL) UserRepo {
	return UserRepo{
		userDAL: userDAL,
	}
}

func (r *UserRepo) CreateUser(name, phoneNumber string) error {
	return r.userDAL.CreateUser(name, phoneNumber)
}

func (r *UserRepo) GetUserByID(id int) (*models.User, error) {
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
