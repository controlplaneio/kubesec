#!/usr/bin/env bash

load './bin/bats-support/load'
load './bin/bats-assert/load'

TEST_DIR="."

_global_setup() {
    [ ! -f ${BATS_PARENT_TMPNAME}.skip ] || skip "skip remaining tests"
}

_global_teardown() {
    if [ ! -n "$BATS_TEST_COMPLETED" ]; then
      touch ${BATS_PARENT_TMPNAME}.skip
    fi
}

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
    local ARGS="${@:-}"
    ARGS=$(echo "${ARGS}" | sed 's,--json,,g')
    curl --fail -v "$(_get_remote_url)" -F file=@"${FILE}" "${ARGS}"
  }

  assert_gt_zero_points() {
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
    assert_output --regexp ".*couldn't open file \"somefile.yaml\".*" \
     || assert_output --regexp ".*somefile.yaml: no such file or directory.*"
  }

  assert_invalid_input() {
    assert_output --regexp '  "message": "Invalid input"' \
     || assert_output --regexp ".*Invalid input.*" \
     || assert_output --regexp ".*Kubernetes kind not found.*"
  }

  assert_failure_local() { :; }

else

  _app() {
    local ARGS="${@:-}"
    if [[ "${BIN_UNDER_TEST}" != "" ]]; then
      # remove --json flag for golang rewrite
      ARGS=$(echo "${ARGS}" | sed -E 's,--json,,g')
    fi
    echo "# DEBUG: ARGS ${ARGS}" >&3
    ./../${BIN_UNDER_TEST:-kseccheck.sh} "${ARGS}";
  }

  assert_gt_zero_points() {
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
    assert_output --regexp ".*File somefile.yaml does not exist.*" \
     || assert_output --regexp ".*somefile.yaml: no such file or directory.*"  \
     || assert_output --regexp ".*Invalid input.*"
  }

  assert_invalid_input() {
    assert_output --regexp '  "message": "Invalid input"' \
     || assert_output --regexp ".*Kubernetes kind not found.*" \
     || assert_output --regexp ".*Invalid input.*"
  }

  assert_failure_local() {
    assert_failure
  }

fi
