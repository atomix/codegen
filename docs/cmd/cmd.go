// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "docs",
		Short:   "Generates documentation from Protobuf sources",
		Aliases: []string{"doc"},
		Args:    cobra.NoArgs,
		RunE:    run,
	}
	cmd.Flags().StringP("config", "c", "", "the path to the generator configuration")
	cmd.Flags().StringP("proto-path", "p", ".", "the relative path to the Protobuf API root")
	cmd.Flags().StringSliceP("proto-pattern", "f", []string{"**/*.proto"}, "a pattern by which to filter Protobuf sources")
	cmd.Flags().StringP("docs-path", "d", ".", "the relative path to the documentation root")
	cmd.Flags().String("docs-format", "markdown", "the documentation format")
	_ = cmd.MarkFlagFilename("config")
	return cmd
}
