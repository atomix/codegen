// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"github.com/atomix/codegen/internal/generator/proto"
	"github.com/atomix/codegen/internal/generator/template"
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
	if err := template.Generate(g.Config.Config, values); err != nil {
		return err
	}
	if g.Config.Proto != nil {
		if err := proto.Generate(*g.Config.Proto, values); err != nil {
			return err
		}
	}
	return nil
}
