# kube-sec-check

Validate the security parameters of Kubernetes YAML resources.

Currently supported types: Pod, Deployment, StatefulSet, DaemonSet

## Usage

```bash
kseccheck.sh [options] <k8s resource file>
```


## TODO

1. short form output behind `-o` option (default)
1. long form output (details of problem, recommended fix, link to docs)
1. JSON output 
1. More test cases
