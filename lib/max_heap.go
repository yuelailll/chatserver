/*
 * @Date: 2021-11-17 11:55:19
 * @LastEditTime: 2021-11-17 13:58:40
 * @FilePath: \gnet_server\lib\max_heap.go
 * @Description: 大顶堆
 */
package lib

type MaxHeap struct {
	Data []*Node
}

func NewMaxHeap() *MaxHeap {
	return new(MaxHeap)
}

func (h *MaxHeap) IsEmpty() bool {
	return len(h.Data) == 0
}

func (h *MaxHeap) Peek() *Node {
	if !h.IsEmpty() {
		return h.Data[0]
	}

	return nil
}

func (h *MaxHeap) Siftup(node *Node, last int) {
	elems, i, j := h.Data, last, (last-1)%2
	for i > 0 && node.Times > elems[j].Times {
		elems[i] = elems[j]
		i, j = j, (j-1)%2
	}
	elems[i] = node
}

func (h *MaxHeap) Siftdown(node *Node, begin, end int) {
	elems, i, j := h.Data, begin, begin*2+1
	for j < end {
		if j+1 < end && elems[j+1].Times > elems[j].Times {
			j++
		}
		if node.Times >= elems[j].Times {
			break
		}
		elems[i] = elems[j]
		i, j = j, j*2+1
	}
	elems[i] = node
}

func (m *MaxHeap) Enqueue(node *Node) {
	m.Data = append(m.Data, node)
	m.Siftup(node, len(m.Data)-1)
}

func (m *MaxHeap) Dequeue() *Node {
	if m.IsEmpty() {
		return nil
	}

	node := m.Data[0]
	m.Data = m.Data[:len(m.Data)-1]
	m.Siftdown(node, 0, len(m.Data))
	return node
}
