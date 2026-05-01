// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package encode

import (
	"errors"
	"fmt"

	"github.com/basecomplextech/baselibrary/buffer"
	"github.com/basecomplextech/baseproto/internal/format"
)

func EncodeBytes(b buffer.Buffer, v []byte) (int, error) {
	size := len(v)
	if size > format.MaxSize {
		return 0, fmt.Errorf("encode: bytes too large, max size=%d, actual size=%d",
			format.MaxSize, size)
	}

	p := b.Grow(size)
	copy(p, v)
	n := size

	n += encodeSizeKind(b, uint32(size), format.KindBytes)
	return n, nil
}

// Decode

func DecodeBytes(b []byte) (_ format.Bytes, size int, err error) {
	if len(b) == 0 {
		return nil, 0, nil
	}

	// format.Kind
	kind, n := decodeKind(b)
	if n < 0 {
		err = errors.New("decode bytes: invalid data")
		return
	}
	if kind != format.KindBytes {
		err = fmt.Errorf("decode bytes: invalid kind, kind=%v", kind)
		return
	}

	size = n
	end := len(b) - size

	// Data size
	dataSize, n := decodeSize(b[:end])
	if n < 0 {
		err = errors.New("decode bytes: invalid data size")
		return
	}
	size += n
	end -= n

	// Data
	data, err := decodeBytesData(b[:end], dataSize)
	if err != nil {
		return nil, 0, err
	}

	size += int(dataSize)
	return data, size, nil
}

// private

func decodeBytesData(b []byte, size uint32) ([]byte, error) {
	off := len(b) - int(size)
	if off < 0 {
		return nil, errors.New("decode bytes: invalid data size")
	}
	return b[off:], nil
}
