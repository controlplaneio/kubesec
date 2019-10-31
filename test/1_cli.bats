#!/usr/bin/env bash

load '_helper'

setup() {
  _global_setup
}

teardown() {
  _global_teardown
}

@test "fails Pod with unconfined seccomp" {
  run _app ${TEST_DIR}/asset/score-0-pod-seccomp-unconfined.yml
  assert_lt_zero_points
}

@test "fails with CAP_SYS_ADMIN" {
  run _app ${TEST_DIR}/asset/score-0-cap-sys-admin.yml
  assert_lt_zero_points
}

@test "fails with CAP_CHOWN" {
  run _app ${TEST_DIR}/asset/score-0-cap-chown.yml
  assert_zero_points
}

@test "fails with CAP_SYS_ADMIN and CAP_CHOWN" {
  run _app ${TEST_DIR}/asset/score-0-cap-sys-admin-and-cap-chown.yml
  assert_lt_zero_points
}

@test "passes with securityContext capabilities drop all" {
  run _app ${TEST_DIR}/asset/score-1-cap-drop-all.yml
  assert_gt_zero_points
}

@test "passes deployment with securitycontext readOnlyRootFilesystem" {
  run _app ${TEST_DIR}/asset/score-1-dep-ro-root-fs.yml
  assert_gt_zero_points
}

@test "passes deployment with securitycontext runAsNonRoot" {
  run _app ${TEST_DIR}/asset/score-1-dep-seccon-run-as-non-root.yml
  assert_gt_zero_points
}

@test "fails deployment with securitycontext runAsUser 1" {
  run _app ${TEST_DIR}/asset/score-1-dep-seccon-run-as-user-1.yml
  assert_zero_points
}

@test "passes deployment with securitycontext runAsUser > 10000" {
  run _app ${TEST_DIR}/asset/score-1-dep-seccon-run-as-user-10001.yml
  assert_gt_zero_points
}

@test "fails deployment with empty security context" {
  run _app ${TEST_DIR}/asset/score-1-dep-empty-security-context.yml
  assert_zero_points
}

@test "passes deployment with cgroup resource limits" {
  run _app ${TEST_DIR}/asset/score-1-dep-resource-limit-cpu.yml
  assert_gt_zero_points
}

@test "passes deployment with cgroup memory limits" {
  run _app ${TEST_DIR}/asset/score-1-dep-resource-limit-memory.yml
  assert_gt_zero_points
}

@test "passes StatefulSet with volumeClaimTemplate" {
  run _app ${TEST_DIR}/asset/score-1-statefulset-volumeclaimtemplate.yml
  assert_gt_zero_points
}

@test "fails StatefulSet with no security" {
  run _app ${TEST_DIR}/asset/score-0-statefulset-no-sec.yml
  assert_zero_points
}

@test "fails DaemonSet with securityContext.privileged = true" {
  run _app ${TEST_DIR}/asset/score-0-daemonset-securitycontext-privileged.yml
  assert_lt_zero_points
}

@test "fails DaemonSet with mounted host docker.sock" {
  run _app ${TEST_DIR}/asset/score-0-daemonset-mount-docker-socket.yml
  assert_lt_zero_points
}

@test "passes Pod with apparmor annotation" {
  run _app ${TEST_DIR}/asset/score-3-pod-apparmor.yaml
  assert_gt_zero_points
}

@test "fails Pod with unconfined seccomp for all containers" {
  run _app ${TEST_DIR}/asset/score-0-pod-seccomp-unconfined.yml
  assert_lt_zero_points
}

@test "passes Pod with non-unconfined seccomp for all containers" {
  run _app ${TEST_DIR}/asset/score-0-pod-seccomp-non-unconfined.yml
  assert_gt_zero_points
}

@test "fails DaemonSet with hostNetwork" {
  run _app ${TEST_DIR}/asset/score-0-daemonset-host-network.yml
  assert_lt_zero_points
}

@test "fails DaemonSet with hostPid" {
  run _app ${TEST_DIR}/asset/score-0-daemonset-host-pid.yml
  assert_lt_zero_points
}

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

@test "fails deployment with allowPrivilegeEscalation" {
  run _app ${TEST_DIR}/asset/allowPrivilegeEscalation.yaml

  assert_lt_zero_points
}

@test "returns a unordered point score for specific response lines" {
  # NB response from use of parallel results in different permutations of rule order
  run \
    _app ${TEST_DIR}/asset/score-2-pod-serviceaccount.yml
  for LINE in 11 16 21 26 31 36 41 46 51 56 61
  do
    assert_line --index ${LINE} --regexp '^.*"points": [0-9]+$'
  done
}

@test "returns a ordered point score present" {
  # for #44 (later)
  skip

  run \
    _app ${TEST_DIR}/asset/score-2-pod-serviceaccount.yml
  assert_line --index 11 --regexp '^.*\"points\": 3$'
  for LINE in 16 21 26 31 36 41 46 51 56 61
  do 
    assert_line --index ${LINE} --regexp '^.*\"points\": 1$'
  done
}
