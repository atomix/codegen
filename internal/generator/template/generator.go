// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"fmt"
	"os"
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
	for _, template := range g.Config.Templates {
		if err := NewTemplate(g, template).Generate(values); err != nil {
			return err
		}
	}
	return nil
}

func NewTemplate(parent *Generator, template TemplateConfig) *TemplateGenerator {
	return &TemplateGenerator{
		Generator: parent,
		Template:  template,
	}
}

type TemplateGenerator struct {
	*Generator
	Template TemplateConfig
}

func (g *TemplateGenerator) Generate(values interface{}) error {
	template, err := New(g.Template.Name).ParseFiles(g.Template.Path)
	if err != nil {
		return err
	}
	file, err := os.Create(g.Template.Output.Path)
	if err != nil {
		return err
	}
	return template.Execute(file, values)
}

type Spec struct {
	Template  string
	Output    string
	Atom      string
	Component string
	Args      map[string]string
}

func (s Spec) String() string {
	var elems []string
	elems = append(elems, formatArg("template", s.Template))
	elems = append(elems, formatArg("output", s.Output))
	if s.Args != nil {
		for key, value := range s.Args {
			elems = append(elems, formatArg(key, value))
		}
	}
	if s.Atom != "" {
		elems = append(elems, formatArg("atom", s.Atom))
	}
	if s.Component != "" {
		elems = append(elems, formatArg("component", s.Component))
	}
	return strings.Join(elems, ",")
}

func formatArg(key, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}
