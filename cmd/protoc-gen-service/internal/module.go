// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"github.com/atomix/codegen/pkg/generator/template"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"path/filepath"
	gotemplate "text/template"
)

const moduleName = "primitive"

// NewModule creates a new proto module
func NewModule() pgs.Module {
	return &Module{
		ModuleBase: &pgs.ModuleBase{},
	}
}

// Module is the code generation module
type Module struct {
	*pgs.ModuleBase
	ctx *Context
}

// Name returns the module name
func (m *Module) Name() string {
	return moduleName
}

// InitContext initializes the module context
func (m *Module) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.ctx = newContext(pgsgo.InitContext(c.Parameters()))
	for key, value := range c.Parameters() {
		c.Logf("%s=%s", key, value)
	}
}

// Execute executes the code generator
func (m *Module) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	for _, target := range targets {
		for _, service := range target.Services() {
			m.generateService(service)
		}
	}
	return m.Artifacts()
}

func (m *Module) generateService(service pgs.Service) {
	descriptor, err := m.getDescriptor(service)
	if err != nil {
		panic(err)
	}

	values, err := m.ctx.Values()
	if err != nil {
		panic(err)
	}

	// Generate the store metadata.
	params := Params{
		Service: descriptor,
		Values:  values,
	}

	outputPath := m.ctx.OutputPath(params)
	m.Logf("%s => ", service.Name().String(), outputPath)
	tpl := gotemplate.Must(template.New(filepath.Base(m.ctx.TemplatePath())).ParseFiles(m.ctx.TemplatePath()))
	m.OverwriteGeneratorTemplateFile(outputPath, tpl, params)
}

func (m *Module) getDescriptor(service pgs.Service) (ServiceParams, error) {
	// Iterate through the methods on the service and construct method metadata for the template.
	methods := make([]MethodParams, 0)
	for _, method := range service.Methods() {
		requestParams := RequestParams{
			MessageParams: MessageParams{
				Type: m.ctx.MessageTypeParams(method.Input()),
			},
			IsUnary:  !method.ClientStreaming(),
			IsStream: method.ClientStreaming(),
		}

		// Generate output metadata from the output type.
		responseParams := ResponseParams{
			MessageParams: MessageParams{
				Type: m.ctx.MessageTypeParams(method.Output()),
			},
			IsUnary:  !method.ServerStreaming(),
			IsStream: method.ServerStreaming(),
		}

		methodParams := MethodParams{
			Name:     method.Name().UpperCamelCase().String(),
			Comment:  method.SourceCodeInfo().LeadingComments(),
			Request:  requestParams,
			Response: responseParams,
		}

		methods = append(methods, methodParams)
	}

	return ServiceParams{
		EntityParams: m.ctx.EntityParams(service),
		Name:         pgsgo.PGGUpperCamelCase(service.Name()).String(),
		Comment:      service.SourceCodeInfo().LeadingComments(),
		Methods:      methods,
	}, nil
}
