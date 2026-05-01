// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package encode

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/basecomplextech/baselibrary/buffer"
	"github.com/basecomplextech/baseproto/internal/format"
)

func EncodeMessageTable(b buffer.Buffer, dataSize int, table []format.MessageField) (int, error) {
	if dataSize > format.MaxSize {
		return 0, fmt.Errorf("encode: message too large, max size=%d, actual size=%d",
			format.MaxSize, dataSize)
	}

	// format.Kind
	big := format.IsBigMessage(table)
	kind := format.KindMessage
	if big {
		kind = format.KindMessageBig
	}

	// Write table
	tableSize, err := encodeMessageTable(b, table, big)
	if err != nil {
		return 0, err
	}
	n := tableSize

	// Write data size
	n += encodeSize(b, uint32(dataSize))

	// Write table size and
	n += encodeSizeKind(b, uint32(tableSize), kind)
	return n, nil
}

func DecodeMessageTable(b []byte) (_ format.MessageTable, size int, err error) {
	if len(b) == 0 {
		return
	}

	// Decode type
	kind, n := decodeKind(b)
	if n < 0 {
		err = errors.New("decode message: invalid type")
		return
	}
	switch kind {
	case format.KindMessage, format.KindMessageBig:
	default:
		err = fmt.Errorf("decode message: invalid kind, kind=%v", kind)
		return
	}

	// Start
	size = n
	end := len(b) - size
	big := kind == format.KindMessageBig

	// Table size
	tableSize, m := decodeSize(b[:end])
	if m < 0 {
		err = errors.New("decode message: invalid table size")
		return
	}
	end -= m
	size += m

	// Data size
	dataSize, m := decodeSize(b[:end])
	if m < 0 {
		err = fmt.Errorf("decode message: invalid data size")
		return
	}
	end -= m
	size += m

	// Table
	table, err := decodeMessageTable(b[:end], tableSize, big)
	if err != nil {
		return
	}
	end -= int(tableSize) + int(dataSize)
	size += int(tableSize)

	// Data
	if end < 0 {
		err = errors.New("decode message: invalid data")
		return
	}
	size += int(dataSize)

	// Done
	t := format.NewMessageTable(table, dataSize, big)
	return t, size, nil
}

// private

func encodeMessageTable(b buffer.Buffer, table []format.MessageField, big bool) (int, error) {
	// Field size
	var fieldSize int
	if big {
		fieldSize = format.MessageFieldSize_Big
	} else {
		fieldSize = format.MessageFieldSize_Small
	}

	// Check table size
	size := len(table) * fieldSize
	if size > format.MaxSize {
		return 0, fmt.Errorf("encode: message table too large, max size=%d, actual size=%d",
			format.MaxSize, size)
	}

	// Write table
	p := b.Grow(size)
	off := 0

	// Put fields
	for _, field := range table {
		q := p[off : off+fieldSize]

		if big {
			binary.BigEndian.PutUint16(q, field.Tag)
			binary.BigEndian.PutUint32(q[2:], field.Offset)
		} else {
			q[0] = byte(field.Tag)
			binary.BigEndian.PutUint16(q[1:], uint16(field.Offset))
		}

		off += fieldSize
	}

	return size, nil
}

func decodeMessageTable(b []byte, size uint32, big bool) (_ []byte, err error) {
	// Field size
	fieldSize := format.MessageFieldSize_Small
	if big {
		fieldSize = format.MessageFieldSize_Big
	}

	// Check offset
	start := len(b) - int(size)
	if start < 0 {
		err = errors.New("decode message: invalid table")
		return
	}

	// Check divisible
	if size%uint32(fieldSize) != 0 {
		err = errors.New("decode message: invalid table")
		return
	}

	table := b[start:]
	return table, nil
}
