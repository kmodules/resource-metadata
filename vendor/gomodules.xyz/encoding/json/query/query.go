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

package query

import (
	"errors"
	"fmt"
	"strings"

	"gomodules.xyz/encoding/json"

	"gomodules.xyz/jsonpath"
)

func QueryFieldNoCopy(d map[string]interface{}, expr string) (interface{}, bool, error) {
	if strings.HasPrefix(expr, "{") {
		return QueryJSONPathNoCopy(d, expr)
	}
	return NestedFieldNoCopy(d, fields(expr)...)
}

func QuerySlice(d map[string]interface{}, expr string) ([]interface{}, bool, error) {
	val, found, err := QueryFieldNoCopy(d, expr)
	if !found || err != nil {
		return nil, found, err
	}
	_, ok := val.([]interface{})
	if !ok {
		return nil, false, fmt.Errorf("%v accessor error: %v is of the type %T, expected []interface{}", expr, val, val)
	}
	return json.DeepCopyJSONValue(val).([]interface{}), true, nil
}

func fields(path string) []string {
	return strings.Split(strings.Trim(path, "."), ".")
}

func QueryJSONPathNoCopy(d interface{}, expr string) (interface{}, bool, error) {
	enableJSONOutput := false

	jp := jsonpath.New("")
	if err := jp.Parse(expr); err != nil {
		return nil, false, err
	}
	jp.AllowMissingKeys(true)
	jp.EnableJSONOutput(enableJSONOutput)

	fullResults, err := jp.FindResults(d)
	if err != nil {
		return nil, false, err
	}
	switch len(fullResults) {
	case 0:
		return nil, false, nil
	case 1:
		if len(fullResults[0]) > 1 {
			return nil, false, errors.New("expr returned multiple results")
		}
		return fullResults[0][0].Interface(), true, nil
	default:
		return nil, false, errors.New("expr returned multiple results")
	}
}

func QueryJSONPathCopy(obj map[string]interface{}, expr string) (interface{}, bool, error) {
	val, found, err := QueryJSONPathNoCopy(obj, expr)
	if !found || err != nil {
		return nil, found, err
	}
	return json.DeepCopyJSONValue(val), true, nil
}
