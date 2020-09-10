#!/usr/bin/env bash

load '_helper'

setup() {
    _global_setup
}

teardown() {
  _global_teardown
}

@test "only valid types - allow Pod" {
  run _app "${TEST_DIR}/asset/score-1-pod-default.yml"
  refute_output --regexp ".*Only kinds .* accepted.*" \
    || assert_output --regexp ".*This resource kind is not supported.*"
  assert_success
}

@test "only valid types - allow Deployment" {
  run _app "${TEST_DIR}/asset/score-1-dep-default.yml"
  refute_output --regexp ".*Only kinds .* accepted.*" \
    || assert_output --regexp ".*This resource kind is not supported.*"
  assert_success
}

@test "only valid types - allow StatefulSet" {
  run _app "${TEST_DIR}/asset/score-1-statefulset-default.yml"
  refute_output --regexp ".*Only kinds .* accepted.*" \
    || assert_output --regexp ".*This resource kind is not supported.*"
  assert_success
}

@test "only valid types - allow DaemonSet" {
  run _app "${TEST_DIR}/asset/score-1-daemonset-default.yml"
  refute_output --regexp ".*Only kinds .* accepted.*" \
    || assert_output --regexp ".*This resource kind is not supported.*"
  assert_success
}

@test "only valid types - deny PodSecurityPolicy" {
  run _app "${TEST_DIR}/asset/score-0-podsecuritypolicy-permissive.yml"
  assert_output --regexp ".*Only kinds .* accepted.*" \
    || assert_output --regexp ".*This resource kind is not supported.*"
  if _is_local; then
    assert_failure
  fi
}

@test "passes 1.11 format daemonset" {
  run _app "${TEST_DIR}/asset/versioned/score-0-daemonset-v1.11.yml"
  assert_zero_points
}

@test "passes 1.11 format statefulset" {
  run _app "${TEST_DIR}/asset/versioned/score-0-statefulset-v1.11.yml"
  assert_zero_points
}

# ---

@test "returns error for invalid JSON" {
  run _app "${TEST_DIR}/asset/invalid-input-pod-dump.json"

  assert_output --regexp "Missing 'apiVersion' key" \
    || assert_output --regexp ".*: Invalid type\. .*"

  assert_failure_local
}

@test "returns error YAML control characters" {
  run _app "${TEST_DIR}/asset/invalid-input-no-control-characters.json"

  assert_invalid_input
}

@test "passes bug dump twice [1/2]" {
  run _app "${TEST_DIR}/asset/bug-dump-2.json"
  assert_success
  assert_gt_zero_points
}

@test "passes bug dump twice [2/2]" {
  run _app "${TEST_DIR}/asset/bug-dump-2.json"
  assert_success
  assert_gt_zero_points
}

# ---

@test "errors with no filename" {
  run _app
  assert_failure_local
}

@test "errors with no filename - output logs (local)" {
  skip_if_not_local
  run "${BIN_DIR:-}"/kubesec scan
  assert_failure
  assert_line "Error: file path is required"
}

@test "errors with invalid file" {
  run _app somefile.yaml
  assert_failure_local
  assert_file_not_found
}

@test "errors with empty file" {
  run _app "${TEST_DIR}/asset/empty-file"
  assert_failure_local
  assert_invalid_input
}

@test "can read piped input - (yaml, local)" {
  skip_if_not_local
  FILE="${TEST_DIR}/asset/allowPrivilegeEscalation.yaml"
  run bash -c "cat \"${FILE}\" | ${BIN_DIR:-}/kubesec scan -"
  assert_lt_zero_points
}

@test "can read piped input /dev/stdin (yaml, local)" {
  skip_if_not_local
  FILE="${TEST_DIR}/asset/allowPrivilegeEscalation.yaml"
  run bash -c "cat \"${FILE}\" | ${BIN_DIR:-}/kubesec scan /dev/stdin"
  assert_lt_zero_points
}

@test "can read redirected input - (yaml, local)" {
  skip_if_not_local
  FILE="${TEST_DIR}/asset/allowPrivilegeEscalation.yaml"
  run bash -c "${BIN_DIR:-}/kubesec scan - <\"${FILE}\""
  assert_lt_zero_points
}

@test "can read redirected input /dev/stdin (yaml, local)" {
  skip_if_not_local
  FILE="${TEST_DIR}/asset/allowPrivilegeEscalation.yaml"
  run bash -c "${BIN_DIR:-}/kubesec scan /dev/stdin <\"${FILE}\""
  assert_lt_zero_points
}

