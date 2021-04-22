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
	"fmt"
	"reflect"
	"testing"
)

func Test_cellForJSONValue(t *testing.T) {
	tests := []struct {
		headerType string
		value      string
		want       interface{}
		wantErr    bool
	}{
		{"integer", "42", int64(42), false},
		{"integer", "3.14", nil, true},
		{"integer", "true", nil, true},
		{"integer", "foo", nil, true},

		{"number", "42", float64(42), false},
		{"number", "3.14", float64(3.14), false},
		{"number", "true", nil, true},
		{"number", "foo", nil, true},

		{"boolean", "42", nil, true},
		{"boolean", "3.14", nil, true},
		{"boolean", "true", true, false},
		{"boolean", "foo", nil, true},

		{"string", "42", "42", false},
		{"string", "3.14", "3.14", false},
		{"string", "true", "true", false},
		{"string", "foo", "foo", false},

		{"object", `{"app": "xyz"}`, map[string]interface{}{"app": "xyz"}, false},

		{"unknown", "foo", nil, true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v of type %s", tt.value, tt.headerType), func(t *testing.T) {
			got, err := cellForJSONValue("", tt.headerType, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("cellForJSONValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cellForJSONValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}
