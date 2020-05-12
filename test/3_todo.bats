#!/usr/bin/env bash

load '_helper'

setup() {
  _global_setup
}

teardown() {
  _global_teardown
}

# TODO(ajm) BEHAVIOURAL CHANGE (previous scan didn't account for all containers) - FIX BEFORE RELEASE
#@test "does not error on very long file" {
#  run _app ${TEST_DIR}/asset/very-long-file
#
#  assert_success
#}

# TODO(ajm) BEHAVIOURAL CHANGE (previous scan didn't account for all containers) - FIX BEFORE RELEASE
@test "passes production dump" {
  run _app ${TEST_DIR}/asset/score-1-prod-dump.yaml

  assert_lt_zero_points
}

# TODO(ajm) v2 fail - FIX BEFORE RELEASE
#@test "succeeds with valid file (json, local)" {
#  if _is_remote; then
#    skip
#  fi
#
#  run _app ${TEST_DIR}/asset/score-1-pod-default.yml --json
#  assert_output --regexp '  "score": [0-9]+.*'
#  assert_success
#}

# TODO(ajm) v2 new behaviour- FIX BEFORE RELEASE
#@test "returns content-type text/plain for failure" {
#  if _is_local; then
#    skip
#  fi
#
#   run _app \
#    arse \
#    -w '%{content_type}' \
#    -o /dev/null
#
#  assert_output --regexp "text/plain.*"
#}

# TODO: tests for apparmor loaders
@test "TODO: passes DaemonSet with apparmor loader" {
  skip
  https://github.com/kubernetes/contrib/blob/master/apparmor/loader/example-daemon.yaml
  run _app "${TEST_DIR}/asset/score-0-daemonset-"
  assert_zero_points
}

# TODO: acceptance test for pod-specific seccomp

# TODO: deployment serviceAccountName pass
@test "TODO: passes Deployment with serviceaccountname" {
  skip
  run _app "${TEST_DIR}/asset/score-2-dep-serviceaccount.yml"
  assert_zero_points
}

# TODO: tests for all the permutations of this file
@test "TODO: fails DaemonSet with loads o' permutations" {
  skip
  run _app "${TEST_DIR}/asset/score-0-daemonset-"
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
