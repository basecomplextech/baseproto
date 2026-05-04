// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package generator

import (
	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type fileWriter struct {
	writer.Writer

	skipRPC bool
}

func newFileWriter(w writer.Writer, skipRPC bool) *fileWriter {
	return &fileWriter{
		Writer:  w,
		skipRPC: skipRPC,
	}
}

func (w *fileWriter) file(file *model.File) error {
	// Package
	w.Line("package ", file.Package.Name)
	w.Line()

	// Imports
	w.Line("import (")
	w.Line(`"github.com/basecomplextech/baselibrary/alloc"`)
	w.Line(`"github.com/basecomplextech/baselibrary/async"`)
	w.Line(`"github.com/basecomplextech/baselibrary/bin"`)
	w.Line(`"github.com/basecomplextech/baselibrary/buffer"`)
	w.Line(`"github.com/basecomplextech/baselibrary/pools"`)
	w.Line(`"github.com/basecomplextech/baselibrary/ref"`)
	w.Line(`"github.com/basecomplextech/baselibrary/status"`)
	w.Line(`"github.com/basecomplextech/baseproto"`)

	if !w.skipRPC {
		w.Line(`"github.com/basecomplextech/baseproto/baserpc"`)
		w.Line(`"github.com/basecomplextech/baseproto/proto/prpc"`)
	}

	for _, imp := range file.Imports {
		pkg := importPackage(imp)
		w.Linef(`"%v"`, pkg)
	}
	w.Line(")")
	w.Line()

	// Empty values for imports
	w.Line(`var (`)
	w.Line(`_ alloc.Buffer`)
	w.Line(`_ async.Context`)
	if !w.skipRPC {
		w.Line(`_ baserpc.Client`)
	}
	w.Line(`_ baseproto.Kind`)
	w.Line(`_ bin.Bin64`)
	w.Line(`_ buffer.Buffer`)
	w.Line(`_ pools.Pool[any]`)
	if !w.skipRPC {
		w.Line(`_ prpc.Request`)
	}
	w.Line(`_ ref.Ref`)
	w.Line(`_ status.Status`)
	w.Line(`)`)

	// Definitions
	return w.definitions(file)
}

func (w *fileWriter) definitions(file *model.File) error {
	// Types
	for _, def := range file.Definitions {
		switch def.Type {
		case model.DefinitionEnum:
			if err := w.enum(def); err != nil {
				return err
			}
		case model.DefinitionMessage:
			if err := w.message(def); err != nil {
				return err
			}
		case model.DefinitionStruct:
			if err := w.struct_(def); err != nil {
				return err
			}
		case model.DefinitionService:
			if w.skipRPC {
				continue
			}

			if err := w.service(def); err != nil {
				return err
			}
			if err := w.client(def); err != nil {
				return err
			}
		}
	}

	// Message writers
	for _, def := range file.Definitions {
		if def.Type != model.DefinitionMessage {
			continue
		}
		if err := w.messageWriter(def); err != nil {
			return err
		}
	}

	// Service impls
	if !w.skipRPC {
		for _, def := range file.Definitions {
			if def.Type != model.DefinitionService {
				continue
			}
			if err := w.serviceImpl(def); err != nil {
				return err
			}
		}

		for _, def := range file.Definitions {
			if def.Type != model.DefinitionService {
				continue
			}
			if err := w.clientImpl(def); err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *fileWriter) enum(def *model.Definition) error {
	return newEnumWriter(w.Writer).write(def)
}

func (w *fileWriter) message(def *model.Definition) error {
	return newMessageWriter(w.Writer).write(def)
}

func (w *fileWriter) messageWriter(def *model.Definition) error {
	return newMessageWriterWriter(w.Writer).write(def)
}

func (w *fileWriter) struct_(def *model.Definition) error {
	return newStructWriter(w.Writer).write(def)
}

func (w *fileWriter) client(def *model.Definition) error {
	return newClientWriter(w.Writer).write(def)
}

func (w *fileWriter) clientImpl(def *model.Definition) error {
	return newClientImplWriter(w.Writer).write(def)
}

func (w *fileWriter) service(def *model.Definition) error {
	return newServiceWriter(w.Writer).write(def)
}

func (w *fileWriter) serviceImpl(def *model.Definition) error {
	return newServiceImplWriter(w.Writer).write(def)
}
