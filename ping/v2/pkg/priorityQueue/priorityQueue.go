package priorityQueue

import "time"

type DurationHeap []time.Duration

func (h DurationHeap) Len() int           { return len(h) }
func (h DurationHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h DurationHeap) Swap(i, j int) 	 { h[i], h[j] = h[j], h[i] }

func (h *DurationHeap) Push (x any) {
	*h = append(*h, x.(time.Duration))
}

func (h *DurationHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n - 1]		// Mark on the last element
	*h = old[0 : n- 1]	// Assign the heap except for last element
	return x
}