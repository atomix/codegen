// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "atomix-gen-example",
		Args: cobra.NoArgs,
		RunE: run,
	}
	cmd.Flags().StringP("input", "i", ".", "the input path")
	cmd.Flags().String("repo-url", "", "the input repo URL")
	cmd.Flags().String("repo-tag", "", "the input repo tag")
	cmd.Flags().StringP("output", "o", ".", "the output path")
	return cmd
}
