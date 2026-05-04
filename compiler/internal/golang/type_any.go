// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import (
	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type AnyType interface {
	Type

	any()
}

// internal

var _ AnyType = (*anyType)(nil)

type anyType struct{}

func newAnyType(typ *model.Type) (*anyType, error) {
	kind := typ.Kind
	if kind != model.KindAny {
		panic("not any kind")
	}

	t := &anyType{}
	return t, nil
}

func (t *anyType) any() {}

// Kind returns a type kind.
func (t *anyType) Kind() model.Kind {
	return model.KindAny
}

// Name returns a type name.
func (t *anyType) Name() string {
	return "baseproto.Value"
}

// InputName returns an input type, i.e. string (not baseproto.String).
func (t *anyType) InputName() string {
	return "baseproto.Value"
}

// OutputName returns an output type, i.e. baseproto.String (not string).
func (t *anyType) OutputName() string {
	return "baseproto.Value"
}

// Funcs

// ParseFunc returns a parse func.
func (t *anyType) ParseFunc() string {
	return "baseproto.ParseValue"
}

// List

// GetListElem returns a decode func for a list element.
func (t *anyType) GetListElem() string {
	return "baseproto.OpenValue"
}

// AddListElem returns an encode func for a list element.
func (t *anyType) AddListElem() string {
	return "baseproto.WriteValue"
}

// Message

// GetField writes a field get.
func (t *anyType) GetField(w writer.Writer, tag int) error {
	w.Writef(`return m.msg.Field(%d)`, tag)
	return nil
}
