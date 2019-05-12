#!/usr/bin/env bash

load '_helper'

setup() {
    _global_setup
}

teardown() {
  _global_teardown
}

@test "fix multipart form bug" {
  run curl -sSX POST --fail "${REMOTE_URL}" -H 'authority: kubesec.io' -H 'cache-control: max-age=0' -H 'origin: https://kubesec.io' -H 'upgrade-insecure-requests: 1' -H 'dnt: 1' -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundaryjh6hXNUDWrzlmi5c' -H 'user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36' -H 'accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8' -H 'referer: https://kubesec.io/' -H 'accept-encoding: gzip, deflate, br' -H 'accept-language: en-GB,en;q=0.9' -H 'cookie: gsScrollPos-239=0' --data-binary $'------WebKitFormBoundaryjh6hXNUDWrzlmi5c\r\nContent-Disposition: form-data; name="file"\r\n\r\napiVersion: v1\r\nkind: Pod\r\nmetadata:\r\n  annotations:\r\n    scheduler.alpha.kubernetes.io/critical-pod: ""\r\n  creationTimestamp: 2018-07-16T01:48:24Z\r\n  generateName: calico-node-vertical-autoscaler-664dc78496-\r\n  labels:\r\n    k8s-app: calico-node-autoscaler\r\n    pod-template-hash: "2208734052"\r\n  name: calico-node-vertical-autoscaler-664dc78496-qnvn5\r\n  namespace: kube-system\r\n  ownerReferences:\r\n  - apiVersion: extensions/v1beta1\r\n    blockOwnerDeletion: true\r\n    controller: true\r\n    kind: ReplicaSet\r\n    name: calico-node-vertical-autoscaler-664dc78496\r\n    uid: 7d82ff8e-7e06-11e8-8e86-42010a9a0138\r\n  resourceVersion: "1691101"\r\n  selfLink: /api/v1/namespaces/kube-system/pods/calico-node-vertical-autoscaler-664dc78496-qnvn5\r\n  uid: 53b407cc-889a-11e8-8e86-42010a9a0138\r\nspec:\r\n  containers:\r\n  - command:\r\n    - /cpvpa\r\n    - --target=daemonset/calico-node\r\n    - --namespace=kube-system\r\n    - --logtostderr=true\r\n    - --poll-period-seconds=30\r\n    - --v=2\r\n    - --config-file=/etc/config/node-autoscaler\r\n    image: gcr.io/google_containers/cpvpa-amd64:v0.6.0\r\n    imagePullPolicy: IfNotPresent\r\n    name: autoscaler\r\n    resources: {}\r\n    terminationMessagePath: /dev/termination-log\r\n    terminationMessagePolicy: File\r\n    volumeMounts:\r\n    - mountPath: /etc/config\r\n      name: config\r\n    - mountPath: /var/run/secrets/kubernetes.io/serviceaccount\r\n      name: calico-cpva-token-kwjj4\r\n      readOnly: true\r\n  dnsPolicy: ClusterFirst\r\n  nodeName: gke-netpol3-default-pool-a06ebdd6-0gq0\r\n  restartPolicy: Always\r\n  schedulerName: default-scheduler\r\n  securityContext: {}\r\n  serviceAccount: calico-cpva\r\n  serviceAccountName: calico-cpva\r\n  terminationGracePeriodSeconds: 30\r\n  tolerations:\r\n  - effect: NoExecute\r\n    key: node.kubernetes.io/not-ready\r\n    operator: Exists\r\n    tolerationSeconds: 300\r\n  - effect: NoExecute\r\n    key: node.kubernetes.io/unreachable\r\n    operator: Exists\r\n    tolerationSeconds: 300\r\n  volumes:\r\n  - configMap:\r\n      defaultMode: 420\r\n      name: calico-node-vertical-autoscaler\r\n    name: config\r\n  - name: calico-cpva-token-kwjj4\r\n    secret:\r\n      defaultMode: 420\r\n      secretName: calico-cpva-token-kwjj4\r\nstatus:\r\n  conditions:\r\n  - lastProbeTime: null\r\n    lastTransitionTime: 2018-07-16T02:19:40Z\r\n    status: "True"\r\n    type: Initialized\r\n  - lastProbeTime: null\r\n    lastTransitionTime: 2018-07-20T00:40:35Z\r\n    status: "True"\r\n    type: Ready\r\n  - lastProbeTime: null\r\n    lastTransitionTime: 2018-07-16T02:19:40Z\r\n    status: "True"\r\n    type: PodScheduled\r\n  containerStatuses:\r\n  - containerID: docker://f32db0abaf32e67604a6fd5b2f60574b98633ef16f0dbf6c01b59f4fac08d758\r\n    image: asia.gcr.io/google_containers/cpvpa-amd64:v0.6.0\r\n    imageID: docker-pullable://asia.gcr.io/google_containers/cpvpa-amd64@sha256:cfe7b0a11c9c8e18c87b1eb34fef9a7cbb8480a8da11fc2657f78dbf4739f869\r\n    lastState: {}\r\n    name: autoscaler\r\n    ready: true\r\n    restartCount: 0\r\n    state:\r\n      running:\r\n        startedAt: 2018-07-20T00:40:32Z\r\n  hostIP: 10.154.0.2\r\n  phase: Running\r\n  podIP: 10.36.52.10\r\n  qosClass: BestEffort\r\n  startTime: 2018-07-16T02:19:40Z\r\n\r\n------WebKitFormBoundaryjh6hXNUDWrzlmi5c\r\nContent-Disposition: form-data; name="filename"\r\n\r\nwebfile\r\n------WebKitFormBoundaryjh6hXNUDWrzlmi5c--\r\n' --compressed

  # TODO(ajm) v2 responds: 400 Bad Request
  assert_gt_zero_points || true
}

