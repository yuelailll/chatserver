/*
 * @Date: 2021-11-17 10:47:26
 * @LastEditTime: 2021-11-17 11:47:56
 * @FilePath: \gnet_server\lib\min_heap.go
 * @Description: 小顶堆热词统计
 */
package lib

type Node struct {
	Word  string
	Times int
}

type MinHeap struct {
	Data []*Node
}

func NewMinHeap() *MinHeap {
	return new(MinHeap)
}

func (h *MinHeap) IsEmpty() bool {
	return len(h.Data) == 0
}

func (h *MinHeap) Peek() *Node {
	if !h.IsEmpty() {
		return h.Data[0]
	}

	return nil
}

func (h *MinHeap) Siftup(node *Node, last int) {
	elems, i, j := h.Data, last, (last-1)%2
	for i > 0 && node.Times < elems[j].Times {
		elems[i] = elems[j]
		i, j = j, (j-1)%2
	}
	elems[i] = node
}

func (h *MinHeap) Siftdown(node *Node, begin, end int) {
	elems, i, j := h.Data, begin, begin*2+1
	for j < end {
		if j+1 < end && elems[j+1].Times < elems[j].Times {
			j++
		}
		if node.Times <= elems[j].Times {
			break
		}
		elems[i] = elems[j]
		i, j = j, j*2+1
	}
	elems[i] = node
}

func (m *MinHeap) Enqueue(node *Node) {
	m.Data = append(m.Data, node)
	m.Siftup(node, len(m.Data)-1)
}

func (m *MinHeap) Dequeue() *Node {
	if m.IsEmpty() {
		return nil
	}

	node := m.Data[0]
	m.Data = m.Data[:len(m.Data)-1]
	m.Siftdown(node, 0, len(m.Data))
	return node
}
