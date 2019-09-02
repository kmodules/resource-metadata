package graph

import (
	"testing"
)

func TestExtractName(t *testing.T) {
	tests := []struct {
		name     string
		selector string
		result   string
	}{
		{
			name:     "kubedb-sample",
			selector: "kubedb-{.metadata.name}",
			result:   "sample",
		},
		{
			name:     "sample-kubedb",
			selector: "{.metadata.name}-kubedb",
			result:   "sample",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, ok := ExtractName(test.name, test.selector)
			if !ok {
				t.FailNow()
			}
			if test.result != r {
				t.FailNow()
			}
		})
	}
}
