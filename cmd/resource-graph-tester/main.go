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
	"context"
	"fmt"
	"log"
	"path/filepath"

	"kmodules.xyz/resource-metadata/hub"
	"kmodules.xyz/resource-metadata/pkg/graph"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2/klogr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

func NewClient(cfg *rest.Config) (client.Client, error) {
	scheme := runtime.NewScheme()

	_ = clientgoscheme.AddToScheme(scheme)
	ctrl.SetLogger(klogr.New())

	mapper, err := apiutil.NewDynamicRESTMapper(cfg)
	if err != nil {
		return nil, err
	}

	return client.New(cfg, client.Options{
		Scheme: scheme,
		Mapper: mapper,
		Opts: client.WarningHandlerOptions{
			SuppressWarnings:   false,
			AllowDuplicateLogs: false,
		},
	})
}

func main() {
	masterURL := ""
	// kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	kubeconfigPath := "/home/tamal/Downloads/kubedb-demo-ui-kubeconfig.yaml"

	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}

	reg := hub.NewRegistryOfKnownResources()
	err = reg.DiscoverResources(config)
	if err != nil {
		log.Fatalln(err)
	}

	gvr := schema.GroupVersionKind{
		Group:   "kubedb.com",
		Version: "v1alpha2",
		Kind:    "MongoDB",
	}
	//gvr := schema.GroupVersionKind{
	//	Group:    "",
	//	Version:  "v1",
	//	Resource: "pods",
	//}
	//edges, err := graph.GetConnectedGraph(config, reg, gvr, "kube-apiserver-kind-control-plane", "kube-system")

	f, err := NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}
	edges, err := graph.GetConnectedGraph(config, f, reg, gvr, types.NamespacedName{Namespace: "default", Name: "mongo-rs"})
	if err != nil {
		log.Fatalln(err)
	}
	for _, edge := range edges {
		fmt.Printf("%v -> %v\n", edge.Src, edge.Dst)
	}
}

func main_list() {
	masterURL := ""
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	// kubeconfigPath := "$HOME/Downloads/ui-builder-demo-kubeconfig.yaml"

	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}

	g, err := graph.LoadGraphOfKnownResources()
	if err != nil {
		log.Fatalln(err)
	}

	f, err := NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	if err := CheckNodeToPod(f, g); err != nil {
		log.Fatalln(err)
	}
}

func CheckNodeToPod(f client.Client, g *graph.Graph) error {
	podGVR := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Pod",
	}
	var pod unstructured.Unstructured
	pod.SetGroupVersionKind(podGVR)
	err := f.Get(context.TODO(), client.ObjectKey{
		Name:      "kube-apiserver-kind-control-plane",
		Namespace: "kube-system",
	}, &pod)
	if err != nil {
		return err
	}

	nodeGVR := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Node",
	}

	nodes, err := g.List(f, &pod, nodeGVR)
	if err != nil {
		return err
	}
	for _, obj := range nodes {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetName())
	}

	return nil
}

