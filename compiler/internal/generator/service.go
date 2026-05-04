// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package generator

import (
	"github.com/basecomplextech/baseproto/compiler/internal/golang"
	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type serviceWriter struct {
	writer.Writer
}

func newServiceWriter(w writer.Writer) *serviceWriter {
	return &serviceWriter{w}
}

func (w *serviceWriter) write(def *model.Definition) error {
	srv, err := golang.NewService(def)
	if err != nil {
		return err
	}

	if err := w.iface(srv); err != nil {
		return err
	}
	if err := w.newHandler(srv); err != nil {
		return err
	}
	if err := w.channels(srv); err != nil {
		return err
	}
	return nil
}

func (w *serviceWriter) iface(srv *golang.Service) error {
	w.Linef(`// %v`, srv.Name)
	w.Line()
	w.Linef(`type %v interface {`, srv.Name)

	for _, m := range srv.Methods {
		if err := w.method(m); err != nil {
			return err
		}
	}

	w.Line(`}`)
	w.Line()
	return nil
}

func (w *serviceWriter) method(m *golang.Method) error {
	if err := w.methodInput(m); err != nil {
		return err
	}
	if err := w.methodOutput(m); err != nil {
		return err
	}
	w.Line()
	return nil
}

func (w *serviceWriter) methodInput(m *golang.Method) error {
	w.Writef(`%v`, m.Name)

	if m.Oneway {
		w.Write(`(ctx baserpc.ConnContext`)
	} else {
		w.Write(`(ctx baserpc.Context`)
	}

	switch {
	case m.Type == model.MethodType_Channel:
		w.Writef(`, ch %v`, m.Channel.Name)
	case m.Request != nil:
		w.Writef(`, req %v`, m.Request.Name())
	}

	if m.Type == model.MethodType_Subservice {
		w.Writef(`, next baserpc.NextHandler[%v]`, m.Subservice.Name())
	}

	w.Write(`) `)
	return nil
}

func (w *serviceWriter) methodOutput(m *golang.Method) error {
	if m.Response != nil {
		w.Writef(`(ref.R[%v], status.Status)`, m.Response.Name())
	} else {
		w.Write(`status.Status`)
	}
	return nil
}

// newHandler

func (w *serviceWriter) newHandler(srv *golang.Service) error {
	if srv.Sub {
		w.Linef(`func New%v(ctx baserpc.Context, ch baserpc.ServerChannel, `, srv.Handler)
		w.Linef(`index int) baserpc.Subhandler1[%v] {`, srv.Name)
		w.Linef(`return new%v(ctx, ch, index)`, srv.Handler)
		w.Line(`}`)
	} else {
		w.Linef(`func New%v(s %v) baserpc.Handler {`, srv.Handler, srv.Name)
		w.Linef(`return &%v{service: s}`, srv.HandlerImpl)
		w.Line(`}`)
	}

	w.Line()
	return nil
}

// channels

func (w *serviceWriter) channels(srv *golang.Service) error {
	for _, m := range srv.Methods {
		if m.Type != model.MethodType_Channel {
			continue
		}

		if err := w.channel(m); err != nil {
			return err
		}
	}
	return nil
}

func (w *serviceWriter) channel(m *golang.Method) error {
	name := m.Channel.Name
	w.Linef(`type %v interface {`, name)

	// Request method
	switch {
	case m.Request != nil:
		w.Linef(`Request() (%v, status.Status)`, m.Request.Name())
	}

	// Receive methods
	if in := m.Channel.In; in != nil {
		typeName := in.Name()
		w.Linef(`Receive(ctx async.Context) (%v, status.Status)`, typeName)
		w.Linef(`ReceiveAsync(ctx async.Context) (%v, bool, status.Status)`, typeName)
		w.Line(`ReceiveWait() <-chan struct{}`)
	}

	// Send methods
	if out := m.Channel.Out; out != nil {
		typeName := out.Name()
		w.Linef(`Send(ctx async.Context, msg %v) status.Status`, typeName)
		w.Line(`SendEnd(ctx async.Context) status.Status`)
	}

	w.Line(`}`)
	w.Line()
	return nil
}
