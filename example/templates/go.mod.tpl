module github.com/atomix/{{ .Values.Foo }}

go 1.18

require (
    github.com/atomix/sdk {{ .Values.Bar }}
)