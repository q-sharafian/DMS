package dal

import (
	"DMS/internal/db"
	l "DMS/internal/logger"
	m "DMS/internal/models"
)

type DocDAL interface {
	// Create doc and return its id
	CreateDoc(doc *m.Doc) (*m.ID, error)
	// Get n "last" docs by the event id.
	GetNLastDocByEventID(eventID m.ID, n int) (*[]m.Doc, error)
	// Get latest created document of event with event_id by user_id. Then return that
	// document together with the name of event and user.
	GetLastEventDocByUserID(event_id m.ID, user_id m.ID) (doc *m.Doc, event_name string, user_name string, err error)
}

const (
	cacheKeyUserByID = "doc:id:%s"
)

// It's an implementaion of DocDAL interface
type psqlDocDAL struct {
	cache  *cache
	db     *db.PSQLDB
	logger l.Logger
}

func newPsqlDocDAL(db *db.PSQLDB, cache *cache, logger l.Logger) DocDAL {
	return &psqlDocDAL{cache, db, logger}
}

func (d *psqlDocDAL) CreateDoc(doc *m.Doc) (*m.ID, error) {
	newDoc := db.Doc{
		CreatedByID: *modelID2DBID(&doc.CreatedBy),
		EventID:     *modelID2DBID(&doc.EventID),
		Context:     doc.Context,
		Multimedia:  modelMultimedias2DBMultimedias(&doc.Paths, d.logger),
	}
	result := d.db.Create(&newDoc)

	if result.Error != nil {
		d.logger.Debugf("Failed to create doc for user-id %s (%s)", newDoc.CreatedByID.ToString(), result.Error.Error())
		return nil, result.Error
	}
	return dbID2ModelID(&newDoc.ID), nil
}

// If n be equals to -1, then return all docs
func (d *psqlDocDAL) GetNLastDocByEventID(eventID m.ID, n int) (*[]m.Doc, error) {
	var docs *[]db.Doc
	result := d.db.Order("created_at desc").Limit(n).Where(
		&db.Doc{EventID: *modelID2DBID(&eventID)},
	).Find(&docs)

	if result.Error != nil {
		d.logger.Debugf("Failed to get last %d docs in the  event-id %s (%s)", n, eventID.String(), result.Error.Error())
		return nil, result.Error
	}

	return dbDocs2modelDocs(docs, d.logger), nil
}

func (d *psqlDocDAL) GetLastEventDocByUserID(event_id m.ID, user_id m.ID) (doc *m.Doc, event_name string, user_name string, err error) {
	return nil, "", "", nil
}

func modelMediaType2DBMediaType(media m.MediaType, logger l.Logger) db.MediaType {
	switch media {
	case m.MediaImage:
		return db.MediaImage
	case m.MediaVideo:
		return db.MediaVideo
	case m.MediaAudio:
		return db.MediaAudio
	default:
		logger.Panicf("Unknown media type in the model: %s", media)
	}
	return db.MediaImage
}

func dbMediaType2ModelMediaType(media db.MediaType, logger l.Logger) m.MediaType {
	switch media {
	case db.MediaImage:
		return m.MediaImage
	case db.MediaVideo:
		return m.MediaVideo
	case db.MediaAudio:
		return m.MediaAudio
	default:
		logger.Panicf("Unknown media type in the database: %s", media)
	}
	return m.MediaImage
}

func modelMultimedia2DBMultimedia(m *m.MediaPath, logger l.Logger) *db.Multimedia {
	return &db.Multimedia{
		Type:     modelMediaType2DBMediaType(m.Type, logger),
		Src:      m.Src,
		FileName: m.FileName,
	}
}

func dbMultimedia2ModelMultimedia(mum *db.Multimedia, logger l.Logger) *m.MediaPath {
	return &m.MediaPath{
		Type:     dbMediaType2ModelMediaType(mum.Type, logger),
		Src:      mum.Src,
		FileName: mum.FileName,
	}
}

func modelMultimedias2DBMultimedias(m *[]m.MediaPath, logger l.Logger) *[]db.Multimedia {
	var dbMultimedias []db.Multimedia
	for _, media := range *m {
		dbMultimedias = append(dbMultimedias, *modelMultimedia2DBMultimedia(&media, logger))
	}
	return &dbMultimedias
}

// If it be nil, return empty slice
func dbMultimedias2ModelMultimedias(mum *[]db.Multimedia, logger l.Logger) *[]m.MediaPath {
	if mum == nil {
		return &[]m.MediaPath{}
	}
	var modelMultimedias []m.MediaPath
	for _, media := range *mum {
		modelMultimedias = append(modelMultimedias, *dbMultimedia2ModelMultimedia(&media, logger))
	}
	return &modelMultimedias
}

func dbDoc2modelDoc(doc *db.Doc, logger l.Logger) *m.Doc {
	return &m.Doc{
		ID:        *dbID2ModelID(&doc.ID),
		CreatedBy: *dbID2ModelID(&doc.CreatedByID),
		EventID:   *dbID2ModelID(&doc.EventID),
		Context:   doc.Context,
		Paths:     *dbMultimedias2ModelMultimedias(doc.Multimedia, logger),
	}
}

func dbDocs2modelDocs(docs *[]db.Doc, logger l.Logger) *[]m.Doc {
	var modelDocs []m.Doc
	for _, doc := range *docs {
		modelDocs = append(modelDocs, *dbDoc2modelDoc(&doc, logger))
	}
	return &modelDocs
}
