package main

import (
	goflag "flag"
	"fmt"
	"log"

	"github.com/emicklei/dot"
	flag "github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kmodules.xyz/resource-metadata/pkg/graph"
)

func main() {
	var src metav1.TypeMeta
	flag.StringVar(&src.APIVersion, "apiVersion", "", "apiVersion of resource who paths will be rendered")
	flag.StringVar(&src.Kind, "kind", "", "Kind of resource who paths will be rendered")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	if src.APIVersion == "" || src.Kind == "" {
		log.Fatalln("--apiVersion, --kind can't bt empty")
	}

	g, err := graph.LoadGraph()
	if err != nil {
		log.Fatalln(err)
	}
	_, prev := graph.Dijkstra(g, src)

	gv := dot.NewGraph(dot.Directed)
	nodes := make(map[metav1.TypeMeta]dot.Node)

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
