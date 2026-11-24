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

package clusterprofiles

import (
	"testing"
)

// SanitizeFeatureSets does not expand a feature's requirements, so the profile
// has to list every prerequisite of what Patch Console needs itself.
func TestPatchConsoleProfileListsPrerequisites(t *testing.T) {
	cp, err := LoadInternalByName("patch-console")
	if err != nil {
		t.Fatal(err)
	}

	want := map[string][]string{
		"opscenter-core":              {"license-proxyserver"},
		"opscenter-observability":     {"grafana-operator", "monitoring-operator"},
		"opscenter-security":          {"cert-manager", "patch-console-agent", "scanner", "trivy-operator"},
		"opscenter-policy-management": {"kyverno", "kyverno-policies"},
		"opscenter-backup":            {"kubestash"},
	}
	if len(cp.Spec.RequiredFeatureSets) != len(want) {
		t.Errorf("got feature sets %v, want %v", cp.Spec.RequiredFeatureSets, want)
	}
	for set, features := range want {
		got := map[string]bool{}
		for _, f := range cp.Spec.RequiredFeatureSets[set] {
			got[f] = true
		}
		for _, f := range features {
			if !got[f] {
				t.Errorf("feature set %s is missing %s", set, f)
			}
		}
		if len(got) != len(features) {
			t.Errorf("feature set %s has %v, want %v", set, cp.Spec.RequiredFeatureSets[set], features)
		}
	}
}
