package hierarchy

import (
	"DMS/internal/graph"
	l "DMS/internal/logger"
)

type HierarchyTree struct {
	graph  *graph.DynamicGraph
	logger l.Logger
}

func NewHierarchyTree(graph *graph.DynamicGraph, logger l.Logger) *HierarchyTree {
	return &HierarchyTree{
		graph:  graph,
		logger: logger,
	}
}

// Check if there's a path from claimed ancestorID to nodeID.
// If claimed ancestorID be "NilVertex", return true anyway.
func (h *HierarchyTree) IsAncestor(ancestorID, nodeID graph.Vertex) (bool, error) {
	if ancestorID.Equals(graph.NilVertex) {
		return true, nil
	}
	return h.graph.HasPath(ancestorID, nodeID)
}

func (h *HierarchyTree) Graph() *graph.DynamicGraph {
	return h.graph
}

// Get all nested children of the input vertex with the self vertex.
func (h *HierarchyTree) GetNestedChilds(nodeID graph.Vertex) ([]graph.Vertex, error) {
	return h.graph.GetAllNestedChildren(nodeID), nil
}

// Return true if the given vertex is a source vertex. Means it has no parents.
func (h *HierarchyTree) IsSourceVertex(nodeID graph.Vertex) (bool, error) {
	parents := h.graph.GetParents(nodeID)
	return len(*parents) == 0, nil
}
