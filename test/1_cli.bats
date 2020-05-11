#!/usr/bin/env bash

load '_helper'

setup() {
  _global_setup
}

teardown() {
  _global_teardown
}

@test "fails Pod with unconfined seccomp" {
  run _app "${TEST_DIR}/asset/score-0-pod-seccomp-unconfined.yml"
  assert_lt_zero_points
}

@test "fails with CAP_SYS_ADMIN" {
  run _app "${TEST_DIR}/asset/score-0-cap-sys-admin.yml"
  assert_lt_zero_points
}

@test "fails with CAP_CHOWN" {
  run _app "${TEST_DIR}/asset/score-0-cap-chown.yml"
  assert_zero_points
}

@test "fails with CAP_SYS_ADMIN and CAP_CHOWN" {
  run _app "${TEST_DIR}/asset/score-0-cap-sys-admin-and-cap-chown.yml"
  assert_lt_zero_points
}

@test "passes with securityContext capabilities drop all" {
  run _app "${TEST_DIR}/asset/score-1-cap-drop-all.yml"
  assert_gt_zero_points
}

@test "passes deployment with securitycontext readOnlyRootFilesystem" {
  run _app "${TEST_DIR}/asset/score-1-dep-ro-root-fs.yml"
  assert_gt_zero_points
}

@test "passes deployment with securitycontext runAsNonRoot" {
  run _app "${TEST_DIR}/asset/score-1-dep-seccon-run-as-non-root.yml"
  assert_gt_zero_points
}

@test "fails deployment with securitycontext runAsUser 1" {
  run _app "${TEST_DIR}/asset/score-1-dep-seccon-run-as-user-1.yml"
  assert_zero_points
}

@test "passes deployment with securitycontext runAsUser > 10000" {
  run _app "${TEST_DIR}/asset/score-1-dep-seccon-run-as-user-10001.yml"
  assert_gt_zero_points
}

@test "fails deployment with empty security context" {
  run _app "${TEST_DIR}/asset/score-1-dep-empty-security-context.yml"
  assert_zero_points
}

@test "fails deployment with invalid security context" {
  run _app "${TEST_DIR}/asset/score-1-dep-invalid-security-context.yml"
  assert_line --index 4 --regexp 'fake: Additional property fake is not allowed'
}

@test "passes deployment with cgroup resource limits" {
  run _app "${TEST_DIR}/asset/score-1-dep-resource-limit-cpu.yml"
  assert_gt_zero_points
}

@test "passes deployment with cgroup memory limits" {
  run _app "${TEST_DIR}/asset/score-1-dep-resource-limit-memory.yml"
  assert_gt_zero_points
}

@test "passes StatefulSet with volumeClaimTemplate" {
  run _app "${TEST_DIR}/asset/score-1-statefulset-volumeclaimtemplate.yml"
  assert_gt_zero_points
}

@test "fails StatefulSet with no security" {
  run _app "${TEST_DIR}/asset/score-0-statefulset-no-sec.yml"
  assert_zero_points
}

@test "fails DaemonSet with securityContext.privileged = true" {
  run _app "${TEST_DIR}/asset/score-0-daemonset-securitycontext-privileged.yml"
  assert_lt_zero_points
}

@test "fails DaemonSet with mounted host docker.sock" {
  run _app "${TEST_DIR}/asset/score-0-daemonset-mount-docker-socket.yml"
  assert_lt_zero_points
}

@test "passes Pod with apparmor annotation" {
  run _app "${TEST_DIR}/asset/score-3-pod-apparmor.yaml"
  assert_gt_zero_points
}

@test "fails Pod with unconfined seccomp for all containers" {
  run _app "${TEST_DIR}/asset/score-0-pod-seccomp-unconfined.yml"
  assert_lt_zero_points
}

@test "passes Pod with non-unconfined seccomp for all containers" {
  run _app "${TEST_DIR}/asset/score-0-pod-seccomp-non-unconfined.yml"
  assert_gt_zero_points
}

@test "fails DaemonSet with hostNetwork" {
  run _app "${TEST_DIR}/asset/score-0-daemonset-host-network.yml"
  assert_lt_zero_points
}

