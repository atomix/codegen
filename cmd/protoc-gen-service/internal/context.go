// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/atomix/codegen/pkg/generator/template"
	"github.com/golang/protobuf/proto"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"path/filepath"
)

const (
	templateParamKey = "template"
	outputParamKey   = "output"
	valuesParamKey   = "values"
)

// newContext creates a new metadata context
func newContext(ctx pgsgo.Context) *Context {
	return &Context{
		ctx: ctx,
	}
}

// Context is the code generation context
type Context struct {
	ctx pgsgo.Context
}

func (c *Context) TemplatePath() string {
	return c.ctx.Params().Str(templateParamKey)
}

func (c *Context) OutputPath(params Params) string {
	outputTemplate := c.ctx.Params().Str(outputParamKey)
	decodedTemplate, err := base64.RawURLEncoding.DecodeString(outputTemplate)
	if err != nil {
		panic(err)
	}

	template, err := template.New(outputParamKey).Parse(string(decodedTemplate))
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err := template.Execute(&buf, params); err != nil {
		panic(err)
	}
	return buf.String()
}

func (c *Context) ImportPath(entity pgs.Entity) string {
	return c.ctx.ImportPath(entity).String()
}

func (c *Context) Values() (map[string]interface{}, error) {
	encodedValues := c.ctx.Params().Str(valuesParamKey)
	decodedValues, err := base64.RawURLEncoding.DecodeString(encodedValues)
	if err != nil {
		return nil, err
	}
	values := make(map[string]interface{})
	if err := json.Unmarshal(decodedValues, &values); err != nil {
		return nil, err
	}
	return values, nil
}

// FilePath returns the output path for the given entity
func (c *Context) FilePath(entity pgs.Entity, file string) string {
	path := c.ctx.Params().OutputPath()
	if path == "" {
		path = c.ctx.OutputPath(entity).Dir().String()
	}
	return filepath.Join(path, file)
}

// FileParams extracts the file parameters for the given entity
func (c *Context) FileParams(entity pgs.Entity) FileParams {
	return FileParams{
		Name:     entity.File().InputPath().Base(),
		BaseName: entity.File().InputPath().BaseName(),
		Path:     entity.File().InputPath().String(),
	}
}

func (c *Context) PackageParams(entity pgs.Entity) PackageParams {
	return PackageParams{
		Name: c.ctx.PackageName(entity).String(),
	}
}

func (c *Context) EntityParams(entity pgs.Entity) EntityParams {
	return EntityParams{
		File:    c.FileParams(entity),
		Package: c.PackageParams(entity),
	}
}

func (c *Context) findMessage(typeName string, packages map[string]pgs.Package) (pgs.Message, bool) {
	for _, pkg := range packages {
		for _, file := range pkg.Files() {
			for _, msg := range file.Messages() {
				if msg.FullyQualifiedName() == typeName {
					return msg, true
				}
			}
		}
	}
	return nil, false
}

func (c *Context) FieldParams(field pgs.Field) FieldParams {
	params := FieldParams{
		Type: c.FieldTypeParams(field),
		Path: []PathParams{
			{
				Name: c.FieldName(field),
				Type: c.FieldTypeParams(field),
			},
		},
	}
	if field.Type().IsEmbed() {
		message := c.MessageParams(field.Message())
		params.Message = &message
	}
	return params
}

func (c *Context) MessageParams(message pgs.Message) MessageParams {
	fields := make(map[string]FieldParams)
	for _, field := range message.Fields() {
		fields[field.Name().String()] = c.FieldParams(field)
	}
	return MessageParams{
		Type:   c.MessageTypeParams(message),
		Fields: fields,
	}
}

// MessageTypeParams extracts the type metadata for the given message
func (c *Context) MessageTypeParams(message pgs.Message) TypeParams {
	return TypeParams{
		EntityParams: c.EntityParams(message),
		Name:         pgsgo.PGGUpperCamelCase(message.Name()).String(),
		IsMessage:    true,
	}
}

