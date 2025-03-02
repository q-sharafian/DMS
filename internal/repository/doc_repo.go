package repository

import (
	"DMS/internal/dal"
	"DMS/internal/models"
	"path/filepath"
)

// An implementation of document operations that could be used by the application
type DocRepo struct {
	docDAL dal.DocDAL
}

func newDocRepo(docDAL dal.DocDAL) DocRepo {
	return DocRepo{
		docDAL,
	}
}

// Convert media path (from model) such that be suitable for DAL.
func convertMediaPathsModelToDAL(paths []models.MediaPath) []dal.MediaPath {
	var dalPaths []dal.MediaPath
	for _, path := range paths {
		dalPaths = append(dalPaths, dal.MediaPath{
			Type:     dal.MediaType(path.Type),
			Src:      path.Src[:len(path.Src)-len(path.FileName)],
			FileName: path.FileName,
		})
	}
	return dalPaths
}

// Convert media path (from DAL) such that be suitable for the model package.
func convertMediaPathsDALToModel(paths []dal.MediaPath) []models.MediaPath {
	var modelPaths []models.MediaPath
	for _, path := range paths {
		modelPaths = append(modelPaths, models.MediaPath{
			Type:     models.MediaType(path.Type),
			Src:      filepath.Join(path.Src, path.FileName),
			FileName: path.FileName,
		})
	}
	return modelPaths
}

// Create a new document for the `event` have event_id.
func (r *DocRepo) CreateDoc(doc *models.Doc, event_id models.ID) error {
	var err = r.docDAL.CreateDoc(&dal.Doc{
		ID:        toDALID(doc.ID),
		CreatedBy: toDALID(doc.CreatedByID),
		AtEvent:   toDALID(event_id),
		Context:   doc.Context,
		Paths:     convertMediaPathsModelToDAL(doc.Paths),
	})
	return err
}

// Get latest created document of event withevent_id by user_id
func (r *DocRepo) GetLastEventDocByUserID(event_id models.ID, user_id models.ID) (*models.Doc, error) {
	var doc, event_name, user_name, err = r.docDAL.GetLastEventDocByUserID(toDALID(event_id), toDALID(user_id))
	if err != nil {
		return nil, err
	}
	return &models.Doc{
		ID:        toModelID(doc.ID),
		CreatedBy: user_name,
		AtEvent:   event_name,
		Context:   doc.Context,
		Paths:     convertMediaPathsDALToModel(doc.Paths),
	}, nil
}
