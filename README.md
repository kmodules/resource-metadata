[![Go Report Card](https://goreportcard.com/badge/kmodules.xyz/resource-metadata)](https://goreportcard.com/report/kmodules.xyz/resource-metadata)
[![GoDoc](https://godoc.org/kmodules.xyz/resource-metadata?status.svg "GoDoc")](https://godoc.org/kmodules.xyz/resource-metadata)
[![Build Status](https://travis-ci.org/kmodules/resource-metadata.svg?branch=master)](https://travis-ci.org/kmodules/resource-metadata)
[![Slack](https://shields.io/badge/Join_Slack-salck?color=4A154B&logo=slack)](https://slack.appscode.com)
[![Twitter](https://img.shields.io/twitter/follow/appscodehq.svg?style=social&logo=twitter&label=Follow)](https://twitter.com/intent/follow?screen_name=AppsCodeHQ)

# resource-metadata

API for defining metadata about Kubernetes resources. You can read about the motivation and design of `ResourceDescriptors` in our blog post: https://blog.byte.builders/post/resourcedescriptor/

### Test GraphFinder api

```
curl -X POST \
  https://api.crd.builders/apis/meta.k8s.appscode.com/v1alpha1/graphfinders \
  -H 'content-type: application/json' \
  -d '{
   "apiVersion": "meta.k8s.appscode.com/v1alpha1",
   "kind": "GraphFinder",
   "request": {
      "source": {
         "group": "apps",
         "version": "v1",
         "resource": "deployments"
      }
   }
}'
```

### Test PathFinder api

```
# find path from deployments -> services

curl -X POST \
  https://api.crd.builders/apis/meta.k8s.appscode.com/v1alpha1/pathfinders \
  -H 'content-type: application/json' \
  -d '{
   "apiVersion": "meta.k8s.appscode.com/v1alpha1",
   "kind": "PathFinder",
   "request": {
      "source": {
         "group": "apps",
         "version": "v1",
         "resource": "deployments"
      },
      "target": {
         "group": "",
         "version": "v1",
         "resource": "services"
      }
   }
}'

# find all paths from deployments

curl -X POST \
  https://api.crd.builders/apis/meta.k8s.appscode.com/v1alpha1/pathfinders \
  -H 'content-type: application/json' \
  -d '{
   "apiVersion": "meta.k8s.appscode.com/v1alpha1",
   "kind": "PathFinder",
   "request": {
      "source": {
         "group": "apps",
         "version": "v1",
         "resource": "deployments"
      }
   }
}'
```

### Generate Path Diagram using Graphviz

```console
$ go run cmd/pathviz/main.go --group=kubedb.com --version=v1alpha1 --resource=postgreses | dot -Tpng > postgres.png
$ go run cmd/pathviz/main.go --group=apps --version=v1 --resource=deployments | dot -Tpng > deployment.png
```
