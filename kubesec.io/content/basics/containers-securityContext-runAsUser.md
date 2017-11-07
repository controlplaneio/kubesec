+++
title = "containers[] .securityContext .runAsUser > 10000"
weight = 5
+++

## Run as a high-UID user to avoid conflicts with the host's user table

RunAsUser is the UID to run the entrypoint of the container process. The user id that runs the first process in the container. 

# Notes
- `MustRunAs` - Requires a range to be configured. Uses the first value of the range as the default. Validates against the configured range.
- `MustRunAsNonRoot` - Requires that the pod be submitted with a non-zero runAsUser or have the USER directive defined in the image. No default provided.
- `RunAsAny` - No default provided. Allows any runAsUser to be specified.

## External Links

- [Kubernetes Docs: Pod Security Policy](https://kubernetes.io/docs/concepts/policy/pod-security-policy/#runasuser)


{{% katacoda %}}
