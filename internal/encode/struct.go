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

func EncodeStruct(b buffer.Buffer, dataSize int) (int, error) {
	if dataSize > format.MaxSize {
		return 0, fmt.Errorf("encode: struct too large, max size=%d, actual size=%d",
			format.MaxSize, dataSize)
	}

	n := encodeSizeKind(b, uint32(dataSize), format.KindStruct)
	return n, nil
}

// Decode

func DecodeStruct(b []byte) (dataSize int, size int, err error) {
	if len(b) == 0 {
		return 0, 0, nil
	}

	// Decode type
	kind, n := decodeKind(b)
	if n < 0 {
		err = errors.New("decode struct: invalid type")
		return
	}
	if kind != format.KindStruct {
		err = fmt.Errorf("decode struct: invalid kind, kind=%v", kind)
		return
	}

	size = n
	end := len(b) - size

	// Data size
	dataSize_, n := decodeSize(b[:end])
	if n < 0 {
		err = errors.New("decode struct: invalid data size")
		return
	}
	size += n + int(dataSize_)

	return int(dataSize_), size, nil
}
