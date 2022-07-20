// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import "github.com/spf13/cobra"

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "atomix-gen-go",
		Short:   "Generates Go sources from Protobuf sources",
		Aliases: []string{"golang"},
		Args:    cobra.NoArgs,
		RunE:    run,
	}
	cmd.Flags().StringP("config", "c", "", "the path to the generator configuration")
	cmd.Flags().StringSliceP("proto-path", "p", []string{"."}, "the relative path to the Protobuf API root")
	cmd.Flags().StringSliceP("proto-files", "f", []string{"**/*.proto"}, "file name patterns by which to filter Protobuf sources")
	cmd.Flags().StringP("go-path", "d", ".", "the relative path to the documentation root")
	cmd.Flags().StringP("import-path", "i", "", "the base Go path for generated sources")
	_ = cmd.MarkFlagFilename("config")
	return cmd
}
