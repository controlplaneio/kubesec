+++
title = "containers[] .securityContext .runAsNonRoot == true"
weight = 2
+++

## Force the running image to run as a non-root user to ensure least privilege

Indicates that containers should run as non-root user. 


## Notes

- Container level security context settings are applied to the specific container and override settings made at the pod level where there is
overlap
- Container level settings are not applied to the pod's volumes.

## External Links

- [Kubernertes Blog: Security Best Practices](http://blog.kubernetes.io/2016/08/security-best-practices-kubernetes-deployment.html
)


{{% katacoda %}}
