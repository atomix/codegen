// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/atomix/codegen/internal/exec"
	"github.com/atomix/codegen/internal/generator"
	"github.com/atomix/codegen/internal/generator/proto"
	"github.com/atomix/codegen/internal/generator/template"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

const protoExt = ".proto"

func run(cmd *cobra.Command, args []string) error {
	var context Context

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	context.Driver.Name = name

	apiVersion, err := cmd.Flags().GetString("api-version")
	if err != nil {
		return err
	}
	context.Driver.APIVersion = apiVersion

	inputPath, err := cmd.Flags().GetString("input")
	if err != nil {
		return err
	}

	outputPath, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	path, err := cmd.Flags().GetString("module-path")
	if err != nil {
		return err
	}
	context.Module.Path = path

	repoOwner, err := cmd.Flags().GetString("github-owner")
	if err != nil {
		return err
	}
	context.Repo.Owner = repoOwner

	repoName, err := cmd.Flags().GetString("github-repo")
	if err != nil {
		return err
	}
	context.Repo.Name = repoName

	config := generator.Config{
		Generator: "driver",
		Config: template.Config{
			Templates: []template.TemplateConfig{
				{
					Name: ".gitignore",
					Path: getTemplatePath(".gitignore.tpl"),
					Output: template.OutputConfig{
						Path: filepath.Join(outputPath, ".gitignore"),
					},
				},
				{
					Name: ".goreleaser.yaml",
					Path: getTemplatePath(".goreleaser.yaml.tpl"),
					Output: template.OutputConfig{
						Path: filepath.Join(outputPath, ".goreleaser.yaml"),
					},
				},
				{
					Name: "Makefile",
					Path: getTemplatePath("Makefile.tpl"),
					Output: template.OutputConfig{
						Path: filepath.Join(outputPath, "Makefile"),
					},
				},
				{
					Name: "go.mod",
					Path: getTemplatePath("go.mod.tpl"),
					Output: template.OutputConfig{
						Path: filepath.Join(outputPath, "go.mod"),
					},
				},
				{
					Name: "driver.go",
					Path: getTemplatePath("driver.go.tpl"),
					Output: template.OutputConfig{
						Path: filepath.Join(outputPath, "driver/driver.go"),
					},
				},
			},
		},
		Proto: &proto.Config{
			Input: proto.InputConfig{
				Path: inputPath,
			},
			Templates: []proto.TemplateConfig{
				{
					Name: "atom.go",
					Path: getTemplatePath("atom.go.tpl"),
					Output: proto.OutputConfig{
						PathTemplate: "driver/{{ .Atom | toSnake }}.go",
					},
				},
			},
		},
	}

	err = generator.Generate(config, context)
	if err != nil {
		return err
	}

	err = exec.RunIn(outputPath, "go", "mod", "tidy")
	if err != nil {
		return err
	}
	return nil
}

type Context struct {
	Driver DriverContext
	Module ModuleContext
	Repo   RepoContext
}

type DriverContext struct {
	Name       string
	APIVersion string
}

type ModuleContext struct {
	Path string
}

type RepoContext struct {
	Owner string
	Name  string
}

func apply(templatePath string, outputPath string, context Context) error {
	tpl := template.New(templatePath)
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	return tpl.Execute(file, context)
}

func getTemplatePath(name string) string {
	return filepath.Join("/templates/", name)
}
