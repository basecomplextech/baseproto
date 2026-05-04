// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package generator

import (
	"fmt"
	"strings"

	"github.com/basecomplextech/baseproto/compiler/internal/model"
	"github.com/basecomplextech/baseproto/compiler/internal/writer"
)

type clientImplWriter struct {
	writer.Writer
}

func newClientImplWriter(w writer.Writer) *clientImplWriter {
	return &clientImplWriter{w}
}

func (w *clientImplWriter) write(def *model.Definition) error {
	if err := w.def(def); err != nil {
		return err
	}
	if err := w.methods(def); err != nil {
		return err
	}
	if err := w.unwrap(def); err != nil {
		return err
	}
	if err := w.free(def); err != nil {
		return err
	}
	if err := w.channels(def); err != nil {
		return err
	}
	return nil
}

func (w *clientImplWriter) def(def *model.Definition) error {
	name := clientImplName(def)
	w.Linef(`// %v`, name)
	w.Line()

	if def.Service.Sub {
		w.Linef(`var %vPool = pools.NewPoolFunc(`, name)
		w.Linef(`func() *%v {`, name)
		w.Linef(`return &%v{}`, name)
		w.Line(`},`)
		w.Line(`)`)
		w.Line()
		w.Linef(`type %v struct {`, name)
		w.Line(`client rpc.Client`)
		w.Line(`req *rpc.Request`)
		w.Line(`st status.Status`)
		w.Line(`}`)
		w.Line()
		w.Linef(`func new%vCall(client rpc.Client, req *rpc.Request) %vCall {`, def.Name, def.Name)
		w.Linef(`c := %vPool.New()`, name)
		w.Line(`c.client = client`)
		w.Line(`c.req = req`)
		w.Line(`c.st = status.OK`)
		w.Line(`return c`)
		w.Line(`}`)
		w.Line()
	} else {
		w.Linef(`type %v struct {`, name)
		w.Line(`client rpc.Client`)
		w.Line(`}`)
		w.Line()
	}
	return nil
}

func (w *clientImplWriter) free(def *model.Definition) error {
	if !def.Service.Sub {
		return nil
	}

	name := clientImplName(def)
	w.Linef(`func (c *%v) free() {`, name)
	w.Linef(`*c = %v{}`, name)
	w.Linef(`%vPool.Put(c)`, name)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *clientImplWriter) methods(def *model.Definition) error {
	for _, m := range def.Service.Methods {
		if err := w.method(def, m); err != nil {
			return err
		}
	}
	return nil
}

func (w *clientImplWriter) method(def *model.Definition, m *model.Method) error {
	name := clientImplName(def)
	methodName := toUpperCamelCase(m.Name)
	w.Writef(`func (c *%v) %v`, name, methodName)

	if err := w.method_input(def, m); err != nil {
		return err
	}
	if err := w.method_output(def, m); err != nil {
		return err
	}
	w.Line(`{`)

	if err := w.method_call(def, m); err != nil {
		return err
	}

	switch {
	case m.Subservice != nil:
		if err := w.method_subservice(def, m); err != nil {
			return err
		}
	case m.Channel != nil:
		if err := w.method_channel(def, m); err != nil {
			return err
		}
	default:
		if err := w.method_request(def, m); err != nil {
			return err
		}
		if err := w.method_response(def, m); err != nil {
			return err
		}
	}

	w.Line(`}`)
	w.Line()
	return nil
}

