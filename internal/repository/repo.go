package repository

import (
  "DMS/internal/dal"
  "DMS/internal/models"
)

// Mapping given id (from models package) to the ID type that is acceptable by the DAL
func toDALID(id models.ID) dal.ID {
  return dal.ID(id)
}

// Mapping given id (from DAL package) to the ID type that is acceptable by the Model
func toModelID(id dal.ID) models.ID {
  return models.ID(id)
}

type Repository struct {
  // It's local
  dal dal.DAL
  // Contains operations related with user entity
  User userRepo
  // Contains operations related with document entity
  Doc docRepo
  // Contains operations related with event entity
  Event eventRepo
}

// Create New repository that represent operations on data entities in database.
func NewRepository(dal dal.DAL) *Repository {
  return &Repository{
    dal:   dal,
    Doc:   newDocRepo(dal.Doc),
    User:  newUserRepo(dal.User),
    Event: newEventRepo(dal.Event),
  }
}
