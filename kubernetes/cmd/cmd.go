// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "atomix-gen-kubernetes",
		Args: cobra.NoArgs,
		RunE: run,
	}
	cmd.Flags().StringP("input-path", "p", ".", "the relative path to the API root")
	cmd.Flags().StringP("output-path", "o", "", "the relative path to the output directory")
	cmd.Flags().StringP("group-version", "g", "", "the group:version tuple")
	cmd.Flags().Bool("deepcopy", false, "generate deepcopy files")
	cmd.Flags().Bool("client", false, "generate API clients")
	cmd.Flags().String("boilerplate", "", "the path to a boilerplate file")
	_ = cmd.MarkFlagRequired("input-path")
	_ = cmd.MarkFlagRequired("input-path")
	_ = cmd.MarkFlagRequired("group-version")
	return cmd
}
