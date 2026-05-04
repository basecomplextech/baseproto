// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import (
	"fmt"

	"github.com/basecomplextech/baseproto/compiler/internal/model"
)

type Message struct {
	Name   string
	Writer string

	Fields []*MessageField
}

func NewMessage(def *model.Definition) (*Message, error) {
	msg := def.Message
	name := toUpperCamelCase(def.Name)
	writerName := fmt.Sprintf("%vWriter", name)

	// Make message
	m := &Message{
		Name:   name,
		Writer: writerName,
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
