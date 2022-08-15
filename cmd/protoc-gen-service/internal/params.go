// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package internal

// Params is the parameters for the code generator
type Params struct {
	Imports []PackageParams
	Service ServiceParams
	Values  map[string]any
}

type EntityParams struct {
	File    FileParams
	Package PackageParams
}

// FileParams is the Protobuf file parameters
type FileParams struct {
	Name     string
	BaseName string
	Path     string
}

// PackageParams is the package for a code file
type PackageParams struct {
	Name string
}

// TypeParams is the metadata for a store type
type TypeParams struct {
	EntityParams
	Name        string
	IsPointer   bool
	IsScalar    bool
	IsCast      bool
	IsMessage   bool
	IsMap       bool
	IsRepeated  bool
	IsEnum      bool
	IsEnumValue bool
	IsBytes     bool
	IsString    bool
	IsInt32     bool
	IsInt64     bool
	IsUint32    bool
	IsUint64    bool
	IsFloat     bool
	IsDouble    bool
	IsBool      bool
	KeyType     *TypeParams
	ValueType   *TypeParams
	Values      []TypeParams
}

// ServiceParams is the metadata for a service
type ServiceParams struct {
	EntityParams
	Name    string
	Comment string
	Methods map[string]MethodParams
}

// FieldRefParams is metadata for a field reference
type FieldRefParams struct {
	Field FieldParams
}

// FieldParams is metadata for a field
type FieldParams struct {
	Type    TypeParams
	Path    []PathParams
	Message *MessageParams
}

// PathParams is metadata for a field path
type PathParams struct {
	Name string
	Type TypeParams
}

// MethodParams is the metadata for a primitive method
type MethodParams struct {
	Name     string
	Comment  string
	Request  RequestParams
	Response ResponseParams
}

// MessageParams is the metadata for a message
type MessageParams struct {
	Type   TypeParams
	Fields map[string]FieldParams
}

// RequestParams is the type metadata for a message
type RequestParams struct {
	MessageParams
	IsUnary  bool
	IsStream bool
}

// ResponseParams is the type metadata for a message
type ResponseParams struct {
	MessageParams
	IsUnary  bool
	IsStream bool
}
