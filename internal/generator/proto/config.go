// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package proto

type Config struct {
	Input     InputConfig      `yaml:"input,omitempty"`
	Templates []TemplateConfig `yaml:"templates,omitempty"`
}

type InputConfig struct {
	Path  string   `yaml:"path,omitempty"`
	Files []string `yaml:"files,omitempty"`
}

type ModuleConfig struct {
	Path    string `yaml:"path,omitempty"`
	Version string `yaml:"version,omitempty"`
}

type TemplateType string

type TemplateConfig struct {
	Name   string       `yaml:"name,omitempty"`
	Path   string       `yaml:"path,omitempty"`
	Output OutputConfig `yaml:"output,omitempty"`
}

type OutputConfig struct {
	PathTemplate string `yaml:"pathTemplate,omitempty"`
}
