+++
title = "kubesec.io"
chapter = true
weight = 5
pre = "<b>1. </b>"
+++

# kubesec.io

Score Kubernetes resources for using security features 

## Usage

Define a BASH function
```bash
kubesec () 
{ 
    local FILE="${1:-}";
    [[ ! -f "${FILE}" ]] && { 
        echo "kubesec: ${FILE}: No such file" >&2;
        return 1
    };
    curl --silent \
      --compressed \
      --connect-timeout 5 \
      -F file=@"${FILE}" \
      https://kubesec.io/
}
```

POST a Kubernetes resource to kubesec.io   
```bash
kubesec test/asset/score-1-pod-default.yml
```

Return non-zero status code is the score is not greater than 10
```bash
kubesec test/asset/score-1-pod-default.yml  | jq --exit-status '.score > 10' >/dev/null 
# status code 1

```
## Example output
```json
{
  "score": 2,
  "scoring": {
    "critical": [],
    "advise": [
      {
        "selector": "containers[] .securityContext .capabilities .drop",
        "reason": "Reducing kernel capabilities available to a container limits its attack surface",
        "href": "https://kubernetes.io/docs/tasks/configure-pod-container/security-context/"
      },
      {
        "selector": "containers[] .securityContext .capabilities .drop | index(\"ALL\")",
        "reason": "Drop all capabilities and add only those required to reduce syscall attack surface",
      },
      {
        "selector": "containers[] .securityContext .runAsNonRoot == true",
        "reason": "Force the running image to run as a non-root user to ensure least privilege",
      },
      ...
    ],
    "positive": [
      {
        "selector": ".spec .volumeClaimTemplates[] .spec .accessModes | index(\"ReadWriteOnce\")",
        "kind": "StatefulSet"
      },
      {
        "selector": ".spec .volumeClaimTemplates[] .spec .resources .requests .storage",
        "kind": "StatefulSet"
      },
      ...
    ],
    "negative": []
  }
}
```


{{% children  %}}
