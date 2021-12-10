/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package graphfinder

import (
	"context"

	apiv1 "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/discovery"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"kmodules.xyz/resource-metadata/pkg/graph"

	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/rest"
)

type Storage struct {
	mapper discovery.ResourceMapper
}

var _ rest.GroupVersionKindProvider = &Storage{}
var _ rest.Scoper = &Storage{}
var _ rest.Creater = &Storage{}

func NewStorage(mapper discovery.ResourceMapper) *Storage {
	return &Storage{mapper: mapper}
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

	g, err := graph.LoadGraphOfKnownResources()
	if err != nil {
		return nil, kerr.NewInternalError(err)
	}

	srcGVK, err := r.mapper.GVK(apiv1.FromMetaGVR(gf.Request.Source))
	if err != nil {
		return nil, kerr.NewInternalError(err)
	}
	dist, prev := graph.Dijkstra(g, srcGVK)

	out := make([]*v1alpha1.Edge, 0, len(prev))

	for target, edge := range prev {
		if target != srcGVK && edge != nil {
			srcGVR, err := r.mapper.GVR(edge.Src)
			if err != nil {
				return nil, kerr.NewInternalError(err)
			}
			dstGVR, err := r.mapper.GVR(edge.Dst)
			if err != nil {
				return nil, kerr.NewInternalError(err)
			}
			out = append(out, &v1alpha1.Edge{
				Src:        apiv1.ToMetaGVR(srcGVR),
				Dst:        apiv1.ToMetaGVR(dstGVR),
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
