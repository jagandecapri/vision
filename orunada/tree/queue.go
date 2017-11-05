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

type QueueUnit []*Unit

func (q *QueueUnit) Push(n *Unit) {
	*q = append(*q, n)
}

func (q *QueueUnit) Pop() (n *Unit) {
	if len := q.Len(); len > 0 {
		n = (*q)[0]
		*q = append(QueueUnit(nil), (*q)[1:]...)
	}
	return n
}

func (q *QueueUnit) Len() int {
	return len(*q)
}