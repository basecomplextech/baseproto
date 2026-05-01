// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package encode

import (
	"testing"

	"github.com/basecomplextech/baselibrary/buffer"
	"github.com/basecomplextech/baselibrary/tests"
	"github.com/basecomplextech/baseproto/internal/format"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testEncodeListTable(t tests.T, dataSize int, elements []format.ListElement) []byte {
	buf := buffer.New()
	buf.Grow(dataSize)

	_, err := EncodeListTable(buf, dataSize, elements)
	if err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// List

func TestListTable__should_encode_decode_list(t *testing.T) {
	elems := format.TestElements()
	dataSize := 100
	b := testEncodeListTable(t, dataSize, elems)

	table, n, err := DecodeListTable(b)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(b), n)
	assert.Equal(t, uint32(dataSize), table.DataSize())
	assert.Equal(t, len(elems), table.Len())

	kind, size, err := DecodeKindSize(b)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, format.KindList, kind)
	assert.Equal(t, size, len(b))
}

func TestListTable__should_encode_decode_list_table(t *testing.T) {
	elems := format.TestElements()

	for i := 0; i <= len(elems); i++ {
		b := buffer.New()
		elems1 := elems[i:]

		_, err := EncodeListTable(b, 0, elems1)
		if err != nil {
			t.Fatal(err)
		}
		p := b.Bytes()

		table1, _, err := DecodeListTable(p)
		if err != nil {
			t.Fatal(err)
		}

		elems2 := table1.Elements()
		require.Equal(t, elems1, elems2)
	}
}

func TestListTable__should_return_error_when_invalid_kind(t *testing.T) {
	elems := format.TestElements()
	dataSize := 100

	b := testEncodeListTable(t, dataSize, elems)
	b[len(b)-1] = byte(format.KindMessage)

	_, _, err := DecodeListTable(b)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid kind")
}

func TestListTable__should_return_error_when_invalid_table_size(t *testing.T) {
	b := []byte{}
	b = append(b, 0xff)
	b = append(b, byte(format.KindList))

	_, _, err := DecodeListTable(b)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid table size")
}

func TestListTable__should_return_error_when_invalid_data_size(t *testing.T) {
	big := false
	b := []byte{}
	b = append(b, 0xff)
	b = appendSize(b, big, 1000)
	b = append(b, byte(format.KindList))

	_, _, err := DecodeListTable(b)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid data size")
}

func TestListTable__should_return_error_when_invalid_table(t *testing.T) {
	buf := buffer.New()
	_, err := EncodeListTable(buf, 0, nil) // TODO: big(true)
	if err != nil {
		t.Fatal(err)
	}

	big := false
	b := buf.Bytes()
	b = appendSize(b, big, 0)    // data size
	b = appendSize(b, big, 1000) // table size
	b = append(b, byte(format.KindList))

	_, _, err = DecodeListTable(b)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid table")
}

func TestListTable__should_return_error_when_invalid_data(t *testing.T) {
	buf := buffer.New()
	_, err := EncodeListTable(buf, 0, nil) // TODO: big(true)
	if err != nil {
		t.Fatal(err)
	}

	big := false
	b := buf.Bytes()
	b = appendSize(b, big, 1000) // data size
	b = appendSize(b, big, 0)    // table size
	b = append(b, byte(format.KindList))

	_, _, err = DecodeListTable(b)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid data")
}
