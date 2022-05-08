// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package proto

import (
	"encoding/json"
	"fmt"
	"github.com/atomix/codegen/internal/exec"
	"os"
	"path/filepath"
	"strings"
)

func Generate(config Config, values interface{}) error {
	return NewGenerator(config).Generate(values)
}

func NewGenerator(config Config) *Generator {
	if config.Input.Files == nil {
		config.Input.Files = []string{"**/*.proto"}
	}
	return &Generator{
		Config: config,
	}
}

type Generator struct {
	Config Config
}

func (g *Generator) Generate(values interface{}) error {
	for _, pattern := range g.Config.Input.Files {
		if err := NewPath(g, pattern).Generate(values); err != nil {
			return err
		}
	}
	return nil
}

func NewPath(parent *Generator, path string) *PathGenerator {
	return &PathGenerator{
		Generator: parent,
		Path:      path,
	}
}

type PathGenerator struct {
	*Generator
	Path string
}

func (g *PathGenerator) Generate(values interface{}) error {
	for _, template := range g.Config.Templates {
		if err := NewTemplate(g, template).Generate(values); err != nil {
			return err
		}
	}
	return nil
}

func NewTemplate(parent *PathGenerator, template TemplateConfig) *TemplateGenerator {
	return &TemplateGenerator{
		PathGenerator: parent,
		Template:      template,
	}
}

type TemplateGenerator struct {
	*PathGenerator
	Template TemplateConfig
}

func (g *TemplateGenerator) Generate(values interface{}) error {
	var protoPath []string
	protoPath = append(protoPath, ".")
	protoPath = append(protoPath, g.Config.Input.Path)
	protoPath = append(protoPath, filepath.Join(os.Getenv("GOPATH"), "src/github.com/gogo/protobuf"))

	bytes, err := json.Marshal(values)
	if err != nil {
		return err
	}

	var specArgs []string
	specArgs = append(specArgs, fmt.Sprintf("template=%s", g.Template.Path))
	specArgs = append(specArgs, fmt.Sprintf("output=%s", g.Template.Output.PathTemplate))
	specArgs = append(specArgs, fmt.Sprintf("values='%s'", string(bytes)))
	spec := strings.Join(specArgs, ",")

	var protoArgs []string
	protoArgs = append(protoArgs, "-I", strings.Join(protoPath, ":"))
	protoArgs = append(protoArgs, "--template_out=%s", spec)
	protoArgs = append(protoArgs, g.Path)

	return exec.Run("protoc", protoArgs...)
}
