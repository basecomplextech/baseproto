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

func EncodeListTable(b buffer.Buffer, dataSize int, table []format.ListElement) (int, error) {
	if dataSize > format.MaxSize {
		return 0, fmt.Errorf("encode: list too large, max size=%d, actual size=%d",
			format.MaxSize, dataSize)
	}

	// format.Kind
	big := format.IsBigList(table)
	kind := format.KindList
	if big {
		kind = format.KindListBig
	}

	// Write table
	tableSize, err := encodeListTable(b, table, big)
	if err != nil {
		return int(tableSize), err
	}
	n := tableSize

	// Write data size
	n += encodeSize(b, uint32(dataSize))

	// Write table size and kind
	n += encodeSizeKind(b, uint32(tableSize), kind)
	return n, nil
}

// Decode

func DecodeListTable(b []byte) (_ format.ListTable, size int, err error) {
	if len(b) == 0 {
		return
	}

	// Decode type
	kind, n := decodeKind(b)
	if n < 0 {
		n = 0
		err = errors.New("decode list: invalid data")
		return
	}
	if kind != format.KindList && kind != format.KindListBig {
		err = fmt.Errorf("decode list: invalid kind, kind=%v", kind)
		return
	}

	// Start
	size = n
	end := len(b) - n
	big := kind == format.KindListBig

	// Table size
	tableSize, n := decodeSize(b[:end])
	if n < 0 {
		err = errors.New("decode list: invalid table size")
		return
	}
	end -= n
	size += n

	// Data size
	dataSize, n := decodeSize(b[:end])
	if n < 0 {
		err = errors.New("decode list: invalid data size")
		return
	}
	end -= n
	size += n

	// Table
	table, err := decodeListTable(b[:end], tableSize, big)
	if err != nil {
		return
	}
	end -= int(tableSize) + int(dataSize)
	size += int(tableSize)

	// Data
	if end < 0 {
		err = errors.New("decode list: invalid data")
		return
	}
	size += int(dataSize)

	// Done
	t := format.NewListTable(table, dataSize, big)
	return t, size, nil
}

// private

func encodeListTable(b buffer.Buffer, table []format.ListElement, big bool) (int, error) {
	// Element size
	elemSize := format.ListElementSize_Small
	if big {
		elemSize = format.ListElementSize_Big
	}

	// Check table size
	size := len(table) * elemSize
	if size > format.MaxSize {
		return 0, fmt.Errorf("encode: list table too large, max size=%d, actual size=%d",
			format.MaxSize, size)
	}

	// Write table
	p := b.Grow(size)
	off := 0

	// Put elements
	for _, elem := range table {
		q := p[off : off+elemSize]

		if big {
			binary.BigEndian.PutUint32(q, elem.Offset)
		} else {
			binary.BigEndian.PutUint16(q, uint16(elem.Offset))
		}

		off += elemSize
	}

	return size, nil
}

func decodeListTable(b []byte, size uint32, big bool) (_ []byte, err error) {
	// Element size
	elemSize := format.ListElementSize_Small
	if big {
		elemSize = format.ListElementSize_Big
	}

	// Check offset
	start := len(b) - int(size)
	if start < 0 {
		err = errors.New("decode list: invalid table")
		return
	}

	// Check divisible
	if size%uint32(elemSize) != 0 {
		err = errors.New("decode list: invalid table")
		return
	}

	table := b[start:]
	return table, nil
}
