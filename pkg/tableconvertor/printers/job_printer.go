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
	"time"

	batch "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/duration"
)

func init() {
	Register(JobPrinter{})
}

// ref: https://github.com/kubernetes/kubernetes/blob/v1.21.0/pkg/printers/internalversion/printers.go#L164-L174

type JobPrinter struct{}

var _ ColumnConverter = JobPrinter{}

func (_ JobPrinter) GVK() schema.GroupVersionKind {
	return batch.SchemeGroupVersion.WithKind("Job")
}

func (p JobPrinter) Convert(o runtime.Object) (map[string]interface{}, error) {
	obj, ok := o.(*batch.Job)
	if !ok {
		return nil, fmt.Errorf("expected %v, received %v", p.GVK().Kind, reflect.TypeOf(o))
	}

	row := map[string]interface{}{}

	var completions string
	if obj.Spec.Completions != nil {
		completions = fmt.Sprintf("%d/%d", obj.Status.Succeeded, *obj.Spec.Completions)
	} else {
		parallelism := int32(0)
		if obj.Spec.Parallelism != nil {
			parallelism = *obj.Spec.Parallelism
		}
		if parallelism > 1 {
			completions = fmt.Sprintf("%d/1 of %d", obj.Status.Succeeded, parallelism)
		} else {
			completions = fmt.Sprintf("%d/1", obj.Status.Succeeded)
		}
	}
	var jobDuration string
	switch {
	case obj.Status.StartTime == nil:
	case obj.Status.CompletionTime == nil:
		jobDuration = duration.HumanDuration(time.Since(obj.Status.StartTime.Time))
	default:
		jobDuration = duration.HumanDuration(obj.Status.CompletionTime.Sub(obj.Status.StartTime.Time))
	}

	row["Name"] = obj.Name
	row["Completions"] = completions
	row["Duration"] = jobDuration
	row["Age"] = translateTimestampSince(obj.CreationTimestamp)

	names, images := layoutContainerCells(obj.Spec.Template.Spec.Containers)
	row["Containers"] = names
	row["Images"] = images
	row["Selector"] = metav1.FormatLabelSelector(obj.Spec.Selector)

	return row, nil
}
