+++
title = "kubesec.io"
chapter = true
weight = 5
pre = "<b>1. </b>"
+++

# kubesec.io

Score Kubernetes resources for using security features 

## Usage

Try it here

<textarea rows=11 name="file" form="usrform" >
apiVersion: v1
kind: Pod
metadata:
  name: kubesec-demo
spec:
  containers:
  - name: kubesec-demo
    image: gcr.io/google-samples/node-hello:1.0
    securityContext:
      readOnlyRootFilesystem: true
</textarea>
<form action="/"  method="post" enctype="multipart/form-data" id="usrform">
  <input type="hidden" value="webfile" name="filename" />
  <input type="submit" value="submit" id="submit" class="btn btn-default">
</form>

---

<!--
-->

Define a BASH function

{{< highlight bash >}}
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
{{< /highlight >}}

POST a Kubernetes resource to kubesec.io   
{{< highlight bash >}}
kubesec ./deployment.yml
{{< /highlight >}}

Return non-zero status code is the score is not greater than 10
{{< highlight bash >}}
kubesec ./score-9-deployment.yml | jq --exit-status '.score > 10' >/dev/null 
# status code 1
{{< /highlight >}}

## Example output
{{< highlight json >}}
{
  "score": -30,
  "scoring": {
    "critical": [
      {
        "selector": "containers[] .securityContext .capabilities .add | index(\"SYS_ADMIN\")",
        "reason": "CAP_SYS_ADMIN is the most privileged capability and should always be avoided"
      }
    ],
    "advise": [
      {
        "selector": "containers[] .securityContext .runAsNonRoot == true",
        "reason": "Force the running image to run as a non-root user to ensure least privilege"
      },
      {
        "selector": "containers[] .securityContext .capabilities .drop",
        "reason": "Reducing kernel capabilities available to a container limits its attack surface"
      },
      {
        "selector": "containers[] .securityContext .readOnlyRootFilesystem == true",
        "reason": "An immutable root filesystem can prevent malicious binaries being added to PATH and increase attack cost"
      },
      {
        "selector": "containers[] .securityContext .runAsUser > 10000",
        "reason": "Run as a high-UID user to avoid conflicts with the host's user table"
      },
      {
        "selector": "containers[] .securityContext .capabilities .drop | index(\"ALL\")",
        "reason": "Drop all capabilities and add only those required to reduce syscall attack surface"
      }
    ]
  }
}
{{< /highlight >}}

---

{{% children  %}}
