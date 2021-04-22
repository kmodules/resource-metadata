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

package printers

import (
	"fmt"
	"reflect"
	"strings"

	networking "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	Register(IngressPrinter{})
}

// ref: https://github.com/kubernetes/kubernetes/blob/v1.21.0/pkg/printers/internalversion/printers.go#L203-L212

type IngressPrinter struct{}

var _ ColumnConverter = IngressPrinter{}

func (_ IngressPrinter) GVK() schema.GroupVersionKind {
	return networking.SchemeGroupVersion.WithKind("Ingress")
}

func (p IngressPrinter) Convert(o runtime.Object) (map[string]interface{}, error) {
	obj := new(networking.Ingress)
	switch to := o.(type) {
	case *unstructured.Unstructured:
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(to.UnstructuredContent(), obj); err != nil {
			return nil, err
		}
	case *networking.Ingress:
		obj = to
	default:
		return nil, fmt.Errorf("expected %v, received %v", p.GVK().Kind, reflect.TypeOf(o))
	}

	row := map[string]interface{}{}

	className := None
	if obj.Spec.IngressClassName != nil {
		className = *obj.Spec.IngressClassName
	}
	hosts := formatHosts(obj.Spec.Rules)
	address := loadBalancerStatusStringer(obj.Status.LoadBalancer)
	ports := formatPorts(obj.Spec.TLS)
	createTime := translateTimestampSince(obj.CreationTimestamp)

	row["_Name"] = obj.Name
	row["_Class"] = className
	row["Hosts"] = hosts
	row["Address"] = address
	row["Ports"] = ports
	row["_Age"] = createTime

	return row, nil
}

func formatHosts(rules []networking.IngressRule) string {
	list := []string{}
	max := 3
	more := false
	for _, rule := range rules {
		if len(list) == max {
			more = true
		}
		if !more && len(rule.Host) != 0 {
			list = append(list, rule.Host)
		}
	}
	if len(list) == 0 {
		return "*"
	}
	ret := strings.Join(list, ",")
	if more {
		return fmt.Sprintf("%s + %d more...", ret, len(rules)-max)
	}
	return ret
}

func formatPorts(tls []networking.IngressTLS) string {
	if len(tls) != 0 {
		return "80, 443"
	}
	return "80"
}
