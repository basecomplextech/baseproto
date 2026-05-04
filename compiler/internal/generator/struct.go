// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package generator

import (
	"fmt"

	"github.com/basecomplextech/baseproto/compiler/internal/golang"
	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type structWriter struct {
	writer.Writer
}

func newStructWriter(w writer.Writer) *structWriter {
	return &structWriter{w}
}

func (w *structWriter) write(def *model.Definition) error {
	str, err := golang.NewStruct(def)
	if err != nil {
		return err
	}

	if err := w.def(str); err != nil {
		return err
	}
	if err := w.open(str); err != nil {
		return err
	}
	if err := w.decode(str); err != nil {
		return err
	}
	if err := w.encode(str); err != nil {
		return err
	}
	return nil
}

func (w *structWriter) def(str *golang.Struct) error {
	w.Linef(`// %v`, str.Name)
	w.Line()
	w.Linef("type %v struct {", str.Name)

	for _, field := range str.Fields {
		goTag := fmt.Sprintf("`json:\"%v\"`", field.Name)
		w.Linef("%v %v %v", field.Name, field.Type.Name(), goTag)
	}

	w.Line("}")
	w.Line()
	return nil
}

func (w *structWriter) open(str *golang.Struct) error {
	w.Linef(`func Open%v(b []byte) %v {`, str.Name, str.Name)
	w.Linef(`s, _, _ := Decode%v(b)`, str.Name)
	w.Line(`return s`)
	w.Line(`}`)

	w.Linef(`func Decode%v(b []byte) (s %v, size int, err error) {`, str.Name, str.Name)
	w.Line(`size, err = s.Decode(b)`)
	w.Line(`return s, size, err`)
	w.Line(`}`)

	w.Linef(`func Encode%v(b buffer.Buffer, s %v) (int, error) {`, str.Name, str.Name)
	w.Line(`return s.Encode(b)`)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *structWriter) decode(str *golang.Struct) error {
	w.Linef(`func (s *%v) Decode(b []byte) (size int, err error) {`, str.Name)
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

	fields := str.Fields
	for i := len(fields) - 1; i >= 0; i-- {
		field := fields[i]
		decode := field.Type.DecodeCloneFunc()

		w.Linef(`s.%v, n, err = %v(b[:off])`, field.Name, decode)
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

func (w *structWriter) encode(str *golang.Struct) error {
	w.Linef(`func (s %v) Encode(b buffer.Buffer) (int, error) {`, str.Name)
	w.Line(`var dataSize, n int`)
	w.Line(`var err error`)
	w.Line()

	for _, field := range str.Fields {
		encode := field.Type.EncodeFunc()

		w.Linef(`n, err = %v(b, s.%v)`, encode, field.Name)
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
