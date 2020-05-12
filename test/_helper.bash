#!/usr/bin/env bash

load './bin/bats-support/load'
load './bin/bats-assert/load'

export TEST_DIR="."

BIN_UNDER_TEST='./dist/kubesec scan'

_global_setup() {
    [ ! -f ${BATS_PARENT_TMPNAME}.skip ] || skip "skip remaining tests"
}

_global_teardown() {
    if [ ! -n "$BATS_TEST_COMPLETED" ]; then
      touch ${BATS_PARENT_TMPNAME}.skip
    fi
}

_get_remote_url() {
  echo "${REMOTE_URL:-https://v2.kubesec.io/scan}"
}

_is_local() {
  [[ "${REMOTE_URL:-}" == "" ]]
}

_is_remote() {
  ! _is_local
}

_test_description_matches_regex() {
  [[ "${BATS_TEST_DESCRIPTION}" =~ ${1} ]]
}

skip_if_not_local() {
  assert _test_description_matches_regex "[ |(]local[)|,]"
  _is_local || skip
}

skip_if_not_remote() {
  assert _test_description_matches_regex "[ |(]remote[)|,]"
  _is_remote || skip
}

if _is_remote; then

  _app() {
    local FILE="${1:-}"
    shift
    local ARGS="${@:-}"
    ARGS=$(echo "${ARGS}" | sed 's,--json,,g')
    # echo \# curl -sSX POST --data-binary @"${FILE}" ${ARGS} "$(_get_remote_url)" >&3
    curl -sSX POST --data-binary @"${FILE}" ${ARGS} "$(_get_remote_url)"
  }

  assert_gt_zero_points() {
    assert_output --regexp ".*\"score\": [1-9]+,.*"
    assert_success
  }

  assert_zero_points() {
    assert_output --regexp ".*\"score\": 0,.*"
    assert_success
  }

  assert_lt_zero_points() {
    assert_output --regexp ".*\"score\": \-[1-9][0-9]*,.*"
    assert_success
  }

  assert_file_not_found() {
    assert_output --regexp ".*couldn't open file \"somefile.yaml\".*" \
     || assert_output --regexp ".*no such file or directory.*" \
     || assert_output --regexp ".*Invalid input.*"
  }

  assert_invalid_input() {
    assert_output --regexp '  "message": "Invalid input"' \
     || assert_output --regexp ".*Invalid input.*" \
     || assert_output --regexp ".*no such file or directory.*" \
     || assert_output --regexp ".*Missing 'kind' key.*"
  }

  assert_failure_local() { :; }

else

  _app() {
    local ARGS="${@:-}"
    if [[ "${BIN_UNDER_TEST}" != "" ]]; then
      # remove --json flags
      ARGS=$(echo "${ARGS}" | sed -E 's,--json,,g')
    fi
    ./../${BIN_UNDER_TEST:-false} "${ARGS}";
  }

  assert_gt_zero_points() {
    assert_output --regexp ".*with a score of [1-9]+ points.*"
    assert_success
  }

  assert_zero_points() {
    assert_output --regexp ".*with a score of 0 points.*"
    assert_failure
  }

  assert_lt_zero_points() {
    assert_output --regexp ".*\with a score of \-[1-9][0-9]* points.*"
    assert_failure
  }

  assert_file_not_found() {
    assert_output --regexp ".*File somefile.yaml does not exist.*" \
     || assert_output --regexp ".*no such file or directory.*"  \
     || assert_output --regexp ".*Invalid input.*"
  }

  assert_invalid_input() {
    assert_output --regexp '  "message": "Invalid input"' \
     || assert_output --regexp ".*Kubernetes kind not found.*" \
     || assert_output --regexp ".*no such file or directory.*" \
     || assert_output --regexp ".*Invalid input.*"
  }

  assert_failure_local() {
    assert_failure
  }

fi
