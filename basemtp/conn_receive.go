// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package basemtp

import (
	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/status"
	"github.com/basecomplextech/baseproto/proto/pmtp"
)

func (c *conn) receiveLoop(ctx async.Context) status.Status {
	for {
		msg, st := c.reader.readMessage()
		if !st.OK() {
			return st
		}

		if st := c.receiveMessage(msg, false /* not inside batch */); !st.OK() {
			return st
		}
	}
}

func (c *conn) receiveMessage(msg pmtp.Message, insideBatch bool) status.Status {
	code := msg.Code()

	switch code {
	case pmtp.Code_Batch:
		if insideBatch {
			return mtpErrorf("received nested batch messages")
		}
		return c.receiveBatch(msg)
	case pmtp.Code_ChannelOpen:
		return c.receiveOpen(msg)
	case pmtp.Code_ChannelClose:
		return c.receiveClose(msg)
	case pmtp.Code_ChannelData:
		return c.receiveData(msg)
	case pmtp.Code_ChannelWindow:
		return c.receiveWindow(msg)
	}

	return mtpErrorf("unexpected message, code=%v", code)
}

func (c *conn) receiveOpen(msg pmtp.Message) status.Status {
	m := msg.ChannelOpen()
	id := m.Id()

	// Add channel
	// Duplicates are impossible, but still check for them.
	ch := openChannel(c, c.client, m)
	_, exists := c.channels.GetOrSet(id, ch)
	if exists {
		ch.Free()
		ch.free()
		return mtpErrorf("received open message for existing channel, channel=%v", id)
	}

	// Start handler
	h := newChannelHandler(c, ch)
	workerPool.Run(h)

	c.maybeChannelsReached()
	return status.OK
}

func (c *conn) receiveClose(msg pmtp.Message) status.Status {
	m := msg.ChannelClose()
	id := m.Id()

	ch, ok := c.channels.Delete(id)
	if !ok {
		return status.OK
	}
	defer ch.free()

	return ch.receive(msg)
}

func (c *conn) receiveData(msg pmtp.Message) status.Status {
	m := msg.ChannelData()
	id := m.Id()

	ch, ok := c.channels.Get(id)
	if !ok {
		return status.OK
	}
	return ch.receive(msg)
}

func (c *conn) receiveWindow(msg pmtp.Message) status.Status {
	m := msg.ChannelWindow()
	id := m.Id()

	ch, ok := c.channels.Get(id)
	if !ok {
		return status.OK
	}
	return ch.receive(msg)
}

func (c *conn) receiveBatch(msg pmtp.Message) status.Status {
	batch := msg.Batch()
	list := batch.List()
	num := list.Len()

	for i := 0; i < num; i++ {
		m1, err := list.GetErr(i)
		if err != nil {
			return status.WrapError(err)
		}

		if st := c.receiveMessage(m1, true /* inside batch */); !st.OK() {
			return st
		}
	}
	return status.OK
}
