// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"fmt"
	"strings"
	"text/template"
)

// New creates a new Template for the given template file
func New(name string) *template.Template {
	names := make(map[string]string)
	aliases := make(map[string]string)

	t := template.New(name)
	funcs := template.FuncMap{
		"dir":          dir,
		"base":         base,
		"abs":          abs,
		"ext":          ext,
		"rel":          rel,
		"toCamel":      toCamelCase,
		"toLowerCamel": toLowerCamelCase,
		"lower":        toLowerCase,
		"upper":        toUpperCase,
		"upperFirst":   upperFirst,
		"quote":        quote,
		"isLast":       isLast,
		"split":        split,
		"trim":         trim,
		"ternary":      ternary,
		"alias": func(name, proto string) string {
			if alias, ok := aliases[name]; ok {
				return alias
			}

			i := 0
			for {
				alias := fmt.Sprintf("%s%d", proto, i)
				if _, ok := names[alias]; !ok {
					names[alias] = name
					aliases[name] = alias
					return alias
				}
				i++
			}
		},
		"include": func(name string, data interface{}) (string, error) {
			var buf strings.Builder
			err := t.ExecuteTemplate(&buf, name, data)
			return buf.String(), err
		},
	}
	return t.Funcs(funcs)
}
