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

	protoPatterns, err := cmd.Flags().GetStringSlice("proto-pattern")
	if err != nil {
		return err
	}
	config.Proto.Files = protoPatterns

	docsPath, err := cmd.Flags().GetString("docs-path")
	if err != nil {
		return err
	}
	err = os.MkdirAll(docsPath, 0755)
	if err != nil {
		return err
	}
	config.Docs.Path = docsPath

	format, err := cmd.Flags().GetString("docs-format")
	if err != nil {
		return err
	}
	config.Docs.Format = format
	return Generate(config)
}
