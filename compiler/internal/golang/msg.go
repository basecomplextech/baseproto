// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import "github.com/basecomplextech/baseproto/compiler/internal/model"

type Message struct {
	Name   string
	Fields []*MessageField
}

func NewMessage(def *model.Definition) (*Message, error) {
	msg := def.Message

	// Make message
	m := &Message{
		Name:   toUpperCamelCase(def.Name),
		Fields: make([]*MessageField, 0, msg.Fields.Len()),
	}

	// Make fields
	for _, field := range msg.Fields.List {
		f, err := newMessageField(field)
		if err != nil {
			return nil, err
		}
		m.Fields = append(m.Fields, f)
	}

	return m, nil
}
