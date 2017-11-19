+++
title = "securityContext capabilities"
weight = 2
+++

## Reducing kernel capabilities available to a container limits its attack surface

Capabilities permit certain named `root` actions without giving full `root` access. They are a more fine-grained permissions model, and all capabilities should be dropped from a pod, with only those required added back. 

There are a large number of capabilities, with `CAP_SYS_ADMIN` bounding most. Never enable this capability - it's equivalent to `root`.

## Example

```yaml

---
apiVersion: extensions/v1beta1
kind: Deployment
...
      containers:
      - name: payment
        image: nginx
        securityContext:
          capabilities:
            drop:
              - all
            add:
              - NET_BIND_SERVICE
```

## Notes

- Drop all capabilities from a pod as above
- Add only those required
- Run a comprehensive test suite to ensure security extensions have not blocked functionality that your containers or pods require

## External Links
- [Kubernetes Docs: Set capabilities for a Container](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/#set-capabilities-for-a-container)
- [Commands and Capabilities](https://lukemarsden.github.io/docs/user-guide/containers/)
- [Removing Setuid Binaries with Capabilities](https://linux-audit.com/linux-capabilities-hardening-linux-binaries-by-removing-setuid/)


{{% katacoda %}}
