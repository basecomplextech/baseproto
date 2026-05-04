// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package generator

import (
	"fmt"

	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type structWriter struct {
	writer.Writer
}

func newStructWriter(w writer.Writer) *structWriter {
	return &structWriter{w}
}

func (w *structWriter) struct_(def *model.Definition) error {
	if err := w.def(def); err != nil {
		return err
	}
	if err := w.open(def); err != nil {
		return err
	}
	if err := w.decode_method(def); err != nil {
		return err
	}
	if err := w.encode_method(def); err != nil {
		return err
	}
	return nil
}

func (w *structWriter) def(def *model.Definition) error {
	w.Linef(`// %v`, def.Name)
	w.Line()
	w.Linef("type %v struct {", def.Name)

	fields := def.Struct.Fields.Values()
	for _, field := range fields {
		name := structFieldName(field)
		typ := typeName(field.Type)
		goTag := fmt.Sprintf("`json:\"%v\"`", field.Name)
		w.Linef("%v %v %v", name, typ, goTag)
	}

	w.Line("}")
	w.Line()
	return nil
}

func (w *structWriter) open(def *model.Definition) error {
	w.Linef(`func Open%v(b []byte) %v {`, def.Name, def.Name)
	w.Linef(`s, _, _ := Decode%v(b)`, def.Name)
	w.Line(`return s`)
	w.Line(`}`)
	w.Line()

	w.Linef(`func Decode%v(b []byte) (s %v, size int, err error) {`, def.Name, def.Name)
	w.Line(`size, err = s.Decode(b)`)
	w.Line(`return s, size, err`)
	w.Line(`}`)
	w.Line()

	w.Linef(`func Encode%vTo(b buffer.Buffer, s %v) (int, error) {`, def.Name, def.Name)
	w.Line(`return s.EncodeTo(b)`)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *structWriter) decode_method(def *model.Definition) error {
	w.Linef(`func (s *%v) Decode(b []byte) (size int, err error) {`, def.Name)
	w.Line(`dataSize, size, err := baseproto.DecodeStruct(b)`)
	w.Line(`if err != nil || size == 0 {
		return
	}`)
	w.Line()

	w.Line(`b = b[len(b)-size:]
	n := size - dataSize
	off := len(b) - n
	`)
	w.Line()

	w.Line(`// Decode in reverse order`)
	w.Line()

	fields := def.Struct.Fields.Values()
	for i := len(fields) - 1; i >= 0; i-- {
		field := fields[i]
		fieldName := structFieldName(field)
		decodeName := typeDecodeFunc(field.Type)
		if field.Type.Kind == model.KindString {
			decodeName = "baseproto.DecodeStringClone"
		}

		w.Linef(`s.%v, n, err = %v(b[:off])`, fieldName, decodeName)
		w.Line(`if err != nil {
			return
		}`)
		w.Line(`off -= n`)
		w.Line()
	}

	w.Line(`return size, err`)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *structWriter) encode_method(def *model.Definition) error {
	w.Linef(`func (s %v) EncodeTo(b buffer.Buffer) (int, error) {`, def.Name)
	w.Line(`var dataSize, n int`)
	w.Line(`var err error`)
	w.Line()

	fields := def.Struct.Fields.Values()
	for _, field := range fields {
		fieldName := structFieldName(field)
		writeFunc := typeWriteFunc(field.Type)

		w.Linef(`n, err = %v(b, s.%v)`, writeFunc, fieldName)
		w.Line(`if err != nil {
			return 0, err
		}`)
		w.Line(`dataSize += n`)
		w.Line()
	}

	w.Line(`n, err = baseproto.EncodeStruct(b, dataSize)`)
	w.Line(`if err != nil {
			return 0, err
		}`)
	w.Line(`return dataSize + n, nil`)
	w.Line(`}`)
	w.Line()
	return nil
}

func structFieldName(field *model.StructField) string {
	return toUpperCamelCase(field.Name)
}
