package dal

import "DMS/internal/db"

// Representation of a document entity
type Doc struct {
	ID        ID
	CreatedBy ID
	// The event the document is for that
	AtEvent ID
	Context string
	// Contains path of multimedia files in the document. (If there's in the document)
	Paths []MediaPath
}

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
	CreateDoc(doc *Doc) error
	GetLastDocByUserID(id int) (*Doc, error)
	// Get latest created document of event with event_id by user_id. Then return that
	// document together with the name of event and user.
	GetLastEventDocByUserID(event_id ID, user_id ID) (doc *Doc, event_name string, user_name string, err error)
}

// It's an implementaion of DocDAL interface
type psqlDocDAL struct {
	db *db.PSQLDB
}

func newPsqlDocDAL(db *db.PSQLDB) *psqlDocDAL {
	return &psqlDocDAL{db}
}

func (d *psqlDocDAL) CreateDoc(doc *Doc) error {
	return nil
}

func (d *psqlDocDAL) GetLastDocByUserID(id int) (*Doc, error) {
	return nil, nil
}

func (d *psqlDocDAL) GetLastEventDocByUserID(event_id ID, user_id ID) (doc *Doc, event_name string, user_name string, err error) {
	return nil, "", "", nil
}
