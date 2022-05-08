module github.com/atomix/codegen/client

go 1.18

require (
	github.com/spf13/cobra v1.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require github.com/bmatcuk/doublestar/v4 v4.0.2

require (
	github.com/atomix/codegen v0.0.0-20220508094714-cc2cae885ff9
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

replace github.com/atomix/codegen => ../
