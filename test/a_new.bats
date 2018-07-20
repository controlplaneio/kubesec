#!/usr/bin/env bash

load '_helper'

# assert
# refute
# assert_equal
# assert_success
# assert_failure
# assert_output
# refute_output
# assert_line
# refute_line

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
#  assert_negative_points
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
#  assert_negative_points
#
#  REMOTE_URL="${REMOTE_URL_STASH}"
#}

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
  assert_negative_points
}

@test "passes Deployment with serviceaccountname" {
  run _app ${TEST_DIR}/asset/score-2-dep-serviceaccount.yml

  assert_non_zero_points
}

@test "passes pod with serviceaccountname" {
  run _app ${TEST_DIR}/asset/score-2-pod-serviceaccount.yml

  assert_non_zero_points
}



# TODO: tests for apparmor loaders
@test "passes DaemonSet with apparmor loader" {
  skip
  https://github.com/kubernetes/contrib/blob/master/apparmor/loader/example-daemon.yaml
  run _app ${TEST_DIR}/asset/score-0-daemonset-
  assert_zero_points
}

# TODO: test for pod-specific seccomp

# TODO: case sensitive check (use jq's ascii_downcase)

# TODO: tests for all the permutations of this file
@test "fails DaemonSet with loads o' permutations" {
  skip
  run _app ${TEST_DIR}/asset/score-0-daemonset-
  assert_zero_points
}



# ---


# TODO deprecated alpha feature
# https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#opaque-integer-resources-alpha-feature
