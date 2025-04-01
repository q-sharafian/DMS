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
			SEDBError, ancestorID.ToString(), nodeID.ToString(), err.Error())
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
