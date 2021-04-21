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

	batch "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	Register(CronJobPrinter{})
}

// ref: https://github.com/kubernetes/kubernetes/blob/v1.21.0/pkg/printers/internalversion/printers.go#L176-L188

type CronJobPrinter struct{}

var _ ColumnConverter = CronJobPrinter{}

func (_ CronJobPrinter) GVK() schema.GroupVersionKind {
	return batch.SchemeGroupVersion.WithKind("CronJob")
}

func (p CronJobPrinter) Convert(o runtime.Object) (map[string]interface{}, error) {
	obj := new(batch.CronJob)
	switch to := o.(type) {
	case *unstructured.Unstructured:
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(to.UnstructuredContent(), obj); err != nil {
			return nil, err
		}
	case *batch.CronJob:
		obj = to
	default:
		return nil, fmt.Errorf("expected %v, received %v", p.GVK().Kind, reflect.TypeOf(o))
	}

	row := map[string]interface{}{}

	lastScheduleTime := None
	if obj.Status.LastScheduleTime != nil {
		lastScheduleTime = translateTimestampSince(*obj.Status.LastScheduleTime)
	}

	row["Name"] = obj.Name
	row["Schedule"] = obj.Spec.Schedule
	row["Suspend"] = printBoolPtr(obj.Spec.Suspend)
	row["Active"] = int64(len(obj.Status.Active))
	row["Last Schedule"] = lastScheduleTime
	row["Age"] = translateTimestampSince(obj.CreationTimestamp)

	names, images := layoutContainerCells(obj.Spec.JobTemplate.Spec.Template.Spec.Containers)
	row["Containers"] = names
	row["Images"] = images
	row["Selector"] = metav1.FormatLabelSelector(obj.Spec.JobTemplate.Spec.Selector)

	return row, nil
}
