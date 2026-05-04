// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package generator

import (
	"github.com/basecomplextech/baseproto/compiler/internal/golang"
	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type enumWriter struct {
	writer.Writer
}

func newEnumWriter(w writer.Writer) *enumWriter {
	return &enumWriter{w}
}

func (w *enumWriter) write(def *model.Definition) error {
	enum, err := golang.NewEnum(def)
	if err != nil {
		return err
	}

	if err := w.def(enum); err != nil {
		return err
	}
	if err := w.values(enum); err != nil {
		return err
	}
	if err := w.parse(enum); err != nil {
		return err
	}
	if err := w.string(enum); err != nil {
		return err
	}
	return nil
}

func (w *enumWriter) def(enum *golang.Enum) error {
	w.Linef(`// %v`, enum.Name)
	w.Line()
	w.Linef("type %v int32", enum.Name)
	w.Line()
	return nil
}

func (w *enumWriter) values(enum *golang.Enum) error {
	w.Line("const (")

	for _, val := range enum.Values {
		// EnumValue Enum = 1
		w.Linef("%v %v = %d", val.Name, enum.Name, val.Tag)
	}

	w.Line(")")
	w.Line()
	return nil
}

func (w *enumWriter) parse(enum *golang.Enum) error {
	w.Linef(`func Open%v(b []byte) %v {`, enum.Name, enum.Name)
	w.Linef(`v, _, _ := baseproto.DecodeInt32(b)`)
	w.Linef(`return %v(v)`, enum.Name)
	w.Line(`}`)

	w.Linef(`func Decode%v(b []byte) (v %v, size int, err error) {`, enum.Name, enum.Name)
	w.Linef(`k, size, err := baseproto.DecodeInt32(b)`)
	w.Linef(`if err != nil || size == 0 {
		return
	}`)
	w.Linef(`v = %v(k)`, enum.Name)
	w.Line(`return`)
	w.Line(`}`)

	w.Linef(`func Encode%v(b buffer.Buffer, v %v) (int, error) {`, enum.Name, enum.Name)
	w.Linef(`return baseproto.EncodeInt32(b, int32(v))`)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *enumWriter) string(enum *golang.Enum) error {
	w.Linef("func (e %v) String() string {", enum.Name)
	w.Line("switch e {")

	for _, val := range enum.Values {
		w.Linef("case %v:", val.Name)
		w.Linef(`return "%v"`, val.String)
	}

	w.Line("}")
	w.Line(`return ""`)
	w.Line("}")
	w.Line()
	return nil
}
