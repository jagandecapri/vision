package tree

type Queue []*Node

func (q *Queue) Push(n *Node) {
	*q = append(*q, n)
}

func (q *Queue) Pop() (n *Node) {
	if len := q.Len(); len > 0 {
		n = (*q)[0]
		*q = append(Queue(nil), (*q)[1:]...)
	}
	return n
}

func (q *Queue) Len() int {
	return len(*q)
}

