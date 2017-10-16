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
  run "${APP}"
  assert_failure
}

@test "errors with invalid file" {
  run ${APP} somefile.yaml
  assert_output --regexp ".*File somefile.yaml does not exist.*"
  assert_failure
}

# ---

@test "only valid types - deny PodSecurityPolicy" {
  run ${APP} ${TEST_DIR}/asset/score-0-podsecuritypolicy-permissive.yml
  assert_output --regexp ".*Only kinds .* accepted.*"
  assert_failure
}

@test "only valid types - allow Pod" {
  run ${APP} ${TEST_DIR}/asset/score-1-pod-default.yml
  refute_output --regexp ".*Only kinds .* accepted.*"
  assert_success
}

@test "only valid types - allow Deployment" {
  run ${APP} ${TEST_DIR}/asset/score-1-dep-default.yml
  refute_output --regexp ".*Only kinds .* accepted.*"
  assert_success
}

@test "only valid types - allow StatefulSet" {
  run ${APP} ${TEST_DIR}/asset/score-1-statefulset-default.yml
  refute_output --regexp ".*Only kinds .* accepted.*"
  assert_success
}

@test "only valid types - allow DaemonSet" {
  run ${APP} ${TEST_DIR}/asset/score-1-daemonset-default.yml
  refute_output --regexp ".*Only kinds .* accepted.*"
  assert_success
}

# ---

@test "fails with CAP_SYS_ADMIN" {
  run ${APP} ${TEST_DIR}/asset/score-0-cap-sys-admin.yml
  assert_zero_points
}

@test "fails with CAP_CHOWN" {
  run ${APP} ${TEST_DIR}/asset/score-0-cap-chown.yml
  assert_zero_points
}

@test "fails with CAP_SYS_ADMIN and CAP_CHOWN" {
  run ${APP} ${TEST_DIR}/asset/score-0-cap-sys-admin-and-cap-chown.yml
  assert_zero_points
}

@test "passes with securityContext capabilities drop all" {
  run ${APP} ${TEST_DIR}/asset/score-1-cap-drop-all.yml
  assert_non_zero_points
}

@test "passes deployment with securitycontext readOnlyRootFilesystem" {
  run ${APP} ${TEST_DIR}/asset/score-1-dep-ro-root-fs.yml
  assert_non_zero_points
}

@test "passes deployment with securitycontext runAsNonRoot" {
  run ${APP} ${TEST_DIR}/asset/score-1-dep-seccon-run-as-non-root.yml
  assert_non_zero_points
}

@test "fails deployment with securitycontext runAsUser 1" {
  run ${APP} ${TEST_DIR}/asset/score-1-dep-seccon-run-as-user-1.yml
  assert_zero_points
}

@test "passes deployment with securitycontext runAsUser > 10000" {
  run ${APP} ${TEST_DIR}/asset/score-1-dep-seccon-run-as-user-10001.yml
  assert_non_zero_points
}

@test "fails deployment with empty security context" {
  run ${APP} ${TEST_DIR}/asset/score-1-dep-empty-security-context.yml
  assert_zero_points
}

@test "passes deployment with cgroup resource limits" {
  run ${APP} ${TEST_DIR}/asset/score-1-dep-resource-limit-cpu.yml
  assert_non_zero_points
}

@test "passes deployment with cgroup memory limits" {
  run ${APP} ${TEST_DIR}/asset/score-1-dep-resource-limit-memory.yml
  assert_non_zero_points
}

@test "passes StatefulSet with volumeClaimTemplate" {
  run ${APP} ${TEST_DIR}/asset/score-1-statefulset-volumeclaimtemplate.yml
  assert_non_zero_points
}

@test "fails StatefulSet with no security" {
  run ${APP} ${TEST_DIR}/asset/score-0-statefulset-no-sec.yml
  assert_zero_points
}

@test "fails DaemonSet with securityContext.privileged = true" {
  run ${APP} ${TEST_DIR}/asset/score-0-daemonset-securitycontext-privileged.yml
  assert_zero_points
}

# TODO: tests for all the permutations of this file
@test "fails DaemonSet with loads o' permutations" {
  skip
  run ${APP} ${TEST_DIR}/asset/score-0-daemonset-securitycontext-privileged.yml
  assert_zero_points
}

