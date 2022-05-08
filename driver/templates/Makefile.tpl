GOLANG_CROSS_VERSION := v1.18.1

.PHONY: build release

build:
	docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/build \
		-w /build \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --snapshot --rm-dist

api-go:
	@cd api && (rm -r **/*.pb.go || true) && cd ..
	atomix gen go \
	    --proto-path ./api \
	    --package github.com/{{ .Repo.Owner }}/{{ .Repo.Name }}/api \
	    --output-path ./api

api-docs:
	@cd api && (rm -r **/*.md || true) && cd ..
	atomix gen docs \
	    --proto-path ./api \
	    --package github.com/{{ .Repo.Owner }}/{{ .Repo.Name }}/api \
	    --output-path ./api

api: api-go api-docs

reuse-tool: # @HELP install reuse if not present
	command -v reuse || python3 -m pip install reuse

license: reuse-tool # @HELP run license checks
	reuse lint