package graph

import (
	e "DMS/internal/error"
	l "DMS/internal/logger"
	"container/list"
	"errors"
	"fmt"
	"sync"
)

// DynamicGraph represents a directed graph with caching capabilities
type DynamicGraph struct {
	graph map[string]map[string]struct{} // adjacency list using maps for O(1) lookups
	cache storage                        // interface for cache storage
	mu    sync.RWMutex                   // mutex for thread safety
}

// NewDynamicGraph creates a new instance of DynamicGraph
func NewDynamicGraph(storage storage, logger l.Logger) *DynamicGraph {
	if storage == nil {
		storage = NewMemoryStorage(logger)
	}
	graph := &DynamicGraph{
		graph: make(map[string]map[string]struct{}),
		cache: storage,
	}
	return graph
}

// HasPath checks if there's a path from start to end using BFS. If there's not a path,
// returns (false, nil).
func (g *DynamicGraph) HasPath(start, end Vertex) (bool, error) {
	pair := Edge{Start: start, End: end}

	// Check cache first
	if hasPath, err := g.cache.Get(pair); err != nil && !errors.Is(err, e.ErrNotFound) {
		return false, fmt.Errorf("error in checking existing path from %s to %s: %s",
			start.String(), end.String(), err.Error())
	} else if err == nil {
		return hasPath, nil
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	// Check cache again after acquiring write lock
	if hasPath, err := g.cache.Get(pair); err != nil && !errors.Is(err, e.ErrNotFound) {
		return false, fmt.Errorf("error in checking existing path from %s to %s: %s",
			start.String(), end.String(), err.Error())
	} else if err == nil {
		return hasPath, nil
	}

	// BFS implementation
	visited := make(map[string]struct{})
	queue := list.New()
	queue.PushBack(start)
	for queue.Len() > 0 {
		vertex := queue.Remove(queue.Front()).(Vertex)
		if vertex.Equals(end) {
			if err := g.cache.Set(pair, true); err != nil {
				return false, err
			}
			return true, nil
		}

		if _, seen := visited[vertex.String()]; !seen {
			visited[vertex.String()] = struct{}{}
			// Cache intermediate results
			if vertex.Equals(start) {
				if err := g.cache.Set(Edge{Start: start, End: vertex}, true); err != nil {
					return false, err
				}
			}

			// Add unvisited neighbors to queue
			if neighbors, exists := g.graph[vertex.String()]; exists {
				for neighbor := range neighbors {
					if _, seen := visited[neighbor]; !seen {
						queue.PushBack(Vertex{}.str2Vertex(neighbor))
					}
				}
			}
		}
	}

	if err := g.cache.Set(pair, false); err != nil {
		return false, err
	}
	return false, nil
}

// addEdge adds a directed edge from u to v.
func (g *DynamicGraph) addEdge(u, v Vertex) error {
	u_str := u.String()
	v_str := v.String()
	g.mu.Lock()
	defer g.mu.Unlock()

	// Initialize adjacency set if not exists
	if _, exists := g.graph[u_str]; !exists {
		g.graph[u_str] = make(map[string]struct{})
	}
	// Check if edge already exists
	if _, exists := g.graph[u_str][v_str]; !exists {
		g.graph[u_str][v_str] = struct{}{}
		// Invalidate cache entries that start from u
		return g.cache.DeleteByPrefix(u)
	}
	return nil
}

// removeEdge removes a directed edge from u to v
func (g *DynamicGraph) removeEdge(u, v Vertex) error {
	u_str := u.String()
	v_str := v.String()
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.graph[u_str]; exists {
		if _, exists := g.graph[u_str][v_str]; exists {
			delete(g.graph[u_str], v_str)
			// Invalidate cache entries that start from u
			return g.cache.DeleteByPrefix(u)
		}
	}
	return nil
}

// ClearCache clears the reachability cache
func (g *DynamicGraph) ClearCache() error {
	return g.cache.Clear()
}

// LimitCacheSize limits the cache size by removing oldest entries
func (g *DynamicGraph) LimitCacheSize(maxSize int) error {
	size, err := g.cache.Size()
	if err != nil {
		return err
	}
	if size > maxSize {
		return g.ClearCache() // Simple implementation - clear all when limit reached
		// For more sophisticated implementation, implement it in the storage
	}
	return nil
}

// Return number of vertices of the graph
func (g *DynamicGraph) Size() (int, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	count := 0
	for _, neighbors := range g.graph {
		count += len(neighbors)
	}
	return count, nil
}

type GraphChangeType int

const (
	AddEdge GraphChangeType = iota
	RemoveEdge
)

type GraphChange struct {
	Type        GraphChangeType
	Edge        Edge
	ResponseErr chan error
}

// ProcessChanges processes and aplies changes to the graph
func (g *DynamicGraph) ProcessChanges(changes <-chan GraphChange) {
	for change := range changes {
		// if change.Edge.isNil() {
		// 	change.ResponseErr <- nil
		// 	continue
		// }
		switch change.Type {
		case AddEdge:
			change.ResponseErr <- g.addEdge(change.Edge.Start, change.Edge.End)
		case RemoveEdge:
			change.ResponseErr <- g.removeEdge(change.Edge.Start, change.Edge.End)
		}

	}
}

// Returns a string representation of the graph
func (g *DynamicGraph) String() string {
	var strGraph string
	for u, neighbors := range g.graph {
		strGraph += fmt.Sprintf("  %s -> %v\n", u, func() []string {
			neighborsSlice := make([]string, 0, len(neighbors))
			for neighbor := range neighbors {
				neighborsSlice = append(neighborsSlice, neighbor)
			}
			return neighborsSlice
		}())
	}
	return strGraph
}

// GetAllNestedChildren returns a slice of vertices that includes the given vertex and
// all its nested children.
func (g *DynamicGraph) GetAllNestedChildren(vertex Vertex) []Vertex {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var result []Vertex
	queue := list.New()
	visited := make(map[string]struct{})

	vertexStr := vertex.String()
	queue.PushBack(vertexStr)
	visited[vertexStr] = struct{}{}

	for queue.Len() > 0 {
		current := queue.Remove(queue.Front()).(string)
		result = append(result, Vertex{}.str2Vertex(current))

		if neighbors, exists := g.graph[current]; exists {
			for neighbor := range neighbors {
				if _, seen := visited[neighbor]; !seen {
					visited[neighbor] = struct{}{}
					queue.PushBack(neighbor)
				}
			}
		}
	}

	return result
}

func (g *DynamicGraph) GetParents(vertex Vertex) *[]Vertex {
	g.mu.RLock()
	defer g.mu.RUnlock()
	parents := make([]Vertex, 0)
	for parent := range g.graph {
		if _, exists := g.graph[parent][vertex.String()]; exists {
			parents = append(parents, Vertex{}.str2Vertex(parent))
		}
	}
	return &parents
}
