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

@test "does not error on very long file" {
  run _app ${TEST_DIR}/asset/very-long-file

  assert_success
}

@test "returns error for invalid JSON" {
  run _app ${TEST_DIR}/asset/invalid-input-pod-dump.json

  assert_output --regexp "'api_version': invalid key, expected 'apiVersion'"
  assert_failure_local
}

@test "returns error YAML control characters" {
  run _app ${TEST_DIR}/asset/invalid-input-no-control-characters.json

  assert_failure
}

@test "passes production dump" {
  run _app ${TEST_DIR}/asset/score-1-prod-dump.yaml

  assert_non_zero_points
}

@test "passes bug dump twice [1/2]" {
  run _app ${TEST_DIR}/asset/bug-dump-2.json
  assert_success
  assert_non_zero_points
}

@test "passes bug dump twice [2/2]" {
  run _app ${TEST_DIR}/asset/bug-dump-2.json
  assert_success
  assert_non_zero_points
}

@test "returns content-type application/json" {
  if _is_local; then
    skip
  fi

   run _app ${TEST_DIR}/asset/score-0-daemonset-volume-host-docker-socket.yml -w '%{content_type}' -o /dev/null

  assert_output --regexp "application/json"
}

@test "fails DaemonSet with host docker.socket" {
  skip
  run _app ${TEST_DIR}/asset/score-0-daemonset-volume-host-docker-socket.yml
  assert_negative_points
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

# TODO: REQURIRES 1.8 deployment serviceAccountName pass
@test "passes Deployment with serviceaccountname" {
  skip
  run _app ${TEST_DIR}/asset/score-2-dep-serviceaccount.yml
  assert_non_zero_points
}


# TODO deprecated alpha feature
# https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#opaque-integer-resources-alpha-feature
