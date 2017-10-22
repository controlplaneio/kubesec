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
        "reason": "gremlins in the system",
        "href": "http:/",
        "points": 1
      }
    ],
    "positive": [
      {
        "reason": "you're really nice",
        "points": 100,
        "href": "https://more-info.com"
      }
    ]
  }
}
```
