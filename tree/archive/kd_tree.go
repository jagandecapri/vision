package tree

type KDTree struct{
	root *Node
	len int
}

type Node struct {
	left  *Node
	data  PointInterface
	right *Node
}

func (kd *KDTree) newNode(p PointInterface) *Node {
	kd.len++
	tmp := &Node{data:p}
	return tmp
}

func (kd *KDTree) insertRec(root *Node, depth int, p PointInterface) *Node {
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

func (kd *KDTree) Insert(p ...PointInterface){
	for _, v := range p{
		kd.root = kd.insertRec(kd.root, 0, v)
	}
}

func (kd *KDTree) arePointsSame(p1 PointInterface, p2 PointInterface) bool {
	k := p1.Dim()
	for i := 0; i < k; i++ {
		if p1.GetValue(i) != p2.GetValue(i) {
			return false
		}
	}
	return true
}

func (kd *KDTree) searchRec(root *Node, depth int, p PointInterface) bool {
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

func (kd *KDTree) Search(p PointInterface) bool {
	return kd.searchRec(kd.root, 0, p)
}

func (kd *KDTree) BFSTraverse() []PointInterface {
	queue := Queue{}
	val := []PointInterface{}
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

func (kd *KDTree) BFSTraverseChan(out chan<- PointInterface){
	go func(out chan<- PointInterface){
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