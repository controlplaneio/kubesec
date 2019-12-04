#!/usr/bin/env bash

load '_helper'

setup() {
  _global_setup
   ./kubesec http 8089 &
   KSI_PID=$!
}

teardown() {
  _global_teardown
  kill -9 ${KSI_PID}
}


@test "website example processed ok" {
  run \
    curl -sSX POST --data-binary @"/home/go/kubesec-test.yaml" http://localhost:8089/scan

  assert_success
}

@test "website example with file= prefix processed ok" {
  run \
    curl -sSX POST --data-binary @"/home/go/kubesec-test-prefixed.yaml" http://localhost:8089/scan

  assert_success
}
