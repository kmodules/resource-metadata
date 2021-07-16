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

package tableconvertor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"kmodules.xyz/resource-metadata/pkg/tableconvertor/printers"

	"github.com/Masterminds/sprig/v3"
	prom_op "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gomodules.xyz/jsonpath"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	metatable "k8s.io/apimachinery/pkg/api/meta/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/duration"
)

var templateFns = sprig.TxtFuncMap()

func init() {
	templateFns["jp"] = jsonpathFn
	templateFns["k8s_fmt_selector"] = formatLabelSelectorFn
	templateFns["k8s_fmt_label"] = formatLabelsFn
	templateFns["k8s_age"] = ageFn
	templateFns["k8s_svc_ports"] = servicePortsFn
	templateFns["k8s_container_ports"] = containerPortFn
	templateFns["k8s_container_images"] = containerImagesFn
	templateFns["k8s_volumes"] = volumesFn
	templateFns["k8s_volumeMounts"] = volumeMountsFn
	templateFns["k8s_duration"] = durationFn
	templateFns["fmt_list"] = fmtListFn
	templateFns["prom_ns_selector"] = promNamespaceSelectorFn
	templateFns["map_key_count"] = mapKeyCountFn
	templateFns["kubedb_db_mode"] = kubedbDBModeFn
	templateFns["kubedb_db_replicas"] = kubedbDBReplicasFn
	templateFns["kubedb_db_resources"] = kubedbDBResourcesFn
	templateFns["rbac_subjects"] = rbacSubjects
	templateFns["cert_validity"] = certificateValidity
}

// TxtFuncMap returns a 'text/template'.FuncMap
func TxtFuncMap() template.FuncMap {
	gfm := make(map[string]interface{}, len(templateFns))
	for k, v := range templateFns {
		gfm[k] = v
	}
	return gfm
}

func jsonpathFn(expr string, data interface{}, jsonoutput ...bool) (interface{}, error) {
	enableJSONoutput := len(jsonoutput) > 0 && jsonoutput[0]

	jp := jsonpath.New("jp")
	if err := jp.Parse(expr); err != nil {
		return nil, fmt.Errorf("unrecognized column definition %q", expr)
	}
	jp.AllowMissingKeys(true)
	jp.EnableJSONOutput(enableJSONoutput)

	var buf bytes.Buffer
	err := jp.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	if enableJSONoutput {
		var v []interface{}
		err = json.Unmarshal(buf.Bytes(), &v)
		return v, err
	}
	return buf.String(), err
}

func formatLabelSelectorFn(data string) (string, error) {
	if strings.TrimSpace(data) == "" {
		return "", nil
	}
	var sel metav1.LabelSelector
	err := json.Unmarshal([]byte(data), &sel)
	if err != nil {
		return "", err
	}
	return metav1.FormatLabelSelector(&sel), nil
}

func formatLabelsFn(data string) (string, error) {
	if strings.TrimSpace(data) == "" {
		return "", nil
	}
	var label map[string]string
	err := json.Unmarshal([]byte(data), &label)
	if err != nil {
		return "", err
	}
	return labels.FormatLabels(label), nil
}

func ageFn(data string) (string, error) {
	if data == "" {
		return "", nil
	}
	var timestamp metav1.Time
	err := timestamp.UnmarshalQueryParameter(data)
	if err != nil {
		return "", err
	}
	return metatable.ConvertToHumanReadableDateType(timestamp), nil
}

func servicePortsFn(data string) (string, error) {
	if strings.TrimSpace(data) == "" {
		return "", nil
	}
	var ports []core.ServicePort
	err := json.Unmarshal([]byte(data), &ports)
	if err != nil {
		return "", err
	}
	return printers.MakeServicePortString(ports), nil
}

func containerPortFn(data string) (string, error) {
	if strings.TrimSpace(data) == "" {
		return "", nil
	}
	var ports []core.ContainerPort
	err := json.Unmarshal([]byte(data), &ports)
	if err != nil {
		return "", err
	}
	pieces := make([]string, len(ports))
	for ix := range ports {
		port := &ports[ix]
		pieces[ix] = fmt.Sprintf("%d/%s", port.ContainerPort, port.Protocol)
		if port.HostPort > 0 {
			pieces[ix] = fmt.Sprintf("%d:%d/%s", port.ContainerPort, port.HostPort, port.Protocol)
		}
	}
	return strings.Join(pieces, ","), nil
}

