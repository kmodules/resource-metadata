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

		src := rd.Spec.Resource.TypeMeta()

		for _, conn := range rd.Spec.Connections {
			dst := conn.Target

			var w uint64 = 1
			if conn.ResourceConnectionSpec.Type == v1alpha1.MatchSelector &&
				conn.TargetLabelPath != "" &&
				strings.Trim(conn.TargetLabelPath, ".") != "metadata.labels" {
				w = 1 + CostFactorOfInAppFiltering
			}

			graph.AddEdge(&Edge{
				Src:        src,
				Dst:        dst,
				W:          w,
				Connection: conn.ResourceConnectionSpec,
				Forward:    true,
			})

			if conn.Type == v1alpha1.MatchSelector || conn.Type == v1alpha1.OwnedBy {
				graph.AddEdge(&Edge{
					Src:        dst,
					Dst:        src,
					W:          1 + CostFactorOfInAppFiltering,
					Connection: conn.ResourceConnectionSpec,
					Forward:    false,
				})
			} else if conn.Type == v1alpha1.MatchName {
				graph.AddEdge(&Edge{
					Src:        dst,
					Dst:        src,
					W:          1,
					Connection: conn.ResourceConnectionSpec,
					Forward:    false,
				})
			} else if conn.Type == v1alpha1.MatchRef {
				graph.AddEdge(&Edge{
					Src:        dst,
					Dst:        src,
					W:          1 + CostFactorOfInAppFiltering<<1,
					Connection: conn.ResourceConnectionSpec,
					Forward:    false,
				})
			}
		}
	}
	return graph, nil
}