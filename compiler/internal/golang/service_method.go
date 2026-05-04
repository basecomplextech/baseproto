// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package golang

import (
	"fmt"

	"github.com/basecomplextech/baseproto/compiler/internal/model"
)

type Method struct {
	Path   string
	Type   model.MethodType
	Oneway bool // Oneway method

	Public  string
	Private string

	Request    MessageType // Message type
	Response   MessageType // Message type
	Channel    *MethodChannel
	Subservice ServiceType // Subservice type
}

// MethodChannel defines in/out channel messages.
type MethodChannel struct {
	Name string

	Client  string
	Handler string

	In  Type
	Out Type
}

func newMethod(srv *model.Service, m *model.Method) (_ *Method, err error) {
	path := m.Name
	public := toUpperCamelCase(m.Name)
	private := toLowerCamelCase(m.Name)

	m1 := &Method{
		Path:   path,
		Oneway: m.Oneway,
		Type:   m.Type,

		Public:  public,
		Private: private,
	}

	if m.Request != nil {
		m1.Request, err = newMessageType(m.Request)
		if err != nil {
			return nil, err
		}
	}

	if m.Response != nil {
		m1.Response, err = newMessageType(m.Response)
		if err != nil {
			return nil, err
		}
	}

	if m.Channel != nil {
		m1.Channel, err = newMethodChannel(srv, m, m.Channel)
		if err != nil {
			return nil, err
		}
	}

	if m.Subservice != nil {
		m1.Subservice, err = newServiceType(m.Subservice)
		if err != nil {
			return nil, err
		}
	}
	return m1, nil
}

func newMethodChannel(srv *model.Service, m *model.Method, ch *model.MethodChannel) (
	_ *MethodChannel, err error) {

	name := fmt.Sprintf("%v%vChannel", srv.Def.Name, toUpperCamelCase(m.Name))
	client := fmt.Sprintf("%v%vClientChannel", m.Service.Def.Name, toUpperCamelCase(m.Name))
	handler := fmt.Sprintf("%v%vChannel", toLowerCamelCase(srv.Def.Name), toUpperCamelCase(m.Name))

	var in Type
	if ch.In != nil {
		in, err = newType(ch.In)
		if err != nil {
			return nil, err
		}
	}

	var out Type
	if ch.Out != nil {
		out, err = newType(ch.Out)
		if err != nil {
			return nil, err
		}
	}

	ch1 := &MethodChannel{
		Name:    name,
		Client:  client,
		Handler: handler,

		In:  in,
		Out: out,
	}
	return ch1, nil
}
