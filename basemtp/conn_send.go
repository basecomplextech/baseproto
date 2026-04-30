// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package basemtp

import (
	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/status"
	"github.com/basecomplextech/baseproto/proto/pmtp"
)

func (c *conn) sendLoop(ctx async.Context) status.Status {
	for {
		// Write pending messages
		b, ok, st := c.writeq.Read()
		switch {
		case !st.OK():
			return st
		case ok:
			if st := c.sendMessage(b); !st.OK() {
				return st
			}
			continue
		}

		// Flush buffered writes
		if st := c.writer.flush(); !st.OK() {
			return st
		}

		// Wait for more messages
		select {
		case <-ctx.Wait():
			return ctx.Status()
		case <-c.writeq.ReadWait():
		}
	}
}

func (c *conn) sendMessage(b []byte) status.Status {
	msg, err := pmtp.OpenMessageErr(b)
	if err != nil {
		return mtpError(err)
	}

	// Handle message
	if st := c.sendHandle(msg); !st.OK() {
		return st
	}

	// Write message
	return c.writer.write(msg)
}

func (c *conn) sendHandle(msg pmtp.Message) status.Status {
	code := msg.Code()

	switch code {
	case pmtp.Code_Batch:
		// Handle batch messages
		batch := msg.Batch()
		list := batch.List()
		num := list.Len()

		for i := 0; i < num; i++ {
			m1, err := list.GetErr(i)
			if err != nil {
				return status.WrapError(err)
			}
			if st := c.sendHandle(m1); !st.OK() {
				return st
			}
		}

	case pmtp.Code_ChannelClose:
		// Remove and free channel
		id := msg.ChannelClose().Id()

		ch, ok := c.channels.Delete(id)
		if ok {
			ch.free()
		}
	}

	return status.OK
}
