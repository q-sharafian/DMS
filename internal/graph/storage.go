package graph

type Vertex []byte

// Represent nil-value for Vertex
var NilVertex Vertex = Vertex([]byte("00000000"))

func (v Vertex) Equals(other Vertex) bool     { return string(v) == string(other) }
func (v Vertex) String() string               { return string(v) }
func (v Vertex) str2Vertex(str string) Vertex { return Vertex(str) }

// Represents an Edge in directed graph. Start in the begining of the Edge, end is
// the end of the Edge.
type Edge struct{ Start, End Vertex }

// Represent nil-value
var NilEdge Edge = Edge{Start: nil, End: nil}

func (e *Edge) equals(other Edge) bool { return e.Start.Equals(other.Start) && e.End.Equals(other.End) }
func (e *Edge) isNil() bool            { return e.equals(NilEdge) }
func (e *Edge) string() string {
	return string(e.Start) + ":" + string(e.End)
}

// storage defines the interface for cache storage implementations
type storage interface {
	// Return (value, error). If there's not such key, the error type would be "e.ErrNotFound"
	Get(edge Edge) (bool, error)
	Set(edge Edge, value bool) error
	Delete(edge Edge)
	// Clear graph cache
	Clear() error
	// Return number of edges
	Size() (int, error)
	// Deletes all entries/edges with matching start vertex
	DeleteByPrefix(start Vertex) error
}
