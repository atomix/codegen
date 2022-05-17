// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package internal

// Params is the parameters for the code generator
type Params struct {
	Imports []PackageParams
	Atom    AtomParams
	Values  map[string]interface{}
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

// AtomParams is the metadata for an primitive
type AtomParams struct {
	ServiceParams
	Manager ManagerParams
	Name    string
}

// ManagerParams is the metadata for an primitive manager
type ManagerParams struct {
	ServiceParams
}

// ServiceParams is the metadata for a service
type ServiceParams struct {
	Type    ServiceTypeParams
	Comment string
	Methods []MethodParams
}

// ServiceTypeParams is metadata for a service type
type ServiceTypeParams struct {
	EntityParams
	Name string
}

// FieldRefParams is metadata for a field reference
type FieldRefParams struct {
	Field FieldParams
}

// FieldParams is metadata for a field
type FieldParams struct {
	Type TypeParams
	Path []PathParams
}

// PathParams is metadata for a field path
type PathParams struct {
	Name string
	Type TypeParams
}

// MethodParams is the metadata for a primitive method
type MethodParams struct {
	ID       uint32
	Name     string
	Type     MethodTypeParams
	Comment  string
	Request  RequestParams
	Response ResponseParams
}

// MessageParams is the metadata for a message
type MessageParams struct {
	Type TypeParams
}

// RequestParams is the type metadata for a message
type RequestParams struct {
	MessageParams
	Headers  FieldRefParams
	IsUnary  bool
	IsStream bool
}

// ResponseParams is the type metadata for a message
type ResponseParams struct {
	MessageParams
	Headers  FieldRefParams
	IsUnary  bool
	IsStream bool
}

// MethodTypeParams is the metadata for a store method type
type MethodTypeParams struct {
	IsCommand bool
	IsQuery   bool
	IsCreate  bool
	IsClose   bool
}
