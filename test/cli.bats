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

@test "errors with no filename" {
  run _app
  assert_failure
}

@test "errors with invalid file" {
  run _app somefile.yaml
  assert_file_not_found
  assert_failure
}

@test "accepts --json flag locally" {
  if _is_remote; then
    skip
  fi

  run _app --json ${TEST_DIR}/asset/score-1-pod-default.yml
  assert_output --regexp '  "score": [0-9]+.*'
  assert_success
}

# ---

@test "only valid types - deny PodSecurityPolicy" {
  run _app ${TEST_DIR}/asset/score-0-podsecuritypolicy-permissive.yml
  assert_output --regexp ".*Only kinds .* accepted.*"
  if _is_local; then
    assert_failure
  fi
}

@test "only valid types - allow Pod" {
  run _app ${TEST_DIR}/asset/score-1-pod-default.yml
  refute_output --regexp ".*Only kinds .* accepted.*"
  assert_success
}

@test "only valid types - allow Deployment" {
  run _app ${TEST_DIR}/asset/score-1-dep-default.yml
  refute_output --regexp ".*Only kinds .* accepted.*"
  assert_success
}

@test "only valid types - allow StatefulSet" {
  run _app ${TEST_DIR}/asset/score-1-statefulset-default.yml
  refute_output --regexp ".*Only kinds .* accepted.*"
  assert_success
}

@test "only valid types - allow DaemonSet" {
  run _app ${TEST_DIR}/asset/score-1-daemonset-default.yml
  refute_output --regexp ".*Only kinds .* accepted.*"
  assert_success
}

# ---

@test "fails with CAP_SYS_ADMIN" {
  run _app ${TEST_DIR}/asset/score-0-cap-sys-admin.yml
  assert_zero_points
}

@test "fails with CAP_CHOWN" {
  run _app ${TEST_DIR}/asset/score-0-cap-chown.yml
  assert_zero_points
}

@test "fails with CAP_SYS_ADMIN and CAP_CHOWN" {
  run _app ${TEST_DIR}/asset/score-0-cap-sys-admin-and-cap-chown.yml
  assert_zero_points
}

@test "passes with securityContext capabilities drop all" {
  run _app ${TEST_DIR}/asset/score-1-cap-drop-all.yml
  assert_non_zero_points
}

@test "passes deployment with securitycontext readOnlyRootFilesystem" {
  run _app ${TEST_DIR}/asset/score-1-dep-ro-root-fs.yml
  assert_non_zero_points
}

@test "passes deployment with securitycontext runAsNonRoot" {
  run _app ${TEST_DIR}/asset/score-1-dep-seccon-run-as-non-root.yml
  assert_non_zero_points
}

@test "fails deployment with securitycontext runAsUser 1" {
  run _app ${TEST_DIR}/asset/score-1-dep-seccon-run-as-user-1.yml
  assert_zero_points
}

@test "passes deployment with securitycontext runAsUser > 10000" {
  run _app ${TEST_DIR}/asset/score-1-dep-seccon-run-as-user-10001.yml
  assert_non_zero_points
}

@test "fails deployment with empty security context" {
  run _app ${TEST_DIR}/asset/score-1-dep-empty-security-context.yml
  assert_zero_points
}

@test "passes deployment with cgroup resource limits" {
  run _app ${TEST_DIR}/asset/score-1-dep-resource-limit-cpu.yml
  assert_non_zero_points
}

@test "passes deployment with cgroup memory limits" {
  run _app ${TEST_DIR}/asset/score-1-dep-resource-limit-memory.yml
  assert_non_zero_points
}

@test "passes StatefulSet with volumeClaimTemplate" {
  run _app ${TEST_DIR}/asset/score-1-statefulset-volumeclaimtemplate.yml
  assert_non_zero_points
}

@test "fails StatefulSet with no security" {
  run _app ${TEST_DIR}/asset/score-0-statefulset-no-sec.yml
  assert_zero_points
}

@test "fails DaemonSet with securityContext.privileged = true" {
  run _app ${TEST_DIR}/asset/score-0-daemonset-securitycontext-privileged.yml
  assert_zero_points
}

@test "fails DaemonSet with mounted host docker.sock" {
  skip
  run _app ${TEST_DIR}/asset/score-0-daemonset-mount-docker-socket.yml
  assert_zero_points
}




# TODO: tests for all the permutations of this file
@test "fails DaemonSet with loads o' permutations" {
  skip
  run _app ${TEST_DIR}/asset/score-0-daemonset-
  assert_zero_points
}

