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
	Docs  DocsConfig  `yaml:"docs,omitempty"`
}

type ProtoConfig struct {
	Path  string   `yaml:"path,omitempty"`
	Files []string `yaml:"files,omitempty"`
}

type DocsConfig struct {
	Path   string `yaml:"path,omitempty"`
	Format string `yaml:"format,omitempty"`
}
