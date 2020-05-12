#!/usr/bin/env bash

load '_helper'

setup() {
  _global_setup
}

teardown() {
  _global_teardown
}

@test "test dep - jq is installed" {
  run command -v jq

  assert_success
}

@test "test dep - curl is installed (remote)" {
  skip_if_not_remote

  run command -v curl

  assert_success
}
