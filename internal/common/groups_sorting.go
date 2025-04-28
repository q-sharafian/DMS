package common

// MergeOrderedGroups combines multiple ordered slices into a single ordered slice.
// It uses a high-performance merge algorithm.
// comparator function declare the strategy to sort items. (e.g., ascending or descending)
func MergeOrderedGroups[T any](groups *[]*[]T, comparator func(a, b T) bool) *[]T {
	// Create a min-heap to keep track of the smallest elements from each group
	type heapNode struct {
		value T
		group int
		index int
	}

	// Create a min-heap using a slice
	heap := make([]heapNode, 0)

	// Initialize the heap with the first element from each group
	for i, group := range *groups {
		if len(*group) > 0 {
			heap = append(heap, heapNode{value: (*group)[0], group: i, index: 0})
		}
	}

	// Heapify the initial heap
	heapify(heap, func(a, b heapNode) bool {
		return comparator(a.value, b.value)
	})

	// Result slice
	result := make([]T, 0)

	// While the heap is not empty
	for len(heap) > 0 {
		// Extract the smallest element from the heap
		smallest := heap[0]
		result = append(result, smallest.value)

		// Move to the next element in the same group
		if nextIndex := smallest.index + 1; nextIndex < len(*((*groups)[smallest.group])) {
			heap[0] = heapNode{
				value: (*(*groups)[smallest.group])[nextIndex],
				group: smallest.group,
				index: nextIndex,
			}
		} else {
			// Remove the group from the heap
			heap[0] = heap[len(heap)-1]
			heap = heap[:len(heap)-1]
		}

		// Re-heapify
		heapify(heap, func(a, b heapNode) bool {
			return comparator(a.value, b.value)
		})
	}

	return &result
}

// Helper function to heapify a slice
func heapify[T any](heap []T, less func(a, b T) bool) {
	n := len(heap)
	for i := (n - 1) / 2; i >= 0; i-- {
		downHeap(heap, i, n, less)
	}
}

func downHeap[T any](heap []T, i, n int, less func(a, b T) bool) {
	for {
		left := 2*i + 1
		right := 2*i + 2
		smallest := i

		if left < n && less(heap[left], heap[smallest]) {
			smallest = left
		}
		if right < n && less(heap[right], heap[smallest]) {
			smallest = right
		}
		if smallest == i {
			break
		}
		heap[i], heap[smallest] = heap[smallest], heap[i]
		i = smallest
	}
}
