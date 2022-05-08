// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	runtimev1 "github.com/atomix/api/pkg/atomix/runtime/v1"
	"github.com/atomix/codegen/internal/generator/template"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"path/filepath"
	"strings"
	gotemplate "text/template"
)

const moduleName = "atom"

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
}

// Execute executes the code generator
func (m *Module) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	for _, target := range targets {
		m.executeTarget(target)
	}
	return m.Artifacts()
}

func (m *Module) executeTarget(target pgs.File) {
	println(target.File().InputPath())
	for _, service := range target.Services() {
		m.executeService(service)
	}
}

// executeService generates a store from a Protobuf service
//nolint:gocyclo
func (m *Module) executeService(service pgs.Service) {
	atomName, err := getAtomName(service)
	if err != nil {
		return
	}

	componentType, err := getComponentType(service)
	if err != nil {
		return
	}

	if filter, enabled := m.ctx.AtomFilter(); enabled {
		if atomName != filter {
			return
		}
	}
	if filter, enabled := m.ctx.ComponentFilter(); enabled {
		switch componentType {
		case runtimev1.ComponentType_ATOM:
			if strings.ToLower(filter) != "atom" {
				return
			}
		case runtimev1.ComponentType_MANAGER:
			if strings.ToLower(filter) != "manager" {
				return
			}
		}
	}

	// Iterate through the methods on the service and construct method metadata for the template.
	methods := make([]MethodParams, 0)
	for _, method := range service.Methods() {
		operationID, err := getOperationID(method)
		if err != nil {
			panic(err)
		}

		// Get the operation type for the method.
		operationType, err := getOperationType(method)
		if err != nil {
			panic(err)
		}

		methodTypeParams := MethodTypeParams{
			IsCommand: operationType == runtimev1.OperationType_COMMAND,
			IsQuery:   operationType == runtimev1.OperationType_QUERY,
		}

		requestHeaders, err := m.ctx.HeadersFieldParams(method.Input())
		if err != nil {
			panic(err)
		} else if requestHeaders == nil {
			panic("no request headers found on method input " + method.Input().FullyQualifiedName())
		}

		requestParams := RequestParams{
			MessageParams: MessageParams{
				Type: m.ctx.MessageTypeParams(method.Input()),
			},
			Headers:  *requestHeaders,
			IsUnary:  !method.ClientStreaming(),
			IsStream: method.ClientStreaming(),
		}

		responseHeaders, err := m.ctx.HeadersFieldParams(method.Output())
		if err != nil {
			panic(err)
		} else if responseHeaders == nil {
			panic("no request headers found on method input " + method.Output().FullyQualifiedName())
		}

		// Generate output metadata from the output type.
		responseParams := ResponseParams{
			MessageParams: MessageParams{
				Type: m.ctx.MessageTypeParams(method.Output()),
			},
			Headers:  *responseHeaders,
			IsUnary:  !method.ServerStreaming(),
			IsStream: method.ServerStreaming(),
		}

		methodParams := MethodParams{
			ID:       operationID,
			Name:     method.Name().UpperCamelCase().String(),
			Comment:  method.SourceCodeInfo().LeadingComments(),
			Type:     methodTypeParams,
			Request:  requestParams,
			Response: responseParams,
		}

		methods = append(methods, methodParams)
	}

	atomParams := AtomParams{
		Name: atomName,
		ServiceParams: ServiceParams{
			Type: ServiceTypeParams{
				EntityParams: m.ctx.EntityParams(service),
				Name:         pgsgo.PGGUpperCamelCase(service.Name()).String(),
			},
			Comment: service.SourceCodeInfo().LeadingComments(),
			Methods: methods,
		},
	}

	// Generate the store metadata.
	params := Params{
		File: m.ctx.FileParams(service),
		Atom: atomParams,
		Args: m.ctx.Args(),
	}

	tpl := gotemplate.Must(template.New(filepath.Base(m.ctx.TemplatePath())).ParseFiles(m.ctx.TemplatePath()))
	m.OverwriteGeneratorTemplateFile(m.OutputPath(), tpl, params)
}
