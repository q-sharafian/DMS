package dal

import (
	"DMS/internal/db"
	l "DMS/internal/logger"
	m "DMS/internal/models"
)

type PermissionDAL interface {
	// Get list of permissions of the given job position.
	GetPermissionsByJPID(jpID m.ID) (*m.Permission, error)
	// Get some records of hierarchy tree. Each record contains all childs of a job position.
	// For each job position we must have one corresponding record in the tree.
	// For example if you set offset to 20 and limit to 10, you'll retrieve records 21 through 30.
	//
	// TODO: What happens if new data is imported or removed during data retrieval?
	// It may happen that some records won't be read and some will be read multiple times.
	GetSomeHierarchyTree(offset, limit int) ([]m.HierarchyTree, error)
	// Get all records of hierarchy tree and format them as adjacency list.
	GetHierarchyTree() (m.Graph, error)
	GetHierarchyTreeCount() (uint64, error)
}

type psqlPermissionDAL struct {
	db     *db.PSQLDB
	logger l.Logger
}

func newPsqlPermissionDAL(db *db.PSQLDB, logger l.Logger) *psqlPermissionDAL {
	return &psqlPermissionDAL{db, logger}
}

func (d *psqlPermissionDAL) GetPermissionsByJPID(jpID m.ID) (*m.Permission, error) {
	panic("GetPermissionsByJPID doesn't implemented")
}

func (d *psqlPermissionDAL) GetSomeHierarchyTree(offset, limit int) ([]m.HierarchyTree, error) {
	var hierarchyTrees *[]db.HierarchyTree
	result := d.db.Limit(int(limit)).Offset(int(offset)).Find(&hierarchyTrees)

	if result.Error != nil {
		d.logger.Debugf("Failed to get hierarchy tree. offset: %d, limit: %d- (%s)", offset, limit, result.Error.Error())
		return nil, result.Error
	}
	return *dbHierarchyTree2ModelHierarchyTree(hierarchyTrees), nil
}

func (d *psqlPermissionDAL) GetHierarchyTreeCount() (uint64, error) {
	var count int64
	result := d.db.Model(&db.HierarchyTree{}).Count(&count)
	if result.Error != nil {
		d.logger.Debugf("Failed to get hierarchy tree count (%s)", result.Error.Error())
		return 0, result.Error
	}
	return uint64(count), nil
}

func (d *psqlPermissionDAL) GetHierarchyTree() (m.Graph, error) {
	count, err := d.GetHierarchyTreeCount()
	if err != nil {
		d.logger.Debugf("Failed to get hierarchy tree count (%s)", err.Error())
		return nil, err
	}
	graph := make(m.Graph)
	// TODO: Optimize
	limit := 10

	for i := 0; i < (int(count)/limit)+1; i += 1 {
		hierarchyTree, err := d.GetSomeHierarchyTree(i, limit)
		if err != nil {
			d.logger.Debugf("Failed to get hierarchy tree elements from %d with limit %d (%s)", i, limit, err.Error())
			return nil, err
		}
		// for _, element := range *hierarchyTree {
		// 	graph[element.JPID] = element.JPChildsID,
		// }
		for _, element := range hierarchyTree {
			graph[element.JPID] = element.JPChildsID
		}
	}
	return graph, nil
}

func dbHierarchyTree2ModelHierarchyTree(hierarchyTrees *[]db.HierarchyTree) *[]m.HierarchyTree {
	var result []m.HierarchyTree
	for _, hierarchyTree := range *hierarchyTrees {
		result = append(result, m.HierarchyTree{
			ID:         *dbID2ModelID(&hierarchyTree.ID),
			JPID:       *dbID2ModelID(&hierarchyTree.JPID),
			JPChildsID: dbID2ModelIDSlice(hierarchyTree.ChildJPsID),
		})
	}
	return &result
}
