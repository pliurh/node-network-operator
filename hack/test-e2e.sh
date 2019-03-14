#!/bin/bash -x
set -e
operator-sdk test local ./test/e2e --go-test-flags "-v"
