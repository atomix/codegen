// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/atomix/codegen/pkg/exec"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func run(cmd *cobra.Command, _ []string) error {
	inputPath, err := cmd.Flags().GetString("input-path")
	if err != nil {
		return err
	}
	outputPath, err := cmd.Flags().GetString("output-path")
	if err != nil {
		return err
	}
	if outputPath == "" {
		outputPath = inputPath
	}
	groupVersion, err := cmd.Flags().GetString("group-version")
	if err != nil {
		return err
	}
	deepcopy, err := cmd.Flags().GetBool("deepcopy")
	if err != nil {
		return err
	}
	client, err := cmd.Flags().GetBool("client")
	if err != nil {
		return err
	}
	boilerplate, err := cmd.Flags().GetString("boilerplate")
	if err != nil {
		return err
	}

	tmpDir, err := ioutil.TempDir("", "code-generator")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	err = exec.Run("git", "clone", "https://github.com/kubernetes/code-generator.git", tmpDir)
	if err != nil {
		return err
	}
	err = exec.RunIn(tmpDir, "git", "checkout", "release-1.24")
	if err != nil {
		return err
	}

	var generators []string
	if deepcopy {
		generators = append(generators, "deepcopy")
	}
	if client {
		generators = append(generators, "client")
	}

	var args []string
	args = append(args, filepath.Join(tmpDir, "generate-groups.sh"))
	args = append(args, strings.Join(generators, ","))
	args = append(args, inputPath, outputPath, groupVersion)
	if boilerplate != "" {
		args = append(args, "--go-header-file", boilerplate)
	}
	err = exec.Run("bash", args...)
	if err != nil {
		return err
	}
	return nil
}
