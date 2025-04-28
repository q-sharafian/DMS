package services

import (
	"DMS/internal/graph"
	"DMS/internal/models"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestVertex2ID(t *testing.T) {
	tests := []struct {
		name        string
		vertex      graph.Vertex
		expected    models.ID
		expectedErr error
	}{
		{
			name:        "NilVertex returns NilID",
			vertex:      graph.NilVertex,
			expected:    models.NilID,
			expectedErr: nil,
		},
		{
			name:        "valid UUID returns valid ID",
			vertex:      graph.Vertex([]byte("123e4567-e89b-12d3-a456-426614174000")),
			expected:    models.ID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")),
			expectedErr: nil,
		},
		{
			name:        "invalid UUID returns error",
			vertex:      graph.Vertex([]byte(" invalid-uuid")),
			expected:    models.NilID,
			expectedErr: fmt.Errorf("invalid uuid %s: %s", graph.Vertex([]byte(" invalid-uuid")), "invalid UUID length: 12"),
		},
		{
			name:        "empty UUID returns error",
			vertex:      graph.Vertex([]byte{}),
			expected:    models.NilID,
			expectedErr: fmt.Errorf("invalid uuid %s: %s", graph.Vertex([]byte{}), "invalid UUID length: 0"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			id, err := vertex2ID(test.vertex)
			if id != test.expected {
				t.Errorf("expected ID %v, got %v", test.expected, id)
			}
			if (err != nil && test.expectedErr == nil) || (err == nil && test.expectedErr != nil) {
				t.Errorf("expected error %v, got %v", test.expectedErr, err)
			}
		})
	}
}