@test "fails DaemonSet with hostPid" {
  run _app "${TEST_DIR}/asset/score-0-daemonset-host-pid.yml"
  assert_lt_zero_points
}

@test "fails DaemonSet with host docker.socket" {
  run _app "${TEST_DIR}/asset/score-0-daemonset-volume-host-docker-socket.yml"
  assert_lt_zero_points
}

@test "passes Deployment with serviceaccountname" {
  run _app "${TEST_DIR}/asset/score-2-dep-serviceaccount.yml"
  assert_gt_zero_points
}

@test "passes pod with serviceaccountname" {
  run _app "${TEST_DIR}/asset/score-2-pod-serviceaccount.yml"
  assert_gt_zero_points
}

@test "fails deployment with allowPrivilegeEscalation" {
  run _app "${TEST_DIR}/asset/allowPrivilegeEscalation.yaml"
  assert_lt_zero_points
}

@test "returns integer point score for specific response lines" {
  run _app "${TEST_DIR}/asset/score-2-pod-serviceaccount.yml"

  for LINE in 11 16 21 26 31 36 41 46 51 56 61; do
    assert_line --index ${LINE} --regexp '^.*"points": [0-9]+$'
  done
}

@test "returns an ordered point score for all responses" {
  run _app "${TEST_DIR}/asset/score-2-pod-serviceaccount.yml"

  assert_line --index 11 --regexp '^.*\"points\": 3$'

  for LINE in 16 21 26 31 36 41 46 51 56 61; do
    assert_line --index ${LINE} --regexp '^.*\"points\": 1$'
  done
}

@test "check critical and advisory points listed by magnitude" {
  run _app "${TEST_DIR}/asset/critical-double.yml"

  # criticals - magnitude sort/lowest number first
  assert_line --index 11 --regexp '^.*\"points\": -30$'
  assert_line --index 16 --regexp '^.*\"points\": -7$'

  # advisories - magnitude sort/highest number first
  assert_line --index 23 --regexp '^.*\"points\": 3$'
  assert_line --index 28 --regexp '^.*\"points\": 3$'
  assert_line --index 33 --regexp '^.*\"points\": 1$'
}

@test "check critical and advisory points as multi-yaml" {
  run _app "${TEST_DIR}/asset/critical-double-multiple.yml"

  # report 1 - criticals - magnitude sort/lowest number first
  assert_line --index 11 --regexp '^.*\"points\": -30$'
  assert_line --index 16 --regexp '^.*\"points\": -7$'

  # report 1 - advisories - magnitude sort/highest number first
  assert_line --index 23 --regexp '^.*\"points\": 3$'
  assert_line --index 28 --regexp '^.*\"points\": 3$'
  assert_line --index 33 --regexp '^.*\"points\": 1$'

  # report 2 - criticals - magnitude sort/lowest number first
  assert_line --index 93 --regexp '^.*\"points\": -30$'
  assert_line --index 98 --regexp '^.*\"points\": -7$'

  # report 2 - advisories - magnitude sort/highest number first
  assert_line --index 105 --regexp '^.*\"points\": 3$'
  assert_line --index 110 --regexp '^.*\"points\": 3$'
  assert_line --index 115 --regexp '^.*\"points\": 1$'
}

@test "returns deterministic report output" {
  run _app "${TEST_DIR}/asset/score-2-pod-serviceaccount.yml"
  assert_success

  RUN_1_SIGNATURE=$(echo "${output}" | sha1sum)

  run _app "${TEST_DIR}/asset/score-2-pod-serviceaccount.yml"
  assert_success

  RUN_2_SIGNATURE=$(echo "${output}" | sha1sum)

  run _app "${TEST_DIR}/asset/score-2-pod-serviceaccount.yml"
  assert_success

  RUN_3_SIGNATURE=$(echo "${output}" | sha1sum)

  assert [ "${RUN_1_SIGNATURE}" = "${RUN_2_SIGNATURE}" ]
  assert [ "${RUN_1_SIGNATURE}" = "${RUN_3_SIGNATURE}" ]
}