func volumesFn(data string) (string, error) {
	if strings.TrimSpace(data) == "" {
		return "", nil
	}
	var volumes []core.Volume
	err := json.Unmarshal([]byte(data), &volumes)
	if err != nil {
		return "", err
	}
	ss := "["
	for i := range volumes {
		ss += describeVolume(volumes[i])
		if i < len(volumes)-1 {
			ss += ","
		}
	}
	ss += "]"
	return ss, nil
}

func volumeMountsFn(data string) (string, error) {
	if strings.TrimSpace(data) == "" {
		return "", nil
	}
	var mounts []core.VolumeMount
	ss := make([]string, 0)
	err := json.Unmarshal([]byte(data), &mounts)
	if err != nil {
		return "", err
	}

	for i := range mounts {
		mnt := fmt.Sprintf("%s:%s", mounts[i].Name, mounts[i].MountPath)
		if mounts[i].SubPath != "" {
			mnt = fmt.Sprintf("%s:%s:%s", mounts[i].Name, mounts[i].MountPath, mounts[i].SubPath)
		}
		ss = append(ss, mnt)
	}
	return strings.Join(ss, "\n"), nil
}

func fmtListFn(data string) (string, error) {
	// Return empty list if the data is empty. This helps to avoid object parsing error.
	// ref: https://stackoverflow.com/a/18419503
	if len(data) == 0 {
		return "[]", nil
	}
	return data, nil
}

type promNamespaceSelector struct {
	metav1.ObjectMeta `json:"metadata"`
	Spec              promNamespaceSelectorSpec `json:"spec"`
}
type promNamespaceSelectorSpec struct {
	NamespaceSelector *prom_op.NamespaceSelector `json:"namespaceSelector,omitempty"`
}

func promNamespaceSelectorFn(data string) (string, error) {
	if strings.TrimSpace(data) == "" {
		return "", nil
	}
	var selOpts promNamespaceSelector
	err := json.Unmarshal([]byte(data), &selOpts)
	if err != nil {
		return "", err
	}
	// If selector field is empty, then all namespaces are the targets.
	if selOpts.Spec.NamespaceSelector == nil {
		return "All", nil
	} else if len(selOpts.Spec.NamespaceSelector.MatchNames) != 0 {
		// If an array of namespace is provided, then those namespaces are the target
		return strings.Join(selOpts.Spec.NamespaceSelector.MatchNames, ", "), nil
	} else if !selOpts.Spec.NamespaceSelector.Any {
		// If "any: false" is set in the namespace selector field, only the object namespace is the target.
		return selOpts.Namespace, nil
	}
	return "", nil
}

func containerImagesFn(data string) (string, error) {
	var imagesBuffer bytes.Buffer
	if strings.TrimSpace(data) == "" {
		return "", nil
	}

	var containers []core.Container
	err := json.Unmarshal([]byte(data), &containers)
	if err != nil {
		return "", err
	}
	for i, container := range containers {
		imagesBuffer.WriteString(container.Image)
		if i != len(containers)-1 {
			imagesBuffer.WriteString(",")
		}
	}
	return imagesBuffer.String(), nil
}

func durationFn(data string) (string, error) {
	if strings.TrimSpace(data) == "" {
		return "", nil
	}

	// make the data a valid json
	ss := strings.Split(strings.TrimSuffix(data, ","), ",")
	jsonData := "["
	for i := range ss {
		jsonData += fmt.Sprintf("%q", ss[i])
		if i < len(ss)-1 {
			jsonData += ","
		}
	}
	jsonData += "]"

	var tt []metav1.Time
	err := json.Unmarshal([]byte(jsonData), &tt)
	if err != nil {
		return "", err
	}
	if len(tt) == 1 {
		// only start time is does exist
		return duration.HumanDuration(time.Since(tt[0].Time)), nil
	} else {
		return duration.HumanDuration(tt[1].Sub(tt[0].Time)), nil
	}
}

func mapKeyCountFn(data string) (string, error) {
	if strings.TrimSpace(data) == "" {
		return "", nil
	}

	var m map[string]string
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(len(m)), nil
}