func getProtoTypeName(protoType pgs.ProtoType) string {
	switch protoType {
	case pgs.BytesT:
		return "[]byte"
	case pgs.StringT:
		return "string"
	case pgs.Int32T:
		return "int32"
	case pgs.Int64T:
		return "int64"
	case pgs.UInt32T:
		return "uint32"
	case pgs.UInt64T:
		return "uint64"
	case pgs.FloatT:
		return "float32"
	case pgs.DoubleT:
		return "float64"
	case pgs.BoolT:
		return "bool"
	}
	return ""
}

// FieldName computes the name for the given field
func (c *Context) FieldName(field pgs.Field) string {
	customName, err := getCustomName(field)
	if err != nil {
		panic(err)
	} else if customName != nil {
		return *customName
	}
	embed, err := getEmbed(field)
	if err != nil {
		panic(err)
	} else if embed != nil && *embed {
		return pgsgo.PGGUpperCamelCase(field.Type().Embed().Name()).String()
	}
	name := field.Name()
	if name == "size" {
		name = "size_"
	}
	return pgsgo.PGGUpperCamelCase(name).String()
}

// RawFieldTypeParams extracts the raw type metadata for the given field
func (c *Context) RawFieldTypeParams(field pgs.Field) TypeParams {
	if field.Type().IsMap() {
		return c.MapFieldTypeParams(field)
	}
	if field.Type().IsRepeated() {
		return c.RepeatedFieldTypeParams(field)
	}
	if field.Type().IsEmbed() {
		return c.MessageFieldTypeParams(field)
	}
	if field.Type().IsEnum() {
		return c.EnumFieldTypeParams(field)
	}
	protoType := field.Type().ProtoType()
	return TypeParams{
		EntityParams: c.EntityParams(field),
		Name:         getProtoTypeName(field.Type().ProtoType()),
		IsScalar:     true,
		IsBytes:      protoType == pgs.BytesT,
		IsString:     protoType == pgs.StringT,
		IsInt32:      protoType == pgs.Int32T,
		IsInt64:      protoType == pgs.Int64T,
		IsUint32:     protoType == pgs.UInt32T,
		IsUint64:     protoType == pgs.UInt64T,
		IsFloat:      protoType == pgs.FloatT,
		IsDouble:     protoType == pgs.DoubleT,
		IsBool:       protoType == pgs.BoolT,
	}
}

// FieldTypeParams extracts the type metadata for the given field
func (c *Context) FieldTypeParams(field pgs.Field) TypeParams {
	if field.Type().IsMap() {
		return c.MapFieldTypeParams(field)
	}
	if field.Type().IsRepeated() {
		return c.RepeatedFieldTypeParams(field)
	}
	if field.Type().IsEmbed() {
		return c.MessageFieldTypeParams(field)
	}
	if field.Type().IsEnum() {
		return c.EnumFieldTypeParams(field)
	}

	protoType := field.Type().ProtoType()
	castType, err := getCastType(field)
	if err != nil {
		panic(err)
	} else if castType != nil {
		return TypeParams{
			EntityParams: c.EntityParams(field),
			Name:         *castType,
			IsScalar:     true,
			IsCast:       true,
			IsBytes:      protoType == pgs.BytesT,
			IsString:     protoType == pgs.StringT,
			IsInt32:      protoType == pgs.Int32T,
			IsInt64:      protoType == pgs.Int64T,
			IsUint32:     protoType == pgs.UInt32T,
			IsUint64:     protoType == pgs.UInt64T,
			IsFloat:      protoType == pgs.FloatT,
			IsDouble:     protoType == pgs.DoubleT,
			IsBool:       protoType == pgs.BoolT,
		}
	}
	return TypeParams{
		EntityParams: c.EntityParams(field),
		Name:         getProtoTypeName(field.Type().ProtoType()),
		IsScalar:     true,
		IsBytes:      protoType == pgs.BytesT,
		IsString:     protoType == pgs.StringT,
		IsInt32:      protoType == pgs.Int32T,
		IsInt64:      protoType == pgs.Int64T,
		IsUint32:     protoType == pgs.UInt32T,
		IsUint64:     protoType == pgs.UInt64T,
		IsFloat:      protoType == pgs.FloatT,
		IsDouble:     protoType == pgs.DoubleT,
		IsBool:       protoType == pgs.BoolT,
	}
}

