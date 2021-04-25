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

//nolint:deadcode,unused
package main

import (
	"fmt"
	"log"
	"path/filepath"

	dynamicfactory "kmodules.xyz/client-go/dynamic/factory"
	"kmodules.xyz/resource-metadata/pkg/graph"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	masterURL := ""
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	// kubeconfigPath := "$HOME/Downloads/ui-builder-demo-kubeconfig.yaml"

	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}

	dc, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	g, err := graph.LoadGraphOfKnownResources()
	if err != nil {
		log.Fatalln(err)
	}

	f := dynamicfactory.New(dc)

	if err := CheckPodToNode(f, g); err != nil {
		log.Fatalln(err)
	}
}

func CheckNodeToPod(f dynamicfactory.Factory, g *graph.Graph) error {
	podGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}
	pod, err := f.ForResource(podGVR).Namespace("kube-system").Get("kube-apiserver-kind-control-plane")
	if err != nil {
		return err
	}

	nodeGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "nodes",
	}

	nodes, err := g.List(f, pod, nodeGVR)
	if err != nil {
		return err
	}
	for _, obj := range nodes {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetName())
	}

	return nil
}

func CheckPodToNode(f dynamicfactory.Factory, g *graph.Graph) error {
	nodeGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "nodes",
	}
	node, err := f.ForResource(nodeGVR).Get("kind-control-plane")
	if err != nil {
		return err
	}

	podGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}

	pods, err := g.List(f, node, podGVR)
	if err != nil {
		return err
	}
	for _, obj := range pods {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func CheckDeployment(f dynamicfactory.Factory, g *graph.Graph) error {
	depGVR := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}
	busyDep, err := f.ForResource(depGVR).Namespace("default").Get("busy-dep")
	if err != nil {
		return err
	}

	podGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}
	pods, err := g.List(f, busyDep, podGVR)
	if err != nil {
		return err
	}
	for _, obj := range pods {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	svcGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "services",
	}
	services, err := g.List(f, busyDep, svcGVR)
	if err != nil {
		return err
	}
	for _, obj := range services {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func checkPodToPVC(f dynamicfactory.Factory, g *graph.Graph) error {
	podGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}
	pod, err := f.ForResource(podGVR).Namespace("default").Get("mypod")
	if err != nil {
		return err
	}

	pvcGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "persistentvolumeclaims",
	}
	pvcs, err := g.List(f, pod, pvcGVR)
	if err != nil {
		return err
	}
	for _, obj := range pvcs {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func checkPVCToPod(f dynamicfactory.Factory, g *graph.Graph) error {
	pvcGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "persistentvolumeclaims",
	}
	pvc, err := f.ForResource(pvcGVR).Namespace("default").Get("myclaim")
	if err != nil {
		return err
	}

	podGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}

	pods, err := g.List(f, pvc, podGVR)
	if err != nil {
		return err
	}
	for _, obj := range pods {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func checkDepToConfigMap(f dynamicfactory.Factory, g *graph.Graph) error {
	depGVR := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}
	pod, err := f.ForResource(depGVR).Namespace("default").Get("busy-dep")
	if err != nil {
		return err
	}

	cfgGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "configmaps",
	}
	cfgs, err := g.List(f, pod, cfgGVR)
	if err != nil {
		return err
	}
	for _, obj := range cfgs {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func checkConfigMapToDep(f dynamicfactory.Factory, g *graph.Graph) error {
	cfgGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "configmaps",
	}
	cfg, err := f.ForResource(cfgGVR).Namespace("default").Get("omni")
	if err != nil {
		return err
	}

	depGVR := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}
	deps, err := g.List(f, cfg, depGVR)
	if err != nil {
		return err
	}
	for _, obj := range deps {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func checkKubeDBToService(f dynamicfactory.Factory, g *graph.Graph) error {
	pgGVR := schema.GroupVersionResource{
		Group:    "kubedb.com",
		Version:  "v1alpha1",
		Resource: "postgreses",
	}
	pg, err := f.ForResource(pgGVR).Namespace("demo").Get("quick-postgres")
	if err != nil {
		return err
	}

	svcGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "services",
	}
	services, err := g.List(f, pg, svcGVR)
	if err != nil {
		return err
	}
	for _, obj := range services {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func checkKubeDBToStatefulset(f dynamicfactory.Factory, g *graph.Graph) error {
	pgGVR := schema.GroupVersionResource{
		Group:    "kubedb.com",
		Version:  "v1alpha1",
		Resource: "postgreses",
	}
	pg, err := f.ForResource(pgGVR).Namespace("demo").Get("quick-postgres")
	if err != nil {
		return err
	}

	ssGVR := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "statefulsets",
	}
	statefulsets, err := g.List(f, pg, ssGVR)
	if err != nil {
		return err
	}
	for _, obj := range statefulsets {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func CheckBackupConfigToAppBinding(f dynamicfactory.Factory, g *graph.Graph) error {
	bkcfgGVR := schema.GroupVersionResource{
		Group:    "stash.appscode.com",
		Version:  "v1beta1",
		Resource: "backupconfigurations",
	}
	node, err := f.ForResource(bkcfgGVR).Namespace("test-namespace").Get("demo-backup-config")
	if err != nil {
		return err
	}

	appBindingGVR := schema.GroupVersionResource{
		Group:    "appcatalog.appscode.com",
		Version:  "v1alpha1",
		Resource: "appbindings",
	}

	pods, err := g.List(f, node, appBindingGVR)
	if err != nil {
		return err
	}
	for _, obj := range pods {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}
