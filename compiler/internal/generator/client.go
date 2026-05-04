// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package generator

import (
	"fmt"

	"github.com/basecomplextech/baseproto/compiler/internal/golang"
	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type clientWriter struct {
	writer.Writer
}

func newClientWriter(w writer.Writer) *clientWriter {
	return &clientWriter{w}
}

func (w *clientWriter) write(def *model.Definition) error {
	srv, err := golang.NewService(def)
	if err != nil {
		return err
	}

	if err := w.iface(srv); err != nil {
		return err
	}
	if err := w.methods(srv); err != nil {
		return err
	}
	if err := w.ifaceEnd(); err != nil {
		return err
	}
	if err := w.newClient(srv); err != nil {
		return err
	}
	if err := w.channels(srv); err != nil {
		return err
	}
	return nil
}

// iface

func (w *clientWriter) iface(srv *golang.Service) error {
	w.Linef(`// %v`, srv.Client)
	w.Line()
	w.Linef(`type %v interface {`, srv.Client)
	w.Line()
	return nil
}

// newClient

func (w *clientWriter) newClient(srv *golang.Service) error {
	if srv.Sub {
		w.Linef(`func New%v(client baserpc.Client, req *baserpc.Request) %v {`,
			srv.Client, srv.Client)
		w.Linef(`return new%v(client, req)`, srv.Client)
		w.Line(`}`)
		w.Line()
		w.Linef(`func New%vErr(st status.Status) %v {`, srv.Client, srv.Client)
		w.Linef(`return &%v{st: st}`, srv.ClientImpl)
		w.Line(`}`)
		w.Line()
	} else {
		w.Linef(`func New%v(client baserpc.Client) %v {`, srv.Client, srv.Client)
		w.Linef(`return &%v{client: client}`, srv.ClientImpl)
		w.Line(`}`)
		w.Line()
	}

	return nil
}

// methods

func (w *clientWriter) methods(srv *golang.Service) error {
	for _, m := range srv.Methods {
		if err := w.method(srv, m); err != nil {
			return err
		}
	}
	return nil
}

func (w *clientWriter) method(srv *golang.Service, m *golang.Method) error {
	w.Write(m.Public)

	if err := w.methodInput(srv, m); err != nil {
		return err
	}
	if err := w.methodOutput(srv, m); err != nil {
		return err
	}
	return nil
}

func (w *clientWriter) methodInput(srv *golang.Service, m *golang.Method) error {
	ctx := "ctx async.Context, "
	switch {
	case m.Subservice != nil:
		ctx = ""
	}

	switch {
	default:
		w.Writef(`(%v) `, ctx)

	case m.Request != nil:
		typeName := m.Request.InputName()
		w.Writef(`(%v req %v) `, ctx, typeName)
	}
	return nil
}

func (w *clientWriter) methodOutput(_ *golang.Service, m *golang.Method) error {
	switch {
	default:
		w.Line(`status.Status`)

	case m.Subservice != nil:
		name := m.Subservice.Name()
		w.Linef(`%vCall`, name)

	case m.Channel != nil:
		name := m.Channel.Client
		w.Linef(`(%v, status.Status)`, name)

	case m.Response != nil:
		typeName := m.Response.Name()
		w.Linef(`(ref.R[%v], status.Status)`, typeName)
	}
	return nil
}

// ifaceEnd

func (w *clientWriter) ifaceEnd() error {
	w.Linef(`Unwrap() baserpc.Client`)
	w.Line(`}`)
	w.Line()
	return nil
}

// channel

func (w *clientWriter) channels(srv *golang.Service) error {
	for _, m := range srv.Methods {
		if m.Channel == nil {
			continue
		}

		if err := w.channel(srv, m); err != nil {
			return err
		}
	}
	return nil
}

func (w *clientWriter) channel(srv *golang.Service, m *golang.Method) error {
	channel := m.Channel.Client
	w.Linef(`type %v interface {`, channel)

	// Send methods
	if in := m.Channel.In; in != nil {
		typeName := in.Name()
		w.Linef(`Send(ctx async.Context, msg %v) status.Status `, typeName)
		w.Line(`SendEnd(ctx async.Context) status.Status `)
	}

	// Receive methods
	if out := m.Channel.Out; out != nil {
		typeName := out.Name()
		w.Linef(`Receive(ctx async.Context) (%v, status.Status)`, typeName)
		w.Linef(`ReceiveAsync(ctx async.Context) (%v, bool, status.Status)`, typeName)
		w.Line(`ReceiveWait() <-chan struct{}`)
	}

	// Response method
	{
		w.Write(`Response(ctx async.Context) `)

		if m.Response != nil {
			typeName := m.Response.Name()
			w.Linef(`(%v, status.Status)`, typeName)
		} else {
			w.Line(`status.Status`)
		}
	}

	// Free method
	w.Line(`Free()`)
	w.Line(`}`)
	w.Line()
	return nil
}

func clientChannel_name(m *model.Method) string {
	return fmt.Sprintf("%v%vClientChannel", m.Service.Def.Name, toUpperCamelCase(m.Name))
}
