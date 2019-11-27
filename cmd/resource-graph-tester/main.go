/*
Copyright The Kmodules Authors.

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
	"fmt"
	"log"
	"path/filepath"

	"kmodules.xyz/resource-metadata/pkg/graph"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	masterURL := ""
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

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

	if err := checkKubeDBToStatefulset(dc, g); err != nil {
		log.Fatalln(err)
	}
}

func CheckNodeToPod(dc dynamic.Interface, g *graph.Graph) error {
	podGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}
	pod, err := dc.Resource(podGVR).Namespace("kube-system").Get("kube-apiserver-minikube", metav1.GetOptions{})
	if err != nil {
		return err
	}

	nodeGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "nodes",
	}

	nodes, err := g.List(dc, *pod, nodeGVR)
	if err != nil {
		return err
	}
	for _, obj := range nodes {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetName())
	}

	return nil
}

func CheckPodToNode(dc dynamic.Interface, g *graph.Graph) error {
	nodeGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "nodes",
	}
	node, err := dc.Resource(nodeGVR).Get("minikube", metav1.GetOptions{})
	if err != nil {
		return err
	}

	podGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}

	pods, err := g.List(dc, *node, podGVR)
	if err != nil {
		return err
	}
	for _, obj := range pods {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func CheckDeployment(dc dynamic.Interface, g *graph.Graph) error {
	depGVR := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}
	busyDep, err := dc.Resource(depGVR).Namespace("default").Get("busy-dep", metav1.GetOptions{})
	if err != nil {
		return err
	}

	podGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}
	pods, err := g.List(dc, *busyDep, podGVR)
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
	services, err := g.List(dc, *busyDep, svcGVR)
	if err != nil {
		return err
	}
	for _, obj := range services {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func checkPodToPVC(dc dynamic.Interface, g *graph.Graph) error {
	podGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}
	pod, err := dc.Resource(podGVR).Namespace("default").Get("mypod", metav1.GetOptions{})
	if err != nil {
		return err
	}

	pvcGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "persistentvolumeclaims",
	}
	pvcs, err := g.List(dc, *pod, pvcGVR)
	if err != nil {
		return err
	}
	for _, obj := range pvcs {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func checkPVCToPod(dc dynamic.Interface, g *graph.Graph) error {
	pvcGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "persistentvolumeclaims",
	}
	pvc, err := dc.Resource(pvcGVR).Namespace("default").Get("myclaim", metav1.GetOptions{})
	if err != nil {
		return err
	}

	podGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}

	pods, err := g.List(dc, *pvc, podGVR)
	if err != nil {
		return err
	}
	for _, obj := range pods {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func checkDepToConfigMap(dc dynamic.Interface, g *graph.Graph) error {
	depGVR := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}
	pod, err := dc.Resource(depGVR).Namespace("default").Get("busy-dep", metav1.GetOptions{})
	if err != nil {
		return err
	}

	cfgGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "configmaps",
	}
	cfgs, err := g.List(dc, *pod, cfgGVR)
	if err != nil {
		return err
	}
	for _, obj := range cfgs {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func checkConfigMapToDep(dc dynamic.Interface, g *graph.Graph) error {
	cfgGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "configmaps",
	}
	cfg, err := dc.Resource(cfgGVR).Namespace("default").Get("omni", metav1.GetOptions{})
	if err != nil {
		return err
	}

	depGVR := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}
	deps, err := g.List(dc, *cfg, depGVR)
	if err != nil {
		return err
	}
	for _, obj := range deps {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func checkKubeDBToService(dc dynamic.Interface, g *graph.Graph) error {
	pgGVR := schema.GroupVersionResource{
		Group:    "kubedb.com",
		Version:  "v1alpha1",
		Resource: "postgreses",
	}
	pg, err := dc.Resource(pgGVR).Namespace("demo").Get("quick-postgres", metav1.GetOptions{})
	if err != nil {
		return err
	}

	svcGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "services",
	}
	services, err := g.List(dc, *pg, svcGVR)
	if err != nil {
		return err
	}
	for _, obj := range services {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func checkKubeDBToStatefulset(dc dynamic.Interface, g *graph.Graph) error {
	pgGVR := schema.GroupVersionResource{
		Group:    "kubedb.com",
		Version:  "v1alpha1",
		Resource: "postgreses",
	}
	pg, err := dc.Resource(pgGVR).Namespace("demo").Get("quick-postgres", metav1.GetOptions{})
	if err != nil {
		return err
	}

	ssGVR := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "statefulsets",
	}
	statefulsets, err := g.List(dc, *pg, ssGVR)
	if err != nil {
		return err
	}
	for _, obj := range statefulsets {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}
