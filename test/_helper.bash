#!/usr/bin/env bash

load './bin/bats-support/load'
load './bin/bats-assert/load'

TEST_DIR="."

_is_local() {
  [[ "${TEST_REMOTE:-0}" != 1 ]]
}
_is_remote() {
  ! _is_local
}
_get_remote_url() {
  echo "${REMOTE_URL:-https://kubesec.io/}"
}
if _is_remote; then

  _app() {
    local FILE="${1:-}"
    shift
    curl --fail -v "$(_get_remote_url)" -F file=@"${FILE}" "${@}"
  }

  assert_non_zero_points() {
    assert_output --regexp ".*\"score\": [1-9]+,.*"
    assert_success
  }

  assert_zero_points() {
    assert_output --regexp ".*\"score\": (0|\-[1-9][0-9]*),.*"
    assert_success
  }

  assert_negative_points() {
    assert_output --regexp ".*\"score\": (\-[1-9][0-9]*),.*"
    assert_success
  }

  assert_file_not_found() {
    assert_output --regexp ".*couldn't open file \"somefile.yaml\".*"
  }

  assert_failure_local() { :; }

else

  _app() { ./../kseccheck.sh "${@}"; }

  assert_non_zero_points() {
    assert_output --regexp ".*with a score of [1-9]+ points.*"
    assert_success
  }

  assert_zero_points() {
    assert_output --regexp ".*with a score of 0 points.*"
    assert_failure
  }

  assert_negative_points() {
    assert_output --regexp ".*\with a score of \-[1-9][0-9]* points.*"
    assert_failure
  }

  assert_file_not_found() {
    assert_output --regexp ".*File somefile.yaml does not exist.*"
  }

  assert_failure_local() {
    assert_failure
  }

fi
