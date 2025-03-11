package services

import (
	"DMS/internal/dal"
	e "DMS/internal/error"
	l "DMS/internal/logger"
	m "DMS/internal/models"
)

// Contains interface for all functionalities related to permissions.
type PermissionService interface {
	// Return true if ancestorID is ancestor of nodeID in the hierarchy tree. Being an
	// ancestor means there is a path from the ancestor to the specified node.
	IsAncestor(ancestorID, nodeID m.ID) bool
	// In a graph, transitive closure means which nodes could be reachable from each other.
	// Even if they're not connected by a single edge.
	// Return a matrix. If there's a path from i to the j, then matrix[i][j] = 1. else 0.
	GetTransitiveClosure() m.Matrix
}

// It's a simple implementation of PermissionService interface.
// This implementation has minimum functionalities.
type sPermissionService struct {
	permission dal.PermissionDAL
	logger     l.Logger
	// Adjacency matrix represents the path between each two job positions. If there's
	// a path from i to j, then matrix[i][j] >= 1 and it means job position i is an
	// ancestor of job position j. Else it would be 0.
	jpPaths m.Matrix
	// Because of the id of job positions are not consecutive, we need to map them to
	// consecutive integers.
	mappedJPID2Int map[m.ID]int
}

// Create a new simple permission service and find transitive closure and mapping
// job position ids to consecutive integers. If occured error, panic.
func newSPermissionService(permission dal.PermissionDAL, logger l.Logger) PermissionService {
	sPermission := &sPermissionService{
		permission,
		logger,
		nil,
		nil,
	}
	err := sPermission.findTransitiveClosure()
	if err != nil {
		sPermission.logger.Panicf("Failed to find transitive closure and mapping job position ids (%s)", err.Error())
	}
	return sPermission
}

func (s *sPermissionService) IsAncestor(ancestorID, nodeID m.ID) bool {
	ancestor := s.mappedJPID2Int[ancestorID]
	node := s.mappedJPID2Int[nodeID]
	return s.jpPaths[ancestor][node] >= 1
}

// Fetch the hierarchy tree. Create adjacency matrix and mapping between job position
// ids and consecutive integers and finally calculate the transitive closure.
// The transitive closure and mapping are stored in the struct immimmediately.
//
// Possible error codes:
// DBError
func (s *sPermissionService) findTransitiveClosure() *e.Error {
	rawHierarchyTree, err := s.permission.GetHierarchyTree()
	if err != nil {
		return e.NewErrorP(err.Error(), SEDBError)
	}
	s.jpPaths, s.mappedJPID2Int = adjacencyList2AdjacencyMatrix(rawHierarchyTree)
	floydWarshall(s.jpPaths, len(s.jpPaths))
	return nil
}

func (s *sPermissionService) GetTransitiveClosure() m.Matrix {
	return s.jpPaths
}

// Find Shortest path between each two verices of given graph. Given graph is as
// adjacency matrix.
func floydWarshall(distance m.Matrix, numVertices int) {
	for k := 0; k < numVertices; k++ {
		for i := 0; i < numVertices; i++ {
			for j := 0; j < numVertices; j++ {
				if distance[i][k]+distance[k][j] < distance[i][j] {
					distance[i][j] = distance[i][k] + distance[k][j]
				}
			}
		}
	}
}

// Map each vertex in the graph to an integer such that these integers be consecutive
// from 0 to n.
func mapVertexIDs2ConInt(adjacencyList m.Graph) map[m.ID]int {
	mapped := make(map[m.ID]int)
	lastInt := 0
	for parent, childrens := range adjacencyList {
		if _, ok := mapped[parent]; !ok {
			mapped[parent] = lastInt
			lastInt++
		}
		for _, child := range *childrens {
			if _, ok := mapped[child]; !ok {
				mapped[child] = lastInt
				lastInt++
			}
		}
	}
	return mapped
}

// Create adjacency matrix represents the hierarchy tree.
// If there's a path from vertex i to j, set matrix[i][j] to 1, otherwise set it to 0.
//
// Return two element. First, the adjacency matrix. Second, a map of vertex IDs to integer.
func adjacencyList2AdjacencyMatrix(adjacencyList m.Graph) (m.Matrix, map[m.ID]int) {
	mapped := mapVertexIDs2ConInt(adjacencyList)
	adjacencyMatrix := make(m.Matrix, 0)
	for parent, childrens := range adjacencyList {
		for _, child := range *childrens {
			adjacencyMatrix[mapped[parent]][mapped[child]] = 1
		}
	}
	return adjacencyMatrix, mapped
}
