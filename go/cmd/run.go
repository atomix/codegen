// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

func run(cmd *cobra.Command, args []string) error {
	var config Config
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	if configPath != "" {
		config, err = ParseConfigFile(configPath)
		if err != nil {
			return err
		}
	}

	protoPath, err := cmd.Flags().GetString("proto-path")
	if err != nil {
		return err
	}
	config.Proto.Path = protoPath

	protoFiles, err := cmd.Flags().GetStringSlice("proto-files")
	if err != nil {
		return err
	}
	config.Proto.Files = protoFiles

	goPath, err := cmd.Flags().GetString("go-path")
	if err != nil {
		return err
	}
	err = os.MkdirAll(goPath, 0755)
	if err != nil {
		return err
	}
	config.Go.Path = goPath

	importPath, err := cmd.Flags().GetString("import-path")
	if err != nil {
		return err
	}
	config.Go.ImportPath = importPath
	return Generate(config)
}
