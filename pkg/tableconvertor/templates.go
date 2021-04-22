package tableconvertor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"kmodules.xyz/resource-metadata/pkg/tableconvertor/printers"

	"github.com/Masterminds/sprig/v3"
	"gomodules.xyz/jsonpath"
	core "k8s.io/api/core/v1"
	metatable "k8s.io/apimachinery/pkg/api/meta/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

var templateFns = sprig.TxtFuncMap()

func init() {
	templateFns["jp"] = jsonpathFn
	templateFns["k8s_fmt_selector"] = formatLabelSelectorFn
	templateFns["k8s_fmt_label"] = formatLabelsFn
	templateFns["k8s_age"] = ageFn
	templateFns["k8s_svc_ports"] = servicePortsFn
	templateFns["k8s_container_ports"] = containerPortFn
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
	var sel metav1.LabelSelector
	err := json.Unmarshal([]byte(data), &sel)
	if err != nil {
		return "", err
	}
	return metav1.FormatLabelSelector(&sel), nil
}

func formatLabelsFn(data string) (string, error) {
	var sel map[string]string
	err := json.Unmarshal([]byte(data), &sel)
	if err != nil {
		return "", err
	}
	return labels.FormatLabels(sel), nil
}

func ageFn(data string) (string, error) {
	var timestamp metav1.Time
	err := timestamp.UnmarshalQueryParameter(data)
	if err != nil {
		return "", err
	}
	return metatable.ConvertToHumanReadableDateType(timestamp), nil
}

func servicePortsFn(data string) (string, error) {
	var ports []core.ServicePort
	err := json.Unmarshal([]byte(data), &ports)
	if err != nil {
		return "", err
	}
	return printers.MakeServicePortString(ports), nil
}

func containerPortFn(data string) (string, error) {
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
