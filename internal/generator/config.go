// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"github.com/atomix/codegen/internal/generator/proto"
	"github.com/atomix/codegen/internal/generator/template"
)

type Config struct {
	Generator       string `yaml:"generator"`
	template.Config `yaml:",inline"`
	Proto           *proto.Config `yaml:"proto,omitempty"`
}
