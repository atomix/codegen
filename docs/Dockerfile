# SPDX-FileCopyrightText: 2022-present Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0

ARG VERSION=latest

FROM atomix/codegen:$VERSION

COPY atomix-gen-docs /usr/local/bin/atomix-gen-docs

ENTRYPOINT ["atomix-gen-docs"]
