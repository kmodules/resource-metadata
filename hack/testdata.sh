# Copyright AppsCode Inc. and Contributors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# https://github.com/kubernetes/kubernetes/blob/master/pkg/printers/internalversion/printers.go

# CSV parser: https://github.com/mholt/PapaParse

$ kubectl create -f busy-dep.yaml
$ kubectl get deploy busy-dep -o=jsonpath='{.status.readyReplicas}/{.spec.replicas}'
$ kubectl get deploy busy-dep -o=jsonpath='{range .spec.template.spec.containers[*]}{.image}{"\n"}{end}'

$ kubectl get deploy busy-dep -o=jsonpath='{.spec.template.spec.containers[*].image}'

$ kubectl get deploy busy-dep -o=jsonpath='{.status.updatedReplicas} updated, {.status.replicas} total, {.status.availableReplicas} available, {.status.unavailableReplicas} unavailable'

$ kubectl create -f nginx-rc.yaml
replicationcontroller/nginx created

kubectl get rc nginx -o=jsonpath='{.status.readyReplicas}/{.spec.replicas}'
kubectl get rc nginx -o=jsonpath='{range .spec.template.spec.containers[*]}{.image}{"\n"}{end}'

https://console.byte.builders/kubernetes/cluster-admin@gke_ackube_us-central1-f_demo.pharmer/replicationcontroller/nginx?namespace=default
