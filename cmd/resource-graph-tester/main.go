package main

import (
	"fmt"
	"log"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"kmodules.xyz/resource-metadata/pkg/graph"
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

	g, err := graph.LoadGraph()
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

	nodeType := metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Node",
	}
	nodes, err := g.List(dc, pod, nodeType)
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

	podType := metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Pod",
	}
	pods, err := g.List(dc, node, podType)
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

	podType := metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Pod",
	}
	pods, err := g.List(dc, busyDep, podType)
	if err != nil {
		return err
	}
	for _, obj := range pods {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	svcType := metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	}
	services, err := g.List(dc, busyDep, svcType)
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

	pvcType := metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "PersistentVolumeClaim",
	}
	pvcs, err := g.List(dc, pod, pvcType)
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

	podType := metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Pod",
	}
	pods, err := g.List(dc, pvc, podType)
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

	cfgType := metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ConfigMap",
	}
	cfgs, err := g.List(dc, pod, cfgType)
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

	depType := metav1.TypeMeta{
		APIVersion: "apps/v1",
		Kind:       "Deployment",
	}
	deps, err := g.List(dc, cfg, depType)
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

	svcType := metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Service",
	}
	services, err := g.List(dc, pg, svcType)
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

	ssType := metav1.TypeMeta{
		APIVersion: "apps/v1",
		Kind:       "StatefulSet",
	}
	statefulsets, err := g.List(dc, pg, ssType)
	if err != nil {
		return err
	}
	for _, obj := range statefulsets {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}
