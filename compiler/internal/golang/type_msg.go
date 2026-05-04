// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import (
	"fmt"

	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type MessageType interface {
	Type

	// Funcs

	// NewFunc returns a new message function.
	NewFunc() string

	// Write

	// Writer returns a message writer type.
	Writer() string

	// NewWriter returns a new message writer func.
	NewWriter() string
}

// internal

var _ MessageType = (*messageType)(nil)

type messageType struct {
	name string
	imp  string
}

func newMessageType(typ *model.Type) (*messageType, error) {
	kind := typ.Kind
	if kind != model.KindMessage {
		panic("not message kind")
	}

	t := &messageType{
		name: typ.Name,
		imp:  typ.ImportName,
	}
	return t, nil
}

// Kind returns a type kind.
func (t *messageType) Kind() model.Kind {
	return model.KindMessage
}

// Name returns a type name.
func (t *messageType) Name() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.%v", t.imp, t.name)
	}
	return t.name
}

// InputName returns an input type, i.e. string (not baseproto.String).
func (t *messageType) InputName() string {
	return t.Name()
}

// OutputName returns an output type, i.e. baseproto.String (not string).
func (t *messageType) OutputName() string {
	return t.Name()
}

// Funcs

// NewFunc returns a new message function.
func (t *messageType) NewFunc() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.New%v", t.imp, t.name)
	}
	return fmt.Sprintf("New%v", t.name)
}

// ParseFunc returns a parse func.
func (t *messageType) ParseFunc() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.Parse%v", t.imp, t.name)
	}
	return fmt.Sprintf("Parse%v", t.name)
}

// List

// AddListElem returns an encode func for a list element.
func (t *messageType) AddListElem() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.New%vWriterTo", t.imp, t.name)
	}
	return fmt.Sprintf("New%vWriterTo", t.name)
}

// GetListElem returns a decode func for a list element.
func (t *messageType) GetListElem() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.Open%vErr", t.imp, t.name)
	}
	return fmt.Sprintf("Open%vErr", t.name)
}

// Message

// GetField writes a field get.
func (t *messageType) GetField(w writer.Writer, tag int) error {
	new := t.NewFunc()
	w.Writef(`return %v(m.msg.Message(%d))`, new, tag)
	return nil
}

// Write

// Writer returns a message writer type.
func (t *messageType) Writer() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.%vWriter", t.imp, t.name)
	}
	return fmt.Sprintf("%vWriter", t.name)
}

// NewWriter returns a new message writer func.
func (t *messageType) NewWriter() string {
	if t.imp != "" {
		return fmt.Sprintf("%v.New%vWriterTo", t.imp, t.name)
	}
	return fmt.Sprintf("New%vWriterTo", t.name)
}
