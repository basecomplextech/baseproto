// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package compiler

import (
	"github.com/basecomplextech/baseproto/compiler/internal/compiler"
	"github.com/basecomplextech/baseproto/compiler/internal/generator"
)

type Compiler struct {
	importPath []string
	skipRPC    bool
}

// New returns a new compiler.
func New(importPath []string, skipRPC bool) *Compiler {
	return &Compiler{
		importPath: importPath,
		skipRPC:    skipRPC,
	}
}

// Generate generates code for the given source path.
func (s *Compiler) Generate(srcPath string, dstPath string) error {
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
