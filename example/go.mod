module github.com/atomix/codegen/example

go 1.18

require github.com/spf13/cobra v1.4.0

require github.com/iancoleman/strcase v0.2.0 // indirect

require (
	github.com/atomix/codegen v0.0.0-20220508094714-cc2cae885ff9
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

replace github.com/atomix/codegen => ../
