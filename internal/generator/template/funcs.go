// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"github.com/iancoleman/strcase"
	"path/filepath"
	"reflect"
	"strings"
)

func dir(path string) string {
	return filepath.Dir(path)
}

func base(path string) string {
	return filepath.Base(path)
}

func abs(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return absPath
}

func ext(path string) string {
	return filepath.Ext(path)
}

func rel(basepath, targetpath string) string {
	path, err := filepath.Rel(basepath, targetpath)
	if err != nil {
		panic(err)
	}
	return path
}

func toCamelCase(value string) string {
	return strcase.ToCamel(value)
}

func toLowerCamelCase(value string) string {
	return strcase.ToLowerCamel(value)
}

func toLowerCase(value string) string {
	return strings.ToLower(value)
}

func toUpperCase(value string) string {
	return strings.ToUpper(value)
}

func upperFirst(value string) string {
	bytes := []byte(value)
	first := strings.ToUpper(string([]byte{bytes[0]}))
	return string(append([]byte(first), bytes[1:]...))
}

func quote(value string) string {
	return "\"" + value + "\""
}

func isLast(values interface{}, index int) bool {
	t := reflect.ValueOf(values)
	return index == t.Len()-1
}

func split(value, sep string) []string {
	return strings.Split(value, sep)
}

func trim(value string) string {
	return strings.Trim(value, " ")
}

func ternary(v1, v2 interface{}, b bool) interface{} {
	if b {
		return v1
	}
	return v2
}
