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

package graph

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"

	"github.com/mitchellh/mapstructure"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/util/jsonpath"
)

func (g *Graph) List(dc dynamic.Interface, src unstructured.Unstructured, dstGVR schema.GroupVersionResource) ([]unstructured.Unstructured, error) {
	srcGVR, err := g.r.GVR(src.GroupVersionKind())
	if err != nil {
		return nil, err
	}
	dist, prev := Dijkstra(g, srcGVR)

	paths := GeneratePaths(srcGVR, dist, prev)
	path, ok := paths[dstGVR]
	if !ok {
		return nil, nil
	}

	in := []unstructured.Unstructured{src}
	var out []unstructured.Unstructured
	for _, e := range path.Edges {
		out = nil
		for _, inObj := range in {
			result, err := g.ResourcesFor(dc, inObj, *e)
			if err != nil {
				return nil, err
			}
			out = appendObjects(out, result...)
		}
		in = out
	}

	return out, nil
}

type objectKey struct {
	name      string
	namespace string
}

func appendObjects(arr []unstructured.Unstructured, items ...unstructured.Unstructured) []unstructured.Unstructured {
	m := make(map[objectKey]unstructured.Unstructured)

	for i := range arr {
		m[objectKey{namespace: arr[i].GetNamespace(), name: arr[i].GetName()}] = arr[i]
	}
	for i := range items {
		m[objectKey{namespace: items[i].GetNamespace(), name: items[i].GetName()}] = items[i]
	}

	out := make([]unstructured.Unstructured, 0, len(m))
	for _, obj := range m {
		out = append(out, obj)
	}
	return out
}

