# SPDX-FileCopyrightText: 2022-present Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0

ARG VERSION=latest

FROM atomix/proto-build:$VERSION

COPY protoc-gen-atom /usr/local/bin/protoc-gen-atom
