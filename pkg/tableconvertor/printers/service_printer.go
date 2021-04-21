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

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
)

func init() {
	Register(ServicePrinter{})
}

// ref: https://github.com/kubernetes/kubernetes/blob/v1.21.0/pkg/printers/internalversion/printers.go#L190-L198

type ServicePrinter struct{}

var _ ColumnConverter = ServicePrinter{}

func (_ ServicePrinter) GVK() schema.GroupVersionKind {
	return core.SchemeGroupVersion.WithKind("Service")
}

func (p ServicePrinter) Convert(o runtime.Object) (map[string]interface{}, error) {
	obj := new(core.Service)
	switch to := o.(type) {
	case *unstructured.Unstructured:
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(to.UnstructuredContent(), obj); err != nil {
			return nil, err
		}
	case *core.Service:
		obj = to
	default:
		return nil, fmt.Errorf("expected %v, received %v", p.GVK().Kind, reflect.TypeOf(o))
	}

	row := map[string]interface{}{}

	svcType := obj.Spec.Type
	// internalIP := None
	/*
		if len(obj.Spec.ClusterIPs) > 0 {
			internalIP = obj.Spec.ClusterIPs[0]
		}
	*/
	internalIP := obj.Spec.ClusterIP

	externalIP := getServiceExternalIP(obj)
	svcPorts := MakePortString(obj.Spec.Ports)
	if len(svcPorts) == 0 {
		svcPorts = None
	}

	row["Name"] = obj.Name
	row["Type"] = string(svcType)
	row["Cluster-IP"] = internalIP
	row["External-IP"] = externalIP
	row["Port(s)"] = svcPorts
	row["Age"] = translateTimestampSince(obj.CreationTimestamp)
	row["Selector"] = labels.FormatLabels(obj.Spec.Selector)

	return row, nil
}

func getServiceExternalIP(svc *core.Service) string {
	switch svc.Spec.Type {
	case core.ServiceTypeClusterIP:
		if len(svc.Spec.ExternalIPs) > 0 {
			return strings.Join(svc.Spec.ExternalIPs, ",")
		}
		return None
	case core.ServiceTypeNodePort:
		if len(svc.Spec.ExternalIPs) > 0 {
			return strings.Join(svc.Spec.ExternalIPs, ",")
		}
		return None
	case core.ServiceTypeLoadBalancer:
		lbIps := loadBalancerStatusStringer(svc.Status.LoadBalancer)
		if len(svc.Spec.ExternalIPs) > 0 {
			results := []string{}
			if len(lbIps) > 0 {
				results = append(results, strings.Split(lbIps, ",")...)
			}
			results = append(results, svc.Spec.ExternalIPs...)
			return strings.Join(results, ",")
		}
		if len(lbIps) > 0 {
			return lbIps
		}
		return "<pending>"
	case core.ServiceTypeExternalName:
		return svc.Spec.ExternalName
	}
	return "<unknown>"
}

// loadBalancerStatusStringer behaves mostly like a string interface and converts the given status to a string.
// `wide` indicates whether the returned value is meant for --o=wide output. If not, it's clipped to 16 bytes.
func loadBalancerStatusStringer(s core.LoadBalancerStatus) string {
	ingress := s.Ingress
	result := sets.NewString()
	for i := range ingress {
		if ingress[i].IP != "" {
			result.Insert(ingress[i].IP)
		} else if ingress[i].Hostname != "" {
			result.Insert(ingress[i].Hostname)
		}
	}

	return strings.Join(result.List(), ",")
}

func MakePortString(ports []core.ServicePort) string {
	pieces := make([]string, len(ports))
	for ix := range ports {
		port := &ports[ix]
		pieces[ix] = fmt.Sprintf("%d/%s", port.Port, port.Protocol)
		if port.NodePort > 0 {
			pieces[ix] = fmt.Sprintf("%d:%d/%s", port.Port, port.NodePort, port.Protocol)
		}
	}
	return strings.Join(pieces, ",")
}
