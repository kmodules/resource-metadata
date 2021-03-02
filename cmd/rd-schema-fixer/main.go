package main

import (
	"fmt"
	"io/ioutil"

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
)

func main() {
	filename := "/home/tamal/go/src/kmodules.xyz/resource-metadata/hub/resourcedescriptors/kubedb.com/v1alpha2/elasticsearches.yaml"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	data, err = v1alpha1.FormatMetadata(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