# ---

# v1.11 tests

@test "passes 1.11 format daemonset" {
  run _app ${TEST_DIR}/asset/versioned/score-0-daemonset-v1.11.yml
  assert_zero_points
}

@test "passes 1.11 format statefulset" {
  run _app ${TEST_DIR}/asset/versioned/score-0-statefulset-v1.11.yml
  assert_zero_points
}

# ---

@test "fails DaemonSet with host docker.socket" {
  run _app ${TEST_DIR}/asset/score-0-daemonset-volume-host-docker-socket.yml
  assert_lt_zero_points
}

@test "passes Deployment with serviceaccountname" {
  run _app ${TEST_DIR}/asset/score-2-dep-serviceaccount.yml

  assert_gt_zero_points
}

@test "passes pod with serviceaccountname" {
  run _app ${TEST_DIR}/asset/score-2-pod-serviceaccount.yml

  assert_gt_zero_points
}

# TODO: test for pod-specific seccomp

# TODO: case sensitive check (use jq's ascii_downcase)

# TODO: tests for all the permutations of this file
@test "TODO: fails DaemonSet with loads o' permutations" {
  skip
  run _app ${TEST_DIR}/asset/score-0-daemonset-
  assert_zero_points
}

# ---

# TODO deprecated alpha feature
# https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#opaque-integer-resources-alpha-feature

# TODO: tests for all the permutations of this file
# tolerations?
# host mounts /dev /boot /proc /var/run/docker/socket
# env vars - certain names?
##
#@test "fails DaemonSet with loads o' permutations" {
#  skip
#  run _app ${TEST_DIR}/asset/score-0-daemonset-
#  assert_zero_points
#}


#@test "webhook: admission control 100" {
#  if _is_local; then
#    skip
#  fi
#
#  REMOTE_URL_STASH="${REMOTE_URL}"
#
#  REMOTE_URL+="?score=100"
#  run _app ./test/asset/score-0-cap-sys-admin.yml \
#    -w '%{content_type}' \
#    -o /dev/null
#
#  assert_output --regexp "^HTTP/1.1 401 Unauthorized$"
#  assert_lt_zero_points
#
#  REMOTE_URL="${REMOTE_URL_STASH}"
#}
#
#@test "webhook: admission control -100" {
#  if _is_local; then
#    skip
#  fi
#
#  REMOTE_URL_STASH="${REMOTE_URL}"
#
#  REMOTE_URL+="?score=100"
#  run _app ./test/asset/score-0-cap-sys-admin.yml \
#    -w '%{content_type}' \
#    -o /dev/null
#
#  assert_output --regexp "^HTTP/1.1 200"
#  assert_lt_zero_points
#
#  REMOTE_URL="${REMOTE_URL_STASH}"
#}
