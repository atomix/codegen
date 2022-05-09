// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/atomix/codegen/internal/generator"
	"github.com/atomix/codegen/internal/generator/proto"
	"github.com/atomix/codegen/internal/generator/template"
	"github.com/spf13/cobra"
	"path/filepath"
)

func run(cmd *cobra.Command, args []string) error {
	inputPath, err := cmd.Flags().GetString("input")
	if err != nil {
		return err
	}

	inputRepoURL, err := cmd.Flags().GetString("repo-url")
	if err != nil {
		return err
	}

	inputRepoTag, err := cmd.Flags().GetString("repo-tag")
	if err != nil {
		return err
	}

	outputPath, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	config := generator.Config{
		Generator: "driver",
		Config: template.Config{
			Templates: []template.TemplateConfig{
				{
					Name: "go.mod",
					Path: getTemplatePath("go.mod.tpl"),
					Output: template.OutputConfig{
						Path: filepath.Join(outputPath, "go.mod"),
					},
				},
			},
		},
		Proto: &proto.Config{
			Input: proto.InputConfig{
				Repo: proto.InputRepo{
					URL: inputRepoURL,
					Tag: inputRepoTag,
				},
				Path: inputPath,
			},
			Output: proto.OutputConfig{
				Path: outputPath,
			},
			Templates: []proto.TemplateConfig{
				{
					Name: "atom.go",
					Path: getTemplatePath("atom.go.tpl"),
					Output: proto.TemplateOutputConfig{
						PathTemplate: "atoms/{{ .Atom.Name | toSnake }}.go",
					},
				},
			},
		},
	}

	context := Context{
		Foo: "bar",
		Bar: "baz",
		Baz: "foo",
	}

	err = generator.Generate(config, context)
	if err != nil {
		return err
	}
	return nil
}

type Context struct {
	Foo string
	Bar string
	Baz string
}

func getTemplatePath(name string) string {
	return filepath.Join("/templates/", name)
}
