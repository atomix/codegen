// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"os"
	"os/exec"
)

func Run(command string, args ...string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	return RunIn(wd, command, args...)
}

func RunIn(dir string, command string, args ...string) error {
	cmd := &exec.Cmd{
		Dir:    dir,
		Path:   command,
		Args:   args,
		Env:    os.Environ(),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	return cmd.Run()
}
