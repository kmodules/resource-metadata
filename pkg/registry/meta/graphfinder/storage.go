package graphfinder

import (
	"context"

	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/rest"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"kmodules.xyz/resource-metadata/pkg/graph"
)

type Storage struct {
}

var _ rest.GroupVersionKindProvider = &Storage{}
var _ rest.Scoper = &Storage{}
var _ rest.Creater = &Storage{}

func NewStorage() *Storage {
	return &Storage{}
}

func (r *Storage) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.ResourceKindGraphFinder)
}

func (r *Storage) NamespaceScoped() bool {
	return false
}

// Getter
func (r *Storage) New() runtime.Object {
	return &v1alpha1.GraphFinder{}
}

func (r *Storage) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	gf := obj.(*v1alpha1.GraphFinder)

	g, err := graph.LoadGraph()
	if err != nil {
		return nil, kerr.NewInternalError(err)
	}

	dist, prev := graph.Dijkstra(g, gf.Request.Source)

	out := make([]v1alpha1.Edge, 0, len(prev))

	for target, edge := range prev {
		if target != gf.Request.Source && edge != nil {
			out = append(out, v1alpha1.Edge{
				Src:        edge.Src,
				Dst:        edge.Dst,
				W:          dist[target],
				Connection: edge.Connection,
				Forward:    edge.Forward,
			})
		}
	}

	gf.Response = &v1alpha1.GraphResponse{
		Source:      gf.Request.Source,
		Connections: out,
	}
	return gf, nil
}