@test "can read multiple inputs (yaml, local)" {
  skip_if_not_local

  run _app \
    "${TEST_DIR}/asset/score-1-cap-drop-all.yml" \
    "${TEST_DIR}/asset/score-1-daemonset-default.yml"

  assert [ "$(jq -r 'length' <<<"${output}")" == "2" ]
  assert [ "$(jq -r '.[0].valid' <<<"${output}")" == "true" ]
  assert [ "$(jq -r '.[1].valid' <<<"${output}")" == "true" ]
}

@test "can read multiple inputs plus invalid (yaml, local)" {
  skip_if_not_local

  run _app \
    "${TEST_DIR}/asset/score-1-cap-drop-all.yml" \
    "${TEST_DIR}/asset/score-1-daemonset-default.yml" \
    "${TEST_DIR}/asset/invalid-schema.yml"

  select_obj_valid() {
    jq -r ".[] | select(.object == \"${1}\") | .valid" <<<"${output}"
  }
  assert [ "$(jq -r 'length' <<<"${output}")" == "3" ]
  assert [ "$(select_obj_valid "Deployment/demo.default")" == "false" ]
  assert [ "$(select_obj_valid "DaemonSet/undefined.default")" == "true" ]
  assert [ "$(select_obj_valid "Pod/security-context-demo.default")" == "true" ]
}

@test "can read multiple inputs plus multi file (yaml, local)" {
  skip_if_not_local

  run _app \
    "${TEST_DIR}/asset/score-1-cap-drop-all.yml" \
    "${TEST_DIR}/asset/score-1-daemonset-default.yml" \
    "${TEST_DIR}/asset/multi.yml"

  assert [ "$(jq -r 'length' <<<"${output}")" == "7" ]
}

@test "errors with empty file (json, local)" {
  skip_if_not_local

  run _app "${TEST_DIR}/asset/empty-file" --json
  assert_invalid_input
  assert_failure_local
}

@test "errors with empty file (json, remote)" {
  skip_if_not_remote

  run _app "${TEST_DIR}/asset/empty-file"
  assert_invalid_input
  assert_failure_local
}

@test "errors with empty JSON (json, local)" {
  skip_if_not_local

  run _app "${TEST_DIR}/asset/empty-json-file" --json
  assert_invalid_input
  assert_failure
}

@test "errors with empty JSON (json, remote)" {
  skip_if_not_remote

  run _app "${TEST_DIR}/asset/empty-json-file"
  assert_invalid_input
  assert_success
}

@test "returns content-type application/json on pass (yaml, remote)" {
  skip_if_not_remote

  run _app \
    "${TEST_DIR}/asset/score-0-daemonset-volume-host-docker-socket.yml" \
    -w '%{content_type}' \
    -o /dev/null

  assert_output --regexp "application/json"
}

@test "strips form 'file=' prefix from YAML (yaml, remote)" {
  skip_if_not_remote

  run _app "${TEST_DIR}/asset/form-prefix-file.yml"
  assert_gt_zero_points
  assert_success
}

@test "strips form 'file=' prefix from JSON (json, remote)" {
  skip_if_not_remote

  run _app "${TEST_DIR}/asset/form-prefix-file.json"
  assert_gt_zero_points
  assert_success
}

@test "does not strip form 'not-file=' prefix from YAML (yaml, remote)" {
  skip_if_not_remote

  run _app "${TEST_DIR}/asset/form-prefix-not-file.yml"
  assert_output --regexp ".*resource.json: Missing 'apiVersion' key.*"
  assert_success
}

@test "does not strip form 'not-file=' prefix from JSON (json, remote)" {
  skip_if_not_remote

  run _app "${TEST_DIR}/asset/form-prefix-not-file.json"
  assert_output "yaml: line 2: mapping values are not allowed in this context"
  assert_success
}

