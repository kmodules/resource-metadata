package dot

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

type Graph struct {
	AttributesMap
	id        string
	isCluster bool
	graphType string
	seq       int
	nodes     map[string]Node
	edgesFrom map[string][]Edge
	subgraphs map[string]*Graph
	parent    *Graph
	sameRank  map[string][]Node
}

func NewGraph(options ...GraphOption) *Graph {
	graph := &Graph{
		AttributesMap: AttributesMap{attributes: map[string]interface{}{}},
		graphType:     Directed.Name,
		nodes:         map[string]Node{},
		edgesFrom:     map[string][]Edge{},
		subgraphs:     map[string]*Graph{},
		sameRank:      map[string][]Node{},
	}
	for _, each := range options {
		each.Apply(graph)
	}
	return graph
}

// ID sets the identifier of the graph.
func (g *Graph) ID(newID string) *Graph {
	g.id = newID
	return g
}

func (g *Graph) beCluster() {
	g.id = "cluster_" + g.id
}

// Root returns the top-level graph if this was a subgraph.
func (g *Graph) Root() *Graph {
	if g.parent == nil {
		return g
	}
	return g.parent.Root()
}

// Subgraph returns the Graph with the given label ; creates one if absent.
func (g *Graph) Subgraph(label string, options ...GraphOption) *Graph {
	sub, ok := g.subgraphs[label]
	if ok {
		return sub
	}
	sub = NewGraph(Sub)
	sub.Attr("label", label)
	sub.ID(fmt.Sprintf("s%d", len(g.subgraphs)))
	for _, each := range options {
		each.Apply(sub)
	}
	sub.parent = g
	g.subgraphs[label] = sub
	return sub
}

func (g *Graph) findNode(id string) (Node, bool) {
	if n, ok := g.nodes[id]; ok {
		return n, ok
	}
	if g.parent == nil {
		return Node{id: "void"}, false
	}
	return g.parent.findNode(id)
}

// Node returns the node created with this id or creates a new node if absent.
// This method can be used as both a constructor and accessor.
// not thread safe!
func (g *Graph) Node(id string) Node {
	if n, ok := g.findNode(id); ok {
		return n
	}
	// create a new, use root sequence
	root := g.Root()
	root.seq++
	n := Node{
		id:  id,
		seq: root.seq,
		AttributesMap: AttributesMap{attributes: map[string]interface{}{
			"label": id}},
		graph: g,
	}
	// store local
	g.nodes[id] = n
	return n
}

// Edge creates a new edge between two nodes.
// Nodes can be have multiple edges to the same other node (or itself).
// If one or more labels are given then the "label" attribute is set to the edge.
func (g *Graph) Edge(fromNode, toNode Node, labels ...string) Edge {
	// assume fromNode owner == toNode owner
	edgeOwner := g
	if fromNode.graph != toNode.graph { // 1 or 2 are subgraphs
		edgeOwner = commonParentOf(fromNode.graph, toNode.graph)
	}
	e := Edge{
		from:          fromNode,
		to:            toNode,
		AttributesMap: AttributesMap{attributes: map[string]interface{}{}},
		graph:         edgeOwner}
	if len(labels) > 0 {
		e.Attr("label", strings.Join(labels, ","))
	}
	edgeOwner.edgesFrom[fromNode.id] = append(edgeOwner.edgesFrom[fromNode.id], e)
	return e
}

// FindEdges finds all edges in the graph that go from the fromNode to the toNode.
// Otherwise, returns an empty slice.
func (g *Graph) FindEdges(fromNode, toNode Node) (found []Edge) {
	found = make([]Edge, 0)
	edgeOwner := g
	if fromNode.graph != toNode.graph {
		edgeOwner = commonParentOf(fromNode.graph, toNode.graph)
	}
	if edges, ok := edgeOwner.edgesFrom[fromNode.id]; ok {
		for _, e := range edges {
			if e.to.id == toNode.id {
				found = append(found, e)
			}
		}
	}
	return found
}

func commonParentOf(one *Graph, two *Graph) *Graph {
	// TODO
	return one.Root()
}

// AddToSameRank adds the given nodes to the specified rank group, forcing them to be rendered in the same row
func (g *Graph) AddToSameRank(group string, nodes ...Node) {
	g.sameRank[group] = append(g.sameRank[group], nodes...)
}

// String returns the source in dot notation.
func (g Graph) String() string {
	b := new(bytes.Buffer)
	g.Write(b)
	return b.String()
}

func (g Graph) Write(w io.Writer) {
	g.IndentedWrite(NewIndentWriter(w))
}

// IndentedWrite write the graph to a writer using simple TAB indentation.
func (g Graph) IndentedWrite(w *IndentWriter) {
	fmt.Fprintf(w, "%s %s {", g.graphType, g.id)
	w.NewLineIndentWhile(func() {
		if len(g.id) > 0 {
			fmt.Fprintf(w, "ID = %q;", g.id)
			w.NewLine()
		}
		// subgraphs
		for _, key := range g.sortedSubgraphsKeys() {
			each := g.subgraphs[key]
			each.IndentedWrite(w)
		}
		// graph attributes
		appendSortedMap(g.AttributesMap.attributes, false, w)
		w.NewLine()
		// graph nodes
		for _, key := range g.sortedNodesKeys() {
			each := g.nodes[key]
			fmt.Fprintf(w, "n%d", each.seq)
			appendSortedMap(each.attributes, true, w)
			fmt.Fprintf(w, ";")
			w.NewLine()
		}
		// graph edges
		denoteEdge := "->"
		if g.graphType == "graph" {
			denoteEdge = "--"
		}
		for _, each := range g.sortedEdgesFromKeys() {
			all := g.edgesFrom[each]
			for _, each := range all {
				fmt.Fprintf(w, "n%d%sn%d", each.from.seq, denoteEdge, each.to.seq)
				appendSortedMap(each.attributes, true, w)
				fmt.Fprint(w, ";")
				w.NewLine()
			}
		}
		for _, nodes := range g.sameRank {
			str := ""
			for _, n := range nodes {
				str += fmt.Sprintf("n%d;", n.seq)
			}
			fmt.Fprintf(w, "{rank=same; %s};", str)
			w.NewLine()
		}
	})
	fmt.Fprintf(w, "}")
}

func appendSortedMap(m map[string]interface{}, mustBracket bool, b io.Writer) {
	if len(m) == 0 {
		return
	}
	if mustBracket {
		fmt.Fprint(b, "[")
	}
	first := true
	// first collect keys
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	sort.StringSlice(keys).Sort()

	for _, k := range keys {
		if !first {
			if mustBracket {
				fmt.Fprint(b, ",")
			} else {
				fmt.Fprintf(b, ";")
			}
		}
		if html, isHTML := m[k].(HTML); isHTML {
			fmt.Fprintf(b, "%s=<%s>", k, html)
		} else if literal, isLiteral := m[k].(Literal); isLiteral {
			fmt.Fprintf(b, "%s=%s", k, literal)
		} else {
			fmt.Fprintf(b, "%s=%q", k, m[k])
		}
		first = false
	}
	if mustBracket {
		fmt.Fprint(b, "]")
	} else {
		fmt.Fprint(b, ";")
	}
}