func (g *Graph) ResourcesFor(dc dynamic.Interface, src unstructured.Unstructured, e Edge) ([]unstructured.Unstructured, error) {
	gvr, err := g.r.GVR(src.GroupVersionKind())
	if err != nil {
		return nil, err
	}
	if e.Src != gvr {
		return nil, fmt.Errorf("edge src %v does not match ref %v", e.Src, gvr)
	}

	if e.Forward {
		// FIXME: How to handle namespace for Backward direction
		if e.Connection.Type == v1alpha1.MatchSelector {
			var ls string
			var selector labels.Selector
			var err error

			if e.Connection.SelectorPath != "" {
				ls, selector, err = ExtractSelector(src, e.Connection.SelectorPath)
				if err != nil {
					return nil, err
				}
			} else if e.Connection.Selector != nil {
				s2, err := evalLabelSelector(src, e.Connection.Selector)
				if err != nil {
					return nil, err
				}
				selector, err = metav1.LabelSelectorAsSelector(s2)
				if err != nil {
					return nil, err
				}
				ls = selector.String()
			} else {
				return nil, fmt.Errorf("edge %v is missing selectorPath and selector", e)
			}

			namespaces, err := Namespaces(src, e.Connection.NamespacePath)
			if err != nil {
				return nil, err
			}
			if len(namespaces) == 0 {
				namespaces = []string{metav1.NamespaceAll}
			}

			var out []unstructured.Unstructured
			for _, ns := range namespaces {
				var ri dynamic.ResourceInterface
				ri = dc.Resource(e.Dst)
				if namespaced, err := g.r.IsNamespaced(e.Dst); err != nil {
					return nil, err
				} else if namespaced {
					ri = dc.Resource(e.Dst).Namespace(ns)
				}

				selInApp := e.Connection.TargetLabelPath != "" && strings.Trim(e.Connection.TargetLabelPath, ".") != MetadataLabels

				var opts metav1.ListOptions
				if !selInApp {
					opts.LabelSelector = ls
				}
				result, err := ri.List(opts)
				if err != nil {
					return nil, err
				}
				for i := range result.Items {
					rs := result.Items[i]

					if selInApp {
						lbl, ok, err := unstructured.NestedStringMap(rs.Object, fields(e.Connection.TargetLabelPath)...)
						if err != nil {
							return nil, err
						}
						if !ok || !selector.Matches(labels.Set(lbl)) {
							continue
						}
					}

					if isConnected(e.Connection.Level, rs, src) {
						out = append(out, rs)
					}
				}
			}
			return out, nil
		} else if e.Connection.Type == v1alpha1.MatchName {
			if e.Connection.NameTemplate == "" {
				return nil, fmt.Errorf("edge %v is missing nameTemplate", e)
			}
			name := strings.ReplaceAll(e.Connection.NameTemplate, MetadataNameQuery, src.GetName())

			namespaces, err := Namespaces(src, e.Connection.NamespacePath)
			if err != nil {
				return nil, err
			}
			if len(namespaces) == 0 {
				namespaces = []string{metav1.NamespaceAll}
			}

			var out []unstructured.Unstructured
			for _, ns := range namespaces {
				var ri dynamic.ResourceInterface
				ri = dc.Resource(e.Dst)
				if namespaced, err := g.r.IsNamespaced(e.Dst); err != nil {
					return nil, err
				} else if namespaced {
					ri = dc.Resource(e.Dst).Namespace(ns)
				}
				rs, err := ri.Get(name, metav1.GetOptions{})
				if err != nil {
					return nil, err
				}

				if isConnected(e.Connection.Level, *rs, src) {
					out = append(out, *rs)
				}
			}
			return out, nil
		} else if e.Connection.Type == v1alpha1.OwnedBy {
			return g.findOwners(dc, e, src.GetOwnerReferences(), src.GetNamespace())
		} else if e.Connection.Type == v1alpha1.MatchRef {
			// TODO: check that namespacePath must be empty

			var out []unstructured.Unstructured

			for _, reference := range e.Connection.References {
				j := jsonpath.New("jsonpath")
				j.AllowMissingKeys(true)
				err := j.Parse(reference)
				if err != nil {
					return nil, fmt.Errorf("fails to parse reference %q between %s -> %s. err:%v", e.Connection.References, e.Src, e.Dst, err)
				}
				buf := new(bytes.Buffer)
				err = j.Execute(buf, src.Object)
				if err != nil {
					return nil, fmt.Errorf("fails to execute reference %q between %s -> %s. err:%v", e.Connection.References, e.Src, e.Dst, err)
				}
				r := csv.NewReader(buf)
				r.Comma = ';'
				r.Comment = '#'
				records, err := r.ReadAll()
				if err != nil {
					return nil, err
				}
				refs, err := ParseResourceRefs(records)
				if err != nil {
					return nil, err
				}

				var objects []unstructured.Unstructured
				for _, ref := range refs {
					// if apiGroup is set, it must match
					if ref.APIGroup != "" && !equalsGV(ref.APIGroup, e.Dst) {
						continue
					}
					// if apiGroup is set, it must match
					gvk, err := g.r.GVK(e.Dst)
					if err != nil {
						return nil, err
					}
					if ref.Kind != "" && ref.Kind != gvk.Kind {
						continue
					}

					var ri dynamic.ResourceInterface
					ri = dc.Resource(e.Dst)
					if namespaced, err := g.r.IsNamespaced(e.Dst); err != nil {
						return nil, err
					} else if namespaced {
						ns := ref.Namespace
						if ns == "" {
							ns = src.GetNamespace()
						}
						if ns == "" {
							// dst is namespaced &&
							// no namespace is defined in reference &&
							// src is not-namespaced
							return nil, errors.New("namespace must be defined in reference")
						}
						ri = dc.Resource(e.Dst).Namespace(ns)
					}
					rs, err := ri.Get(ref.Name, metav1.GetOptions{})
					if err != nil {
						return nil, err
					}

					if isConnected(e.Connection.Level, *rs, src) {
						objects = append(objects, *rs)
					}
				}
				out = appendObjects(out, objects...)
			}
			return out, nil
		}
	} else {
		namespace := core.NamespaceAll
		if e.Connection.NamespacePath == MetadataNamespace {
			namespace = src.GetNamespace()
		} // else all namespace RETHINK

		if e.Connection.Type == v1alpha1.MatchSelector {
			var out []unstructured.Unstructured

			lbl := src.GetLabels()
			if e.Connection.TargetLabelPath != "" && strings.Trim(e.Connection.TargetLabelPath, ".") != MetadataLabels {
				l2, ok, err := unstructured.NestedStringMap(src.Object, fields(e.Connection.TargetLabelPath)...)
				if err != nil {
					return nil, err
				}
				if !ok {
					return out, nil // empty result
				}
				lbl = l2
			}

			var ri dynamic.ResourceInterface
			ri = dc.Resource(e.Dst)
			if namespaced, err := g.r.IsNamespaced(e.Dst); err != nil {
				return nil, err
			} else if namespaced {
				ri = dc.Resource(e.Dst).Namespace(namespace)
			}
			result, err := ri.List(metav1.ListOptions{})
			if err != nil {
				return nil, err
			}
			for i := range result.Items {
				rs := result.Items[i]

				if e.Connection.NamespacePath != "" && e.Connection.NamespacePath != MetadataNamespace {
					namespaces, err := Namespaces(rs, e.Connection.NamespacePath)
					if err != nil {
						return nil, err
					}
					if len(namespaces) > 0 && !contains(namespaces, src.GetNamespace()) {
						continue
					}
				}

				var ls string
				var selector labels.Selector
				if e.Connection.SelectorPath != "" {
					ls, selector, err = ExtractSelector(rs, e.Connection.SelectorPath)
					if err != nil {
						return nil, err
					}
				} else if e.Connection.Selector != nil {
					s2, err := evalLabelSelector(rs, e.Connection.Selector)
					if err != nil {
						return nil, err
					}
					selector, err = metav1.LabelSelectorAsSelector(s2)
					if err != nil {
						return nil, err
					}
					ls = selector.String()
				} else {
					return nil, fmt.Errorf("edge %v is missing selectorPath and selector", e)
				}

				if ls == labels.Nothing().String() {
					continue
				}

				if selector.Matches(labels.Set(lbl)) {
					if isConnected(e.Connection.Level, src, rs) {
						out = append(out, rs)
					}
				}
			}
			return out, nil
		} else if e.Connection.Type == v1alpha1.MatchName {
			if e.Connection.NameTemplate != "" {
				name, ok := ExtractName(src.GetName(), e.Connection.NameTemplate)
				if !ok {
					return nil, fmt.Errorf("failed to detect name from %s and %s", src.GetName(), e.Connection.NameTemplate)
				}

				var out []unstructured.Unstructured
				var ri dynamic.ResourceInterface
				ri = dc.Resource(e.Dst)
				if namespaced, err := g.r.IsNamespaced(e.Dst); err != nil {
					return nil, err
				} else if namespaced {
					ri = dc.Resource(e.Dst).Namespace(namespace)
				}
				rs, err := ri.Get(name, metav1.GetOptions{})
				if err != nil {
					return nil, err
				}

				if isConnected(e.Connection.Level, src, *rs) {
					out = append(out, *rs)
				}

				return out, nil
			}
		} else if e.Connection.Type == v1alpha1.OwnedBy {
			return g.findChildren(dc, e, src)
		} else if e.Connection.Type == v1alpha1.MatchRef {
			// TODO: check that namespacePath must be empty

			var ri dynamic.ResourceInterface
			ri = dc.Resource(e.Dst)
			if namespaced, err := g.r.IsNamespaced(e.Dst); err != nil {
				return nil, err
			} else if namespaced {
				ns := metav1.NamespaceAll
				if e.Connection.NamespacePath == MetadataNamespace {
					ns = src.GetNamespace()
				}
				ri = dc.Resource(e.Dst).Namespace(ns)
			}
			result, err := ri.List(metav1.ListOptions{})
			if err != nil {
				return nil, err
			}

			var out []unstructured.Unstructured
		NextItem:
			for i := range result.Items {
				rs := result.Items[i]

				for _, reference := range e.Connection.References {

					j := jsonpath.New("jsonpath")
					j.AllowMissingKeys(true)
					err := j.Parse(reference)
					if err != nil {
						return nil, fmt.Errorf("fails to parse reference %q between %s -> %s. err:%v", e.Connection.References, e.Src, e.Dst, err)
					}

					buf := new(bytes.Buffer)
					err = j.Execute(buf, rs.Object)
					if err != nil {
						return nil, fmt.Errorf("fails to execute reference %q between %s -> %s. err:%v", e.Connection.References, e.Src, e.Dst, err)
					}
					r := csv.NewReader(buf)
					r.Comma = ';'
					r.Comment = '#'
					records, err := r.ReadAll()
					if err != nil {
						return nil, err
					}
					refs, err := ParseResourceRefs(records)
					if err != nil {
						return nil, err
					}
					for _, ref := range refs {
						// if apiGroup is set, it must match
						if ref.APIGroup != "" && !equalsGV(ref.APIGroup, e.Src) {
							continue
						}

						// if apiGroup is set, it must match
						gvk, err := g.r.GVK(e.Src)
						if err != nil {
							return nil, err
						}
						if ref.Kind != "" && ref.Kind != gvk.Kind {
							continue
						}

						ns := ref.Namespace
						namespaced, err := g.r.IsNamespaced(e.Src)
						if err != nil {
							return nil, err
						}
						if ns == "" && namespaced {
							ns = rs.GetNamespace()
							if ns == "" {
								// src is namespaced &&
								// no namespace is defined in reference &&
								// rs is not-namespaced
								return nil, errors.New("namespace must be defined in reference")
							}

							if ns != src.GetNamespace() {
								continue
							}
						}

						if ref.Name != src.GetName() {
							continue
						}

						if isConnected(e.Connection.Level, src, rs) {
							out = append(out, rs)
						}
						continue NextItem
					}
				}
			}
			return out, nil
		}
	}

	return nil, nil
}

