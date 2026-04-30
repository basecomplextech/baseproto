// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package basemtp

import (
	"bufio"
	"encoding/binary"
	"io"
	"strings"

	"github.com/basecomplextech/baselibrary/alloc"
	"github.com/basecomplextech/baselibrary/opt"
	"github.com/basecomplextech/baselibrary/status"
	"github.com/basecomplextech/baseproto/proto/pmtp"
	"github.com/pierrec/lz4/v4"
)

type connReader struct {
	src    *bufio.Reader
	comp   opt.Opt[*lz4.Reader] // empty when no compression
	reader io.Reader            // points to src or comp

	client bool
	freed  bool

	head [4]byte
	buf  alloc.Buffer
}

func newConnReader(r io.Reader, client bool, bufferSize int) *connReader {
	src := bufio.NewReaderSize(r, bufferSize)
	return &connReader{
		src:    src,
		client: client,
		reader: src,

		buf: alloc.NewBuffer(),
	}
}

func (r *connReader) free() {
	if r.freed {
		return
	}
	r.freed = true

	r.buf.Free()
	r.buf = nil
}

func (r *connReader) initLZ4() status.Status {
	if r.comp.Valid {
		return status.OK
	}

	reader := lz4.NewReader(r.src)
	r.comp = opt.New(reader)
	r.reader = reader
	return status.OK
}

// readLine reads and returns a single line delimited by \n, includes the delimiter.
func (r *connReader) readLine() (string, status.Status) {
	s, err := r.src.ReadString('\n')
	if err != nil {
		return "", mtpError(err)
	}

	if debug {
		debugPrint(r.client, "<- line\t", strings.TrimSpace(s))
	}
	return s, status.OK
}

// readRequest reads and parses a connect request, the message is valid until the next read call.
func (r *connReader) readRequest() (pmtp.ConnectRequest, status.Status) {
	msg, st := r.readMessage()
	if !st.OK() {
		return pmtp.ConnectRequest{}, st
	}

	code := msg.Code()
	if code != pmtp.Code_ConnectRequest {
		return pmtp.ConnectRequest{}, mtpErrorf(
			"unexpected message, expected connect request, got %v", code)
	}

	req := msg.ConnectRequest()
	return req, status.OK
}

// readResponse reads and parses a connect response, the message is valid until the next read call.
func (r *connReader) readResponse() (pmtp.ConnectResponse, status.Status) {
	msg, st := r.readMessage()
	if !st.OK() {
		return pmtp.ConnectResponse{}, st
	}

	code := msg.Code()
	if code != pmtp.Code_ConnectResponse {
		return pmtp.ConnectResponse{}, mtpErrorf(
			"unexpected message, expected connect response, got %v", code)
	}

	resp := msg.ConnectResponse()
	return resp, status.OK
}

// readMessage reads and parses the next message, the message is valid until the next read call.
func (r *connReader) readMessage() (pmtp.Message, status.Status) {
	buf, st := r.read()
	if !st.OK() {
		return pmtp.Message{}, st
	}

	// Parse message
	msg, _, err := pmtp.ParseMessage(buf)
	if err != nil {
		return pmtp.Message{}, mtpError(err)
	}

	if debug {
		code := msg.Code()
		switch code {
		case pmtp.Code_ConnectRequest:
			debugPrint(r.client, "<- connect_req")

		case pmtp.Code_ConnectResponse:
			debugPrint(r.client, "<- connect_resp")

		case pmtp.Code_Batch:
			m := msg.Batch()
			list := m.List()
			num := list.Len()
			codes := make([]string, 0, num)

			for i := 0; i < num; i++ {
				m1 := list.Get(i)
				c1 := m1.Code().String()
				codes = append(codes, c1)
			}
			debugPrint(r.client, "<- batch\t", num, codes)

		case pmtp.Code_ChannelOpen:
			m := msg.ChannelOpen()
			id := m.Id()
			data := debugString(m.Data())
			cmd := "<- channel_open\t"
			debugPrint(r.client, cmd, id, data)

		case pmtp.Code_ChannelClose:
			m := msg.ChannelClose()
			id := m.Id()
			data := debugString(m.Data())
			debugPrint(r.client, "<- channel_close\t", id, data)

		case pmtp.Code_ChannelData:
			m := msg.ChannelData()
			id := m.Id()
			data := debugString(m.Data())
			debugPrint(r.client, "<- channel_data\t", id, data)

		case pmtp.Code_ChannelWindow:
			m := msg.ChannelWindow()
			id := m.Id()
			delta := m.Delta()
			debugPrint(r.client, "<- channel_window\t", id, delta)

		default:
			debugPrint(r.client, "<- unknown", code)
		}
	}
	return msg, status.OK
}

// read reads the next message bytes, the bytes are valid until the next read call.
func (r *connReader) read() ([]byte, status.Status) {
	head := r.head[:]

	// Read size
	if _, err := io.ReadFull(r.reader, head); err != nil {
		return nil, mtpError(err)
	}
	size := binary.BigEndian.Uint32(head)

	// Read bytes
	r.buf.Reset()
	buf := r.buf.Grow(int(size))
	if _, err := io.ReadFull(r.reader, buf); err != nil {
		return nil, mtpError(err)
	}
	return buf, status.OK
}
