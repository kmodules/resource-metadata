# resource-dashboards

```
$ k get grafanadashboards -A -o custom-columns=T:.spec.model.title | sort

Alertmanager / Overview
CoreDNS
Grafana Overview
KubeDB / Elasticsearch / Database
KubeDB / Elasticsearch / Pod
KubeDB / Elasticsearch / Summary
KubeDB / MariaDB / Database
KubeDB / MariaDB / Pod
KubeDB / MariaDB / Summary
KubeDB / MongoDB / Database / ReplicaSet
KubeDB / MongoDB / Pod
KubeDB / MongoDB / Summary
KubeDB / MySQL / Database
KubeDB / MySQL / Group-Replication-Summary
KubeDB / MySQL / Pod
KubeDB / MySQL / Summary
KubeDB / Postgres / Database
KubeDB / Postgres / Pod
KubeDB / Postgres / Summary
KubeDB / Redis / Pod
KubeDB / Redis / Shard
KubeDB / Redis / Summary
Kubernetes / API server
Kubernetes / Compute Resources / Cluster
Kubernetes / Compute Resources / Namespace (Pods)
Kubernetes / Compute Resources / Namespace (Workloads)
Kubernetes / Compute Resources / Node (Pods)
Kubernetes / Compute Resources / Pod
Kubernetes / Compute Resources / Workload
Kubernetes / Controller Manager
Kubernetes / Kubelet
Kubernetes / Networking / Cluster
Kubernetes / Networking / Namespace (Pods)
Kubernetes / Networking / Namespace (Workload)
Kubernetes / Networking / Pod
Kubernetes / Networking / Workload
Kubernetes / Persistent Volumes
Kubernetes / Proxy
Kubernetes / Scheduler
Node Exporter / Nodes
Node Exporter / USE Method / Cluster
Node Exporter / USE Method / Node
Prometheus / Overview
Stash / Dashboard
etcd
```