func isConnected(conn v1alpha1.OwnershipLevel, obj unstructured.Unstructured, owner unstructured.Unstructured) bool {
	switch conn {
	case v1alpha1.Controller:
		if metav1.IsControlledBy(&obj, &owner) {
			return true
		}
	case v1alpha1.Owner:
		if IsOwnedBy(&obj, &owner) {
			return true
		}
	default:
		return true
	}
	return false
}

func evalLabelSelector(obj unstructured.Unstructured, in *metav1.LabelSelector) (*metav1.LabelSelector, error) {
	out := in.DeepCopy()
	for k, v := range out.MatchLabels {
		if strings.ContainsRune(k, '{') {
			return nil, fmt.Errorf("invalid selector key %v", k)
		}
		if v == MetadataNameQuery {
			out.MatchLabels[k] = obj.GetName()
			continue
		}
		if v == MetadataNamespaceQuery {
			out.MatchLabels[k] = obj.GetNamespace()
			continue
		}
		if v[0] == '{' && v[len(v)-1] == '}' {
			val, err := evalJsonPath(obj, v)
			if err != nil {
				return nil, err
			}
			out.MatchLabels[k] = val
		}
	}
	for i := range out.MatchExpressions {
		expr := out.MatchExpressions[i]
		if strings.ContainsRune(expr.Key, '{') {
			return nil, fmt.Errorf("selector has invalid key %v", expr.Key)
		}
		for vi := range expr.Values {
			v := expr.Values[vi]
			if v[0] == '{' && v[len(v)-1] == '}' {
				val, err := evalJsonPath(obj, v)
				if err != nil {
					return nil, err
				}
				expr.Values[vi] = val
			}
		}
		out.MatchExpressions[i] = expr
	}
	return out, nil
}

