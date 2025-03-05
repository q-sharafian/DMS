package services

import (
	"DMS/internal/dal"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
	"time"
)

type MediaPath struct {
	Type MediaType `json:"media_type"`
	// Full path and file name (contains type too)
	Src string `json:"src"`
	// Just contains filename and its type
	FileName string `json:"file_name"`
}

type MediaType uint8

const (
	Image MediaType = iota
	Video
	Audio
)

// Convert media to type compatible with MediaPath in the model package
func (media MediaPath) toModelMediaPath() m.MediaPath {
	return m.MediaPath{
		Type:     m.MediaType(media.Type),
		Src:      media.Src,
		FileName: media.FileName,
	}
}

// Convert media from MediaPath in the model package to the MediaPath in service packagee
func (MediaPath) fromModelMediaPath(media m.MediaPath) MediaPath {
	return MediaPath{
		Type:     MediaType(media.Type),
		Src:      media.Src,
		FileName: media.FileName,
	}
}

type DocService interface {
	// Create document for specified job position and event in the current time that contains input context and mediaPaths
	CreateDoc(eventID, JPID m.ID, context string, meediaPaths []MediaPath) (m.ID, *e.Error)
}

// It's a simple implementation of DocService interface.
// This implementation has minimum functionalities.
type sDocService struct {
	doc    dal.DocDAL
	logger l.Logger
}

// Possible error codese:
// DBError
func (s *sDocService) CreateDoc(eventID, JPID m.ID, context string, mediaPaths []MediaPath) (m.ID, *e.Error) {
	var paths []m.MediaPath
	for _, p := range mediaPaths {
		paths = append(paths, p.toModelMediaPath())
	}
	doc := m.Doc{
		CreatedBy: JPID,
		AtEvent:   eventID,
		Context:   context,
		Paths:     paths,
		At:        time.Now(),
	}
	eventID, err := s.doc.CreateDoc(&doc)
	if err != nil {
		return m.NilID, e.NewErrorP(err.Error(), DBError)
	}
	return eventID, nil
}

// Create an instance of sDocService struct
func newsDocService(doc dal.DocDAL, logger l.Logger) DocService {
	return &sDocService{
		doc,
		logger,
	}
}
