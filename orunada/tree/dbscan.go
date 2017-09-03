package tree

import (
	"fmt"
	"math"
)

const UNCLASSIFIED  = -1
const NOISE = -2

const CORE_POINT = 1
const NOT_CORE_POINT = 0

const SUCCESS = 0
const FAILURE = -3

type point struct{
 x, y, z float64
 cluster_id int
}

type node struct{
	index int
	next *node
}

type epsilon_neighbours struct{
	num_members int
	head, tail *node
}

func create_node(index int) *node{
	n := &node{}
	n.index = index
	return n
}

func append_at_end(index int, en *epsilon_neighbours) int{
	n := create_node(index)
	if (en.head == nil){
		en.head = n
		en.tail = n
	} else {
		en.tail.next = n
		en.tail = n
	}
	en.num_members++
	return SUCCESS
}

func get_epsilon_neighbours(index int, points []point, num_points int, epsilon float64, dist func (a *point, b *point) float64) epsilon_neighbours{
	en := epsilon_neighbours{}
	for i := 0; i < num_points; i++{
		if i == index{
			continue
		} else if dist(&points[index], &points[i]) > epsilon{
			continue
		} else {
			res := append_at_end(i, &en)
			if res == FAILURE{
				break
			}
		}
	}
	return en
}

func print_epsilon_neightbours(points []point, en *epsilon_neighbours){
	if en != nil{
		h := (*en).head
		for h != nil {
			tmp := points[h.index]
			fmt.Printf("%f, %f, %f\n", tmp.x, tmp.y, tmp.z)
			h = h.next
		}
	}
}

func destroy_epsilon_neighbours(en *epsilon_neighbours){
	//TODO: Find out whether this is needed
}

func dbscan(points []point, num_points int, epsilon float64, minpts int, dist func (a *point, b *point) float64){
	cluster_id := 0
	for i := 0; i < num_points; i++{
		if points[i].cluster_id == UNCLASSIFIED{
			if expand(i, cluster_id, points, num_points, epsilon, minpts, dist) == CORE_POINT{
				cluster_id++
			}
		}
	}
}

func expand(index int, cluster_id int, points []point, num_points int, epsilon float64, minpts int, dist func (a *point, b *point) float64) int{
	return_value := NOT_CORE_POINT
	seeds := get_epsilon_neighbours(index, points, num_points, epsilon, dist)

	if (epsilon_neighbours{}) == seeds{
		return FAILURE
	}
	if (seeds.num_members < minpts){
		points[index].cluster_id = NOISE
	} else {
		points[index].cluster_id = cluster_id
		h := seeds.head
		for h != nil{
			points[h.index].cluster_id = cluster_id
			h = h.next
		}

		h = seeds.head
		for h != nil{
			spread(h.index, &seeds, cluster_id, points, num_points, epsilon, minpts, dist)
			h = h.next
		}
		return_value = CORE_POINT
	}
	return return_value
}

func spread(index int, seeds *epsilon_neighbours, cluster_id int, points []point, num_points int, epsilon float64, minpts int, dist func (a *point, b *point) float64) int{
	spread := get_epsilon_neighbours(index, points, num_points, epsilon, dist)
	if (epsilon_neighbours{}) == spread{
		return FAILURE
	}
	if spread.num_members >= minpts {
		n := spread.head
		for n != nil{
			d := &points[n.index];
			if d.cluster_id == NOISE || d.cluster_id == UNCLASSIFIED {
				if (d.cluster_id == UNCLASSIFIED) {
					if append_at_end(n.index, seeds)== FAILURE {
						destroy_epsilon_neighbours(&spread);
						return FAILURE;
					}
				}
				d.cluster_id = cluster_id;
			}
			n = n.next;
		}
	}
	return SUCCESS
}

func euclidean_dist(a *point, b *point) float64{
	return math.Sqrt(math.Pow(a.x - b.x, 2) + math.Pow(a.y - b.y, 2) + math.Pow(a.z - b.z, 2))
}

func parse_input(){
	//TODO: Do we need this this
}

func print_points (points []point, num_points int){
	fmt.Printf("Number of points: %v\nx\ty\tz\tcluster_id\n----------------------------------------\n", num_points)
	for i := 0; i < num_points; i++{
		tmp := points[i]
		fmt.Printf("%v\t%v\t%v\t%v\n", tmp.x, tmp.y, tmp.z, tmp.cluster_id)
	}
}

func Main_mock(points []point, epsilon float64, minpts int){
	num_points := len(points)
	dbscan(points, num_points, epsilon, minpts, euclidean_dist)
	fmt.Printf("Epsilon: %v\n", epsilon)
	fmt.Printf("Minimum points: %v\n", minpts)
	print_points(points, num_points)
}
