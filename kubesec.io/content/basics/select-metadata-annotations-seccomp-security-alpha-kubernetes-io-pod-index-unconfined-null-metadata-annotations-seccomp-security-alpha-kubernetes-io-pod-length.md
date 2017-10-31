+++
title = "select(.metadata .annotations .\"seccomp.security.alpha.kubernetes.io/pod\" | index(\"unconfined\") == null) | .metadata .annotations .\"seccomp.security.alpha.kubernetes.io/pod\" | length > 0"
weight = 5
+++

## Seccomp profiles set minimum privilege and secure against unknown threats
