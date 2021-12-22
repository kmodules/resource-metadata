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

package graph

import (
	"context"
	"fmt"
	"time"

	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/tools/clientcache"
	"kmodules.xyz/resource-metadata/hub"

	"github.com/gregjones/httpcache"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

/*
- Handle when multiple version of resources are available
- How to handle preferred version path missing
*/

func GetConnectedGraph(config *rest.Config, f client.Client, reg *hub.Registry, srcGVK schema.GroupVersionKind, ref types.NamespacedName) ([]*Edge, error) {
	cfg := clientcache.ConfigFor(config, 5*time.Minute, httpcache.NewMemoryCache())

	srcGVR, err := reg.GVR(srcGVK)
	if err != nil {
		return nil, err
	}
	if err := reg.Register(srcGVR, cfg); err != nil {
		return nil, err
	}
	rd, err := reg.LoadByGVK(srcGVK)
	if err != nil {
		return nil, err
	}

	objkey := client.ObjectKey{Name: ref.Name}
	if rd.Spec.Resource.Scope == kmapi.NamespaceScoped {
		if ref.Namespace == "" {
			return nil, fmt.Errorf("missing namespace query parameter for %s with name %s", srcGVK, ref.Name)
		}
		objkey.Namespace = ref.Namespace
	}
	var src unstructured.Unstructured
	err = f.Get(context.TODO(), objkey, &src)
	if err != nil {
		return nil, err
	}

	g, err := LoadGraph(reg)
	if err != nil {
		return nil, err
	}
	realGraph, err := g.generateRealGraph(f, &src)
	if err != nil {
		return nil, err
	}

	dist, prev := Dijkstra(realGraph, srcGVK)

	out := make([]*Edge, 0, len(prev))
	for target, edge := range prev {
		if target != srcGVK && edge != nil {
			out = append(out, &Edge{
				Src:        edge.Src,
				Dst:        edge.Dst,
				W:          dist[target],
				Connection: edge.Connection,
				Forward:    edge.Forward,
			})
		}
	}
	return out, nil
}

// getRealGraph runs BFS on the original graph and returns a graph that has real connection
// with the source resource.
func (g *Graph) generateRealGraph(f client.Client, src *unstructured.Unstructured) (*Graph, error) {
	srcGVK := schema.FromAPIVersionAndKind(src.GetAPIVersion(), src.GetKind())
	objMap := map[schema.GroupVersionKind][]*unstructured.Unstructured{
		srcGVK: {src},
	}

	visited := map[schema.GroupVersionKind]bool{}
	realGraph := NewGraph(g.r)
	// Queue for the BSF
	q := make([]schema.GroupVersionKind, 0)

	// Push the source node
	q = append(q, srcGVK)
	visited[srcGVK] = true
	finder := ObjectFinder{
		Client: f,
	}
	for {
		// Pop the first item
		u := q[0]
		q = q[1:]
		for v, e := range g.edges[u] {
			if !visited[v] {
				// Find the connected objects. The object might be connected via multiple paths.
				// Hence, we are checking connection from all the child object of u.
				srcObjects := objMap[u]
				var dstObjects []*unstructured.Unstructured
				for _, srcObj := range srcObjects {
					objects, err := finder.ResourcesFor(srcObj, e)
					if err != nil && !kerr.IsNotFound(err) {
						return nil, err
					}
					dstObjects = appendObjects(dstObjects, objects...)
				}
				if len(dstObjects) > 0 {
					// Real edge exists, we need to traverse. So, add it to the queue.
					q = append(q, v)
					realGraph.AddEdge(e)
					objMap[v] = dstObjects
					visited[v] = true
				}
			}
		}
		if len(q) == 0 {
			break
		}
	}
	return realGraph, nil
}
