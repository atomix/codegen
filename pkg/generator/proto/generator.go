// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package proto

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/atomix/codegen/pkg/exec"
	"github.com/bmatcuk/doublestar/v4"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Generate(config Config, values interface{}) error {
	return NewGenerator(config).Generate(values)
}

func NewGenerator(config Config) *Generator {
	if config.Input.Path == "" {
		config.Input.Path = "."
	}
	if config.Input.Files == nil {
		config.Input.Files = []string{"**/*.proto"}
	}
	if config.Output.Path == "" {
		config.Output.Path = "."
	}
	return &Generator{
		Config: config,
	}
}

type Generator struct {
	Config Config
}

func (g *Generator) Generate(values interface{}) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	if g.Config.Input.Repo.URL != "" {
		dir, err = ioutil.TempDir("", "input")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		err = exec.Run("git", "clone", g.Config.Input.Repo.URL, dir)
		if err != nil {
			return err
		}
		if g.Config.Input.Repo.Branch != "" {
			err = exec.RunIn(dir, "git", "checkout", g.Config.Input.Repo.Branch)
			if err != nil {
				return err
			}
		}
		if g.Config.Input.Repo.Tag != "" {
			err = exec.RunIn(dir, "git", "checkout", g.Config.Input.Repo.Tag)
			if err != nil {
				return err
			}
		}
	}
	return NewDir(g, dir).Generate(values)
}

func NewDir(parent *Generator, dir string) *DirGenerator {
	return &DirGenerator{
		Generator: parent,
		Dir:       dir,
	}
}

type DirGenerator struct {
	*Generator
	Dir string
}

func (g *DirGenerator) Generate(values interface{}) error {
	return NewPath(g, g.Config.Input.Path).Generate(values)
}

func NewPath(parent *DirGenerator, path string) *PathGenerator {
	return &PathGenerator{
		DirGenerator: parent,
		Path:         path,
	}
}

type PathGenerator struct {
	*DirGenerator
	Path string
}

func (g *PathGenerator) Generate(values interface{}) error {
	for _, pattern := range g.Config.Input.Files {
		if err := NewGlob(g, pattern).Generate(values); err != nil {
			return err
		}
	}
	return nil
}

func NewGlob(parent *PathGenerator, pattern string) *GlobGenerator {
	return &GlobGenerator{
		PathGenerator: parent,
		Pattern:       pattern,
	}
}

type GlobGenerator struct {
	*PathGenerator
	Pattern string
}

func (g *GlobGenerator) Generate(values interface{}) error {
	files, err := doublestar.Glob(os.DirFS(filepath.Join(g.Dir, g.Path)), g.Pattern)
	if err != nil {
		return err
	}
	return NewFiles(g, files...).Generate(values)
}

func NewFiles(parent *GlobGenerator, files ...string) *FilesGenerator {
	return &FilesGenerator{
		GlobGenerator: parent,
		Files:         files,
	}
}

type FilesGenerator struct {
	*GlobGenerator
	Files []string
}

func (g *FilesGenerator) Generate(values interface{}) error {
	for _, template := range g.Config.Templates {
		if err := NewTemplate(g, template).Generate(values); err != nil {
			return err
		}
	}
	return nil
}

func NewTemplate(parent *FilesGenerator, template TemplateConfig) *TemplateGenerator {
	return &TemplateGenerator{
		FilesGenerator: parent,
		Template:       template,
	}
}

type TemplateGenerator struct {
	*FilesGenerator
	Template TemplateConfig
}

func (g *TemplateGenerator) Generate(values interface{}) error {
	var protoPath []string
	protoPath = append(protoPath, filepath.Join(g.Dir, g.Config.Input.Path))
	protoPath = append(protoPath, filepath.Join(os.Getenv("GOPATH"), "src/github.com/gogo/protobuf"))

	bytes, err := json.Marshal(values)
	if err != nil {
		return err
	}

	var specArgs []string
	specArgs = append(specArgs, fmt.Sprintf("template=%s", g.Template.Path))
	specArgs = append(specArgs, fmt.Sprintf("output=%s", base64.RawURLEncoding.EncodeToString([]byte(g.Template.Output.PathTemplate))))
	specArgs = append(specArgs, fmt.Sprintf("values=%s", base64.RawURLEncoding.EncodeToString(bytes)))
	spec := strings.Join(specArgs, ",")

	var protoArgs []string
	protoArgs = append(protoArgs, fmt.Sprintf("-I=%s", strings.Join(protoPath, ":")))
	protoArgs = append(protoArgs, fmt.Sprintf("--atom_out=%s:%s", spec, g.Config.Output.Path))
	protoArgs = append(protoArgs, g.Files...)

	return exec.Run("protoc", protoArgs...)
}
