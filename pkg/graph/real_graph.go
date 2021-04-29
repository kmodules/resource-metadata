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
	dynamicfactory "kmodules.xyz/client-go/dynamic/factory"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"kmodules.xyz/resource-metadata/hub"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var empty = struct{}{}

/*
- Handle when multiple version of resources are available
- How to handle preferred version path missing
*/
func LoadFromCluster(f dynamicfactory.Factory, r *hub.Registry, src *unstructured.Unstructured) (*Graph, error) {
	g := NewGraph(r)

	srcGVK := schema.FromAPIVersionAndKind(src.GetAPIVersion(), src.GetKind())
	srcGVR, err := g.r.GVR(srcGVK)
	if err != nil {
		return nil, err
	}

	travered := map[schema.GroupResource]struct{}{}
	toBeTravered := []schema.GroupVersionResource{srcGVR}
	objMap := map[schema.GroupVersionResource][]*unstructured.Unstructured{
		srcGVR: {src},
	}

	g.r.Visit(func(key string, rd *v1alpha1.ResourceDescriptor) {
		gvr := rd.Spec.Resource.GroupVersionResource()
		if _, exists := travered[gvr.GroupResource()]; exists {
			return
		}

		for _, conn := range rd.Spec.Connections {
			dst := conn.Target
			dstGVR, err := g.r.GVR(dst.GroupVersionKind())
			if err != nil {
				// TODO: should panic ?
				panic(err)
			}
			if dstGVR != srcGVR {
				continue
			}

			backEdge := &Edge{
				Src:        srcGVR, // == dstGVR
				Dst:        gvr,
				W:          getWeight(conn.Type),
				Connection: conn.ResourceConnectionSpec,
				Forward:    false,
			}
			objects, err := g.ResourcesFor(f, src, backEdge)
			if err != nil {
				panic(err)
			}
			if len(objects) > 0 {
				// real edge exists, so need to traverse
				toBeTravered = append(toBeTravered, gvr)
				objMap[gvr] = objects
			}

			break // since we found the edge to the srcGVR
		}
	})

	for {
		var gvr schema.GroupVersionResource

		// https://github.com/golang/go/wiki/SliceTricks
		gvr, toBeTravered = toBeTravered[0], toBeTravered[1:]
		srcObjects := objMap[gvr]

		rd, err := g.r.LoadByGVR(gvr)
		if err != nil {
			return nil, err
		}
		for _, conn := range rd.Spec.Connections {
			dst := conn.Target
			dstGVR, err := g.r.GVR(dst.GroupVersionKind())
			if err != nil {
				return nil, err
			}

			if dstGVR == srcGVR {
				continue // already added to graph in g.r.Visit(...)
			}

			edge := &Edge{
				Src:        srcGVR,
				Dst:        dstGVR,
				W:          getWeight(conn.Type),
				Connection: conn.ResourceConnectionSpec,
				Forward:    true,
			}

			var dstObjects []*unstructured.Unstructured
			for _, srcObj := range srcObjects {
				objects, err := g.ResourcesFor(f, srcObj, edge)
				if err != nil {
					return nil, err
				}
				dstObjects = appendObjects(dstObjects, objects...)
			}
			if len(dstObjects) > 0 {
				g.AddEdge(edge)
				backEdge := &Edge{
					Src:        dstGVR,
					Dst:        srcGVR,
					W:          getWeight(conn.Type),
					Connection: conn.ResourceConnectionSpec,
					Forward:    false,
				}
				g.AddEdge(backEdge)

				if _, exists := travered[dstGVR.GroupResource()]; !exists {
					toBeTravered = append(toBeTravered, dstGVR)
					objMap[dstGVR] = dstObjects
				}
			}
		}
		travered[gvr.GroupResource()] = empty

		if len(toBeTravered) == 0 {
			break
		}
	}
	return g, nil
}