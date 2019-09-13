# https://github.com/kubernetes/kubernetes/blob/master/pkg/printers/internalversion/printers.go

# CSV parser: https://github.com/mholt/PapaParse

$ kubectl create -f busy-dep.yaml
$ kubectl get deploy busy-dep -o=jsonpath='{.status.readyReplicas}/{.spec.replicas}'
$ kubectl get deploy busy-dep -o=jsonpath='{range .spec.template.spec.containers[*]}{.image}{"\n"}{end}'

$ kubectl get deploy busy-dep -o=jsonpath='{.spec.template.spec.containers[*].image}'



$ kubectl create -f nginx-rc.yaml
replicationcontroller/nginx created

kubectl get rc nginx -o=jsonpath='{.status.readyReplicas}/{.spec.replicas}'
kubectl get rc nginx -o=jsonpath='{range .spec.template.spec.containers[*]}{.image}{"\n"}{end}'


https://console.byte.builders/kubernetes/cluster-admin@gke_ackube_us-central1-f_demo.pharmer/replicationcontroller/nginx?namespace=default

