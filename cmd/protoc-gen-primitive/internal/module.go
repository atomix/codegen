// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"errors"
	"github.com/atomix/codegen/pkg/generator/template"
	runtimev1 "github.com/atomix/runtime/api/atomix/runtime/v1"
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
	atoms := make(map[string]*protoAtom)
	for _, target := range targets {
		for _, service := range target.Services() {
			atomName, err := getAtomName(service)
			if err != nil {
				continue
			}
			atom, ok := atoms[atomName]
			if !ok {
				atom = &protoAtom{}
				atoms[atomName] = atom
			}

			atomService := service
			componentType, err := getComponentType(service)
			if err != nil {
				continue
			}
			switch componentType {
			case runtimev1.ComponentType_ATOM:
				atom.service = atomService
			case runtimev1.ComponentType_MANAGER:
				atom.manager = atomService
			}
		}
	}

	for _, atom := range atoms {
		if atom.isComplete() {
			m.generateAtom(atom)
		}
	}
	return m.Artifacts()
}

func (m *Module) getAtomParams(atom *protoAtom) (AtomParams, error) {
	var atomParams AtomParams

	atomName, err := getAtomName(atom.service)
	if err != nil {
		return atomParams, err
	}
	atomParams.Name = atomName

	// Iterate through the methods on the service and construct method metadata for the template.
	methods := make([]MethodParams, 0)
	for _, method := range atom.service.Methods() {
		operationID, err := getOperationID(method)
		if err != nil {
			return atomParams, err
		}

		// Get the operation type for the method.
		operationType, err := getOperationType(method)
		if err != nil {
			return atomParams, err
		}

		var methodTypeParams MethodTypeParams
		switch operationType {
		case runtimev1.OperationType_COMMAND:
			methodTypeParams.IsCommand = true
		case runtimev1.OperationType_QUERY:
			methodTypeParams.IsQuery = true
		default:
			return atomParams, errors.New("primitive service can only contain COMMAND or QUERY type operations")
		}

		requestHeaders, err := m.ctx.HeadersFieldParams(method.Input())
		if err != nil {
			return atomParams, err
		} else if requestHeaders == nil {
			return atomParams, errors.New("no request headers found on method input " + method.Input().FullyQualifiedName())
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
			return atomParams, err
		} else if responseHeaders == nil {
			return atomParams, errors.New("no request headers found on method input " + method.Output().FullyQualifiedName())
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

	atomParams.ServiceParams = ServiceParams{
		Type: ServiceTypeParams{
			EntityParams: m.ctx.EntityParams(atom.service),
			Name:         pgsgo.PGGUpperCamelCase(atom.service.Name()).String(),
		},
		Comment: atom.service.SourceCodeInfo().LeadingComments(),
		Methods: methods,
	}

	managerParams, err := m.getManagerParams(atom)
	if err != nil {
		return atomParams, err
	}
	atomParams.Manager = managerParams
	return atomParams, nil
}

func (m *Module) getManagerParams(atom *protoAtom) (ManagerParams, error) {
	var managerParams ManagerParams

	// Iterate through the methods on the service and construct method metadata for the template.
	methods := make([]MethodParams, 0)
	for _, method := range atom.manager.Methods() {
		operationID, err := getOperationID(method)
		if err != nil {
			return managerParams, err
		}

		// Get the operation type for the method.
		operationType, err := getOperationType(method)
		if err != nil {
			return managerParams, err
		}

		var methodTypeParams MethodTypeParams
		switch operationType {
		case runtimev1.OperationType_CREATE:
			methodTypeParams.IsCreate = true
		case runtimev1.OperationType_CLOSE:
			methodTypeParams.IsClose = true
		default:
			panic("manager service can only contain CREATE or CLOSE type operations")
		}

		requestHeaders, err := m.ctx.HeadersFieldParams(method.Input())
		if err != nil {
			return managerParams, err
		} else if requestHeaders == nil {
			return managerParams, errors.New("no request headers found on method input " + method.Input().FullyQualifiedName())
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
			return managerParams, err
		} else if responseHeaders == nil {
			return managerParams, errors.New("no request headers found on method input " + method.Output().FullyQualifiedName())
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

	managerParams.ServiceParams = ServiceParams{
		Type: ServiceTypeParams{
			EntityParams: m.ctx.EntityParams(atom.manager),
			Name:         pgsgo.PGGUpperCamelCase(atom.manager.Name()).String(),
		},
		Comment: atom.manager.SourceCodeInfo().LeadingComments(),
		Methods: methods,
	}
	return managerParams, nil
}

func (m *Module) generateAtom(atom *protoAtom) {
	atomParams, err := m.getAtomParams(atom)
	if err != nil {
		panic(err)
	}

	values, err := m.ctx.Values()
	if err != nil {
		panic(err)
	}

	// Generate the store metadata.
	params := Params{
		Atom:   atomParams,
		Values: values,
	}

	outputPath := m.ctx.OutputPath(params)
	m.Logf("%s => ", atomParams.Name, outputPath)
	tpl := gotemplate.Must(template.New(filepath.Base(m.ctx.TemplatePath())).ParseFiles(m.ctx.TemplatePath()))
	m.OverwriteGeneratorTemplateFile(outputPath, tpl, params)
}

type protoAtom struct {
	service pgs.Service
	manager pgs.Service
}

func (a *protoAtom) isComplete() bool {
	return a.service != nil && a.manager != nil
}
