When Containers have resource requests specified the scheduler can make better decisions about which nodes to place Pods on and how to deal with resource contention.

Limits and requests for memory are measured in bytes. You can express memory as a plain integer or as a fixed-point integer using one of these suffixes: E, P, T, G, M, K. You can also use the power-of-two equivalents: Ei, Pi, Ti, Gi, Mi, Ki. For example, the following represent roughly the same value:

```
128974848, 129e6, 129M, 123Mi
```

## Example

```yaml
...
resources:
  limits:
    memory: 200Mi
  requests:
    memory: 100Mi
...
```

## Notes

- A Container can exceed its memory request if the Node has memory available, although this is not allowed
- If a Container allocates more memory than its limit, the Container becomes a candidate for termination. It will be terminated if it continues to consume memory beyond its limit
- The memory request for the Pod is the sum of the memory requests for all the Containers in the Pod
- If a memory limit is not applied, the namespace's limit is automatically assigned via a [LimitRange](https://kubernetes.io/docs/tasks/administer-cluster/memory-default-namespace/). If this does not exist there is no upper bound to the memory a container can use

## External Links
- [Manage Container Compute Resources](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/)
- [Kubernetes - Understanding Resources](http://www.noqcks.io/note/kubernetes-resources-limits/)
- [Kubernetes Docs: Assign Memory Resources to Containers and Pods](https://kubernetes.io/docs/tasks/configure-pod-container/assign-memory-resource/)
- [Kubernetes Docs: Configure Memory and CPU Quotas for a Namespace](https://kubernetes.io/docs/tasks/administer-cluster/quota-memory-cpu-namespace/)
- [Configure Quality of Service for Pods](https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod/)
