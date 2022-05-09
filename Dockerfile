# SPDX-FileCopyrightText: 2022-present Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0

FROM golang:1.18

RUN apt-get update && \
    apt-get install -y unzip git

RUN mkdir -p /build

RUN apt-get update && apt-get install -y protobuf-compiler

RUN mkdir -p /go/src/github.com/protocolbuffers && \
    git clone --branch main https://github.com/protocolbuffers/protobuf /go/src/github.com/protocolbuffers/protobuf

RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protoc-3.19.4-linux-x86_64.zip && \
    unzip -o protoc-3.19.4-linux-x86_64.zip -d /usr/local bin/protoc && \
    unzip -o protoc-3.19.4-linux-x86_64.zip -d /usr/local include/* && \
    rm -rf protoc-3.19.4-linux-x86_64.zip

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest && \
    export GO111MODULE=off && \
    go get github.com/gogo/protobuf/proto && \
    go get github.com/gogo/protobuf/gogoproto && \
    go get github.com/gogo/protobuf/protoc-gen-gofast && \
    go get github.com/gogo/protobuf/protoc-gen-gogo && \
    go get github.com/gogo/protobuf/protoc-gen-gogofast && \
    go get github.com/gogo/protobuf/protoc-gen-gogofaster && \
    go get github.com/gogo/protobuf/protoc-gen-gogoslick && \
    go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway && \
    go get github.com/envoyproxy/protoc-gen-validate && \
    go get github.com/favadi/protoc-go-inject-tag

RUN mkdir -p /go/src/github.com/google && \
    git clone --branch main https://github.com/google/protobuf /go/src/github.com/google/protobuf && \
    git clone --branch main https://github.com/googleapis/api-common-protos /go/src/github.com/googleapis/api-common-protos && \
    mkdir -p /go/src/github.com/ && \
    wget "https://github.com/grpc/grpc-web/releases/download/1.2.1/protoc-gen-grpc-web-1.2.1-linux-x86_64" --quiet && \
    mv protoc-gen-grpc-web-1.2.1-linux-x86_64 /usr/local/bin/protoc-gen-grpc-web && \
    chmod +x /usr/local/bin/protoc-gen-grpc-web

WORKDIR /build

COPY protoc-gen-atom /usr/local/bin/protoc-gen-atom
