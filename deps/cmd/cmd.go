// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "atomix-gen-deps",
		Args: cobra.NoArgs,
		RunE: run,
	}
	cmd.Flags().BoolP("check", "c", false, "check module compatibility only")
	cmd.Flags().StringP("version", "v", "", "the target runtime API version")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}
