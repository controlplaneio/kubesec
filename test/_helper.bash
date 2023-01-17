#!/usr/bin/env bash

load './bin/bats-support/load'
load './bin/bats-assert/load'

export TEST_DIR="."

export BIN_DIR='../dist/'

export SUBCOMMAND="${SUBCOMMAND:-scan}"

_global_setup() {
  [ ! -f "${BATS_PARENT_TMPNAME}".skip ] || skip "skip remaining tests"
}

_global_teardown() {
  if [ -z "$BATS_TEST_COMPLETED" ]; then
    touch "${BATS_PARENT_TMPNAME}".skip
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
    local ARGS=("$@")
    # remove --json flags
    ARGS=("${@/--json/}")
    curl -sSX POST --data-binary @"${FILE}" "${ARGS[@]}" "$(_get_remote_url)"
  }

  assert_gt_zero_points() {
    SCORE=$(jq -r .[].score <<<"${output:-}")
    ((SCORE > 0))
    assert_success
  }

  assert_zero_points() {
    SCORE=$(jq -r .[].score <<<"${output:-}")
    ((SCORE == 0))
    assert_success
  }

  assert_lt_zero_points() {
    SCORE=$(jq -r .[].score <<<"${output:-}")
    ((SCORE < 0))
    assert_success
  }

  assert_file_not_found() {
    assert_output --regexp ".*couldn't open file \"somefile.yaml\".*" ||
      assert_output --regexp ".*no such file or directory.*" ||
      assert_output --regexp ".*Invalid input.*"
  }

  assert_invalid_input() {
    assert_output --regexp '  "message": "Invalid input"' ||
      assert_output --regexp ".*Invalid input.*" ||
      assert_output --regexp ".*no such file or directory.*" ||
      assert_output --regexp ".*error while parsing.*" ||
      assert_output --regexp ".*[mM]issing 'kind' key.*"
  }

  assert_failure_local() { :; }

else

  _app() {
    local ARGS=("$@")
    if [[ "${BIN_DIR}" != "" ]]; then
      # remove --json flags
      ARGS=("${@/--json/}")
    fi
    "${BIN_DIR}"/kubesec "${SUBCOMMAND}" "${ARGS[@]}"
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
    assert_output --regexp ".*File somefile.yaml does not exist.*" ||
      assert_output --regexp ".*no such file or directory.*" ||
      assert_output --regexp ".*Invalid input.*"
  }

  assert_invalid_input() {
    assert_output --regexp '  "message": "Invalid input"' ||
      assert_output --regexp ".*Kubernetes kind not found.*" ||
      assert_output --regexp ".*no such file or directory.*" ||
      assert_output --regexp ".*Invalid input.*" ||
      assert_output --regexp ".*error while parsing.*"
  }

  assert_failure_local() {
    assert_failure
  }

  ## Pod Security Standard (PSS) tests

  assert_pss_rsc_valid() {
    local expected_objects
    expected_objects=${1:-1}

    OBJECTS=$(jq -r '[.[] | select(.valid == true)] | length' <<<"${output:-}")
    ((OBJECTS == expected_objects))

    assert_success
  }

  assert_pss_rsc_invalid() {
    local expected_objects
    expected_objects=${1:-1}

    OBJECTS=$(jq -r '[.[] | select(.valid == false)] | length' <<<"${output:-}")
    ((OBJECTS == expected_objects))

    assert_failure
  }

fi