// MessageFieldTypeParams extracts the type metadata for the given message field
func (c *Context) MessageFieldTypeParams(field pgs.Field) TypeParams {
	var fieldType string
	castType, err := getCastType(field)
	if err != nil {
		panic(err)
	} else if castType != nil {
		fieldType = *castType
	}

	customType, err := getCustomType(field)
	if err != nil {
		panic(err)
	} else if customType != nil {
		fieldType = *customType
	} else if fieldType == "" {
		fieldType = pgsgo.PGGUpperCamelCase(field.Type().Embed().Name()).String()
	}

	pointer := true
	nullable, err := getNullable(field)
	if err != nil {
		panic(err)
	} else if nullable != nil {
		pointer = *nullable
	}

	return TypeParams{
		EntityParams: c.EntityParams(field),
		Name:         fieldType,
		IsMessage:    true,
		IsPointer:    pointer,
	}
}

// RepeatedFieldTypeParams extracts the type metadata for the given repeated field
func (c *Context) RepeatedFieldTypeParams(field pgs.Field) TypeParams {
	elementTypeParams := c.FieldElementTypeParams(field)
	elementTypeParams.IsRepeated = true
	return elementTypeParams
}

// MapFieldTypeParams extracts the type metadata for the given map field
func (c *Context) MapFieldTypeParams(field pgs.Field) TypeParams {
	keyTypeParams := c.FieldKeyTypeParams(field)
	valueTypeParams := c.FieldValueTypeParams(field)
	return TypeParams{
		EntityParams: c.EntityParams(field),
		Name:         "map",
		IsMap:        true,
		KeyType:      &keyTypeParams,
		ValueType:    &valueTypeParams,
	}
}

// FieldKeyTypeParams extracts the key type metadata for the given field
func (c *Context) FieldKeyTypeParams(field pgs.Field) TypeParams {
	castKey, err := getCastKey(field)
	if err != nil {
		panic(err)
	} else if castKey != nil {
		return TypeParams{
			Name: *castKey,
		}
	}
	if field.Type().Key().IsEmbed() {
		return c.MessageTypeParams(field.Type().Key().Embed())
	}
	protoType := field.Type().Element().ProtoType()
	return TypeParams{
		EntityParams: c.EntityParams(field),
		Name:         getProtoTypeName(field.Type().Key().ProtoType()),
		IsScalar:     true,
		IsBytes:      protoType == pgs.BytesT,
		IsString:     protoType == pgs.StringT,
		IsInt32:      protoType == pgs.Int32T,
		IsInt64:      protoType == pgs.Int64T,
		IsUint32:     protoType == pgs.UInt32T,
		IsUint64:     protoType == pgs.UInt64T,
		IsFloat:      protoType == pgs.FloatT,
		IsDouble:     protoType == pgs.DoubleT,
		IsBool:       protoType == pgs.BoolT,
	}
}

