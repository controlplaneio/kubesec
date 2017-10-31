+++
title = "select(.metadata .annotations .\"container.seccomp.security.alpha.kubernetes.io/explorer\" | index(\"unconfined\") == null) | .metadata .annotations .\"container.seccomp.security.alpha.kubernetes.io/explorer\" | length > 0"
weight = 5
+++

## Seccomp profiles for OpenShift set minimum privilege and secure against unknown threats
