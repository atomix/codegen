// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"github.com/atomix/codegen/pkg/exec"
	"github.com/rogpeppe/go-internal/modfile"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func run(cmd *cobra.Command, args []string) error {
	var path string
	if len(args) == 1 {
		path = args[0]
	} else {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		path = dir
	}

	version, err := cmd.Flags().GetString("version")
	if err != nil {
		return err
	}

	checkOnly, err := cmd.Flags().GetBool("check")
	if err != nil {
		return err
	}

	if checkOnly {
		fmt.Fprintf(cmd.OutOrStdout(), "Checking plugin module constraints against target API version %s\n", version)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "Updating plugin module constraints for target API version %s\n", version)
	}

	err = exec.Run("go", "mod", "tidy")
	if err != nil {
		return err
	}

	rootPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	srcModPath := filepath.Join(rootPath, "go.mod")
	fmt.Fprintf(cmd.OutOrStdout(), "Parsing %s\n", srcModPath)
	srcModBytes, err := ioutil.ReadFile(srcModPath)
	if err != nil {
		return err
	}

	srcMod, err := modfile.Parse("go.mod", srcModBytes, nil)
	if err != nil {
		return err
	}

	tmpDir, err := ioutil.TempDir(rootPath, "atomix")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	tgtModURL := fmt.Sprintf("https://raw.githubusercontent.com/atomix/runtime/%s/go.mod", version)
	fmt.Fprintf(cmd.OutOrStdout(), "Downloading go.mod for atomix/runtime %s from %s\n", version, tgtModURL)
	resp, err := http.Get(tgtModURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tgtModPath := filepath.Join(tmpDir, "go.mod")
	tgtModFile, err := os.Create(tgtModPath)
	if err != nil {
		return err
	}
	defer tgtModFile.Close()

	_, err = io.Copy(tgtModFile, resp.Body)
	if err != nil {
		return err
	}

	tgtModBytes, err := ioutil.ReadFile(tgtModPath)
	if err != nil {
		return err
	}

	tgtMod, err := modfile.Parse(tgtModPath, tgtModBytes, nil)
	if err != nil {
		return err
	}

	tgtReqs := make(map[string]string)
	for _, tgtReq := range tgtMod.Require {
		tgtReqs[tgtReq.Mod.Path] = tgtReq.Mod.Version
	}

	for _, srcReq := range srcMod.Require {
		if tgtReqVersion, ok := tgtReqs[srcReq.Mod.Path]; ok {
			fmt.Fprintf(cmd.OutOrStdout(), "Evaluating common dependency %s\n", srcReq.Mod.Path)
			if srcReq.Mod.Version != tgtReqVersion {
				if checkOnly {
					fmt.Fprintf(cmd.OutOrStderr(), "Detected incompatible dependency %s: %s <> %s\n", srcReq.Mod.Path, srcReq.Mod.Version, tgtReqVersion)
					os.Exit(1)
				} else {
					fmt.Fprintf(cmd.OutOrStderr(), "Updating dependency %s: %s => %s\n", srcReq.Mod.Path, srcReq.Mod.Version, tgtReqVersion)
					srcReq.Mod.Version = tgtReqVersion
				}
			}
		}
	}

	srcModBytes, err = srcMod.Format()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(srcModPath, srcModBytes, 0755)
	if err != nil {
		return err
	}
	return nil
}
