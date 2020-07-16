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

package main

import (
	goflag "flag"
	"fmt"
	"log"

	"kmodules.xyz/resource-metadata/pkg/graph"

	"github.com/emicklei/dot"
	flag "github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func main() {
	var src schema.GroupVersionResource
	flag.StringVar(&src.Group, "group", "", "Group of resource who paths will be rendered")
	flag.StringVar(&src.Version, "version", "", "Version of resource who paths will be rendered")
	flag.StringVar(&src.Resource, "resource", "", "Name (plural) of resource who paths will be rendered")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	if src.Group == "" || src.Version == "" || src.Resource == "" {
		log.Fatalln("--group, --version, --resource can't bt empty")
	}

	g, err := graph.LoadGraphOfKnownResources()
	if err != nil {
		log.Fatalln(err)
	}
	_, prev := graph.Dijkstra(g, src)

	gv := dot.NewGraph(dot.Directed)
	nodes := make(map[schema.GroupVersionResource]dot.Node)

	for target, edge := range prev {
		if target != src && edge != nil {
			n1, ok := nodes[edge.Src]
			if !ok {
				n1 = gv.Node(edge.Src.String())
			}
			n2, ok := nodes[edge.Dst]
			if !ok {
				n2 = gv.Node(edge.Dst.String())
			}

			e2 := gv.Edge(n1, n2,
				string(edge.Connection.Type),
				fmt.Sprintf("distnace=%d", edge.W))
			if !edge.Forward {
				e2.Attr("color", "red")
			}
		}
	}
	fmt.Println(gv.String())
}
