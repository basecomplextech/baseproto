// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import (
	"fmt"

	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type StructType interface {
	EncodableType

	// Funcs

	// OpenFunc returns an open struct func.
	OpenFunc() string
}

// internal

var _ StructType = (*structType)(nil)

type structType struct {
	name string
	imp  string
}

func newStructType(typ *model.Type) (*structType, error) {
	kind := typ.Kind
	if kind != model.KindStruct {
		panic("not struct kind")
	}

	t := &structType{
		name: typ.Name,
		imp:  typ.ImportName,
	}
	return t, nil
}

// Kind returns a type kind.
func (t *structType) Kind() model.Kind {
	return model.KindStruct
}

// Name returns a type name.
func (t *structType) Name() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.%v", t.imp, t.name)
	}
	return t.name
}

// InputName returns an input type, i.e. string (not baseproto.String).
func (t *structType) InputName() string {
	return t.Name()
}

// OutputName returns an output type, i.e. baseproto.String (not string).
func (t *structType) OutputName() string {
	return t.Name()
}

// Funcs

// OpenFunc returns an open struct func.
func (t *structType) OpenFunc() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.Open%v", t.imp, t.name)
	}
	return "Open" + t.name
}

// EncodeFunc returns an encode struct func.
func (t *structType) EncodeFunc() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.Encode%v", t.imp, t.name)
	}
	return fmt.Sprintf("Encode%v", t.name)
}

// DecodeFunc returns a decode struct func.
func (t *structType) DecodeFunc() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.Decode%v", t.imp, t.name)
	}
	return "Decode" + t.name
}

// DecodeCloneFunc returns a decode func, which returns string clones.
func (t *structType) DecodeCloneFunc() string {
	return t.DecodeFunc()
}

// ParseFunc returns a parse func.
func (t *structType) ParseFunc() string {
	return t.DecodeFunc()
}

// List

// AddListElem returns an encode func for a list element.
func (t *structType) AddListElem() string {
	return t.EncodeFunc()
}

// GetListElem returns a decode func for a list element.
func (t *structType) GetListElem() string {
	return t.DecodeFunc()
}

// Message

// GetField writes a field get.
func (t *structType) GetField(w writer.Writer, tag int) error {
	open := t.OpenFunc()
	w.Writef(`return %v(m.msg.FieldRaw(%d))`, open, tag)
	return nil
}

// WriteField writes a field write.
func (t *structType) WriteField(w writer.Writer, tag int) error {
	return nil
}
