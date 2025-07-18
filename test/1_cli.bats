#!/usr/bin/env bash

load '_helper'

setup() {
  _global_setup
}

teardown() {
  _global_teardown
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

@test "passes deployment with pod securitycontext runAsNonRoot" {
  run _app "${TEST_DIR}/asset/score-1-dep-podseccon-run-as-non-root.yml"
  assert_gt_zero_points
}

@test "passes deployment with securitycontext runAsNonRoot" {
  run _app "${TEST_DIR}/asset/score-1-dep-seccon-run-as-non-root.yml"
  assert_gt_zero_points
}

@test "fails deployment with pod securitycontext runAsUser 1" {
  run _app "${TEST_DIR}/asset/score-1-dep-podseccon-run-as-user-1.yml"
  assert_zero_points
}

@test "fails deployment with securitycontext runAsUser 1" {
  run _app "${TEST_DIR}/asset/score-1-dep-seccon-run-as-user-1.yml"
  assert_zero_points
}

@test "passes deployment with pod securitycontext runAsUser > 10000" {
  run _app "${TEST_DIR}/asset/score-1-dep-podseccon-run-as-user-10001.yml"
  assert_gt_zero_points
}

@test "passes deployment with securitycontext runAsUser > 10000" {
  run _app "${TEST_DIR}/asset/score-1-dep-seccon-run-as-user-10001.yml"
  assert_gt_zero_points
}

@test "fails deployment with pod securitycontext runAsGroup 1" {
  run _app "${TEST_DIR}/asset/score-1-dep-podseccon-run-as-group-1.yml"
  assert_zero_points
}

@test "fails deployment with securitycontext runAsGroup 1" {
  run _app "${TEST_DIR}/asset/score-1-dep-seccon-run-as-group-1.yml"
  assert_zero_points
}

@test "passes deployment with pod securitycontext runAsGroup > 10000" {
  run _app "${TEST_DIR}/asset/score-1-dep-podseccon-run-as-group-10001.yml"
  assert_gt_zero_points
}

@test "passes deployment with securitycontext runAsGroup > 10000" {
  run _app "${TEST_DIR}/asset/score-1-dep-seccon-run-as-group-10001.yml"
  assert_gt_zero_points
}

@test "fails deployment with empty security context" {
  run _app "${TEST_DIR}/asset/score-1-dep-empty-security-context.yml"
  assert_zero_points
}

@test "fails deployment with invalid security context" {
  run _app "${TEST_DIR}/asset/score-1-dep-invalid-security-context.yml"

  run jq -r .[].message <<<"${output}"

  assert_output --partial "additional properties 'fake' not allowed"
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

@test "passes StatefulSet with no volumeClaimTemplate" {
  run _app "${TEST_DIR}/asset/score-1-statefulset-novolumeclaimtemplate.yml"
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

@test "fails Deployment with unconfined apparmor for all containers" {
  run _app "${TEST_DIR}/asset/score-0-dep-apparmor-empty-securitycontext.yml"
  assert_zero_points
}

@test "passes Deployment with non-unconfined apparmor on the .spec section" {
  run _app "${TEST_DIR}/asset/score-1-dep-apparmor-nonunconfined-spec-securitycontext.yml"
  assert_gt_zero_points
}

@test "fails Deployment with unconfined apparmor on the .spec section" {
  run _app "${TEST_DIR}/asset/score-0-dep-apparmor-unconfined-spec-securitycontext.yml"
  assert_lt_zero_points
}

@test "passes Deployment with non-unconfined apparmor on the .spec.containers section" {
  run _app "${TEST_DIR}/asset/score-1-dep-apparmor-nonunconfined-container.yml"
  assert_gt_zero_points
}

@test "fails Deployment with unconfined apparmor on the .spec.containers section" {
  run _app "${TEST_DIR}/asset/score-0-dep-apparmor-unconfined-container.yml"
  assert_lt_zero_points
}

@test "passes Deployment with non-unconfined apparmor on the .spec.initContainers section" {
  run _app "${TEST_DIR}/asset/score-1-dep-apparmor-nonunconfined-initcontainer.yml"
  assert_gt_zero_points
}

@test "fails Deployment with unconfined apparmor on the .spec.initcontainers section" {
  run _app "${TEST_DIR}/asset/score-0-dep-apparmor-unconfined-initcontainer.yml"
  assert_lt_zero_points
}

@test "passes Deployment with non-unconfined apparmor on the .spec.ephemeralContainers section" {
  run _app "${TEST_DIR}/asset/score-1-dep-apparmor-nonunconfined-ephemeralcontainer.yml"
  assert_gt_zero_points
}

@test "fails Deployment with unconfined apparmor on the .spec.ephemeralContainers section" {
  run _app "${TEST_DIR}/asset/score-0-dep-apparmor-unconfined-ephemeralcontainer.yml"
  assert_lt_zero_points
}

@test "fails Deployment with unconfined seccomp for all containers" {
  run _app "${TEST_DIR}/asset/score-0-dep-seccomp-empty-securitycontext.yml"
  assert_zero_points
}

@test "passes Deployment with non-unconfined seccomp on the .spec section" {
  run _app "${TEST_DIR}/asset/score-1-dep-seccomp-nonunconfined-spec-securitycontext.yml"
  assert_gt_zero_points
}

@test "fails Deployment with unconfined seccomp on the .spec section" {
  run _app "${TEST_DIR}/asset/score-0-dep-seccomp-unconfined-spec-securitycontext.yml"
  assert_lt_zero_points
}

@test "passes Deployment with non-unconfined seccomp on the .spec.containers section" {
  run _app "${TEST_DIR}/asset/score-1-dep-seccomp-nonunconfined-container.yml"
  assert_gt_zero_points
}

@test "fails Deployment with unconfined seccomp on the .spec.containers section" {
  run _app "${TEST_DIR}/asset/score-0-dep-seccomp-unconfined-container.yml"
  assert_lt_zero_points
}

@test "passes Deployment with non-unconfined seccomp on the .spec.initContainers section" {
  run _app "${TEST_DIR}/asset/score-1-dep-seccomp-nonunconfined-initcontainer.yml"
  assert_gt_zero_points
}

@test "fails Deployment with unconfined seccomp on the .spec.initcontainers section" {
  run _app "${TEST_DIR}/asset/score-0-dep-seccomp-unconfined-initcontainer.yml"
  assert_lt_zero_points
}

@test "passes Deployment with non-unconfined seccomp on the .spec.ephemeralContainers section" {
  run _app "${TEST_DIR}/asset/score-1-dep-seccomp-nonunconfined-ephemeralcontainer.yml"
  assert_gt_zero_points
}

@test "fails Deployment with unconfined seccomp on the .spec.ephemeralContainers section" {
  run _app "${TEST_DIR}/asset/score-0-dep-seccomp-unconfined-ephemeralcontainer.yml"
  assert_lt_zero_points
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

@test "passes Pod with automountServiceAccountToken set to false" {
  run _app "${TEST_DIR}/asset/score-1-pod-automount-sa-set-to-false.yml"
  assert_gt_zero_points
}

@test "passes Pod with hostUsers set to false" {
  run _app "${TEST_DIR}/asset/score-1-pod-hostUsers-set-to-false.yml"
  assert_gt_zero_points
}


@test "returns integer point score for each advice element" {
  run _app "${TEST_DIR}/asset/score-2-pod-serviceaccount.yml"
  assert_success

  run jq -r .[].scoring.advise[].points <<<"${output}"

  for SCORE in ${output}; do
    assert bash -c "[[ ${SCORE} =~ ^[0-9]+$ ]]"
  done
}

@test "returns an ordered point score for all advice" {
  run _app "${TEST_DIR}/asset/score-2-pod-serviceaccount.yml"
  assert_success

  run jq -r .[].scoring.advise[].points <<<"${output}"

  PREVIOUS=""
  for CURRENT in ${output}; do
    [ "${PREVIOUS}" = "" ] || assert [ "${CURRENT}" -le "${PREVIOUS}" ]
    PREVIOUS="${CURRENT}"
  done
}

@test "returns integer point score for each pass element" {
  run _app "${TEST_DIR}/asset/score-5-pod-serviceaccount.yml"
  assert_success

  run jq -r .[].scoring.passed[].points <<<"${output}"

  for SCORE in ${output}; do
    assert bash -c "[[ ${SCORE} =~ ^[0-9]+$ ]]"
  done
}

@test "returns an ordered point score for all passed" {
  run _app "${TEST_DIR}/asset/score-5-pod-serviceaccount.yml"

  run jq -r .[].scoring.passed[].points <<<"${output}"

  PREVIOUS=""
  for CURRENT in ${output}; do
    [ "${PREVIOUS}" = "" ] || assert [ "${CURRENT}" -le "${PREVIOUS}" ]
    PREVIOUS="${CURRENT}"
  done
}

@test "check critical and advisory points listed by magnitude" {
  run _app "${TEST_DIR}/asset/critical-double.yml"

  # criticals - magnitude sort/lowest number first
  CRITICAL_FIRST=$(jq -r .[].scoring.critical[0].points <<<"${output}")
  CRITICAL_SECOND=$(jq -r .[].scoring.critical[1].points <<<"${output}")
  (( CRITICAL_FIRST <= CRITICAL_SECOND ))

  # advisories - magnitude sort/highest number first
  ADVISE_FIRST=$(jq -r .[].scoring.advise[0].points <<<"${output}")
  ADVISE_SECOND=$(jq -r .[].scoring.advise[1].points <<<"${output}")
  ADVISE_THIRD=$(jq -r .[].scoring.advise[2].points <<<"${output}")
  (( ADVISE_FIRST >= ADVISE_SECOND >= ADVISE_THIRD ))
}

@test "check critical and advisory points as multi-yaml" {
  run _app "${TEST_DIR}/asset/critical-double-multiple.yml"

  # report 1 - criticals - magnitude sort/lowest number first
  CRITICAL_FIRST_FIRST=$(jq -r .[0].scoring.critical[0].points <<<"${output}")
  CRITICAL_FIRST_SECOND=$(jq -r .[0].scoring.critical[1].points <<<"${output}")
  (( CRITICAL_FIRST_FIRST <= CRITICAL_FIRST_SECOND ))

  # report 1 - advisories - magnitude sort/highest number first
  ADVISE_FIRST_FIRST=$(jq -r .[0].scoring.advise[0].points <<<"${output}")
  ADVISE_FIRST_SECOND=$(jq -r .[0].scoring.advise[1].points <<<"${output}")
  ADVISE_FIRST_THIRD=$(jq -r .[0].scoring.advise[2].points <<<"${output}")
  (( ADVISE_FIRST_FIRST >= ADVISE_FIRST_SECOND >= ADVISE_FIRST_THIRD ))

  # report 2 - criticals - magnitude sort/lowest number first
  CRITICAL_SECOND_FIRST=$(jq -r .[1].scoring.critical[0].points <<<"${output}")
  CRITICAL_SECOND_SECOND=$(jq -r .[1].scoring.critical[1].points <<<"${output}")
  (( CRITICAL_SECOND_FIRST <= CRITICAL_SECOND_SECOND ))

  # report 2 - advisories - magnitude sort/highest number first
  ADVISE_SECOND_FIRST=$(jq -r .[1].scoring.advise[0].points <<<"${output}")
  ADVISE_SECOND_SECOND=$(jq -r .[1].scoring.advise[1].points <<<"${output}")
  ADVISE_SECOND_THIRD=$(jq -r .[1].scoring.advise[2].points <<<"${output}")
  (( ADVISE_SECOND_FIRST >= ADVISE_SECOND_SECOND >= ADVISE_SECOND_THIRD ))
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
