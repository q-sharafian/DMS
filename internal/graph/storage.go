package graph

type Vertex []byte

func (v Vertex) equals(other Vertex) bool     { return string(v) == string(other) }
func (v Vertex) string() string               { return string(v) }
func (v Vertex) str2Vertex(str string) Vertex { return Vertex(str) }

// Represents an Edge in directed graph. Start in the begining of the Edge, end is
// the end of the Edge.
type Edge struct{ Start, End Vertex }

// Represent nil-value
var NilEdge Edge = Edge{Start: nil, End: nil}

func (e *Edge) equals(other Edge) bool { return e.Start.equals(other.Start) && e.End.equals(other.End) }
func (e *Edge) isNil() bool            { return e.equals(NilEdge) }
func (e *Edge) string() string {
	return string(e.Start) + ":" + string(e.End)
}

// storage defines the interface for cache storage implementations
type storage interface {
	Get(edge Edge) (bool, bool) // Returns (value, exists)
	Set(edge Edge, value bool) error
	Delete(edge Edge)
	// Clear graph cache
	Clear() error
	// Return number of edges
	Size() (int, error)
	// Deletes all entries/edges with matching start vertex
	DeleteByPrefix(start Vertex) error
}
