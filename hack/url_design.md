https://console.byte.builders/kubernetes/cluster-admin@gke_ackube_us-central1-f_demo.pharmer/daemonset/fluentd-elasticsearch?namespace=default








{
	"basic": {
		"name": "xyz"
	},
	"selfTables": {
		"ports": "",
	}
}




https://console.byte.builders/{user_id}/{cluster_name}/daemonset/fluentd-elasticsearch?namespace=default


https://console.byte.builders/{user_id}/kubernetes/{cluster_name}/{group}/{version}/{resource}/{name}?namespace=default


```json
{
  "sections": [
    {
      "section_name": "Cluster",
      "sub_sections": [
        {
          "ref": null,
          "name": "Basic",
          "editable": false
        },
        {
          "ref": {
            "kind": "Machine",
            "group": "cluster.k8s.io",
            "version": "v1alpha1",
            "resource": "machines",
            "namespaced": false
          },
          "name": "Machines",
          "editable": false
        },
        {
          "ref": {
            "kind": "MachineSet",
            "group": "cluster.k8s.io",
            "version": "v1alpha1",
            "resource": "machinesets",
            "namespaced": true
          },
          "name": "Machine Sets",
          "editable": false
        }
      ]
    },
    {
      "section_name": "Workloads",
      "sub_sections": [
        {
          "ref": {
            "kind": "Deployment",
            "group": "extensions",
            "version": "v1beta1",
            "resource": "deployments",
            "namespaced": true
          },
          "name": "Deployments",
          "editable": false
        },
        {
          "ref": {
            "kind": "Replica Set",
            "group": "apps",
            "version": "v1",
            "resource": "replicasets",
            "namespaced": true
          },
          "name": "Replica Sets",
          "editable": false
        },
        {
          "ref": {
            "kind": "ReplicationController",
            "group": "core",
            "version": "v1",
            "resource": "replicationcontrollers",
            "namespaced": true
          },
          "name": "Replication Controllers",
          "editable": false
        },
        {
          "ref": {
            "kind": "DaemonSet",
            "group": "apps",
            "version": "v1",
            "resource": "daemonsets",
            "namespaced": true
          },
          "name": "Daemon Sets",
          "editable": false
        },
        {
          "ref": {
            "kind": "StatefulSet",
            "group": "apps",
            "version": "v1",
            "resource": "statefulsets",
            "namespaced": true
          },
          "name": "Stateful Sets",
          "editable": false
        },
        {
          "ref": {
            "kind": "Job",
            "group": "batch",
            "version": "v1",
            "resource": "jobs",
            "namespaced": true
          },
          "name": "Jobs",
          "editable": false
        },
        {
          "ref": {
            "kind": "Pod",
            "group": "core",
            "version": "v1",
            "resource": "pods",
            "namespaced": true
          },
          "name": "Pods",
          "editable": false
        }
      ]
    }
  ]
}
```