// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package pcommon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpanID(t *testing.T) {
	sid := SpanID([]byte{1, 2, 3, 4, 4, 3, 2, 1})
	assert.Equal(t, []byte{1, 2, 3, 4, 4, 3, 2, 1}, []byte(sid))
	assert.False(t, sid.IsEmpty())
}

func TestNewSpanIDEmpty(t *testing.T) {
	sid := NewSpanIDEmpty()
	assert.Equal(t, []byte{}, []byte(sid))
	assert.True(t, sid.IsEmpty())
}

func TestSpanIDString(t *testing.T) {
	sid := SpanID([]byte{})
	assert.Equal(t, "", sid.String())

	sid = SpanID([]byte{0x12, 0x23, 0xAD, 0x12, 0x23, 0xAD, 0x12, 0x23})
	assert.Equal(t, "1223ad1223ad1223", sid.String())
}

func TestSpanIDImmutable(t *testing.T) {
	initialBytes := []byte{0x12, 0x23, 0xAD, 0x12, 0x23, 0xAD, 0x12, 0x23}
	sid := SpanID(initialBytes)
	assert.Equal(t, SpanID(initialBytes), sid)

	// Get the bytes and try to mutate.
	sid[4] = 0x89

	// Does not change the already created SpanID.
	assert.NotEqual(t, SpanID(initialBytes), sid)
}
