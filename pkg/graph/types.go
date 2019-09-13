package graph

import (
	"strings"
	"sync"

	jsoniter "github.com/json-iterator/go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
)

var json = jsoniter.Config{
	EscapeHTML:             true,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
	DisallowUnknownFields:  true, // non-standard
}.Froze()

// CostFactorOfInAppFiltering = 4 means, we assume that the cost of listing all resources and
// filtering them in the app (instead of using kube-apiserver) is 5x of that via label based selection
const CostFactorOfInAppFiltering = 4

type Edge struct {
	Src        metav1.TypeMeta
	Dst        metav1.TypeMeta
	W          uint64
	Connection v1alpha1.ResourceConnectionSpec
	Forward    bool
}

type AdjacencyMap map[metav1.TypeMeta]*Edge

type Graph struct {
	edges    map[metav1.TypeMeta]AdjacencyMap
	regTypes map[metav1.TypeMeta]*v1alpha1.ResourceID
	regGVR   map[schema.GroupVersionResource]*v1alpha1.ResourceID

	m sync.Mutex
}

func (g *Graph) GVR(t metav1.TypeMeta) schema.GroupVersionResource {
	return g.regTypes[t].GroupVersionResource()
}

func (g *Graph) TypeMeta(gvr schema.GroupVersionResource) metav1.TypeMeta {
	return g.regGVR[gvr].TypeMeta()
}

func NewGraph() *Graph {
	return &Graph{
		edges:    make(map[metav1.TypeMeta]AdjacencyMap),
		regTypes: make(map[metav1.TypeMeta]*v1alpha1.ResourceID),
		regGVR:   make(map[schema.GroupVersionResource]*v1alpha1.ResourceID),
	}
}

func (g *Graph) AddVertex(v v1alpha1.ResourceID) {
	g.regTypes[v.TypeMeta()] = &v
	g.regGVR[v.GroupVersionResource()] = &v
}

func (g *Graph) AddEdge(e *Edge) {
	if _, ok := g.edges[e.Src]; !ok {
		g.edges[e.Src] = AdjacencyMap{}
	}

	// only keep the shortest edge between 2 vertices
	// example: ReplicaSet -> Dep
	// 1. Backward (Dep -> ReplicaSet)
	// 2. Owner Ref (shorter path)
	if old, ok := g.edges[e.Src][e.Dst]; !ok || old.W > e.W {
		g.edges[e.Src][e.Dst] = e
	}
}

// Types of Selectors

// metav1.LabelSelector
// *metav1.LabelSelector

// map[string]string

// ref: https://github.com/coreos/prometheus-operator/blob/cc584ecfa08d2eb95ba9401f116e3a20bf71be8b/pkg/apis/monitoring/v1/types.go#L578
// NamespaceSelector is a selector for selecting either all namespaces or a
// list of namespaces.
// +k8s:openapi-gen=true
type NamespaceSelector struct {
	// Boolean describing whether all namespaces are selected in contrast to a
	// list restricting them.
	Any bool `json:"any,omitempty"`
	// List of namespace names.
	MatchNames []string `json:"matchNames,omitempty"`

	// TODO(fabxc): this should embed metav1.LabelSelector eventually.
	// Currently the selector is only used for namespaces which require more complex
	// implementation to support label selections.
}

// ResourceRef contains information that points to the resource being used
type ResourceRef struct {
	// Name is the name of resource being referenced
	Name string `json:"name"`
	// Namespace is the namespace of resource being referenced
	Namespace string `json:"namespace,omitempty"`
	// Kind is the type of resource being referenced
	Kind string `json:"kind,omitempty"`
	// APIGroup is the group for the resource being referenced
	APIGroup string `json:"apiGroup,omitempty"`
}

func fields(path string) []string {
	return strings.Split(strings.Trim(path, "."), ".")
}

func contains(arr []string, item string) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

func equalsGV(apiGroup string, t metav1.TypeMeta) bool {
	gv1, err := schema.ParseGroupVersion(apiGroup)
	if err != nil {
		return false
	}
	gv2 := t.GroupVersionKind().GroupVersion()
	if gv1.Version != "" && gv1.Version != gv2.Version {
		// if gv2 has version, than version must match
		return false
	}
	return gv1.Group != gv2.Group
}
