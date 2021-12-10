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

package pathfinder

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
	return v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.ResourceKindPathFinder)
}

func (r *Storage) NamespaceScoped() bool {
	return false
}

// Getter
func (r *Storage) New() runtime.Object {
	return &v1alpha1.PathFinder{}
}

func (r *Storage) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	pf := obj.(*v1alpha1.PathFinder)

	g, err := graph.LoadGraphOfKnownResources()
	if err != nil {
		return nil, kerr.NewInternalError(err)
	}

	srcGVK, err := r.mapper.GVK(apiv1.FromMetaGVR(pf.Request.Source))
	if err != nil {
		return nil, kerr.NewInternalError(err)
	}
	dist, prev := graph.Dijkstra(g, srcGVK)
	paths := graph.GeneratePaths(srcGVK, dist, prev)
	dstGVK, err := r.mapper.GVK(apiv1.FromMetaGVR(*pf.Request.Target))
	if err != nil {
		return nil, kerr.NewInternalError(err)
	}

	out := make([]v1alpha1.Path, 0, len(paths))

	if pf.Request.Target != nil {
		path, ok := paths[dstGVK]
		if ok {
			cp, err := r.convertPath(*path)
			if err != nil {
				return nil, kerr.NewInternalError(err)
			}
			out = append(out, cp)
		}
	} else {
		for i := range paths {
			cp, err := r.convertPath(*paths[i])
			if err != nil {
				return nil, kerr.NewInternalError(err)
			}
			out = append(out, cp)
		}
	}

	pf.Response = &v1alpha1.PathResponse{Paths: out}
	return pf, nil
}

func (r *Storage) convertPath(in graph.Path) (v1alpha1.Path, error) {
	srcGVR, err := r.mapper.GVR(in.Source)
	if err != nil {
		return v1alpha1.Path{}, err
	}
	dstGVR, err := r.mapper.GVR(in.Target)
	if err != nil {
		return v1alpha1.Path{}, err
	}
	out := v1alpha1.Path{
		Source:   apiv1.ToMetaGVR(srcGVR),
		Target:   apiv1.ToMetaGVR(dstGVR),
		Distance: in.Distance,
		Edges:    make([]*v1alpha1.Edge, len(in.Edges)),
	}

	for i := range in.Edges {
		if out.Edges[i], err = r.convertEdge(in.Edges[i]); err != nil {
			return v1alpha1.Path{}, err
		}
	}

	return out, nil
}

func (r *Storage) convertEdge(in *graph.Edge) (*v1alpha1.Edge, error) {
	srcGVR, err := r.mapper.GVR(in.Src)
	if err != nil {
		return nil, err
	}
	dstGVR, err := r.mapper.GVR(in.Dst)
	if err != nil {
		return nil, err
	}
	return &v1alpha1.Edge{
		Src:        apiv1.ToMetaGVR(srcGVR),
		Dst:        apiv1.ToMetaGVR(dstGVR),
		W:          in.W,
		Connection: in.Connection,
		Forward:    in.Forward,
	}, nil
}
