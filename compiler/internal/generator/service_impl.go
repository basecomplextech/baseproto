// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package generator

import (
	"strings"

	"github.com/basecomplextech/baseproto/compiler/internal/golang"
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
	srv, err := golang.NewService(def)
	if err != nil {
		return err
	}

	if err := w.def(srv); err != nil {
		return err
	}
	if err := w.free(srv); err != nil {
		return err
	}
	if err := w.result(srv); err != nil {
		return err
	}
	if err := w.handle(srv); err != nil {
		return err
	}
	if err := w.methods(srv); err != nil {
		return err
	}
	if err := w.channels(srv); err != nil {
		return err
	}
	return nil
}

func (w *serviceImplWriter) def(srv *golang.Service) error {
	handler := srv.HandlerImpl

	w.Linef(`// %v`, handler)
	w.Line()

	if srv.Sub {
		w.Linef(`var %vPool = pools.NewPoolFunc(`, handler)
		w.Linef(`func() *%v {`, handler)
		w.Linef(`return &%v{}`, handler)
		w.Line(`},`)
		w.Line(`)`)
		w.Line()
		w.Linef(`type %v struct {`, handler)
		w.Linef(`ctx baserpc.Context`)
		w.Linef(`channel baserpc.ServerChannel`)
		w.Linef(`index int`)
		w.Linef(`service %v`, srv.Name)
		w.Linef(`result ref.R[[]byte]`)
		w.Line(`}`)
		w.Line()
		w.Linef(`func new%v(ctx baserpc.Context, ch baserpc.ServerChannel, `, srv.Handler)
		w.Linef(`index int) baserpc.Subhandler1[%v] {`, srv.Name)
		w.Linef(`h := %vPool.New()`, handler)
		w.Line(`h.ctx = ctx`)
		w.Line(`h.channel = ch`)
		w.Line(`h.index = index`)
		w.Line(`return h`)
		w.Line(`}`)
		w.Line()

	} else {
		w.Linef(`type %v struct {`, handler)
		w.Linef(`service %v`, srv.Name)
		w.Line(`}`)
		w.Line()
	}

	return nil
}

func (w *serviceImplWriter) free(srv *golang.Service) error {
	if !srv.Sub {
		return nil
	}

	handler := srv.HandlerImpl
	w.Linef(`func (h *%v) Free() {`, handler)
	w.Linef(`*h = %v{}`, handler)
	w.Linef(`%vPool.Put(h)`, handler)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *serviceImplWriter) result(srv *golang.Service) error {
	if !srv.Sub {
		return nil
	}

	handler := srv.HandlerImpl
	w.Linef(`func (h *%v) Result() ref.R[[]byte] {`, handler)
	w.Line(`return h.result`)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *serviceImplWriter) handle(srv *golang.Service) error {
	handler := srv.HandlerImpl

	if srv.Sub {
		w.Linef(`func (h *%v) Handle(service %v) status.Status {`, handler, srv.Name)
		w.Line(`ctx := h.ctx`)
		w.Line(`ch := h.channel`)
		w.Line(`index := h.index`)
		w.Line(`h.service = service`)
		w.Line()
	} else {
		w.Linef(`func (h *%v) Handle(ctx baserpc.Context, ch baserpc.ServerChannel) (`, handler)
		w.Line(`ref.R[[]byte], status.Status) {`)
		w.Line(`index := 0`)
	}

	w.Line(`req, st := ch.Request(ctx)`)
	w.Line(`if !st.OK() {`)
	if srv.Sub {
		w.Line(`return st`)
	} else {
		w.Line(`return nil, st`)
	}
	w.Line(`}`)
	w.Line()

	w.Line(`call, err := req.Calls().GetErr(index)`)
	w.Line(`if err != nil {`)
	if srv.Sub {
		w.Line(`return baserpc.WrapError(err)`)
	} else {
		w.Line(`return nil, baserpc.WrapError(err)`)
	}
	w.Line(`}`)
	w.Line()

	w.Line(`method := call.Method()`)
	w.Line(`switch method {`)
	for _, m := range srv.Methods {
		w.Linef(`case %q:`, m.Path)
		if srv.Sub {
			w.Linef(`h.result, st = h._%v(ctx, ch, call, index)`, m.Private)
			w.Line(`return st`)
		} else {
			w.Linef(`return h._%v(ctx, ch, call, index)`, m.Private)
		}
	}
	w.Line(`}`)
	w.Line()

	if srv.Sub {
		w.Linef(`return baserpc.Errorf("unknown %v method %%q", method)`, srv.Name)
	} else {
		w.Linef(`return nil, baserpc.Errorf("unknown %v method %%q", method)`, srv.Name)
	}
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *serviceImplWriter) methods(srv *golang.Service) error {
	for _, m := range srv.Methods {
		if err := w.method(srv, m); err != nil {
			return err
		}
	}
	return nil
}

