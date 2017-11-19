+++
title = ".spec .hostIPC"
weight = 0
+++

## Sharing the host's IPC namespace allows container processes to communicate with processes on the host

Removing namespaces from pods reduces isolation and allows the processes in the pod to perform tasks as if they were running natively on the host.

This circumvents the protection models that containers are based on and should only be done with absolutely certainty (for example, for low-level observation of other containers).

## External Links


{{% katacoda %}}
