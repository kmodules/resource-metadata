package graph

import (
	"strings"

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	hub "kmodules.xyz/resource-metadata/hub/v1alpha1"

	"sigs.k8s.io/yaml"
)

func LoadGraph() (*Graph, error) {
	graph := NewGraph()

	for _, f := range hub.AssetNames() {
		data, err := hub.Asset(f)
		if err != nil {
			return nil, err
		}

		var rd v1alpha1.ResourceDescriptor
		err = yaml.UnmarshalStrict(data, &rd)
		if err != nil {
			return nil, err
		}
		addRDConnectionsToGraph(graph, &rd)
	}
	cache := hub.GetCachedResourceDescriptor()
	for _, rd := range cache {
		addRDConnectionsToGraph(graph, rd)
	}
	return graph, nil
}

func addRDConnectionsToGraph(graph *Graph, rd *v1alpha1.ResourceDescriptor) {
	src := rd.Spec.Resource.GroupVersionResource()
	for _, conn := range rd.Spec.Connections {
		dst := conn.Target
		dstGVR := hub.GVR(dst.GroupVersionKind())

		var w uint64 = 1
		if conn.ResourceConnectionSpec.Type == v1alpha1.MatchSelector &&
			conn.TargetLabelPath != "" &&
			strings.Trim(conn.TargetLabelPath, ".") != "metadata.labels" {
			w = 1 + CostFactorOfInAppFiltering
		}

		graph.AddEdge(&Edge{
			Src:        src,
			Dst:        dstGVR,
			W:          w,
			Connection: conn.ResourceConnectionSpec,
			Forward:    true,
		})

		backEdge := &Edge{
			Src:        dstGVR,
			Dst:        src,
			Connection: conn.ResourceConnectionSpec,
			Forward:    false,
		}
		switch conn.Type {
		case v1alpha1.MatchName:
			backEdge.W = 1
		case v1alpha1.MatchSelector, v1alpha1.OwnedBy:
			backEdge.W = 1 + CostFactorOfInAppFiltering
		case v1alpha1.MatchRef:
			backEdge.W = 1 + CostFactorOfInAppFiltering<<1
		}

		graph.AddEdge(backEdge)
	}
}