func kubedbDBModeFn(data string) (string, error) {
	if strings.TrimSpace(data) == "" {
		return "", nil
	}
	var obj unstructured.Unstructured
	err := json.Unmarshal([]byte(data), &obj)
	if err != nil {
		return "", err
	}

	switch obj.GetKind() {
	case ResourceKindMongoDB:
		shards, found, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), "spec", "shardTopology")
		if err != nil {
			return "", err
		}
		if found && shards != nil {
			return DBModeSharded, nil
		}
		rs, found, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), "spec", "replicaSet")
		if err != nil {
			return "", err
		}
		if found && rs != nil {
			return DBModeReplicaSet, nil
		}
		return DBModeStandalone, nil
	case ResourceKindPostgres:
		mode, found, err := unstructured.NestedString(obj.UnstructuredContent(), "spec", "standbyMode")
		if err != nil {
			return "", err
		}
		if found && mode != "" {
			return mode, nil
		}
		return "Hot", nil
	case ResourceKindElasticsearch:
		topology, found, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), "spec", "topology")
		if err != nil {
			return "", err
		}
		if found && topology != nil {
			return "Topology", nil
		}
		return "Combined", nil
	case ResourceKindMariaDB:
		replicas, found, err := unstructured.NestedInt64(obj.UnstructuredContent(), "spec", "replicas")
		if err != nil {
			return "", err
		}
		if found && replicas > 1 {
			return DBModeCluster, nil
		}
		return DBModeStandalone, nil
	case ResourceKindMySQL:
		topology, found, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), "spec", "topology")
		if err != nil {
			return "", err
		}
		if found && topology != nil {
			mode, found, err := unstructured.NestedFieldCopy(topology.(map[string]interface{}), "mode")
			if err != nil {
				return "", err
			}
			if found && mode != nil {
				return fmt.Sprintf("%v", mode), nil
			}
		}
		return DBModeStandalone, nil
	case ResourceKindRedis:
		mode, found, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), "spec", "mode")
		if err != nil {
			return "", err
		}
		if found && mode != nil {
			return fmt.Sprintf("%v", mode), nil
		}
		return DBModeStandalone, nil
	}
	return "", fmt.Errorf("failed to detectect database mode. Reason: Unknown database type `%s`", obj.GetKind())
}

func kubedbDBReplicasFn(data string) (string, error) {
	var obj unstructured.Unstructured
	err := json.Unmarshal([]byte(data), &obj)
	if err != nil {
		return "", err
	}

	switch obj.GetKind() {
	case ResourceKindMongoDB:
		return getMongoDBReplicas(obj)
	case ResourceKindPostgres:
		replicas, _, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), "spec", "replicas")
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%v", replicas), nil
	case ResourceKindElasticsearch:
		return getElasticsearchReplicas(obj)
	case ResourceKindMariaDB:
		replicas, _, err := unstructured.NestedFieldCopy(obj.UnstructuredContent(), "spec", "replicas")
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%v", replicas), nil
	case ResourceKindMySQL:
		return getMySQLReplicas(obj)
	case ResourceKindRedis:
		return getRedisReplicas(obj)
	}
	return "", fmt.Errorf("failed to detect replica number. Reason: Unknown database type `%s`", obj.GetKind())
}

func kubedbDBResourcesFn(data string) (string, error) {
	var obj unstructured.Unstructured
	err := json.Unmarshal([]byte(data), &obj)
	if err != nil {
		return "", err
	}

	switch obj.GetKind() {
	case ResourceKindMongoDB:
		return mongoDBResources(obj)
	case ResourceKindPostgres:
		return postgresResources(obj)
	case ResourceKindElasticsearch:
		return elasticsearchDBResources(obj)
	case ResourceKindMariaDB:
		return mariaDBResources(obj)
	case ResourceKindMySQL:
		return mySQLResources(obj)
	case ResourceKindRedis:
		return redisResources(obj)
	}
	return "", fmt.Errorf("failed to extract CPU information. Reason: Unknown database type `%s`", obj.GetKind())
}

func rbacSubjects(data string) (string, error) {
	if strings.TrimSpace(data) == "" {
		return "", nil
	}

	var subjects []rbac.Subject
	err := json.Unmarshal([]byte(data), &subjects)
	if err != nil {
		return "", err
	}
	var ss []string
	for i := range subjects {
		s := fmt.Sprintf("%s %s", subjects[i].Kind, subjects[i].Name)
		if subjects[i].Namespace != "" {
			s = fmt.Sprintf("%s %s/%s", subjects[i].Kind, subjects[i].Namespace, subjects[i].Name)
		}
		ss = append(ss, s)
	}
	return strings.Join(ss, ","), nil
}

func certificateValidity(data string) (string, error) {
	if strings.TrimSpace(data) == "" {
		return "", nil
	}

	certStatus := struct {
		NotBefore metav1.Time `json:"notBefore"`
		NotAfter  metav1.Time `json:"notAfter"`
	}{}
	err := json.Unmarshal([]byte(data), &certStatus)
	if err != nil {
		return "", err
	}

	if certStatus.NotBefore.After(time.Now()) {
		return "Not valid yet", nil
	} else if time.Now().After(certStatus.NotAfter.Time) {
		return "Expired", nil
	}
	return duration.HumanDuration(time.Until(certStatus.NotAfter.Time)), nil
}
