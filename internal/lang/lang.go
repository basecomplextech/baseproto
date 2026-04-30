// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package lang

import (
	"github.com/basecomplextech/baseproto/internal/lang/compiler"
	"github.com/basecomplextech/baseproto/internal/lang/generator"
)

type baseProto struct {
	importPath []string
	skipRPC    bool
}

func New(importPath []string, skipRPC bool) *baseProto {
	return &baseProto{
		importPath: importPath,
		skipRPC:    skipRPC,
	}
}

func (s *baseProto) Generate(srcPath string, dstPath string) error {
	if dstPath == "" {
		dstPath = srcPath
	}

	compiler, err := compiler.New(compiler.Options{
		ImportPath: s.importPath,
	})
	if err != nil {
		return err
	}

	pkg, err := compiler.Compile(srcPath)
	if err != nil {
		return err
	}

	gen := generator.New(s.skipRPC)
	return gen.Package(pkg, dstPath)
}
