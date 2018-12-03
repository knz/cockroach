#!/usr/bin/env bash

set -euo pipefail

source "$(dirname "${0}")/teamcity-support.sh"

tc_prepare

export TMPDIR=$PWD/artifacts/test
mkdir -p "$TMPDIR"

tc_start_block "Maybe stress pull request"
run build/builder.sh go install ./pkg/cmd/github-pull-request-make
run build/builder.sh env BUILD_VCS_NUMBER="$BUILD_VCS_NUMBER" TARGET=stress github-pull-request-make
tc_end_block "Maybe stress pull request"

tc_start_block "Compile C dependencies"
run script -t5 artifacts/test-deps.log \
	build/builder.sh make -Otarget c-deps
tc_end_block "Compile C dependencies"

tc_start_block "Run Go tests"
run script -t5 artifacts/test.log \
	build/builder.sh \
	make test TESTFLAGS='-v' \
	| go-test-teamcity
tc_end_block "Run Go tests"

tc_start_block "Run C++ tests"
run script -t5 artifacts/test-cpp.log \
	build/builder.sh \
	make check-libroach
tc_end_block "Run C++ tests"
