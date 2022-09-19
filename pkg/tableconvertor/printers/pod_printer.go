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

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	Register(PodPrinter{})
}

// ref: https://github.com/kubernetes/kubernetes/blob/v1.21.0/pkg/printers/internalversion/printers.go#L89-L101

type PodPrinter struct{}

var _ ColumnConverter = PodPrinter{}

func (_ PodPrinter) GVK() schema.GroupVersionKind {
	return core.SchemeGroupVersion.WithKind("Pod")
}

/*
"name": "Name",
"name": "Ready",
"name": "Status",
"name": "Restarts",
"name": "Age",
"name": "IP",
"name": "Node",
"name": "Nominated Node",
"name": "Readiness Gates",
*/
func (p PodPrinter) Convert(o runtime.Object) (map[string]interface{}, error) {
	pod := new(core.Pod)
	switch to := o.(type) {
	case *unstructured.Unstructured:
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(to.UnstructuredContent(), pod); err != nil {
			return nil, err
		}
	case *core.Pod:
		pod = to
	default:
		return nil, fmt.Errorf("expected %v, received %v", p.GVK().Kind, reflect.TypeOf(o))
	}

	restarts := 0
	totalContainers := len(pod.Spec.Containers)
	readyContainers := 0

	reason := string(pod.Status.Phase)
	if pod.Status.Reason != "" {
		reason = pod.Status.Reason
	}

	row := map[string]interface{}{}

	initializing := false
	for i := range pod.Status.InitContainerStatuses {
		container := pod.Status.InitContainerStatuses[i]
		restarts += int(container.RestartCount)
		switch {
		case container.State.Terminated != nil && container.State.Terminated.ExitCode == 0:
			continue
		case container.State.Terminated != nil:
			// initialization is failed
			if len(container.State.Terminated.Reason) == 0 {
				if container.State.Terminated.Signal != 0 {
					reason = fmt.Sprintf("Init:Signal:%d", container.State.Terminated.Signal)
				} else {
					reason = fmt.Sprintf("Init:ExitCode:%d", container.State.Terminated.ExitCode)
				}
			} else {
				reason = "Init:" + container.State.Terminated.Reason
			}
			initializing = true
		case container.State.Waiting != nil && len(container.State.Waiting.Reason) > 0 && container.State.Waiting.Reason != "PodInitializing":
			reason = "Init:" + container.State.Waiting.Reason
			initializing = true
		default:
			reason = fmt.Sprintf("Init:%d/%d", i, len(pod.Spec.InitContainers))
			initializing = true
		}
		break
	}
	if !initializing {
		restarts = 0
		hasRunning := false
		for i := len(pod.Status.ContainerStatuses) - 1; i >= 0; i-- {
			container := pod.Status.ContainerStatuses[i]

			restarts += int(container.RestartCount)
			if container.State.Waiting != nil && container.State.Waiting.Reason != "" {
				reason = container.State.Waiting.Reason
			} else if container.State.Terminated != nil && container.State.Terminated.Reason != "" {
				reason = container.State.Terminated.Reason
			} else if container.State.Terminated != nil && container.State.Terminated.Reason == "" {
				if container.State.Terminated.Signal != 0 {
					reason = fmt.Sprintf("Signal:%d", container.State.Terminated.Signal)
				} else {
					reason = fmt.Sprintf("ExitCode:%d", container.State.Terminated.ExitCode)
				}
			} else if container.Ready && container.State.Running != nil {
				hasRunning = true
				readyContainers++
			}
		}

		// change pod status back to "Running" if there is at least one container still reporting as "Running" status
		if reason == "Completed" && hasRunning {
			if hasPodReadyCondition(pod.Status.Conditions) {
				reason = "Running"
			} else {
				reason = "NotReady"
			}
		}
	}

	if pod.DeletionTimestamp != nil && pod.Status.Reason == NodeUnreachablePodReason {
		reason = "Unknown"
	} else if pod.DeletionTimestamp != nil {
		reason = "Terminating"
	}

	row["_Name"] = pod.Name
	row["Ready"] = fmt.Sprintf("%d/%d", readyContainers, totalContainers)
	row["Status"] = reason
	row["Restarts"] = int64(restarts)
	row["_Age"] = translateTimestampSince(pod.CreationTimestamp)

	nodeName := pod.Spec.NodeName
	nominatedNodeName := pod.Status.NominatedNodeName
	podIP := ""
	if len(pod.Status.PodIPs) > 0 {
		podIP = pod.Status.PodIPs[0].IP
	}

	if podIP == "" {
		podIP = None
	}
	if nodeName == "" {
		nodeName = None
	}
	if nominatedNodeName == "" {
		nominatedNodeName = None
	}

	readinessGates := None
	if len(pod.Spec.ReadinessGates) > 0 {
		trueConditions := 0
		for _, readinessGate := range pod.Spec.ReadinessGates {
			conditionType := readinessGate.ConditionType
			for _, condition := range pod.Status.Conditions {
				if condition.Type == conditionType {
					if condition.Status == core.ConditionTrue {
						trueConditions++
					}
					break
				}
			}
		}
		readinessGates = fmt.Sprintf("%d/%d", trueConditions, len(pod.Spec.ReadinessGates))
	}

	row["IP"] = podIP
	row["Node"] = nodeName
	row["Nominated Node"] = nominatedNodeName
	row["Readiness Gates"] = readinessGates

	return row, nil
}

func hasPodReadyCondition(conditions []core.PodCondition) bool {
	for _, condition := range conditions {
		if condition.Type == core.PodReady && condition.Status == core.ConditionTrue {
			return true
		}
	}
	return false
}
