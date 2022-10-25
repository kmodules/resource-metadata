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
	"context"
	"fmt"
	"log"

	"kmodules.xyz/resource-metadata/pkg/tableconvertor"

	"gomodules.xyz/encoding/json"
	crdinstall "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/install"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2/klogr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

var scheme = runtime.NewScheme()

func main() {
	_ = clientgoscheme.AddToScheme(scheme)
	crdinstall.Install(scheme)

	ctrl.SetLogger(klogr.New())
	cfg := ctrl.GetConfigOrDie()

	mapper, err := apiutil.NewDynamicRESTMapper(cfg)
	if err != nil {
		panic(err)
	}

	c, err := client.New(cfg, client.Options{
		Scheme: scheme,
		Mapper: mapper,
		Opts: client.WarningHandlerOptions{
			SuppressWarnings:   false,
			AllowDuplicateLogs: false,
		},
	})
	if err != nil {
		panic(err)
	}

	{
		var list unstructured.UnstructuredList
		list.SetAPIVersion("apps/v1")
		list.SetKind("Deployment")
		err := c.List(context.TODO(), &list, client.InNamespace("default"))
		if err != nil {
			log.Fatalln(err)
		}

		t, err := tableconvertor.TableForAnyList(c, list.Items, "", nil, nil)
		if err != nil {
			log.Fatalln(err)
		}

		data, err := json.MarshalIndent(t, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(data))
	}

	{
		var dep unstructured.Unstructured
		dep.SetAPIVersion("apps/v1")
		dep.SetKind("Deployment")
		err := c.Get(context.TODO(), client.ObjectKey{Namespace: "default", Name: "busy-dep"}, &dep)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(dep.GroupVersionKind().String())

		t, err := tableconvertor.TableForObject(c, &dep, "", nil, nil)
		if err != nil {
			log.Fatalln(err)
		}

		data, err := json.MarshalIndent(t, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(data))
	}
}
