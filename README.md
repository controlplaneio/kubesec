# kubesec [![Build Status](https://travis-ci.com/controlplaneio/kubesec.svg?token=2zTFdbp4Jrcox4MuDKaD&branch=master)](https://travis-ci.com/controlplaneio/kubesec)

<p align="center">
  <img src="http://casual-hosting.s3.amazonaws.com/kubesec-logo.png">
</p>

Validate the security parameters of Kubernetes YAML resources.

Currently supported types: Pod, Deployment, StatefulSet, DaemonSet

# Why?

# Usage

Kubesec can run as a CLI tool or an HTTP web server, and is offered as a hosted service, a Kubernetes Admission Controller, and a `kubectl` plugin.

## CLI

```bash
$ kubesec scan test/asset/score-0-cap-sys-admin.yml
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
          "reason": "CAP_SYS_ADMIN is the most privileged capability and should always be avoided"
        }
      ],
      "advise": [
        {
          "selector": ".spec .hostAliases",
          "reason": "Managing /etc/hosts aliases can prevent Docker from modifying the file after a pod's containers have already been started"
        },
...
```

## HTTP

```bash
$ kubesec http 8080 &
[1] 12345
{"severity":"info","timestamp":"2019-05-12T11:58:34.662+0100","caller":"server/server.go:69","message":"Starting HTTP server on port 8080"}
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
          "reason": "CAP_SYS_ADMIN is the most privileged capability and should always be avoided"
        }
      ],
      "advise": [
        {
          "selector": ".spec .hostAliases",
          "reason": "Managing /etc/hosts aliases can prevent Docker from modifying the file after a pod's containers have already been started"
        },
...
$ kill % # to stop the background process
```

## Hosted Service


# Release Notes
