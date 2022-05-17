# SPDX-FileCopyrightText: 2022-present Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0

ARG VERSION=latest

FROM atomix/codegen:$VERSION

COPY atomix-gen-kubernetes /usr/local/bin/atomix-gen-kubernetes

ENTRYPOINT ["atomix-gen-kubernetes"]
