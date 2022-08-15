// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"github.com/atomix/codegen/pkg/exec"
	"github.com/bmatcuk/doublestar/v4"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const protoExt = ".proto"
const markdownFormat = "markdown"

func Generate(config Config) error {
	return NewGenerator(config).Generate()
}

func NewGenerator(config Config) *Generator {
	return &Generator{
		Config: config,
	}
}

type Generator struct {
	Config Config
}

func (g *Generator) Generate() error {
	for _, pattern := range g.Config.Proto.Files {
		if err := NewGlob(g, pattern).Generate(); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) gen(file string, spec Spec) error {
	var path []string
	path = append(path, ".")
	path = append(path, g.Config.Proto.Path)
	path = append(path, filepath.Join(os.Getenv("GOPATH"), "src/github.com/gogo/protobuf"))

	var args []string
	args = append(args, "-I", strings.Join(path, ":"))
	docDir := filepath.Dir(filepath.Join(g.Config.Docs.Path, file))
	if err := os.MkdirAll(docDir, 0755); err != nil {
		return err
	}
	args = append(args, fmt.Sprintf("--doc_out=%s", docDir))
	args = append(args, fmt.Sprintf("--doc_opt=%s", spec.String()))
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

func (g *GlobGenerator) Generate() error {
	return doublestar.GlobWalk(os.DirFS(g.Config.Proto.Path), g.Pattern, func(path string, info fs.DirEntry) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(info.Name()) != protoExt {
			return nil
		}
		return NewFile(g, path).Generate()
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

func (g *FileGenerator) Generate() error {
	mdFile := fmt.Sprintf("%s.md", g.File[:len(g.File)-len(filepath.Ext(g.File))])
	return g.gen(g.File, Spec{
		FileName: mdFile,
		Format:   markdownFormat,
	})
}

type Spec struct {
	FileName string
	Format   string
}

func (s Spec) String() string {
	var args []string
	args = append(args, s.Format)
	args = append(args, s.FileName)
	return strings.Join(args, ",")
}
