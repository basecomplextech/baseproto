// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import (
	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type AnyMessageType interface {
	Type

	anyMessage()
}

// internal

var _ AnyMessageType = (*anyMessageType)(nil)

type anyMessageType struct{}

func newAnyMessageType(typ *model.Type) (*anyMessageType, error) {
	kind := typ.Kind
	if kind != model.KindAnyMessage {
		panic("not any message kind")
	}

	t := &anyMessageType{}
	return t, nil
}

func (t *anyMessageType) anyMessage() {}

// Kind returns a type kind.
func (t *anyMessageType) Kind() model.Kind {
	return model.KindAnyMessage
}

// Name returns a type name.
func (t *anyMessageType) Name() string {
	return "baseproto.Message"
}

// InputName returns an input type, i.e. string (not baseproto.String).
func (t *anyMessageType) InputName() string {
	return "baseproto.Message"
}

// OutputName returns an output type, i.e. baseproto.String (not string).
func (t *anyMessageType) OutputName() string {
	return "baseproto.Message"
}

// List

// AddListElem returns an encode func for a list element.
func (t *anyMessageType) AddListElem() string {
	return "baseproto.WriteMessage"
}

// GetListElem returns a decode func for a list element.
func (t *anyMessageType) GetListElem() string {
	return "baseproto.ParseMessage"
}

// Fields

// GetField writes a field get.
func (t *anyMessageType) GetField(w writer.Writer, tag int) error {
	w.Writef(`return m.msg.Field(%d).Message()`, tag)
	return nil
}
