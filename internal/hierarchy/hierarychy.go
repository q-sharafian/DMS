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

// Check if there's a path from ancestorID to nodeID
func (h *HierarchyTree) IsAncestor(ancestorID, nodeID graph.Vertex) (bool, error) {
	return h.graph.HasPath(ancestorID, nodeID)
}

func (h *HierarchyTree) Graph() *graph.DynamicGraph {
	return h.graph
}
