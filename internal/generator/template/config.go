// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package template

type Config struct {
	Templates []TemplateConfig `yaml:"templates,omitempty"`
}

type TemplateConfig struct {
	Name   string       `yaml:"name"`
	Path   string       `yaml:"path"`
	Output OutputConfig `yaml:"output,omitempty"`
}

type OutputConfig struct {
	Path string `yaml:"path"`
}