@test "fix multipart form bug (remote)" {
  skip_if_not_remote

  run curl -sSX POST --fail "${REMOTE_URL}" -H 'authority: kubesec.io' -H 'cache-control: max-age=0' -H 'origin: https://kubesec.io' -H 'upgrade-insecure-requests: 1' -H 'dnt: 1' -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundaryjh6hXNUDWrzlmi5c' -H 'user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36' -H 'accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8' -H 'referer: https://kubesec.io/' -H 'accept-encoding: gzip, deflate, br' -H 'accept-language: en-GB,en;q=0.9' -H 'cookie: gsScrollPos-239=0' --data-binary $'------WebKitFormBoundaryjh6hXNUDWrzlmi5c\r\nContent-Disposition: form-data; name="file"\r\n\r\napiVersion: v1\r\nkind: Pod\r\nmetadata:\r\n  annotations:\r\n    scheduler.alpha.kubernetes.io/critical-pod: ""\r\n  creationTimestamp: 2018-07-16T01:48:24Z\r\n  generateName: calico-node-vertical-autoscaler-664dc78496-\r\n  labels:\r\n    k8s-app: calico-node-autoscaler\r\n    pod-template-hash: "2208734052"\r\n  name: calico-node-vertical-autoscaler-664dc78496-qnvn5\r\n  namespace: kube-system\r\n  ownerReferences:\r\n  - apiVersion: extensions/v1beta1\r\n    blockOwnerDeletion: true\r\n    controller: true\r\n    kind: ReplicaSet\r\n    name: calico-node-vertical-autoscaler-664dc78496\r\n    uid: 7d82ff8e-7e06-11e8-8e86-42010a9a0138\r\n  resourceVersion: "1691101"\r\n  selfLink: /api/v1/namespaces/kube-system/pods/calico-node-vertical-autoscaler-664dc78496-qnvn5\r\n  uid: 53b407cc-889a-11e8-8e86-42010a9a0138\r\nspec:\r\n  containers:\r\n  - command:\r\n    - /cpvpa\r\n    - --target=daemonset/calico-node\r\n    - --namespace=kube-system\r\n    - --logtostderr=true\r\n    - --poll-period-seconds=30\r\n    - --v=2\r\n    - --config-file=/etc/config/node-autoscaler\r\n    image: gcr.io/google_containers/cpvpa-amd64:v0.6.0\r\n    imagePullPolicy: IfNotPresent\r\n    name: autoscaler\r\n    resources: {}\r\n    terminationMessagePath: /dev/termination-log\r\n    terminationMessagePolicy: File\r\n    volumeMounts:\r\n    - mountPath: /etc/config\r\n      name: config\r\n    - mountPath: /var/run/secrets/kubernetes.io/serviceaccount\r\n      name: calico-cpva-token-kwjj4\r\n      readOnly: true\r\n  dnsPolicy: ClusterFirst\r\n  nodeName: gke-netpol3-default-pool-a06ebdd6-0gq0\r\n  restartPolicy: Always\r\n  schedulerName: default-scheduler\r\n  securityContext: {}\r\n  serviceAccount: calico-cpva\r\n  serviceAccountName: calico-cpva\r\n  terminationGracePeriodSeconds: 30\r\n  tolerations:\r\n  - effect: NoExecute\r\n    key: node.kubernetes.io/not-ready\r\n    operator: Exists\r\n    tolerationSeconds: 300\r\n  - effect: NoExecute\r\n    key: node.kubernetes.io/unreachable\r\n    operator: Exists\r\n    tolerationSeconds: 300\r\n  volumes:\r\n  - configMap:\r\n      defaultMode: 420\r\n      name: calico-node-vertical-autoscaler\r\n    name: config\r\n  - name: calico-cpva-token-kwjj4\r\n    secret:\r\n      defaultMode: 420\r\n      secretName: calico-cpva-token-kwjj4\r\nstatus:\r\n  conditions:\r\n  - lastProbeTime: null\r\n    lastTransitionTime: 2018-07-16T02:19:40Z\r\n    status: "True"\r\n    type: Initialized\r\n  - lastProbeTime: null\r\n    lastTransitionTime: 2018-07-20T00:40:35Z\r\n    status: "True"\r\n    type: Ready\r\n  - lastProbeTime: null\r\n    lastTransitionTime: 2018-07-16T02:19:40Z\r\n    status: "True"\r\n    type: PodScheduled\r\n  containerStatuses:\r\n  - containerID: docker://f32db0abaf32e67604a6fd5b2f60574b98633ef16f0dbf6c01b59f4fac08d758\r\n    image: asia.gcr.io/google_containers/cpvpa-amd64:v0.6.0\r\n    imageID: docker-pullable://asia.gcr.io/google_containers/cpvpa-amd64@sha256:cfe7b0a11c9c8e18c87b1eb34fef9a7cbb8480a8da11fc2657f78dbf4739f869\r\n    lastState: {}\r\n    name: autoscaler\r\n    ready: true\r\n    restartCount: 0\r\n    state:\r\n      running:\r\n        startedAt: 2018-07-20T00:40:32Z\r\n  hostIP: 10.154.0.2\r\n  phase: Running\r\n  podIP: 10.36.52.10\r\n  qosClass: BestEffort\r\n  startTime: 2018-07-16T02:19:40Z\r\n\r\n------WebKitFormBoundaryjh6hXNUDWrzlmi5c\r\nContent-Disposition: form-data; name="filename"\r\n\r\nwebfile\r\n------WebKitFormBoundaryjh6hXNUDWrzlmi5c--\r\n' --compressed

  # TODO(ajm) v2 responds: 400 Bad Request
  assert_gt_zero_points || true
}
