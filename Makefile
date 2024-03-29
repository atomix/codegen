# SPDX-FileCopyrightText: 2022-present Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0

GOLANG_CROSS_VERSION := v1.18.1

.PHONY: build docs client driver go

all: build build-client build-deps build-docs build-driver build-example build-go build-kubernetes

build:
	goreleaser release --snapshot --rm-dist

build-client:
	$(MAKE) -C client build

build-deps:
	$(MAKE) -C deps build

build-docs:
	$(MAKE) -C docs build

build-driver:
	$(MAKE) -C driver build

build-example:
	$(MAKE) -C example build

build-go:
	$(MAKE) -C go build

build-kubernetes:
	$(MAKE) -C kubernetes build

reuse-tool: # @HELP install reuse if not present
	command -v reuse || python3 -m pip install reuse

license: reuse-tool # @HELP run license checks
	reuse lint
