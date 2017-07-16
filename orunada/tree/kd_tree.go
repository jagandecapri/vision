package tree

import (
	"github.com/golang-collections/go-datastructures/queue"
	"fmt"
)

var k int = 2

type KDTree struct{
	root *Node
	len int
}

type Node struct {
	left  *Node
	point []int
	right *Node
}

func (kd *KDTree) newNode(arr ...int) *Node {
	kd.len++
	tmp := &Node{}

	for i := 0; i < k; i++ {
		tmp.point = append(tmp.point, arr[i])
	}
	return tmp
}
func (kd *KDTree) insertRec(root *Node, depth int, point ...int) *Node {
	if root == nil {
		return kd.newNode(point...)
	}

	cd := depth % k

	if point[cd] < root.point[cd] {
		root.left = kd.insertRec(root.left, depth+1, point...)
	} else {
		root.right = kd.insertRec(root.right, depth+1, point...)
	}
	return root
}

func (kd *KDTree) Insert(point ...int){
	kd.root = kd.insertRec(kd.root, 0, point...)
}

func (kd *KDTree) arePointsSame(point1 []int, point2 []int) bool {
	for i := 0; i < k; i++ {
		if point1[i] != point2[i] {
			return false
		}
	}
	return true
}

func (kd *KDTree) searchRec(root *Node, depth int, point ...int) bool {
	if root == nil {
		return false
	}
	if kd.arePointsSame(root.point, point) {
		return true
	}

	cd := depth % k

	if point[cd] < root.point[cd] {
		return kd.searchRec(root.left, depth+1, point...)
	}

	return kd.searchRec(root.right, depth+1, point...)
}

func (kd *KDTree) Search(point ...int) bool {
	return kd.searchRec(kd.root, 0, point...)
}

func loop(queue *queue.Queue, val [][]int) [][]int{
	elems, _ := queue.Get(1)
	if elems[0] == nil{
		fmt.Println(elems[0])
		return val
	}
	tmp, ok := elems[0].(Node)
	if ok{
		val = append(val, tmp.point)
		queue.Put(tmp.left, tmp.right)
		val = loop(queue, val)
	}
	return val
}

func (kd *KDTree) BFSTraverse() [][]int{
	queue := queue.New(int64(1000))
	queue.Put(kd.root)
	val := [][]int{}
	fmt.Println()
	val = loop(queue, val)
	return val
}