func evalJsonPath(src unstructured.Unstructured, template string) (string, error) {
	j := jsonpath.New("jsonpath")
	j.AllowMissingKeys(true)
	err := j.Parse(template)
	if err != nil {
		return "", fmt.Errorf("fails to parse value of selector key. err:%v", err)
	}
	buf := new(bytes.Buffer)
	err = j.Execute(buf, src.Object)
	if err != nil {
		return "", fmt.Errorf("fails to evaluate value of selector key. err:%v", err)
	}
	return strings.TrimSpace(buf.String()), nil
}

func (g *Graph) findOwners(dc dynamic.Interface, e Edge, srcOwnerRefs []metav1.OwnerReference, namespace string) ([]unstructured.Unstructured, error) {
	var out []unstructured.Unstructured

	var ri dynamic.ResourceInterface
	ri = dc.Resource(e.Dst)
	if namespaced, err := g.r.IsNamespaced(e.Dst); err != nil {
		return nil, err
	} else if namespaced {
		ri = dc.Resource(e.Dst).Namespace(namespace)
	}
	t, err := g.r.TypeMeta(e.Dst)
	if err != nil {
		return nil, err
	}
	for _, ref := range srcOwnerRefs {
		if ref.APIVersion == t.APIVersion && ref.Kind == t.Kind {
			if e.Connection.Level == v1alpha1.Controller {
				if ref.Controller != nil && *ref.Controller {
					rs, err := ri.Get(ref.Name, metav1.GetOptions{})
					if err != nil {
						return nil, err
					}
					out = append(out, *rs)
					break
				}
			} else if e.Connection.Level == v1alpha1.Owner {
				rs, err := ri.Get(ref.Name, metav1.GetOptions{})
				if err != nil {
					return nil, err
				}
				out = append(out, *rs)
			} else {
				return nil, fmt.Errorf("connection level should be Owner or Controller, found %v", e.Connection.Level)
			}
		}
	}

	return out, nil
}

func (g *Graph) findChildren(dc dynamic.Interface, e Edge, src unstructured.Unstructured) ([]unstructured.Unstructured, error) {
	if e.Connection.Level != v1alpha1.Owner && e.Connection.Level != v1alpha1.Controller {
		return nil, fmt.Errorf("connection level should be Owner or Controller, found %v", e.Connection.Level)
	}

	var out []unstructured.Unstructured

	var ri dynamic.ResourceInterface
	ri = dc.Resource(e.Dst)
	if namespaced, err := g.r.IsNamespaced(e.Dst); err != nil {
		return nil, err
	} else if namespaced {
		ri = dc.Resource(e.Dst).Namespace(src.GetNamespace())
	}

	result, err := ri.List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for i := range result.Items {
		rs := result.Items[i]
		if isConnected(e.Connection.Level, rs, src) {
			out = append(out, rs)
		}
	}

	return out, nil
}

