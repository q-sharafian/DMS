package services

import (
	"DMS/internal/dal"
	e "DMS/internal/error"
	"DMS/internal/hierarchy"
	l "DMS/internal/logger"
	m "DMS/internal/models"
)

// Contains interface for all functionalities related to permissions and hierarchy tree.
type AuthorizationService interface {
	// Check if there's a path from ancestorID to nodeID.
	//
	// Possible error codes:
	// SEDBError
	IsAncestor(ancestorID, nodeID m.ID) (bool, *e.Error)
	// List of permissions of the given job position.
	//
	// Possible error codes:
	// SEDBError
	GetJPPermissions(jpID m.ID) (*m.Permission, *e.Error)
	// List of all nested child job positions of the given job position.
	//
	// Possible error codes:
	// SEDBError
	GetNestedChilds(jpID m.ID) ([]m.ID, *e.Error)
	// Return true if the given job position is an admin job position.
	//
	// Possible error codes:
	// SEDBError
	IsAdminJP(jpID m.ID) (bool, *e.Error)
}

// It's a simple implementation of AuthorizationService interface.
// This implementation has minimum functionalities.
type sAuthorizationService struct {
	hierarchy  hierarchy.HierarchyTree
	permission dal.PermissionDAL
	logger     l.Logger
}

// Create a new simple authorization service
func newSAuthorizationService(hierarchy hierarchy.HierarchyTree, permission dal.PermissionDAL,
	logger l.Logger) AuthorizationService {
	sPermission := &sAuthorizationService{
		hierarchy,
		permission,
		logger,
	}
	return sPermission
}

func (s *sAuthorizationService) IsAncestor(ancestorID, nodeID m.ID) (bool, *e.Error) {
	isAncestor, err := s.hierarchy.IsAncestor(id2Vertex(ancestorID), id2Vertex(nodeID))
	if err != nil {
		return false, e.NewErrorP("failed to check if ancestor id %s is an ancestor of node id %s: %s",
			SEDBError, ancestorID.String(), nodeID.String(), err.Error())
	}
	return isAncestor, nil
}

func (s *sAuthorizationService) GetJPPermissions(jpID m.ID) (*m.Permission, *e.Error) {
	permission, err := s.permission.GetPermissionsByJPID(jpID)
	if err != nil {
		return nil, e.NewErrorP(err.Error(), SEDBError)
	}
	return permission, nil
}

func (s *sAuthorizationService) GetNestedChilds(jpID m.ID) ([]m.ID, *e.Error) {
	nestedChilds, err := s.hierarchy.GetNestedChilds(id2Vertex(jpID))
	if err != nil {
		return nil, e.NewErrorP(err.Error(), SEDBError)
	}

	childs := make([]m.ID, 0)
	for _, vertex := range nestedChilds {
		id, err := vertex2ID(vertex)
		if err != nil {
			return nil, e.NewErrorP("failed to convert vertex %s to id: %s", SEDBError, vertex.String(), err.Error())
		}
		if id != m.NilID {
			childs = append(childs, id)
		}
	}
	return childs, nil
}

func (s *sAuthorizationService) IsAdminJP(jpID m.ID) (bool, *e.Error) {
	result, err := s.hierarchy.IsSourceVertex(id2Vertex(jpID))
	if err != nil {
		return false, e.NewErrorP("failed to check if job position id %s is admin: %s",
			SEDBError, jpID.String(), err.Error())
	}
	return result, nil
}