func CheckPodToNode(f client.Client, g *graph.Graph) error {
	nodeGVR := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Node",
	}
	var node unstructured.Unstructured
	node.SetGroupVersionKind(nodeGVR)
	err := f.Get(context.TODO(), client.ObjectKey{
		Name: "kind-control-plane",
	}, &node)
	if err != nil {
		return err
	}

	podGVR := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Pod",
	}

	pods, err := g.List(f, &node, podGVR)
	if err != nil {
		return err
	}
	for _, obj := range pods {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func CheckDeployment(f client.Client, g *graph.Graph) error {
	depGVR := schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "Deployment",
	}
	var busyDep unstructured.Unstructured
	busyDep.SetGroupVersionKind(depGVR)
	err := f.Get(context.TODO(), client.ObjectKey{
		Name:      "busy-dep",
		Namespace: "default",
	}, &busyDep)
	if err != nil {
		return err
	}

	podGVR := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Pod",
	}
	pods, err := g.List(f, &busyDep, podGVR)
	if err != nil {
		return err
	}
	for _, obj := range pods {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	svcGVR := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Service",
	}
	services, err := g.List(f, &busyDep, svcGVR)
	if err != nil {
		return err
	}
	for _, obj := range services {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func checkPodToPVC(f client.Client, g *graph.Graph) error {
	podGVR := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Pod",
	}
	var pod unstructured.Unstructured
	pod.SetGroupVersionKind(podGVR)
	err := f.Get(context.TODO(), client.ObjectKey{
		Name:      "mypod",
		Namespace: "default",
	}, &pod)
	if err != nil {
		return err
	}

	pvcGVR := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "PersistentVolumeClaim",
	}
	pvcs, err := g.List(f, &pod, pvcGVR)
	if err != nil {
		return err
	}
	for _, obj := range pvcs {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func checkPVCToPod(f client.Client, g *graph.Graph) error {
	pvcGVR := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "PersistentVolumeClaim",
	}
	var pvc unstructured.Unstructured
	pvc.SetGroupVersionKind(pvcGVR)
	err := f.Get(context.TODO(), client.ObjectKey{
		Name:      "myclaim",
		Namespace: "default",
	}, &pvc)
	if err != nil {
		return err
	}

	podGVR := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Pod",
	}

	pods, err := g.List(f, &pvc, podGVR)
	if err != nil {
		return err
	}
	for _, obj := range pods {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func checkDepToConfigMap(f client.Client, g *graph.Graph) error {
	depGVR := schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "Deployment",
	}
	var pod unstructured.Unstructured
	pod.SetGroupVersionKind(depGVR)
	err := f.Get(context.TODO(), client.ObjectKey{
		Name:      "busy-dep",
		Namespace: "default",
	}, &pod)
	if err != nil {
		return err
	}

	cfgGVR := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "ConfigMap",
	}
	cfgs, err := g.List(f, &pod, cfgGVR)
	if err != nil {
		return err
	}
	for _, obj := range cfgs {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func checkConfigMapToDep(f client.Client, g *graph.Graph) error {
	cfgGVR := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "ConfigMap",
	}
	var cfg unstructured.Unstructured
	cfg.SetGroupVersionKind(cfgGVR)
	err := f.Get(context.TODO(), client.ObjectKey{
		Name:      "omni",
		Namespace: "default",
	}, &cfg)
	if err != nil {
		return err
	}

	depGVR := schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "Deployment",
	}
	deps, err := g.List(f, &cfg, depGVR)
	if err != nil {
		return err
	}
	for _, obj := range deps {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func checkKubeDBToService(f client.Client, g *graph.Graph) error {
	pgGVR := schema.GroupVersionKind{
		Group:   "kubedb.com",
		Version: "v1alpha1",
		Kind:    "Postgres",
	}
	var pg unstructured.Unstructured
	pg.SetGroupVersionKind(pgGVR)
	err := f.Get(context.TODO(), client.ObjectKey{
		Name:      "quick-postgres",
		Namespace: "demo",
	}, &pg)
	if err != nil {
		return err
	}

	svcGVR := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Service",
	}
	services, err := g.List(f, &pg, svcGVR)
	if err != nil {
		return err
	}
	for _, obj := range services {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func checkKubeDBToStatefulset(f client.Client, g *graph.Graph) error {
	pgGVR := schema.GroupVersionKind{
		Group:   "kubedb.com",
		Version: "v1alpha1",
		Kind:    "Postgrese",
	}
	var pg unstructured.Unstructured
	pg.SetGroupVersionKind(pgGVR)
	err := f.Get(context.TODO(), client.ObjectKey{
		Name:      "quick-postgres",
		Namespace: "demo",
	}, &pg)
	if err != nil {
		return err
	}

	ssGVR := schema.GroupVersionKind{
		Group:   "apps",
		Version: "v1",
		Kind:    "StatefulSet",
	}
	statefulsets, err := g.List(f, &pg, ssGVR)
	if err != nil {
		return err
	}
	for _, obj := range statefulsets {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}

func CheckBackupConfigToAppBinding(f client.Client, g *graph.Graph) error {
	bkcfgGVK := schema.GroupVersionKind{
		Group:   "stash.appscode.com",
		Version: "v1beta1",
		Kind:    "BackupConfiguration",
	}

	var node unstructured.Unstructured
	node.SetGroupVersionKind(bkcfgGVK)
	err := f.Get(context.TODO(), client.ObjectKey{
		Name:      "demo-backup-config",
		Namespace: "test-namespace",
	}, &node)
	if err != nil {
		return err
	}

	appBindingGVK := schema.GroupVersionKind{
		Group:   "appcatalog.appscode.com",
		Version: "v1alpha1",
		Kind:    "AppBinding",
	}

	pods, err := g.List(f, &node, appBindingGVK)
	if err != nil {
		return err
	}
	for _, obj := range pods {
		fmt.Println(obj.GetObjectKind().GroupVersionKind(), ":", obj.GetNamespace()+"/"+obj.GetName())
	}

	return nil
}