func (w *clientImplWriter) method_input(def *model.Definition, m *model.Method) error {
	ctx := "ctx async.Context, "
	if m.Subservice != nil {
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

func (w *clientImplWriter) method_output(def *model.Definition, m *model.Method) error {
	switch {
	default:
		w.Write(`(_st status.Status)`)

	case m.Subservice != nil:
		typeName := typeName(m.Subservice)
		w.Writef(`%vCall`, typeName)

	case m.Channel != nil:
		name := clientChannel_name(m)
		w.Writef(`(_ %v, _st status.Status)`, name)

	case m.Response != nil:
		typeName := typeName(m.Response)
		w.Writef(`(_ ref.R[%v], _st status.Status)`, typeName)
	}
	return nil
}

func (w *clientImplWriter) method_error(def *model.Definition, m *model.Method) error {
	if m.Type != model.MethodType_Subservice {
		w.Line(`return`)
		return nil
	}

	name := clientImplNewErr(m.Subservice)
	w.Linef(`return %v(_st)`, name)
	return nil
}

func (w *clientImplWriter) method_call(def *model.Definition, m *model.Method) error {
	if def.Service.Sub {
		w.Line(`defer c.free()`)
		w.Line()
	}

	// Subservice methods do not return status
	if m.Subservice != nil {
		w.Line(`var _st status.Status`)
		w.Line(``)
	}

	// Begin request
	if def.Service.Sub {
		w.Line(`// Continue request`)
		w.Line(`if _st = c.st; !_st.OK() {`)
		w.method_error(def, m)
		w.Line(`}`)
		w.Line(`req := c.req`)
		w.Line(`c.req = nil`)
	} else {
		w.Line(`// Begin request`)
		w.Line(`req := rpc.NewRequest()`)
	}

	// Free request
	if m.Subservice != nil {
		w.Line(`ok := false`)
		w.Line(`defer func() {`)
		w.Line(`if !ok {`)
		w.Line(`req.Free()`)
		w.Line(`}`)
		w.Line(`}()`)
		w.Line()
	} else {
		w.Line(`defer req.Free()`)
		w.Line()
	}

	// Add call
	w.Line(`// Add call`)
	switch {
	default:
		w.Linef(`st := req.AddEmpty("%v")`, m.Name)
		w.Line(`if !st.OK() {`)
		w.Line(`_st = st`)
		w.method_error(def, m)
		w.Line(`}`)

	case m.Request != nil:
		w.Linef(`st := req.AddMessage("%v", req_.Unwrap())`, m.Name)
		w.Line(`if !st.OK() {`)
		w.Line(`_st = st`)
		w.method_error(def, m)
		w.Line(`}`)
	}

	// End request
	w.Line()
	return nil
}

func (w *clientImplWriter) method_subservice(def *model.Definition, m *model.Method) error {
	// Return subservice
	newFunc := clientImplNew(m.Subservice)

	w.Line(`// Return subservice`)
	w.Linef(`sub := %v(c.client, req)`, newFunc)
	w.Line(`ok = true`)
	w.Linef(`return sub`)
	return nil
}

func (w *clientImplWriter) method_channel(def *model.Definition, m *model.Method) error {
	// Build request
	w.Line(`// Build request`)
	w.Line(`preq, st := req.Build()`)
	w.Line(`if !st.OK() {`)
	w.Line(`_st = st`)
	w.Line(`return`)
	w.Line(`}`)
	w.Line()

	// Open channel
	name := clientChannelImpl_name(m)
	w.Line(`// Open channel`)
	w.Line(`ch, st := c.client.Channel(ctx, preq)`)
	w.Line(`if !st.OK() {`)
	w.Line(`_st = st`)
	w.Line(`return`)
	w.Line(`}`)
	w.Linef(`return new%v(ch), status.OK`, strings.Title(name))
	return nil
}

func (w *clientImplWriter) method_request(def *model.Definition, m *model.Method) error {
	// Build request
	w.Line(`// Build request`)
	w.Line(`preq, st := req.Build()`)
	w.Line(`if !st.OK() {`)
	w.Line(`_st = st`)
	w.Line(`return`)
	w.Line(`}`)
	w.Line()

	// Send request
	if m.Oneway {
		w.Line(`// Send request`)
		w.Line(`return c.client.RequestOneway(ctx, preq)`)
	} else {
		w.Line(`// Send request`)
		w.Line(`resp, st := c.client.Request(ctx, preq)`)
		w.Line(`if !st.OK() {`)
		w.Line(`_st = st`)
		w.Line(`return`)
		w.Line(`}`)
		w.Line(`defer resp.Release()`)
		w.Line(``)
	}
	return nil
}

func (w *clientImplWriter) method_response(def *model.Definition, m *model.Method) error {
	switch {
	default:
		w.Line(`return status.OK`)

	case m.Oneway:
		// pass

	case m.Response != nil:
		parseFunc := typeParseFunc(m.Response)
		w.Line(`// Parse result`)
		w.Linef(`result, _, err := %v(resp.Unwrap())`, parseFunc)
		w.Line(`if err != nil {`)
		w.Line(`_st = status.WrapError(err)`)
		w.Line(`return`)
		w.Line(`}`)
		w.Line(`return ref.NextRetain(result, resp), status.OK`)
	}

	return nil
}

// unwrap

func (w *clientImplWriter) unwrap(def *model.Definition) error {
	name := clientImplName(def)
	w.Linef(`func (c *%v) Unwrap() rpc.Client {`, name)
	w.Line(`return c.client `)
	w.Line(`}`)
	w.Line()
	return nil
}

// channel

func (w *clientImplWriter) channels(def *model.Definition) error {
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

func (w *clientImplWriter) channel(def *model.Definition, m *model.Method) error {
	if err := w.channel_def(def, m); err != nil {
		return err
	}
	if err := w.channel_receive(def, m); err != nil {
		return err
	}
	if err := w.channel_send(def, m); err != nil {
		return err
	}
	if err := w.channel_response(def, m); err != nil {
		return err
	}
	if err := w.channel_free(def, m); err != nil {
		return err
	}
	return nil
}

func (w *clientImplWriter) channel_def(def *model.Definition, m *model.Method) error {
	name := clientChannelImpl_name(m)

	w.Linef(`// %v`, name)
	w.Line()
	w.Linef(`type %v struct {`, name)
	w.Line(`ch rpc.Channel`)
	w.Line(`}`)
	w.Line()
	w.Linef(`func new%v(ch rpc.Channel) *%v {`, strings.Title(name), name)
	w.Linef(`return &%v{ch: ch}`, name)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *clientImplWriter) channel_send(def *model.Definition, m *model.Method) error {
	in := m.Channel.In
	if in == nil {
		return nil
	}

	name := clientChannelImpl_name(m)
	typeName := typeName(in)

	// Send
	w.Linef(`func (c *%v) Send(ctx async.Context, msg %v) status.Status {`, name, typeName)
	switch in.Kind {
	case model.KindList, model.KindMessage:
		w.Line(`return c.ch.Send(ctx, msg.Unwrap().Raw())`)

	case model.KindStruct:
		writeFunc := typeWriteFunc(in)
		w.Line(`buf := alloc.AcquireBuffer()`)
		w.Line(`defer buf.Free()`)
		w.Linef(`if _, err := %v(buf, msg); err != nil {`, writeFunc)
		w.Line(`return status.WrapError(err)`)
		w.Line(`}`)
		w.Line(`return c.ch.Send(ctx, buf.Bytes())`)

	default:
		w.Line(`return c.ch.Send(ctx, msg)`)
	}
	w.Line(`}`)
	w.Line()

	// SendEnd
	w.Linef(`func (c *%v) SendEnd(ctx async.Context) status.Status {`, name)
	w.Line(`return c.ch.SendEnd(ctx)`)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *clientImplWriter) channel_receive(def *model.Definition, m *model.Method) error {
	out := m.Channel.Out
	if out == nil {
		return nil
	}

	name := clientChannelImpl_name(m)
	typeName := typeName(out)
	parseFunc := typeParseFunc(out)

	// Receive
	w.Linef(`func (c *%v) Receive(ctx async.Context) (%v, status.Status) {`, name, typeName)
	w.Line(`b, st := c.ch.Receive(ctx)`)
	w.Line(`if !st.OK() {`)
	w.Linef(`return %v{}, st`, typeName)
	w.Line(`}`)
	w.Linef(`msg, _, err := %v(b)`, parseFunc)
	w.Line(`if err != nil {`)
	w.Linef(`return %v{}, status.WrapError(err)`, typeName)
	w.Line(`}`)
	w.Line(`return msg, status.OK`)
	w.Line(`}`)
	w.Line()

	// ReceiveAsync
	w.Linef(`func (c *%v) ReceiveAsync(ctx async.Context) (%v, bool, status.Status) {`, name, typeName)
	w.Line(`b, ok, st := c.ch.ReceiveAsync(ctx)`)
	w.Line(`switch {`)
	w.Line(`case !st.OK():`)
	w.Linef(`return %v{}, false, st`, typeName)
	w.Line(`case !ok:`)
	w.Linef(`return %v{}, false, status.OK`, typeName)
	w.Line(`}`)
	w.Linef(`msg, _, err := %v(b)`, parseFunc)
	w.Line(`if err != nil {`)
	w.Linef(`return %v{}, false, status.WrapError(err)`, typeName)
	w.Line(`}`)
	w.Line(`return msg, true, status.OK`)
	w.Line(`}`)
	w.Line()

	// ReceiveWait
	w.Linef(`func (c *%v) ReceiveWait() <-chan struct{} {`, name)
	w.Line(`return c.ch.ReceiveWait()`)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *clientImplWriter) channel_response(def *model.Definition, m *model.Method) error {
	if err := w.channel_response_def(m); err != nil {
		return err
	}
	if err := w.channel_response_receive(m); err != nil {
		return err
	}
	if err := w.channel_response_parse(m); err != nil {
		return err
	}
	return nil
}

func (w *clientImplWriter) channel_response_def(m *model.Method) error {
	// Response method
	name := clientChannelImpl_name(m)
	w.Writef(`func (c *%v) Response(ctx async.Context) `, name)

	switch {
	default:
		w.Write(`(_st status.Status)`)

	case m.Response != nil:
		typeName := typeName(m.Response)
		w.Writef(`(_ %v, _st status.Status)`, typeName)
	}

	w.Line(`{`)
	return nil
}

func (w *clientImplWriter) channel_response_receive(m *model.Method) error {
	// Receive response
	w.Line(`// Receive response`)
	w.Line(`resp, st := c.ch.Response(ctx)`)
	w.Line(`if !st.OK() {`)
	w.Line(`_st = st`)
	w.Line(`return`)
	w.Line(`}`)
	w.Line(``)
	return nil
}

func (w *clientImplWriter) channel_response_parse(m *model.Method) error {
	// Parse results
	switch {
	default:
		w.Line(`_ = resp`)
		w.Line(`return status.OK`)

	case m.Response != nil:
		parseFunc := typeParseFunc(m.Response)
		w.Line(`// Parse result`)
		w.Linef(`result, _, err := %v(resp)`, parseFunc)
		w.Line(`if err != nil {`)
		w.Line(`_st = status.WrapError(err)`)
		w.Line(`return`)
		w.Line(`}`)
		w.Line(`return result, status.OK`)
	}

	w.Line(`}`)
	w.Line()
	return nil
}

func (w *clientImplWriter) channel_free(def *model.Definition, m *model.Method) error {
	name := clientChannelImpl_name(m)
	w.Linef(`func (c *%v) Free() {`, name)
	w.Line(`c.ch.Free()`)
	w.Line(`}`)
	w.Line()
	return nil
}

// util

func clientImplName(def *model.Definition) string {
	if def.Service.Sub {
		return fmt.Sprintf("%vCall", toLowerCameCase(def.Name))
	}
	return fmt.Sprintf("%vClient", toLowerCameCase(def.Name))
}

func clientImplNew(typ *model.Type) string {
	var name string
	if typ.Ref.Service.Sub {
		name = fmt.Sprintf("New%vCall", typ.Name)
	} else {
		name = fmt.Sprintf("New%vClient", typ.Name)
	}

	if typ.Import != nil {
		return fmt.Sprintf("%v.%v", typ.ImportName, name)
	}
	return name
}

func clientImplNewErr(typ *model.Type) string {
	var name string
	if typ.Ref.Service.Sub {
		name = fmt.Sprintf("New%vCallErr", typ.Name)
	} else {
		name = fmt.Sprintf("New%vClientErr", typ.Name)
	}

	if typ.Import != nil {
		return fmt.Sprintf("%v.%v", typ.ImportName, name)
	}
	return name
}

func clientChannelImpl_name(m *model.Method) string {
	return fmt.Sprintf("%v%vClientChannel", toLowerCameCase(m.Service.Def.Name), toUpperCamelCase(m.Name))
}
