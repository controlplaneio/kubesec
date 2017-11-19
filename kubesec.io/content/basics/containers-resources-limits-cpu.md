+++
title = "containers[] .resources .limits .cpu"
weight = 2
+++

## Enforcing CPU limits prevents DOS via resource exhaustion

When Containers have resource requests specified the scheduler can make better decisions about which nodes to place Pods on and how to deal with resource contention.

Limits and requests for CPU resources are measured in cpu units. Kubernetes judges these as:

- 1 AWS vCPU
- 1 GCP Core
- 1 Azure vCore
- 1 Hyperthread on a bare-metal Intel processor with Hyperthreading

## Notes
- Fractional requests are allowed. 
- CPU is always requested as an absolute quantity, never as a relative quantity; 0.1 is the same amount of CPU on a single-core, dual-core, or 48-core machine.
- Each node has a maximum capacity for each of the resource types: the amount of CPU and memory it can provide for Pods
- Although actual memory or CPU resource usage on nodes is very low, the scheduler still refuses to place a Pod on a node if the capacity check fails
- If a CPU limit is not applied, the namespace's limit is automatically assigned via a [LimitRange](https://kubernetes.io/docs/tasks/administer-cluster/cpu-default-namespace/). If this does not exist there is no upper bound to the memory a container can use


## External Links
- [Manage Container Compute Resources](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/)
- [Kubernetes - Understanding Resources](http://www.noqcks.io/note/kubernetes-resources-limits/)
- [Kubernetes Docs: Assign CPU Resources to Containers and Pods](https://kubernetes.io/docs/tasks/configure-pod-container/assign-cpu-resource/)
- [Kubernetes Docs: Configure Memory and CPU Quotas for a Namespace](https://kubernetes.io/docs/tasks/administer-cluster/quota-memory-cpu-namespace/)
- [Configure Quality of Service for Pods](https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod/)


{{% katacoda %}}
