// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/atomix/codegen/cmd/protoc-gen-atom/internal"
	"github.com/lyft/protoc-gen-star"
)

func main() {
	pgs.Init(pgs.DebugMode()).
		RegisterModule(internal.NewModule()).
		Render()
}
