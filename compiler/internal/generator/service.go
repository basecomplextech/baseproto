// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package generator

import (
	"fmt"

	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type serviceWriter struct {
	writer.Writer
}

func newServiceWriter(w writer.Writer) *serviceWriter {
	return &serviceWriter{w}
}

func (w *serviceWriter) service(def *model.Definition) error {
	if err := w.iface(def); err != nil {
		return err
	}
	if err := w.new_handler(def); err != nil {
		return err
	}
	if err := w.channels(def); err != nil {
		return err
	}
	return nil
}

func (w *serviceWriter) iface(def *model.Definition) error {
	w.Linef(`// %v`, def.Name)
	w.Line()
	w.Linef(`type %v interface {`, def.Name)

	for _, m := range def.Service.Methods {
		if err := w.method(def, m); err != nil {
			return err
		}
	}

	w.Linef(`}`)
	w.Line()
	return nil
}

func (w *serviceWriter) method(def *model.Definition, m *model.Method) error {
	if err := w.method_input(def, m); err != nil {
		return err
	}
	if err := w.method_output(def, m); err != nil {
		return err
	}
	w.Line()
	return nil
}

func (w *serviceWriter) method_input(def *model.Definition, m *model.Method) error {
	name := toUpperCamelCase(m.Name)
	w.Writef(`%v`, name)

	if m.Oneway {
		w.Write(`(ctx rpc.ConnContext`)
	} else {
		w.Write(`(ctx rpc.Context`)
	}

	switch {
	case m.Type == model.MethodType_Channel:
		channel := serviceChannel_name(m)
		w.Writef(`, ch %v`, channel)
	case m.Request != nil:
		typeName := typeName(m.Request)
		w.Writef(`, req %v`, typeName)
	}

	if m.Type == model.MethodType_Subservice {
		sub := m.Subservice
		typeName := typeName(sub)
		w.Writef(`, next rpc.NextHandler[%v]`, typeName)
	}

	w.Write(`) `)
	return nil
}

func (w *serviceWriter) method_output(def *model.Definition, m *model.Method) error {
	if m.Response != nil {
		typeName := typeName(m.Response)
		w.Writef(`(ref.R[%v], status.Status)`, typeName)
	} else {
		w.Write(`status.Status`)
	}
	return nil
}

// new_handler

func (w *serviceWriter) new_handler(def *model.Definition) error {
	name := handler_name(def)

	if def.Service.Sub {
		w.Linef(`func New%vHandler(ctx rpc.Context, channel rpc.ServerChannel, index int) rpc.Subhandler1[%v] {`,
			def.Name, def.Name)
		w.Linef(`return new%vHandler(ctx, channel, index)`, def.Name)
		w.Linef(`}`)
	} else {
		w.Linef(`func New%vHandler(s %v) rpc.Handler {`, def.Name, def.Name)
		w.Linef(`return &%v{service: s}`, name)
		w.Linef(`}`)
	}

	w.Line()
	return nil
}

// channels

func (w *serviceWriter) channels(def *model.Definition) error {
	for _, m := range def.Service.Methods {
		if m.Type != model.MethodType_Channel {
			continue
		}

		if err := w.channel(def, m); err != nil {
			return err
		}
	}
	return nil
}

func (w *serviceWriter) channel(def *model.Definition, m *model.Method) error {
	name := serviceChannel_name(m)
	w.Linef(`type %v interface {`, name)

	// Request method
	switch {
	case m.Request != nil:
		typeName := typeName(m.Request)
		w.Linef(`Request() (%v, status.Status)`, typeName)
	}

	// Receive methods
	if in := m.Channel.In; in != nil {
		typeName := typeName(in)
		w.Linef(`Receive(ctx async.Context) (%v, status.Status)`, typeName)
		w.Linef(`ReceiveAsync(ctx async.Context) (%v, bool, status.Status)`, typeName)
		w.Line(`ReceiveWait() <-chan struct{}`)
	}

	// Send methods
	if out := m.Channel.Out; out != nil {
		typeName := typeName(out)
		w.Linef(`Send(ctx async.Context, msg %v) status.Status`, typeName)
		w.Line(`SendEnd(ctx async.Context) status.Status`)
	}

	w.Linef(`}`)
	w.Line()
	return nil
}

func serviceChannel_name(m *model.Method) string {
	return fmt.Sprintf("%v%vChannel", m.Service.Def.Name, toUpperCamelCase(m.Name))
}
