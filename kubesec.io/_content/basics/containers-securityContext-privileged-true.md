Privileged containers share namespaces with the host system, eschew cgroup restrictions, and do not offer any security. They should be used exclusively as a bundling and distribution mechanism for the code in the container, and not for isolation.

## Notes

- Processes within the container get almost the same privileges that are available to processes outside a container
- Privileged containers have significantly fewer kernel isolation features
- `root` inside a privileged container is close to `root` on the host as User Namespaces are not enforced 
- Privileged containers shared `/dev` with the host, which allows mounting of the host's filesystem
- They can also interact with the kernel to load kernel and alter settings (including the hostname), interfere with the network stack, and many other subtle permissions

## External Links

- [Kubernetes Docs: Privileged mode for pod containers](https://kubernetes.io/docs/concepts/workloads/pods/pod/#privileged-mode-for-pod-containers)
