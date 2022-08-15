// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package proto

type Config struct {
	Input     InputConfig      `yaml:"input,omitempty"`
	Output    OutputConfig     `yaml:"output,omitempty"`
	Templates []TemplateConfig `yaml:"templates,omitempty"`
}

type InputConfig struct {
	Repo  InputRepo `yaml:"repo,omitempty"`
	Path  string    `yaml:"path,omitempty"`
	Files []string  `yaml:"files,omitempty"`
}

type InputRepo struct {
	URL    string `yaml:"url,omitempty"`
	Branch string `yaml:"branch,omitempty"`
	Tag    string `yaml:"tag,omitempty"`
}

type OutputConfig struct {
	Path string `yaml:"path"`
}

type ModuleConfig struct {
	Path    string `yaml:"path,omitempty"`
	Version string `yaml:"version,omitempty"`
}

type TemplateType string

type TemplateConfig struct {
	Name   string               `yaml:"name,omitempty"`
	Path   string               `yaml:"path,omitempty"`
	Output TemplateOutputConfig `yaml:"output,omitempty"`
}

type TemplateOutputConfig struct {
	PathTemplate string `yaml:"pathTemplate,omitempty"`
}
