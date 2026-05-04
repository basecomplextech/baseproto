// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package generator

import (
	"fmt"
	"strings"

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
	if err := w.def(def); err != nil {
		return err
	}
	if err := w.values(def); err != nil {
		return err
	}
	if err := w.open_method(def); err != nil {
		return err
	}
	if err := w.decode_method(def); err != nil {
		return err
	}
	if err := w.encode_method(def); err != nil {
		return err
	}
	if err := w.string_method(def); err != nil {
		return err
	}
	return nil
}

func (w *enumWriter) def(def *model.Definition) error {
	w.Linef(`// %v`, def.Name)
	w.Line()
	w.Linef("type %v int32", def.Name)
	w.Line()
	return nil
}

func (w *enumWriter) values(def *model.Definition) error {
	w.Line("const (")

	for _, val := range def.Enum.Values {
		// EnumValue Enum = 1
		name := enumValueName(val)
		w.Linef("%v %v = %d", name, def.Name, val.Number)
	}

	w.Line(")")
	w.Line()
	return nil
}

func (w *enumWriter) open_method(def *model.Definition) error {
	name := def.Name
	w.Linef(`func Open%v(b []byte) %v {`, name, name)
	w.Linef(`v, _, _ := baseproto.DecodeInt32(b)`)
	w.Linef(`return %v(v)`, name)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *enumWriter) decode_method(def *model.Definition) error {
	name := def.Name
	w.Linef(`func Decode%v(b []byte) (result %v, size int, err error) {`, name, name)
	w.Linef(`v, size, err := baseproto.DecodeInt32(b)`)
	w.Linef(`if err != nil || size == 0 {
		return
	}`)
	w.Linef(`result = %v(v)`, name)
	w.Line(`return`)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *enumWriter) encode_method(def *model.Definition) error {
	w.Linef(`func Encode%vTo(b buffer.Buffer, v %v) (int, error) {`, def.Name, def.Name)
	w.Linef(`return baseproto.EncodeInt32(b, int32(v))`)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *enumWriter) string_method(def *model.Definition) error {
	w.Linef("func (e %v) String() string {", def.Name)
	w.Line("switch e {")

	for _, val := range def.Enum.Values {
		name := enumValueName(val)
		w.Linef("case %v:", name)
		w.Linef(`return "%v"`, strings.ToLower(val.Name))
	}

	w.Line("}")
	w.Line(`return ""`)
	w.Line("}")
	w.Line()
	return nil
}

func enumValueName(val *model.EnumValue) string {
	name := toUpperCamelCase(val.Name)
	return fmt.Sprintf("%v_%v", val.Enum.Def.Name, name)
}