func IsOwnedBy(obj metav1.Object, owner metav1.Object) bool {
	for _, ref := range obj.GetOwnerReferences() {
		if ref.UID == owner.GetUID() {
			return true
		}
	}
	return false
}

// len([]string) == 0 && err == nil => all namespaces
func Namespaces(ref unstructured.Unstructured, nsSelector string) ([]string, error) {
	if nsSelector == MetadataNamespace {
		return []string{ref.GetNamespace()}, nil
	} else if nsSelector != "" {
		var nsel NamespaceSelector
		ok, err := Extract(ref, nsSelector, &nsel)
		if err != nil {
			return nil, err
		}
		if ok {
			// https://gitg.r.com/coreos/prometheus-operator/blob/cc584ecfa08d2eb95ba9401f116e3a20bf71be8b/pkg/prometheus/promcfg.go#L392
			if !nsel.Any && len(nsel.MatchNames) == 0 {
				return []string{ref.GetNamespace()}, nil
			} else if len(nsel.MatchNames) > 0 {
				return nsel.MatchNames, nil
			}
			return nil, nil
		}
	}
	return nil, nil
}

func Extract(u unstructured.Unstructured, fieldPath string, v interface{}) (bool, error) {
	if fieldPath == "" {
		return false, errors.New("fieldPath can't be empty")
	}
	f, ok, err := unstructured.NestedMap(u.Object, fields(fieldPath)...)
	if !ok || err != nil {
		return false, err
	}
	err = DecodeJSON(f, v)
	return err == nil, err
}

func keyExists(m map[string]interface{}, key string) bool {
	_, ok := m[key]
	return ok
}

func ExtractSelector(u unstructured.Unstructured, fieldPath string) (string, labels.Selector, error) {
	nothing := labels.Nothing().String()

	if fieldPath == "" {
		return nothing, labels.Nothing(), errors.New("fieldPath can't be empty")
	}
	val, found, err := unstructured.NestedFieldNoCopy(u.Object, fields(fieldPath)...)
	if !found || err != nil {
		return nothing, labels.Nothing(), err
	}
	m, ok := val.(map[string]interface{})
	if !ok {
		return nothing, labels.Nothing(), fmt.Errorf("%v accessor error: %v is of the type %T, expected map[string]interface{}", fieldPath, val, val)
	}

	if len(m) <= 2 && (keyExists(m, "matchLabels") || keyExists(m, "matchExpressions")) {
		var ls metav1.LabelSelector
		err = DecodeJSON(m, &ls)
		if err != nil {
			return nothing, labels.Nothing(), err
		}

		sel, err := metav1.LabelSelectorAsSelector(&ls)
		if err != nil {
			return nothing, labels.Nothing(), err
		}
		return sel.String(), sel, nil
	}

	strMap := make(map[string]string, len(m))
	for k, v := range m {
		if str, ok := v.(string); ok {
			strMap[k] = str
		} else {
			return nothing, labels.Nothing(), fmt.Errorf("%v accessor error: contains non-string key in the map: %v is of the type %T, expected string", fieldPath, v, v)
		}
	}
	sel := labels.SelectorFromSet(strMap)
	return sel.String(), sel, nil
}

func DecodeJSON(input map[string]interface{}, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   output,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func ExtractName(name, selector string) (string, bool) {
	re := regexp.MustCompile(`^` + strings.ReplaceAll(selector, MetadataNameQuery, `(.*)`) + `$`)
	matches := re.FindStringSubmatch(name)
	if len(matches) != 2 {
		return "", false
	}
	return matches[1], true
}

func ParseResourceRefs(records [][]string) ([]ResourceRef, error) {
	var refs []ResourceRef

	var cols int
	for i, rec := range records {
		n := len(rec)
		if i == 0 {
			cols = n
		} else if cols != n {
			return nil, errors.New("all rows must have same number of columns")
		}

		switch n {
		case 0:
			return nil, errors.New("must have at least one column")
		case 1:
			refs = append(refs, ResourceRef{
				Name: rec[0],
			})
		case 2:
			refs = append(refs, ResourceRef{
				Name:      rec[0],
				Namespace: rec[1],
			})
		case 3:
			refs = append(refs, ResourceRef{
				Name:      rec[0],
				Namespace: rec[1],
				Kind:      rec[2],
			})
		case 4:
			refs = append(refs, ResourceRef{
				Name:      rec[0],
				Namespace: rec[1],
				Kind:      rec[2],
				APIGroup:  rec[3],
			})
		default:
			return nil, fmt.Errorf("maximum 4 columns can be present, found %d", n)
		}
	}
	return refs, nil
}
