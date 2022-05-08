// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package proto

import (
	"encoding/json"
	"fmt"
	"github.com/atomix/codegen/internal/exec"
	"github.com/bmatcuk/doublestar/v4"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Generate(config Config, values interface{}) error {
	return NewGenerator(config).Generate(values)
}

func NewGenerator(config Config) *Generator {
	return &Generator{
		Config: config,
	}
}

type Generator struct {
	Config Config
}

func (g *Generator) Generate(values interface{}) error {
	for _, pattern := range g.Config.Input.Patterns {
		if err := NewGlob(g, pattern).Generate(values); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) gen(file string, spec Spec) error {
	var path []string
	path = append(path, ".")
	path = append(path, g.Config.Input.Path)
	path = append(path, filepath.Join(os.Getenv("GOPATH"), "src/github.com/gogo/protobuf"))

	var args []string
	args = append(args, "-I", strings.Join(path, ":"))
	args = append(args, "--template_out=%s", spec.String())
	args = append(args, file)

	return exec.Run("protoc", args...)
}

func NewGlob(generator *Generator, pattern string) *GlobGenerator {
	return &GlobGenerator{
		Generator: generator,
		Pattern:   pattern,
	}
}

type GlobGenerator struct {
	*Generator
	Pattern string
}

func (g *GlobGenerator) Generate(values interface{}) error {
	return doublestar.GlobWalk(os.DirFS(g.Config.Input.Path), g.Pattern, func(path string, info fs.DirEntry) error {
		if info.IsDir() {
			return nil
		}
		return NewFile(g, path).Generate(values)
	})
}

func NewFile(parent *GlobGenerator, file string) *FileGenerator {
	return &FileGenerator{
		GlobGenerator: parent,
		File:          file,
	}
}

type FileGenerator struct {
	*GlobGenerator
	File string
}

func (g *FileGenerator) Generate(values interface{}) error {
	for _, template := range g.Config.Templates {
		if err := NewTemplate(g, template).Generate(values); err != nil {
			return err
		}
	}
	return nil
}

func NewTemplate(parent *FileGenerator, template TemplateConfig) *TemplateGenerator {
	return &TemplateGenerator{
		FileGenerator: parent,
		Template:      template,
	}
}

type TemplateGenerator struct {
	*FileGenerator
	Template TemplateConfig
}

func (g *TemplateGenerator) Generate(values interface{}) error {
	bytes, err := json.Marshal(values)
	if err != nil {
		return err
	}
	return g.gen(g.File, Spec{
		Template: g.Template.Path,
		Output:   g.Template.Output.PathTemplate,
		Values:   string(bytes),
	})
}

type Spec struct {
	Template  string
	Output    string
	Atom      string
	Component string
	Values    string
}

func (s Spec) String() string {
	var elems []string
	elems = append(elems, fmt.Sprintf("template=%s", s.Template))
	elems = append(elems, fmt.Sprintf("output=%s", s.Output))
	elems = append(elems, fmt.Sprintf("values='%s'", s.Values))
	return strings.Join(elems, ",")
}