// FieldValueTypeParams extracts the value type metadata for the given field
func (c *Context) FieldValueTypeParams(field pgs.Field) TypeParams {
	castValue, err := getCastValue(field)
	if err != nil {
		panic(err)
	} else if castValue != nil {
		return TypeParams{
			Name: *castValue,
		}
	}
	if field.Type().Element().IsEmbed() {
		return c.MessageTypeParams(field.Type().Element().Embed())
	}
	protoType := field.Type().Element().ProtoType()
	return TypeParams{
		EntityParams: c.EntityParams(field),
		Name:         getProtoTypeName(field.Type().Element().ProtoType()),
		IsScalar:     true,
		IsBytes:      protoType == pgs.BytesT,
		IsString:     protoType == pgs.StringT,
		IsInt32:      protoType == pgs.Int32T,
		IsInt64:      protoType == pgs.Int64T,
		IsUint32:     protoType == pgs.UInt32T,
		IsUint64:     protoType == pgs.UInt64T,
		IsFloat:      protoType == pgs.FloatT,
		IsDouble:     protoType == pgs.DoubleT,
		IsBool:       protoType == pgs.BoolT,
	}
}

// FieldElementTypeParams extracts the element type metadata for the given field
func (c *Context) FieldElementTypeParams(field pgs.Field) TypeParams {
	castValue, err := getCastValue(field)
	if err != nil {
		panic(err)
	} else if castValue != nil {
		return TypeParams{
			Name: *castValue,
		}
	}
	if field.Type().Element().IsEmbed() {
		return c.MessageTypeParams(field.Type().Element().Embed())
	}
	protoType := field.Type().Element().ProtoType()
	return TypeParams{
		EntityParams: c.EntityParams(field),
		Name:         getProtoTypeName(field.Type().Element().ProtoType()),
		IsScalar:     true,
		IsBytes:      protoType == pgs.BytesT,
		IsString:     protoType == pgs.StringT,
		IsInt32:      protoType == pgs.Int32T,
		IsInt64:      protoType == pgs.Int64T,
		IsUint32:     protoType == pgs.UInt32T,
		IsUint64:     protoType == pgs.UInt64T,
		IsFloat:      protoType == pgs.FloatT,
		IsDouble:     protoType == pgs.DoubleT,
		IsBool:       protoType == pgs.BoolT,
	}
}

// EnumFieldTypeParams extracts the type metadata for the given enum field
func (c *Context) EnumFieldTypeParams(field pgs.Field) TypeParams {
	values := make([]TypeParams, 0, len(field.Type().Enum().Values()))
	for _, value := range field.Type().Enum().Values() {
		values = append(values, c.EnumValueTypeParams(value))
	}
	return TypeParams{
		EntityParams: c.EntityParams(field),
		Name:         pgsgo.PGGUpperCamelCase(field.Type().Enum().Name()).String(),
		IsEnum:       true,
		Values:       values,
	}
}

// EnumValueTypeParams extracts the type metadata for the given enum value
func (c *Context) EnumValueTypeParams(enumValue pgs.EnumValue) TypeParams {
	return TypeParams{
		EntityParams: c.EntityParams(enumValue),
		Name:         pgsgo.PGGUpperCamelCase(enumValue.Name()).String(),
		IsEnumValue:  true,
	}
}

func (c *Context) findAnnotatedField(message pgs.Message, extension *proto.ExtensionDesc) (*FieldRefParams, error) {
	for _, field := range message.Fields() {
		var isAnnotatedField bool
		ok, err := field.Extension(extension, &isAnnotatedField)
		if err != nil {
			return nil, err
		} else if ok {
			return &FieldRefParams{
				Field: FieldParams{
					Type: c.FieldTypeParams(field),
					Path: []PathParams{
						{
							Name: c.FieldName(field),
							Type: c.FieldTypeParams(field),
						},
					},
				},
			}, nil
		} else if field.Type().IsEmbed() {
			child, err := c.findAnnotatedField(field.Type().Embed(), extension)
			if err != nil {
				return nil, err
			} else if child != nil {
				return &FieldRefParams{
					Field: FieldParams{
						Type: child.Field.Type,
						Path: append([]PathParams{
							{
								Name: c.FieldName(field),
								Type: c.FieldTypeParams(field),
							},
						}, child.Field.Path...),
					},
				}, nil
			}
		}
	}
	return nil, nil
}
