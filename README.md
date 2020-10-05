# kubesec [![Build Status](https://travis-ci.com/controlplaneio/kubesec.svg?token=2zTFdbp4Jrcox4MuDKaD&branch=master)](https://travis-ci.com/controlplaneio/kubesec)

<!-- markdownlint-disable no-inline-html header-increment -->
<!-- markdownlint-disable line-length -->

#### <center>🚨 v1 API is deprecated, please read the <a href="https://github.com/controlplaneio/kubesec/blob/master/README.md#release-notes" target="_blank">release notes</a> 🚨</center>

<!-- markdownlint-enable line-length -->

### Security risk analysis for Kubernetes resources

<p align="center">
  <img src="https://casual-hosting.s3.amazonaws.com/kubesec-logo.png">
</p>

## Live demo

[Visit Kubesec.io](https://kubesec.io)

This uses ControlPlane's hosted API at [v2.kubesec.io/scan](https://v2.kubesec.io/scan)

---

- [Download Kubesec](#download-kubesec)
  - [Command line usage](#command-line-usage-)
  - [Usage example](#usage-example-)
  - [Docker usage](#docker-usage-)
- [Kubesec HTTP Server](#kubesec-http-server)
  - [CLI usage example](#cli-usage-example-)
  - [Docker usage example](#docker-usage-example-)
- [Kubesec-as-a-Service](#kubesec-as-a-service)
  - [Command line usage](#command-line-usage--1)
  - [Usage example](#usage-example--1)
- [Example output](#example-output)
- [Contributors](#contributors)
- [Contributing](#contributing)
- [Getting Help](#getting-help)
- [Release Notes](#release-notes)
  - [2.5.0](#250)
  - [2.4.0](#240)
  - [2.3.1](#231)
  - [2.3.0](#230)
  - [2.2.0](#220)
  - [2.1.0](#210)
  - [2.0.0](#200)
  - [1.0.0](#100)

## Download Kubesec

Kubesec is available as a:

- [Docker container image](https://hub.docker.com/r/kubesec/kubesec/tags) at `docker.io/kubesec/kubesec:v2`
- Linux/MacOS/Win binary (get the [latest release](https://github.com/controlplaneio/kubesec/releases))
- [Kubernetes Admission Controller](https://github.com/controlplaneio/kubesec-webhook)
- [Kubectl plugin](https://github.com/controlplaneio/kubectl-kubesec)

Or install the latest commit from GitHub with `go get -u github.com/controlplaneio/kubesec/cmd/kubesec`

#### Command line usage:

```bash
$ kubesec scan k8s-deployment.yaml
```

#### Usage example:

```bash
$ cat <<EOF > kubesec-test.yaml
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
EOF
$ kubesec scan kubesec-test.yaml
```

#### Docker usage:

Run the same command in Docker:

```bash
$ docker run -i kubesec/kubesec:512c5e0 scan /dev/stdin < kubesec-test.yaml
```

## Kubesec HTTP Server

Kubesec includes a bundled HTTP server

#### CLI usage example:

Start the HTTP server in the background

<!-- markdownlint-disable line-length -->

```bash
$ kubesec http 8080 &
[1] 12345
{"severity":"info","timestamp":"2019-05-12T11:58:34.662+0100","caller":"server/server.go:69","message":"Starting HTTP server on port 8080"}
```

<!-- markdownlint-enable line-length -->

Use curl to POST a file to the server

```bash
$ curl -sSX POST --data-binary @test/asset/score-0-cap-sys-admin.yml http://localhost:8080/scan
[
  {
    "object": "Pod/security-context-demo.default",
    "valid": true,
    "message": "Failed with a score of -30 points",
    "score": -30,
    "scoring": {
      "critical": [
        {
          "selector": "containers[] .securityContext .capabilities .add == SYS_ADMIN",
          "reason": "CAP_SYS_ADMIN is the most privileged capability and should always be avoided",
          "points": -30
        },
        {
          "selector": "containers[] .securityContext .runAsNonRoot == true",
          "reason": "Force the running image to run as a non-root user to ensure least privilege",
          "points": 1
        },
  // ...
```

Finally, stop the Kubesec server by killing the background process

```bash
$ kill %
```

#### Docker usage example:

Start the HTTP server using Docker

```bash
$ docker run -d -p 8080:8080 kubesec/kubesec:d432be9 http 8080
```

Use curl to POST a file to the server

```bash
$ curl -sSX POST --data-binary @test/asset/score-0-cap-sys-admin.yml http://localhost:8080/scan
...
```

Don't forget to stop the server.

## Kubesec-as-a-Service

Kubesec is also available via HTTPS at [v2.kubesec.io/scan](https://v2.kubesec.io/scan)

#### Command line usage:

```bash
$ curl -sSX POST --data-binary @"k8s-deployment.yaml" https://v2.kubesec.io/scan
```

#### Usage example:

Define a BASH function

```bash
$ kubesec ()
{
    local FILE="${1:-}";
    [[ ! -e "${FILE}" ]] && {
        echo "kubesec: ${FILE}: No such file" >&2;
        return 1
    };
    curl --silent \
      --compressed \
      --connect-timeout 5 \
      -sSX POST \
      --data-binary=@"${FILE}" \
      https://v2.kubesec.io/scan
}
```

POST a Kubernetes resource to v2.kubesec.io/scan

```bash
$ kubesec ./deployment.yml
```

Return non-zero status code is the score is not greater than 10

```bash
$ kubesec ./score-9-deployment.yml | jq --exit-status '.score > 10' >/dev/null
# status code 1
```

## Example output

Kubesec returns a returns a JSON array, and can scan multiple YAML documents in a single input file.

```json
[
  {
    "object": "Pod/security-context-demo.default",
    "valid": true,
    "message": "Failed with a score of -30 points",
    "score": -30,
    "scoring": {
      "critical": [
        {
          "selector": "containers[] .securityContext .capabilities .add == SYS_ADMIN",
          "reason": "CAP_SYS_ADMIN is the most privileged capability and should always be avoided",
          "points": -30
        }
      ],
      "advise": [
        {
          "selector": "containers[] .securityContext .runAsNonRoot == true",
          "reason": "Force the running image to run as a non-root user to ensure least privilege",
          "points": 1
        },
        {
          // ...
        }
      ]
    }
  }
]
```

---

## Contributors

Thanks to our awesome contributors!

- [Andrew Martin](@sublimino)
- [Stefan Prodan](@stefanprodan)

## Contributing

Kubesec is Apache 2.0 licensed and accepts contributions via GitHub pull requests.

When submitting bug reports please include as much details as possible:

- which Kubesec version
- which Kubernetes version
- what happened (Kubesec logs and expected output)

## Getting Help

If you have any questions about Kubesec and Kubernetes security:

- Read the Kubesec docs
- Reach out on Twitter to [@sublimino](https://twitter.com/sublimino) or [@controlplaneio](https://twitter.com/controlplaneio)
- File an issue

Your feedback is always welcome!

# Release Notes

## 2.7.2

- bump go and alpine versions
  - this is also part of making `go mod` happy with `v2`

## 2.7.1

- further fixes to make `go mod` happy with `v2`
  - should resolve issues with tools that use `go list ./...` at the project root

## 2.7.0

- fix go mod issues with `v2`
  - can use `go get` again

## 2.6.0

- allow for piping into `kubesec scan` using `-` or `/dev/stdin`
  - `cat somefile.yml | kubesec scan -`
  - `cat somefile.yml | kubesec scan /dev/stdin`

## 2.5.0

- improved in-toto integration

## 2.4.0

- added passed to the JSON output
- note: repo tests now require `jq` - **only concerns maintainers**

## 2.3.1

- patch to accept form data from the <https://kubesec.io> webpage sample form

## 2.3.0

- moved everything to go modules

## 2.2.0

- added in-toto support

## 2.1.0

- add rule for `allowPrivilegeEscalation: true` with a score of -7
- add `points` field to each recommendation so the values that comprise the total score can be seen
- fix case sensitivity bug in `.capabilities.drop | index("ALL")`
- rules in `critical` and `advise` lists prioritised and returned in same order across runs

## 2.0.0

- first open source release
- passes same acceptance tests as Kubesec v1
- more stringent analysis: scoring for a rule is multiplied by number of matches (previously the score was only applied
  once), initContainers are included in score, new securityContext directive support, seccomp and apparmor pod-targeting
  tighter
- CLI and HTTP server bundled in single binary

## 1.0.0

- initial release at <https://kubesec.io>
- closed source
