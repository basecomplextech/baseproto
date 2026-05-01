// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package encode

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/basecomplextech/baselibrary/buffer"
	"github.com/basecomplextech/baseproto/internal/format"
)

func EncodeString(b buffer.Buffer, s string) (int, error) {
	size := len(s)
	if size > format.MaxSize {
		return 0, fmt.Errorf("encode: string too large, max size=%d, actual size=%d",
			format.MaxSize, size)
	}

	n := size + 1 // plus zero byte
	p := b.Grow(n)
	copy(p, s)

	n += encodeSizeKind(b, uint32(size), format.KindString)
	return n, nil
}

func EncodeStringBytes(b buffer.Buffer, s []byte) (int, error) {
	size := len(s)
	if size > format.MaxSize {
		return 0, fmt.Errorf("encode: string too large, max size=%d, actual size=%d",
			format.MaxSize, size)
	}

	n := size + 1 // plus zero byte
	p := b.Grow(n)
	copy(p, s)

	n += encodeSizeKind(b, uint32(size), format.KindString)
	return n, nil
}

// Decode

func DecodeString(b []byte) (_ format.String, size int, err error) {
	if len(b) == 0 {
		return "", 0, nil
	}

	// format.Kind
	kind, n := decodeKind(b)
	if n < 0 {
		err = errors.New("decode string: invalid data")
		return
	}
	if kind != format.KindString {
		err = fmt.Errorf("decode string: invalid kind, kind=%v", kind)
		return
	}

	size = n
	end := len(b) - size

	// Size
	dataSize, n := decodeSize(b[:end])
	if n < 0 {
		err = fmt.Errorf("decode string: invalid data size")
		return
	}
	size += n + 1
	end -= (n + 1) // null terminator

	// Data
	data, err := decodeStringData(b[:end], dataSize)
	if err != nil {
		return
	}

	size += int(dataSize)
	return format.String(data), size, nil
}

func DecodeStringClone(b []byte) (_ string, size int, err error) {
	s, size, err := DecodeString(b)
	if err != nil {
		return "", size, err
	}
	return s.Clone(), size, nil
}

// private

func decodeStringData(b []byte, size uint32) (string, error) {
	off := len(b) - int(size)
	if off < 0 {
		return "", errors.New("decode string: invalid data size")
	}

	p := b[off:]
	s := *(*string)(unsafe.Pointer(&p))
	return s, nil
}
