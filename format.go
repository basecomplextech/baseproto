// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package baseproto

import (
	"github.com/basecomplextech/baseproto/internal/format"
)

// Kind specifies a value kind.
type Kind = format.Kind

const (
	KindUndefined = format.KindUndefined

	KindTrue  = format.KindTrue
	KindFalse = format.KindFalse
	KindByte  = format.KindByte

	KindInt16 = format.KindInt16
	KindInt32 = format.KindInt32
	KindInt64 = format.KindInt64

	KindUint16 = format.KindUint16
	KindUint32 = format.KindUint32
	KindUint64 = format.KindUint64

	KindFloat32 = format.KindFloat32
	KindFloat64 = format.KindFloat64

	KindBin64  = format.KindBin64
	KindBin128 = format.KindBin128
	KindBin256 = format.KindBin256

	KindBytes  = format.KindBytes
	KindString = format.KindString

	KindList    = format.KindList
	KindListBig = format.KindListBig

	KindMessage    = format.KindMessage
	KindMessageBig = format.KindMessageBig

	KindStruct = format.KindStruct
)

type (
	// Bytes is a baseproto byte slice backed by a buffer.
	// Clone it if you need to keep it around.
	Bytes = format.Bytes

	// String is a baseproto string backed by a buffer.
	// Clone it if you need to keep it around.
	String = format.String

	// ListTable is a serialized array of list element offsets ordered by index.
	ListTable = format.ListTable

	// MessageTable is a table of message fields ordered by tags.
	MessageTable = format.MessageTable
)
