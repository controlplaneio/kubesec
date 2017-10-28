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

## JSON struct

```
{
  "points": 1,
  "scoring": {
    "critical": [
      {
        "reason": "you ran priv!",
        "points": -100,
        "href": "https://more-info.com"
      }
    ],
    "advisory": [
      {
        "reason": "you should fix this",
        "href": "http:/",
        "points": -1
      }
    ],
    "positive": [
      {
        "reason": "well done",
        "points": 100,
        "href": "https://more-info.com"
      }
    ]
  }
}
```

## To add

1. contents of https://kubernetes.io/docs/concepts/policy/pod-security-policy/
1. examples from https://github.com/kubernetes/kubernetes/tree/master/examples
1. https://github.com/weaveworks/weave/blob/master/prog/weave-kube/weave-daemonset.yaml
