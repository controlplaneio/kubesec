#!/usr/bin/env bash

load './bin/bats-support/load'
load './bin/bats-assert/load'

APP="./../kseccheck.sh"
TEST_DIR="."

assert_non_zero_points() {
  assert_output --regexp ".*with [1-9]+ points.*"
  assert_success
}

assert_zero_points() {
  assert_output --regexp ".*with 0 points.*"
  assert_failure
}
