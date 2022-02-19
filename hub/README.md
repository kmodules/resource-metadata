## How to update Hub files in kube-ui-server

```bash
POD_NAMESPACE=kubeops
POD_NAME=$(kubectl get pods -n $POD_NAMESPACE -l app.kubernetes.io/instance=kube-ui-server -o jsonpath={.items[0].metadata.name})

kubectl cp hub $POD_NAME:/tmp -n $POD_NAMESPACE
# verify
kubectl exec -it $POD_NAME -n $POD_NAMESPACE -- ls -l /tmp/hub

LOCATION=hub/resourceeditors/kubedb.com/v1alpha2/elasticsearches.yaml
kubectl cp $LOCATION $POD_NAME:/tmp/$LOCATION -n $POD_NAMESPACE
# verify
kubectl exec -it $POD_NAME -n $POD_NAMESPACE -- cat /tmp/$LOCATION

# trigger reload
kubectl exec -it $POD_NAME -n $POD_NAMESPACE -- sh -c "date > /tmp/hub/resourceeditors/trigger"
# verify
kubectl exec -it $POD_NAME -n $POD_NAMESPACE -- cat /tmp/hub/resourceeditors/trigger
```
