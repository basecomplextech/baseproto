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

type serviceImplWriter struct {
	writer.Writer
}

func newServiceImplWriter(w writer.Writer) *serviceImplWriter {
	return &serviceImplWriter{w}
}

func (w *serviceImplWriter) write(def *model.Definition) error {
	if err := w.def(def); err != nil {
		return err
	}
	if err := w.free(def); err != nil {
		return err
	}
	if err := w.result(def); err != nil {
		return err
	}
	if err := w.handle(def); err != nil {
		return err
	}
	if err := w.methods(def); err != nil {
		return err
	}
	if err := w.channels(def); err != nil {
		return err
	}
	return nil
}

func (w *serviceImplWriter) def(def *model.Definition) error {
	name := handler_name(def)
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
		w.Linef(`ctx rpc.Context`)
		w.Linef(`channel rpc.ServerChannel`)
		w.Linef(`index int`)
		w.Linef(`service %v`, def.Name)
		w.Linef(`result ref.R[[]byte]`)
		w.Line(`}`)
		w.Line()
		w.Linef(`func new%vHandler(ctx rpc.Context, channel rpc.ServerChannel, index int) rpc.Subhandler1[%v] {`,
			def.Name, def.Name)
		w.Linef(`h := %vPool.New()`, name)
		w.Line(`h.ctx = ctx`)
		w.Line(`h.channel = channel`)
		w.Line(`h.index = index`)
		w.Line(`return h`)
		w.Line(`}`)
		w.Line()

	} else {
		w.Linef(`type %v struct {`, name)
		w.Linef(`service %v`, def.Name)
		w.Line(`}`)
		w.Line()
	}

	return nil
}

func (w *serviceImplWriter) free(def *model.Definition) error {
	if !def.Service.Sub {
		return nil
	}

	name := handler_name(def)
	w.Linef(`func (h *%v) Free() {`, name)
	w.Linef(`*h = %v{}`, name)
	w.Linef(`%vPool.Put(h)`, name)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *serviceImplWriter) result(def *model.Definition) error {
	if !def.Service.Sub {
		return nil
	}

	name := handler_name(def)
	w.Linef(`func (h *%v) Result() ref.R[[]byte] {`, name)
	w.Line(`return h.result`)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *serviceImplWriter) handle(def *model.Definition) error {
	name := handler_name(def)

	if def.Service.Sub {
		w.Linef(`func (h *%v) Handle(service %v) status.Status {`, name, def.Name)
		w.Line(`ctx := h.ctx`)
		w.Line(`ch := h.channel`)
		w.Line(`index := h.index`)
		w.Line(`h.service = service`)
		w.Line()
	} else {
		w.Linef(`func (h *%v) Handle(ctx rpc.Context, ch rpc.ServerChannel) (ref.R[[]byte], status.Status) {`,
			name)
		w.Line(`index := 0`)
	}

	w.Line(`req, st := ch.Request(ctx)`)
	w.Line(`if !st.OK() {`)
	if def.Service.Sub {
		w.Line(`return st`)
	} else {
		w.Line(`return nil, st`)
	}
	w.Line(`}`)
	w.Line()

	w.Line(`call, err := req.Calls().GetErr(index)`)
	w.Line(`if err != nil {`)
	if def.Service.Sub {
		w.Line(`return rpc.WrapError(err)`)
	} else {
		w.Line(`return nil, rpc.WrapError(err)`)
	}
	w.Line(`}`)
	w.Line()

	w.Line(`method := call.Method()`)
	w.Line(`switch method {`)
	for _, m := range def.Service.Methods {
		w.Linef(`case %q:`, m.Name)
		if def.Service.Sub {
			w.Linef(`h.result, st = h._%v(ctx, ch, call, index)`, toLowerCameCase(m.Name))
			w.Line(`return st`)
		} else {
			w.Linef(`return h._%v(ctx, ch, call, index)`, toLowerCameCase(m.Name))
		}
	}
	w.Line(`}`)
	w.Line()

	if def.Service.Sub {
		w.Linef(`return rpc.Errorf("unknown %v method %%q", method)`, def.Name)
	} else {
		w.Linef(`return nil, rpc.Errorf("unknown %v method %%q", method)`, def.Name)
	}
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *serviceImplWriter) methods(def *model.Definition) error {
	for _, m := range def.Service.Methods {
		if err := w.method(def, m); err != nil {
			return err
		}
	}
	return nil
}

