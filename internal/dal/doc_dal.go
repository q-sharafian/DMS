package dal

import (
	"DMS/internal/db"
	l "DMS/internal/logger"
	"DMS/internal/models"
)

type MediaPath struct {
	Type MediaType
	// Source path of the file without filename part
	Src string
	// Just contains filename and its type
	FileName string
}

type MediaType int8

const (
	Image MediaType = iota
	Video
	Audio
)

type DocDAL interface {
	// Create doc and return its id
	CreateDoc(doc *models.Doc) (models.ID, error)
	GetLastDocByUserID(id int) (*models.Doc, error)
	// Get latest created document of event with event_id by user_id. Then return that
	// document together with the name of event and user.
	GetLastEventDocByUserID(event_id models.ID, user_id models.ID) (doc *models.Doc, event_name string, user_name string, err error)
}

// It's an implementaion of DocDAL interface
type psqlDocDAL struct {
	db     *db.PSQLDB
	logger l.Logger
}

func newPsqlDocDAL(db *db.PSQLDB, logger l.Logger) *psqlDocDAL {
	return &psqlDocDAL{db, logger}
}

func (d *psqlDocDAL) CreateDoc(doc *models.Doc) (models.ID, error) {
	return models.NilID, nil
}

func (d *psqlDocDAL) GetLastDocByUserID(id int) (*models.Doc, error) {
	return nil, nil
}

func (d *psqlDocDAL) GetLastEventDocByUserID(event_id models.ID, user_id models.ID) (doc *models.Doc, event_name string, user_name string, err error) {
	return nil, "", "", nil
}
