// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"errors"
	"fmt"
	runtimev1 "github.com/atomix/api/pkg/atomix/runtime/v1"
	"github.com/gogo/protobuf/gogoproto"
	gogoprotobuf "github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/proto" //nolint:staticcheck
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	pgs "github.com/lyft/protoc-gen-star"
)

// getAtomName gets the name extension from the given service
func getAtomName(service pgs.Service) (string, error) {
	var name string
	ok, err := service.Extension(runtimev1.E_Name, &name)
	if err != nil {
		return "", err
	} else if !ok {
		return "", errors.New(fmt.Sprintf("extension '%s' is not set", runtimev1.E_Name.Name))
	}
	return name, nil
}

// getComponentType gets the name extension from the given service
func getComponentType(service pgs.Service) (runtimev1.ComponentType, error) {
	var componentType runtimev1.ComponentType
	ok, err := service.Extension(runtimev1.E_Component, &componentType)
	if err != nil {
		return 0, err
	} else if !ok {
		return 0, nil
	}
	return componentType, nil
}

// getOperationID gets the id extension from the given method
func getOperationID(method pgs.Method) (uint32, error) {
	var operationID uint32
	ok, err := method.Extension(runtimev1.E_OperationId, &operationID)
	if err != nil {
		return 0, err
	} else if !ok {
		return 0, errors.New(fmt.Sprintf("extension '%s' is not set", runtimev1.E_OperationId.Name))
	}
	return operationID, nil
}

// getOperationType gets the optype extension from the given method
func getOperationType(method pgs.Method) (runtimev1.OperationType, error) {
	var operationType runtimev1.OperationType
	ok, err := method.Extension(runtimev1.E_OperationType, &operationType)
	if err != nil {
		return 0, err
	} else if !ok {
		return 0, errors.New(fmt.Sprintf("extension '%s' is not set", runtimev1.E_OperationType.Name))
	}
	return operationType, nil
}

// getHeaders gets the headers extension from the given field
func getHeaders(field pgs.Field) (bool, error) {
	var headers bool
	ok, err := field.Extension(runtimev1.E_Headers, &headers)
	if err != nil {
		return false, err
	} else if !ok {
		return false, nil
	}
	return headers, nil
}

// getEmbed gets the embed extension from the given field
func getEmbed(field pgs.Field) (*bool, error) {
	var embed bool
	ok, err := field.Extension(getExtensionDesc(gogoproto.E_Embed), &embed)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &embed, nil
}

// getCastType gets the casttype extension from the given field
func getCastType(field pgs.Field) (*string, error) {
	var castType string
	ok, err := field.Extension(getExtensionDesc(gogoproto.E_Casttype), &castType)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &castType, nil
}

// getCastKey gets the castkey extension from the given field
func getCastKey(field pgs.Field) (*string, error) {
	var castKey string
	ok, err := field.Extension(getExtensionDesc(gogoproto.E_Castkey), &castKey)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &castKey, nil
}

// getCastValue gets the castvalue extension from the given field
func getCastValue(field pgs.Field) (*string, error) {
	var castValue string
	ok, err := field.Extension(getExtensionDesc(gogoproto.E_Castvalue), &castValue)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &castValue, nil
}

// getCustomName gets the customname extension from the given field
func getCustomName(field pgs.Field) (*string, error) {
	var customName string
	ok, err := field.Extension(getExtensionDesc(gogoproto.E_Customname), &customName)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &customName, nil
}

// getCustomType gets the customtype extension from the given field
func getCustomType(field pgs.Field) (*string, error) {
	var customType string
	ok, err := field.Extension(getExtensionDesc(gogoproto.E_Customtype), &customType)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &customType, nil
}

// getNullable gets the nullable extension from the given field
func getNullable(field pgs.Field) (*bool, error) {
	var nullable bool
	ok, err := field.Extension(getExtensionDesc(gogoproto.E_Nullable), &nullable)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &nullable, nil
}

func getExtensionDesc(extension *gogoprotobuf.ExtensionDesc) *proto.ExtensionDesc {
	return &proto.ExtensionDesc{
		ExtendedType:  (*descriptor.FieldOptions)(nil),
		ExtensionType: extension.ExtensionType,
		Field:         extension.Field,
		Name:          extension.Name,
		Tag:           extension.Tag,
		Filename:      extension.Filename,
	}
}
