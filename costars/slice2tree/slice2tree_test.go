package slice2tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_SliceTree tests the functionality of the SliceTree data structure.
func Test_SliceTree(t *testing.T) {
	// Test NewHeap function to ensure it returns a non-nil SliceTree instance.
	t.Run("NewHeap", func(t *testing.T) {
		// Create a new SliceTree instance with a capacity of 10.
		tree := NewHeap(10)

		// Verify that the returned SliceTree instance is not nil.
		require.NotNil(t, tree, "NewHeap() should not return nil")

		// Check that the initial heap length and heapSize are 0.
		assert.Equal(t, 0, len(tree.heap), "heap length should be 0")
		assert.Equal(t, 0, tree.heapSize, "heapSize should be 0")
	})

	// Test IsEmpty function to ensure it correctly identifies empty heaps.
	t.Run("IsEmpty", func(t *testing.T) {
		// Create a new SliceTree instance with a capacity of 10.
		tree := NewHeap(10)

		// Verify that IsEmpty returns true for an empty heap.
		assert.True(t, tree.IsEmpty(), "IsEmpty() should return true for an empty heap")

		// Add an element to the heap and verify that IsEmpty returns false.
		tree.Push(5)
		assert.False(t, tree.IsEmpty(), "IsEmpty() should return false after adding an element")
	})

	// Test Push and Pop functions to ensure they correctly add and remove elements.
	t.Run("Push and Pop", func(t *testing.T) {
		// Create a new SliceTree instance with a capacity of 10.
		tree := NewHeap(10)

		// Define a list of values to push onto the heap.
		values := []int64{3, 10, 5, 6, 2}

		// Define the expected order of popped elements.
		expected := []int64{10, 6, 5, 3, 2}

		// Push each value onto the heap.
		for _, v := range values {
			tree.Push(v)
		}

		// Verify that the heapSize matches the number of pushed elements.
		assert.Equal(t, len(values), tree.heapSize, "heapSize should match number of pushed elements")

		// Pop each element from the heap and verify the correct order.
		for _, exp := range expected {
			assert.Equal(t, exp, tree.Pop(), "Pop() should return the correct value")
		}

		// Verify that the heap is empty after popping all elements.
		assert.True(t, tree.IsEmpty(), "Heap should be empty after popping all elements")
	})

	// Test Push function with duplicate values to ensure correct ordering.
	t.Run("Push with Duplicates", func(t *testing.T) {
		// Create a new SliceTree instance with a capacity of 10.
		tree := NewHeap(10)

		// Push duplicate values onto the heap.
		tree.Push(5)
		tree.Push(5)
		tree.Push(10)
		tree.Push(10)

		// Define the expected order of popped elements.
		expected := []int64{10, 10, 5, 5}

		// Pop each element from the heap and verify the correct order.
		for _, exp := range expected {
			assert.Equal(t, exp, tree.Pop(), "Pop() should handle duplicate values correctly")
		}
	})

	// Test edge case: Pop from an empty heap to ensure it panics.
	t.Run("Pop Empty", func(t *testing.T) {
		// Create a new SliceTree instance with a capacity of 10.
		tree := NewHeap(10)

		// Verify that Pop panics when the heap is empty.
		assert.Panics(t, func() {
			tree.Pop()
		}, "Pop() should panic when the heap is empty")
	})
}
