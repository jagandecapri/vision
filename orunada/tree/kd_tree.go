package tree

type Point interface {
	// Return the total number of dimensions
	Dim() int
	// Return the value X_{dim}, dim is started from 0
	GetValue(dim int) int
	// Return the distance between two points
	//Distance(point Point) float64
	// Return the distance between the point and the plane X_{dim}=val
	//PlaneDistance(val float64, dim int) float64
}

type KDTree struct{
	root *Node
	len int
}

type Node struct {
	left  *Node
	data Point
	right *Node
}

func (kd *KDTree) newNode(p Point) *Node {
	kd.len++
	tmp := &Node{data:p}

	/**for i := 0; i < k; i++ {
		tmp.point = append(tmp.point, arr[i])
	}**/
	return tmp
}

func (kd *KDTree) insertRec(root *Node, depth int, p Point) *Node {
	if root == nil {
		return kd.newNode(p)
	}

	cd := depth % p.Dim()

	if p.GetValue(cd) < root.data.GetValue(cd) {
		root.left = kd.insertRec(root.left, depth+1, p)
	} else {
		root.right = kd.insertRec(root.right, depth+1, p)
	}
	return root
}

func (kd *KDTree) Insert(p Point){
	kd.root = kd.insertRec(kd.root, 0, p)
}

func (kd *KDTree) arePointsSame(p1 Point, p2 Point) bool {
	k := p1.Dim()
	for i := 0; i < k; i++ {
		if p1.GetValue(i) != p2.GetValue(i) {
			return false
		}
	}
	return true
}

func (kd *KDTree) searchRec(root *Node, depth int, p Point) bool {
	if root == nil {
		return false
	}
	if kd.arePointsSame(root.data, p) {
		return true
	}

	cd := depth % p.Dim()

	if p.GetValue(cd) < root.data.GetValue(cd) {
		return kd.searchRec(root.left, depth+1, p)
	}

	return kd.searchRec(root.right, depth+1, p)
}

func (kd *KDTree) Search(p Point) bool {
	return kd.searchRec(kd.root, 0, p)
}

func (kd *KDTree) BFSTraverse() []Point{
	queue := Queue{}
	val := []Point{}
	tmp_node := kd.root
	for tmp_node != nil{
		val = append(val, tmp_node.data)
		if tmp_node.left != nil{
			queue.Push(tmp_node.left)
		}
		if tmp_node.right != nil{
			queue.Push(tmp_node.right)
		}
		tmp_node = queue.Pop()
	}
	return val
}

func (kd *KDTree) BFSTraverseChan(out chan<- Point){
	go func(out chan<- Point){
		queue := Queue{}
		tmp_node := kd.root
		for tmp_node != nil{
			out<- tmp_node.data
			if tmp_node.left != nil{
				queue.Push(tmp_node.left)
			}
			if tmp_node.right != nil{
				queue.Push(tmp_node.right)
			}
			tmp_node = queue.Pop()
		}
		close(out)
	}(out)
}
