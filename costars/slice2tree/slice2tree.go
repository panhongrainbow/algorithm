package slice2tree

type SliceTree struct {
	heap     []int
	heapSize int
}

func New() *SliceTree {
	return &SliceTree{heap: make([]int, 0)}
}

func (h *SliceTree) IsEmpty() bool {
	return h.heapSize == 0
}

func (h *SliceTree) Push(v int) {
	if h.heapSize == len(h.heap) {
		h.heap = append(h.heap, v)
	} else {
		h.heap[h.heapSize] = v
	}
	h.heapInsert(h.heapSize)
	h.heapSize++
}

func (h *SliceTree) Pop() int {
	ans := h.heap[0]
	h.heapSize--
	if h.heapSize > 0 {
		h.swap(0, h.heapSize)
		h.heapify(0, h.heapSize)
	}
	return ans
}

func (h *SliceTree) heapInsert(i int) {
	for h.heap[i] > h.heap[(i-1)/2] {
		h.swap(i, (i-1)/2)
		i = (i - 1) / 2
	}
}

func (h *SliceTree) heapify(i, heapSize int) {
	left := i*2 + 1
	for left < heapSize {
		largest := left
		if left+1 < heapSize && h.heap[left+1] > h.heap[left] {
			largest = left + 1
		}
		if h.heap[largest] <= h.heap[i] {
			largest = i
		}

		if largest == i {
			break
		}

		h.swap(largest, i)
		i = largest
		left = i*2 + 1
	}
}

func (h *SliceTree) swap(i, j int) {
	tmp := h.heap[i]
	h.heap[i] = h.heap[j]
	h.heap[j] = tmp
}
