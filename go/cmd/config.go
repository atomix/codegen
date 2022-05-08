// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func ParseConfigFile(path string) (Config, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	return ParseConfig(bytes)
}

func ParseConfig(bytes []byte) (Config, error) {
	var config Config
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return config, err
	}
	return config, nil
}

type Config struct {
	Proto ProtoConfig `yaml:"proto,omitempty"`
	Go    GoConfig    `yaml:"go,omitempty"`
}

type ProtoConfig struct {
	Path     string   `yaml:"path,omitempty"`
	Patterns []string `yaml:"patterns,omitempty"`
}

type GoConfig struct {
	Path       string `yaml:"path,omitempty"`
	ImportPath string `yaml:"import_path,omitempty"`
}
