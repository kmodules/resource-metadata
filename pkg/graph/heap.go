// This example demonstrates a distance queue built using the heap interface.
package graph

import (
	"container/heap"
	"math"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// An Item is something we manage in a priority queue.
type Item struct {
	vertex metav1.TypeMeta // The value of the item; arbitrary.
	dist   uint64          // The priority of the item in the queue.
	// The index is needed by Update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A Queue implements heap.Interface and holds Items.
type Queue []*Item

func (q Queue) Len() int { return len(q) }

func (q Queue) Less(i, j int) bool {
	return q[i].dist < q[j].dist
}

func (q Queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *Queue) Push(x interface{}) {
	n := len(*q)
	item := x.(*Item)
	item.index = n
	*q = append(*q, item)
}

func (q *Queue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*q = old[0 : n-1]
	return item
}

// Update modifies the priority and value of an Item in the queue.
func (q *Queue) Update(item *Item, dist uint64) {
	item.dist = dist
	heap.Fix(q, item.index)
}

// ref: https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm#Pseudocode

func Dijkstra(graph *Graph, src metav1.TypeMeta) (dist map[metav1.TypeMeta]uint64, prev map[metav1.TypeMeta]*Edge) {
	dist = make(map[metav1.TypeMeta]uint64)
	prev = make(map[metav1.TypeMeta]*Edge)

	q := make(Queue, len(graph.regTypes))
	i := 0
	items := make(map[metav1.TypeMeta]*Item)
	for vertex := range graph.regTypes {
		var d uint64 = math.MaxUint32 // avoid overflow
		if vertex == src {
			d = 0 // dist[src] = 0
		}

		dist[vertex] = d
		prev[vertex] = nil
		item := &Item{
			vertex: vertex,
			dist:   d,
			index:  i,
		}
		items[vertex] = item
		q[i] = item
		i++
	}
	heap.Init(&q)

	for len(q) > 0 {
		u := heap.Pop(&q).(*Item)

		for v, e := range graph.edges[u.vertex] {
			alt := dist[u.vertex] + e.W
			if alt < dist[v] {
				dist[v] = alt
				prev[v] = e
				q.Update(items[v], alt)
			}
		}
	}

	return
}

type Path struct {
	Source   metav1.TypeMeta
	Target   metav1.TypeMeta
	Distance uint64
	Edges    []*Edge
}

// v1 -> v2 -> v3
//func (p Path) String() string {
//	return strings.Join([]string{gvr.Group, "/", gvr.Version, ", Resource=", gvr.Resource}, "")
//}

func GeneratePaths(src metav1.TypeMeta, dist map[metav1.TypeMeta]uint64, prev map[metav1.TypeMeta]*Edge) map[metav1.TypeMeta]*Path {
	paths := make(map[metav1.TypeMeta]*Path)

	for target, d := range dist {
		if d < math.MaxUint32 {
			path := Path{
				Source:   src,
				Target:   target,
				Distance: d,
				Edges:    make([]*Edge, 0),
			}

			u := target
			for prev[u] != nil {
				path.Edges = append(path.Edges, nil)
				copy(path.Edges[1:], path.Edges[0:])
				path.Edges[0] = prev[u]
				u = prev[u].Src
			}

			paths[target] = &path
		}
	}

	return paths
}
