# SPDX-FileCopyrightText: 2022-present Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0

ARG VERSION=latest

FROM atomix/codegen:$VERSION

COPY atomix-gen-deps /usr/local/bin/atomix-gen-deps

ENTRYPOINT ["atomix-gen-deps"]
