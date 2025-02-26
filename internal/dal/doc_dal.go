package dal

// Representation of a document entity
type Doc struct {
  ID        ID `db:"id"`
  CreatedBy ID `db:"created_by"`
  // The event the document is for that
  AtEvent ID     `db:"at_event"`
  Context string `db:"content"`
  // Contains path of multimedia files in the document. (If there's in the document)
  Paths []MediaPath `db:"path"`
}

type MediaPath struct {
  Type MediaType `db:"media_type"`
  // Source path of the file without filename part
  Src string `db:"src"`
  // Just contains filename and its type
  FileName string `db:"file_name"`
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

// It would be an implementation of docDAL interface for PostgreSQL
type psqlDocDAL struct{}

func NewPsqlDocDAL() *psqlDocDAL {
  return &psqlDocDAL{}
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
