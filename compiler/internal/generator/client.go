// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package generator

import (
	"fmt"

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
	if err := w.iface(def); err != nil {
		return err
	}
	if err := w.methods(def); err != nil {
		return err
	}
	if err := w.ifaceEnd(def); err != nil {
		return err
	}
	if err := w.new_client(def); err != nil {
		return err
	}
	if err := w.channels(def); err != nil {
		return err
	}
	return nil
}

// iface

func (w *clientWriter) iface(def *model.Definition) error {
	if def.Service.Sub {
		w.Linef(`// %vCall`, def.Name)
		w.Line()
		w.Linef(`type %vCall interface {`, def.Name)
		w.Line()
	} else {
		w.Linef(`// %vClient`, def.Name)
		w.Line()
		w.Linef(`type %vClient interface {`, def.Name)
		w.Line()
	}
	return nil
}

// new_client

func (w *clientWriter) new_client(def *model.Definition) error {
	name := clientImplName(def)

	if def.Service.Sub {
		w.Linef(`func New%vCall(client baserpc.Client, req *baserpc.Request) %vCall {`, def.Name, def.Name)
		w.Linef(`return new%vCall(client, req)`, def.Name)
		w.Line(`}`)
		w.Line()
		w.Linef(`func New%vCallErr(st status.Status) %vCall {`, def.Name, def.Name)
		w.Linef(`return &%v{st: st}`, name)
		w.Line(`}`)
		w.Line()
	} else {
		w.Linef(`func New%vClient(client baserpc.Client) %vClient {`, def.Name, def.Name)
		w.Linef(`return &%v{client: client}`, name)
		w.Line(`}`)
		w.Line()
	}

	return nil
}

// methods

func (w *clientWriter) methods(def *model.Definition) error {
	for _, m := range def.Service.Methods {
		if err := w.method(def, m); err != nil {
			return err
		}
	}
	return nil
}

func (w *clientWriter) method(def *model.Definition, m *model.Method) error {
	methodName := toUpperCamelCase(m.Name)
	w.Write(methodName)

	if err := w.method_input(def, m); err != nil {
		return err
	}
	if err := w.method_output(def, m); err != nil {
		return err
	}
	return nil
}

func (w *clientWriter) method_input(def *model.Definition, m *model.Method) error {
	ctx := "ctx async.Context, "
	switch {
	case m.Subservice != nil:
		ctx = ""
	}

	switch {
	default:
		w.Writef(`(%v) `, ctx)

	case m.Request != nil:
		typeName := typeName(m.Request)
		w.Writef(`(%v req_ %v) `, ctx, typeName)
	}
	return nil
}

func (w *clientWriter) method_output(def *model.Definition, m *model.Method) error {
	switch {
	default:
		w.Line(`status.Status`)

	case m.Subservice != nil:
		typeName := typeName(m.Subservice)
		w.Linef(`%vCall`, typeName)

	case m.Channel != nil:
		name := clientChannel_name(m)
		w.Linef(`(%v, status.Status)`, name)

	case m.Response != nil:
		typeName := typeName(m.Response)
		w.Linef(`(ref.R[%v], status.Status)`, typeName)
	}
	return nil
}

// ifaceEnd

func (w *clientWriter) ifaceEnd(def *model.Definition) error {
	w.Linef(`Unwrap() baserpc.Client`)
	w.Line(`}`)
	w.Line()
	return nil
}

// channel

func (w *clientWriter) channels(def *model.Definition) error {
	for _, m := range def.Service.Methods {
		if m.Channel == nil {
			continue
		}

		if err := w.channel(def, m); err != nil {
			return err
		}
	}
	return nil
}

func (w *clientWriter) channel(def *model.Definition, m *model.Method) error {
	name := clientChannel_name(m)
	w.Linef(`type %v interface {`, name)

	// Send methods
	if in := m.Channel.In; in != nil {
		typeName := typeName(in)
		w.Linef(`Send(ctx async.Context, msg %v) status.Status `, typeName)
		w.Line(`SendEnd(ctx async.Context) status.Status `)
	}

	// Receive methods
	if out := m.Channel.Out; out != nil {
		typeName := typeName(out)
		w.Linef(`Receive(ctx async.Context) (%v, status.Status)`, typeName)
		w.Linef(`ReceiveAsync(ctx async.Context) (%v, bool, status.Status)`, typeName)
		w.Line(`ReceiveWait() <-chan struct{}`)
	}

	// Response method
	{
		w.Write(`Response(ctx async.Context) `)

		if m.Response != nil {
			typeName := typeName(m.Response)
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