func (w *serviceImplWriter) method(srv *golang.Service, m *golang.Method) error {
	handler := srv.HandlerImpl

	// Declare method
	w.Linef(`func (h *%v) _%v(ctx baserpc.Context, ch baserpc.ServerChannel, `, handler, m.Private)
	w.Line(`call prpc.Call, index int) (ref.R[[]byte], status.Status) {`)

	// Parse input
	switch {
	case m.Channel != nil:
		channelHandler := m.Channel.Handler
		w.Line(`// Make channel`)
		w.Linef(`ch1 := new%v(ch, call.Input())`, strings.Title(channelHandler))
		w.Line()

	case m.Request != nil:
		newFunc := m.Request.NewFunc()
		w.Line(`// Parse input`)
		w.Linef(`in := %v(call.Input())`, newFunc)
		w.Line()
	}

	// Next handler
	if m.Subservice != nil {
		newFunc := m.Subservice.NewHandlerFunc()

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
		w.Linef(`h.service.%v(%v, ch1)`, m.Public, ctx)
	case m.Request != nil:
		w.Writef(`h.service.%v(%v, in`, m.Public, ctx)
		if m.Subservice != nil {
			w.Write(`, next`)
		}
		w.Line(`)`)
	default:
		w.Writef(`h.service.%v(%v`, m.Public, ctx)
		if m.Subservice != nil {
			w.Write(`, next`)
		}
		w.Line(`)`)
	}

	// Handle output
	switch {
	case m.Oneway:
		w.Line(`return nil, baserpc.SkipResponse`)

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

func (w *serviceImplWriter) channels(srv *golang.Service) error {
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

func (w *serviceImplWriter) channel(srv *golang.Service, m *golang.Method) error {
	if err := w.channelDef(srv, m); err != nil {
		return err
	}
	if err := w.channelRequest(srv, m); err != nil {
		return err
	}
	if err := w.channelReceive(srv, m); err != nil {
		return err
	}
	if err := w.channelSend(srv, m); err != nil {
		return err
	}
	return nil
}

func (w *serviceImplWriter) channelDef(srv *golang.Service, m *golang.Method) error {
	handler := m.Channel.Handler

	w.Linef(`// %v`, handler)
	w.Line()
	w.Linef(`type %v struct {`, handler)
	w.Line(`ch baserpc.ServerChannel`)
	w.Line(`req baseproto.Message`)
	w.Line(`}`)
	w.Line()
	w.Linef(`func new%v(ch baserpc.ServerChannel, req baseproto.Message) *%v {`,
		strings.Title(handler), handler)
	w.Linef(`return &%v{ch: ch, req: req}`, handler)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *serviceImplWriter) channelRequest(srv *golang.Service, m *golang.Method) error {
	handler := m.Channel.Handler

	switch {
	case m.Request != nil:
		typeName := m.Request.Name()
		newFunc := m.Request.NewFunc()

		w.Linef(`func (c *%v) Request() (%v, status.Status) {`, handler, typeName)
		w.Linef(`req := %v(c.req)`, newFunc)
		w.Line(`c.req = baseproto.Message{}`)
		w.Line(`return req, status.OK`)
		w.Line(`}`)
		w.Line()
	}
	return nil
}

func (w *serviceImplWriter) channelReceive(srv *golang.Service, m *golang.Method) error {
	in := m.Channel.In
	if in == nil {
		return nil
	}

	handler := m.Channel.Handler
	typeName := in.Name()
	parseFunc := in.ParseFunc()

	// Receive
	w.Linef(`func (c *%v) Receive(ctx async.Context) (%v, status.Status) {`, handler, typeName)
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
	w.Linef(`func (c *%v) ReceiveAsync(ctx async.Context) (%v, bool, status.Status) {`,
		handler, typeName)
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
	w.Linef(`func (c *%v) ReceiveWait() <-chan struct{} {`, handler)
	w.Line(`return c.ch.ReceiveWait()`)
	w.Line(`}`)
	w.Line()
	return nil
}

func (w *serviceImplWriter) channelSend(srv *golang.Service, m *golang.Method) error {
	out := m.Channel.Out
	if out == nil {
		return nil
	}

	name := m.Channel.Handler
	typeName := out.InputName()

	// Send
	w.Linef(`func (c *%v) Send(ctx async.Context, msg %v) status.Status {`, name, typeName)
	switch out := out.(type) {
	case golang.ListType, golang.MessageType:
		w.Line(`return c.ch.Send(ctx, msg.Unwrap().Raw())`)

	case golang.StructType:
		encodeFunc := out.EncodeFunc()
		w.Line(`buf := alloc.AcquireBuffer()`)
		w.Line(`defer buf.Free()`)
		w.Linef(`if _, err := %v(buf, msg); err != nil {`, encodeFunc)
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
