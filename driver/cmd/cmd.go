// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "atomix-gen-driver",
		Args: cobra.NoArgs,
		RunE: run,
	}
	cmd.Flags().StringP("name", "n", "", "the driver name")
	cmd.Flags().StringP("api-verison", "v", "v1", "the driver API version")
	cmd.Flags().StringP("module-path", "p", "", "the driver module path")
	cmd.Flags().String("github-owner", "", "the GitHub user to which to publish release artifacts")
	cmd.Flags().String("github-repo", "", "the GitHub repo to which to publish release artifacts")
	cmd.Flags().StringP("input", "i", ".", "the input path")
	cmd.Flags().StringP("output", "o", ".", "the output path")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("module-path")
	return cmd
}
