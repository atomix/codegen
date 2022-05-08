module {{ .Module.Path }}

go 1.18

require (
    github.com/atomix/sdk {{ .Runtime.Version }}
)