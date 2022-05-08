project_name: {{ .Driver.Name | toSnake }}_driver

before:
  hooks:
    - go mod tidy
    - atomix gen deps --version {{ .Runtime.Version }}
    - go mod tidy

builds:
  - id: plugin
    main: ./cmd/{{ .Driver.Name | toKebab }}
    binary: {{ .Driver.Name | toKebab }}-{{ "{{ .Version }}" }}.{{ .Runtime.Version }}.so
    goos:
      - linux
    goarch:
      - amd64
    env:
      - CC=gcc
      - CXX=g++
    flags:
      - -buildmode=plugin
      - -mod=readonly
      - -trimpath
    gcflags:
      - all=-N -l
    ldflags:
      - -s -w -X ./plugin/{{ .Driver.Name | toKebab }}-driver.version={{ "{{ .Version }}" }} -X ./plugin/{{ .Driver.Name | toKebab }}-driver.commit={{ "{{ .Commit }}" }}

checksum:
  name_template: 'checksums.txt'

snapshot:
  name_template: "{{ "{{ incpatch .Version }}" }}-dev"

changelog:
  sort: asc
  filters:
    exclude:
      - '^docs:'

{{- if .Repo.Name }}
release:
  github:
    owner: {{ .Repo.Owner }}
    name: {{ .Repo.Name }}
  prerelease: auto
  draft: true
{{- end }}
