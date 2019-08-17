package main

import (
	"os"
	"path/filepath"

	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/klog"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
)

func generateCRDDefinitions() {
	v1beta1CRDs := []*crd_api.CustomResourceDefinition{
		v1alpha1.ResourceDescriptor{}.CustomResourceDefinition(),
	}
	genCRD(v1alpha1.SchemeGroupVersion.Version, v1beta1CRDs)

}

func genCRD(version string, crds []*crd_api.CustomResourceDefinition) {
	err := os.MkdirAll(filepath.Join("/src/api/crds", version), 0755)
	if err != nil {
		klog.Fatal(err)
	}

	for _, crd := range crds {
		filename := filepath.Join("/src/api/crds", version, crd.Spec.Names.Singular+".yaml")
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			klog.Fatal(err)
		}
		crdutils.MarshallCrd(f, crd, "yaml")
		f.Close()
	}
}

func main() {
	generateCRDDefinitions()
}
