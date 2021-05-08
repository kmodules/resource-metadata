package main

import (
	"fmt"

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
)

func main() {
	crd := (v1alpha1.ResourceDescriptor{}).CustomResourceDefinition()
	fmt.Println(crd)
}
