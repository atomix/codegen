// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"github.com/atomix/codegen/internal/exec"
	"github.com/bmatcuk/doublestar/v4"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const protoExt = ".proto"

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
	importMappings := make(map[string]string)
	importMappings["google/protobuf/any.proto"] = "github.com/gogo/protobuf/types"
	importMappings["google/protobuf/timestamp.proto"] = "github.com/gogo/protobuf/types"
	importMappings["google/protobuf/duration.proto"] = "github.com/gogo/protobuf/types"
	for _, pattern := range g.Config.Proto.Patterns {
		err := doublestar.GlobWalk(os.DirFS(g.Config.Proto.Path), pattern, func(path string, info fs.DirEntry) error {
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(info.Name()) != protoExt {
				return nil
			}
			importMappings[path] = filepath.Join(g.Config.Go.ImportPath, filepath.Dir(path))
			return nil
		})
		if err != nil {
			return err
		}
	}
	return NewGo(g, importMappings).Generate()
}

func (g *Generator) gen(file string, spec Spec) error {
	var path []string
	path = append(path, ".")
	path = append(path, g.Config.Proto.Path)
	path = append(path, filepath.Join(os.Getenv("GOPATH"), "src/github.com/gogo/protobuf"))

	var args []string
	args = append(args, "-I", strings.Join(path, ":"))
	args = append(args, fmt.Sprintf("--gogofaster_out=%s", spec.String()))
	args = append(args, file)

	return exec.Run("protoc", args...)
}

func NewGo(parent *Generator, imports map[string]string) *GoGenerator {
	return &GoGenerator{
		Generator: parent,
		Imports:   imports,
	}
}

type GoGenerator struct {
	*Generator
	Imports map[string]string
}

func (g *GoGenerator) Generate() error {
	for _, pattern := range g.Config.Proto.Patterns {
		if err := NewGlob(g, pattern).Generate(); err != nil {
			return err
		}
	}
	return nil
}

func NewGlob(generator *GoGenerator, pattern string) *GlobGenerator {
	return &GlobGenerator{
		GoGenerator: generator,
		Pattern:     pattern,
	}
}

type GlobGenerator struct {
	*GoGenerator
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
	return g.gen(g.File, Spec{
		ImportPath:     filepath.Join(g.Config.Go.ImportPath, filepath.Dir(g.File)),
		OutputPath:     g.Config.Go.Path,
		ImportMappings: g.Imports,
	})
}

type Spec struct {
	ImportPath     string
	OutputPath     string
	ImportMappings map[string]string
}

func (s Spec) String() string {
	var args []string
	for key, value := range s.ImportMappings {
		args = append(args, fmt.Sprintf("M%s=%s", key, value))
	}
	args = append(args, fmt.Sprintf("import_path=%s", s.ImportPath))
	args = append(args, fmt.Sprintf("plugins=grpc:%s", s.OutputPath))
	return strings.Join(args, ",")
}
