package models

type Matrix [][]int8

func (m *Matrix) New(rows, cols int) Matrix {
	matrix := make(Matrix, rows)
	for i := range matrix {
		matrix[i] = make([]int8, cols)
	}
	return matrix
}

func (m *Matrix) GetDim() (rows, cols int) {
	return len(*m), len((*m)[0])
}

// It's an adjacency list. Each key is a parent and its value is a list of its children.
// It's possible some values be nil and empty.
type Graph map[ID]*[]ID

// List of some permissions the job position could have
type Permission struct {
	// ID of the job position the permission is for
	JPID ID
	// Does the current job position is allowed to create a job position as child of himself?
	IsAllowCreateJP bool `json:"is_allow_create_jp" validate:"required"`
}

type HierarchyTree struct {
	// Represent ID of this record in database
	ID   ID `json:"id"`
	JPID ID `json:"jp_id"`
	// List of job positions that are child of the job position id
	JPChildsID *[]ID `json:"child_jps_id"`
}