func (w *serviceImplWriter) method(def *model.Definition, m *model.Method) error {
	// Declare method
	name := handler_name(def)
	w.Linef(`func (h *%v) _%v(ctx rpc.Context, ch rpc.ServerChannel, call prpc.Call, index int) (`,
		name, toLowerCameCase(m.Name))
	w.Line(`ref.R[[]byte], status.Status) {`)

	// Parse input
	switch {
	case m.Channel != nil:
		channelName := handlerChannel_name(m)
		w.Line(`// Make channel`)
		w.Linef(`ch1 := new%v(ch, call.Input())`, strings.Title(channelName))
		w.Line()

	case m.Request != nil:
		makeFunc := typeMakeMessageFunc(m.Request)
		w.Line(`// Parse input`)
		w.Linef(`in := %v(call.Input())`, makeFunc)
		w.Line()
	}

	// Next handler
	if m.Subservice != nil {
		newFunc := handler_new(m.Subservice)

		w.Line(`// Next handler`)
		w.Linef(`next := %v(ctx, ch, index+1 /* next call */)`, newFunc)
		w.Line(`defer next.Free()`)
		w.Line()
	}

	// Call context
	w.Line(`// Call method`)
	ctx := "ctx"
	if m.Oneway {
		ctx = "ctx1"
		w.Line("ctx1 := ctx.Conn()")
	}

	// Declare result
	switch {
	case m.Oneway:
		w.Write(`_ = `)
	case m.Subservice != nil:
		w.Write(`st := `)
	case m.Response != nil:
		w.Write(`result, st := `)
	default:
		w.Write(`st := `)
	}

	// Call method
	switch {
	case m.Channel != nil:
		w.Linef(`h.service.%v(%v, ch1)`, toUpperCamelCase(m.Name), ctx)
	case m.Request != nil:
		w.Writef(`h.service.%v(%v, in`, toUpperCamelCase(m.Name), ctx)
		if m.Subservice != nil {
			w.Write(`, next`)
		}
		w.Line(`)`)
	default:
		w.Writef(`h.service.%v(%v`, toUpperCamelCase(m.Name), ctx)
		if m.Subservice != nil {
			w.Write(`, next`)
		}
		w.Line(`)`)
	}

	// Handle output
	switch {
	case m.Oneway:
		w.Line(`return nil, rpc.SkipResponse`)

	case m.Subservice != nil:
		w.Line(`return next.Result(), st`)

	case m.Response != nil:
		w.Line(`if result != nil { `)
		w.Line(`defer result.Release() `)
		w.Line(`}`)
		w.Line(`if !st.OK() {`)
		w.Line(`return nil, st`)
		w.Line(`}`)
		w.Line()
		w.Line(`// Return bytes`)
		w.Line(`bytes := result.Unwrap().Unwrap().Raw()`)
		w.Line(`return ref.NextRetain(bytes, result), status.OK`)

	default:
		w.Line(`return nil, st`)
	}

	w.Line(`}`)
	w.Line()
	return nil
}

// channels

func (w *serviceImplWriter) channels(def *model.Definition) error {
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

func (w *serviceImplWriter) channel(def *model.Definition, m *model.Method) error {
	if err := w.channel_def(def, m); err != nil {
		return err
	}
	if err := w.channel_request(def, m); err != nil {
		return err
	}
	if err := w.channel_receive(def, m); err != nil {
		return err
	}
	if err := w.channel_send(def, m); err != nil {
		return err
	}
	return nil
}

func (w *serviceImplWriter) channel_def(def *model.Definition, m *model.Method) error {
	name := handlerChannel_name(m)

	w.Linef(`// %v`, name)
	w.Line()
	w.Linef(`type %v struct {`, name)
	w.Line(`ch rpc.ServerChannel`)
	w.Line(`req baseproto.Message`)
	w.Line(`}`)
	w.Line()
	w.Linef(`func new%v(ch rpc.ServerChannel, req baseproto.Message) *%v {`, strings.Title(name), name)
	w.Linef(`return &%v{ch: ch, req: req}`, name)
	w.Linef(`}`)
	w.Line()
	return nil
}

func (w *serviceImplWriter) channel_request(def *model.Definition, m *model.Method) error {
	name := handlerChannel_name(m)

	switch {
	case m.Request != nil:
		typeName := typeName(m.Request)
		makeFunc := typeMakeMessageFunc(m.Request)

		w.Linef(`func (c *%v) Request() (%v, status.Status) {`, name, typeName)
		w.Linef(`req := %v(c.req)`, makeFunc)
		w.Line(`c.req = baseproto.Message{}`)
		w.Line(`return req, status.OK`)
		w.Line(`}`)
		w.Line()
	}
	return nil
}

func (w *serviceImplWriter) channel_receive(def *model.Definition, m *model.Method) error {
	in := m.Channel.In
	if in == nil {
		return nil
	}

	name := handlerChannel_name(m)
	typeName := typeName(in)
	parseFunc := typeParseFunc(in)

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

func (w *serviceImplWriter) channel_send(def *model.Definition, m *model.Method) error {
	out := m.Channel.Out
	if out == nil {
		return nil
	}

	name := handlerChannel_name(m)
	typeName := typeName(out)

	// Send
	w.Linef(`func (c *%v) Send(ctx async.Context, msg %v) status.Status {`, name, typeName)
	switch out.Kind {
	case model.KindList, model.KindMessage:
		w.Line(`return c.ch.Send(ctx, msg.Unwrap().Raw())`)

	case model.KindStruct:
		writeFunc := typeWriteFunc(out)
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

func handler_name(def *model.Definition) string {
	return fmt.Sprintf(`%vHandler`, toLowerCameCase(def.Name))
}

func handler_new(typ *model.Type) string {
	if typ.Import != nil {
		return fmt.Sprintf(`%v.New%vHandler`, typ.ImportName, typ.Name)
	}
	return fmt.Sprintf(`New%vHandler`, typ.Name)
}

func handlerChannel_name(m *model.Method) string {
	return fmt.Sprintf("%v%vChannel", toLowerCameCase(m.Service.Def.Name), toUpperCamelCase(m.Name))
}
