// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import "github.com/spf13/cobra"

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "atomix-gen-client",
		RunE: run,
	}
	return cmd
}
