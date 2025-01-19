package slice2tree

// =====================================================================================================================
//                  🧱 Efficient Heap Operations (SliceTree)
// =====================================================================================================================
// 🧩 SliceTree is an implementation of a binary heap designed with efficiency and clarity in mind.
// 🧩 This structure provides a simple yet powerful way to manage dynamic datasets while preserving the heap property.
// 🧩 With optimized insertion (Push) and removal (Pop) operations, it ensures quick access to the maximum element.
// 🧩 Ideal for scenarios requiring priority queue functionality, SliceTree is a versatile and effective solution.

// SliceTree 🧩 represents a binary heap data structure.
type SliceTree struct {
	// heap is the underlying array that stores the heap elements.
	heap []int64
	// heapSize is the current number of elements in the heap.
	heapSize int
}

// NewHeap 🧩 returns a new SliceTree instance with the specified capacity.
func NewHeap(capacity int) *SliceTree {
	return &SliceTree{heap: make([]int64, 0, capacity)}
}

// IsEmpty 🧩 checks if the heap is empty.
func (h *SliceTree) IsEmpty() bool {
	// If the heap size is 0, the heap is empty.
	return h.heapSize == 0
}

// Push 🧩 adds a new element to the heap.
func (h *SliceTree) Push(v int64) {
	// If the heap is full, append a new element to the underlying array.
	if h.heapSize == len(h.heap) {
		h.heap = append(h.heap, v)
	} else {
		// Otherwise, reuse the existing array space.
		h.heap[h.heapSize] = v
	}
	// Insert the new element into the heap.
	h.heapInsert(h.heapSize)
	// Increment the heap size.
	h.heapSize++
}

// heapInsert 🧩 inserts an element at the specified index into the heap.
func (h *SliceTree) heapInsert(i int) {
	// While the element is greater than its parent, swap them.
	for h.heap[i] > h.heap[(i-1)/2] {
		// Swap the element with its parent.
		h.swap(i, (i-1)/2)
		// Move up the heap.
		i = (i - 1) / 2
	}
}

// Pop 🧩 removes and returns the maximum element from the heap.
func (h *SliceTree) Pop() int64 {
	// Save the maximum element.
	ans := h.heap[0]
	// Decrement the heap size.
	h.heapSize--
	// If the heap is not empty, heapify the remaining elements.
	if h.heapSize > 0 {
		// Swap the maximum element with the last element.
		h.swap(0, h.heapSize)
		// Heapify the remaining elements.
		h.heapify(0, h.heapSize)
	}
	// Return the maximum element.
	return ans
}

// heapify 🧩 restores the heap property at the specified index.
func (h *SliceTree) heapify(i, heapSize int) {
	// Initialize the left child index.
	left := i*2 + 1

	// While the left child is within the heap bounds.
	for left < heapSize {
		// Find the largest child.
		largest := left

		// If the right child is larger, update the largest index.
		if left+1 < heapSize && h.heap[left+1] > h.heap[left] {
			largest = left + 1
		}

		// If the largest child is not greater than the current element, break.
		if h.heap[largest] <= h.heap[i] {
			largest = i
		}

		// If the largest child is the current element, break.
		if largest == i {
			break
		}

		// Swap the current element with the largest child.
		h.swap(largest, i)

		// Move down the heap.
		i = largest

		// Update the left child index.
		left = i*2 + 1
	}
}

// swap 🧩 swaps two elements in the heap.
func (h *SliceTree) swap(i, j int) {
	// Save the element at index i.
	tmp := h.heap[i]

	// Swap the elements.
	h.heap[i] = h.heap[j]
	h.heap[j] = tmp
}
