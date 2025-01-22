// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package json

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	otlpcommon "go.opentelemetry.io/collector/pdata/internal/data/protogen/common/v1"
)

func ref[T any](v T) *T {
	return &v
}

func TestReadArray(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    *otlpcommon.ArrayValue
	}{
		{
			name: "values",
			jsonStr: `{"values":[{
"stringValue":"12312"
}]}`,
			want: otlpcommon.ArrayValue_builder{
				Values: []*otlpcommon.AnyValue{
					otlpcommon.AnyValue_builder{
						StringValue: ref("12312"),
					}.Build(),
				},
			}.Build(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iter := jsoniter.ConfigFastest.BorrowIterator([]byte(tt.jsonStr))
			defer jsoniter.ConfigFastest.ReturnIterator(iter)
			got := readArray(iter)
			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestReadKvlistValue(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    *otlpcommon.KeyValueList
	}{
		{
			name: "values",
			jsonStr: `{"values":[{
"key":"testKey",
"value":{
"stringValue": "testValue"
}
}]}`,
			want: otlpcommon.KeyValueList_builder{
				Values: []*otlpcommon.KeyValue{
					otlpcommon.KeyValue_builder{
						Key: "testKey",
						Value: otlpcommon.AnyValue_builder{
							StringValue: ref("testValue"),
						}.Build(),
					}.Build(),
				},
			}.Build(),
		},
		{
			name: "boolValue",
			jsonStr: `{"values":[{
"key":"testKey",
"value":{
"boolValue": true
}
}]}`,
			want: otlpcommon.KeyValueList_builder{
				Values: []*otlpcommon.KeyValue{
					otlpcommon.KeyValue_builder{
						Key: "testKey",
						Value: otlpcommon.AnyValue_builder{
							BoolValue: ref(true),
						}.Build(),
					}.Build(),
				},
			}.Build(),
		},
		{
			name: "intValue",
			jsonStr: `{"values":[{
"key":"testKey",
"value":{
"intValue": 1
}
}]}`,
			want: otlpcommon.KeyValueList_builder{
				Values: []*otlpcommon.KeyValue{
					otlpcommon.KeyValue_builder{
						Key: "testKey",
						Value: otlpcommon.AnyValue_builder{
							IntValue: ref(int64(1)),
						}.Build(),
					}.Build(),
				},
			}.Build(),
		},
		{
			name: "doubleValue",
			jsonStr: `{"values":[{
"key":"testKey",
"value":{
"doubleValue": 1.3
}
}]}`,
			want: otlpcommon.KeyValueList_builder{
				Values: []*otlpcommon.KeyValue{
					otlpcommon.KeyValue_builder{
						Key: "testKey",
						Value: otlpcommon.AnyValue_builder{
							DoubleValue: ref(1.3),
						}.Build(),
					}.Build(),
				},
			}.Build(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iter := jsoniter.ConfigFastest.BorrowIterator([]byte(tt.jsonStr))
			defer jsoniter.ConfigFastest.ReturnIterator(iter)
			got := readKvlistValue(iter)
			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestReadAttributeUnknownField(t *testing.T) {
	jsonStr := `{"extra":""}`
	iter := jsoniter.ConfigFastest.BorrowIterator([]byte(jsonStr))
	defer jsoniter.ConfigFastest.ReturnIterator(iter)
	value := ReadAttribute(iter)
	//  unknown fields should not be an error
	require.NoError(t, iter.Error)
	assert.EqualValues(t, otlpcommon.KeyValue{}, value)
}

func TestReadAttributeValueUnknownField(t *testing.T) {
	// Key after value, to check that we correctly continue to process.
	jsonStr := `{"value": {"unknown": {"extra":""}}, "key":"test"}`
	iter := jsoniter.ConfigFastest.BorrowIterator([]byte(jsonStr))
	defer jsoniter.ConfigFastest.ReturnIterator(iter)
	value := ReadAttribute(iter)
	//  unknown fields should not be an error
	require.NoError(t, iter.Error)
	assert.EqualValues(t, otlpcommon.KeyValue_builder{Key: "test"}.Build(), value)
}

func TestReadValueUnknownField(t *testing.T) {
	jsonStr := `{"extra":""}`
	iter := jsoniter.ConfigFastest.BorrowIterator([]byte(jsonStr))
	defer jsoniter.ConfigFastest.ReturnIterator(iter)
	value := &otlpcommon.AnyValue{}
	ReadValue(iter, value)
	require.NoError(t, iter.Error)
	assert.EqualValues(t, &otlpcommon.AnyValue{}, value)
}

func TestReadValueInvliadBytesValue(t *testing.T) {
	jsonStr := `{"bytesValue": "--"}`
	iter := jsoniter.ConfigFastest.BorrowIterator([]byte(jsonStr))
	defer jsoniter.ConfigFastest.ReturnIterator(iter)

	ReadValue(iter, &otlpcommon.AnyValue{})
	assert.ErrorContains(t, iter.Error, "base64")
}

func TestReadArrayUnknownField(t *testing.T) {
	jsonStr := `{"extra":""}`
	iter := jsoniter.ConfigFastest.BorrowIterator([]byte(jsonStr))
	defer jsoniter.ConfigFastest.ReturnIterator(iter)
	value := readArray(iter)
	require.NoError(t, iter.Error)
	assert.EqualValues(t, &otlpcommon.ArrayValue{}, value)
}

func TestReadKvlistValueUnknownField(t *testing.T) {
	jsonStr := `{"extra":""}`
	iter := jsoniter.ConfigFastest.BorrowIterator([]byte(jsonStr))
	defer jsoniter.ConfigFastest.ReturnIterator(iter)
	value := readKvlistValue(iter)
	require.NoError(t, iter.Error)
	assert.EqualValues(t, &otlpcommon.KeyValueList{}, value)
}

func TestReadArrayValueInvalidArrayValue(t *testing.T) {
	jsonStr := `{"arrayValue": {"extra":""}}`
	iter := jsoniter.ConfigFastest.BorrowIterator([]byte(jsonStr))
	defer jsoniter.ConfigFastest.ReturnIterator(iter)

	value := &otlpcommon.AnyValue{}
	ReadValue(iter, value)
	require.NoError(t, iter.Error)
	assert.EqualValues(t, otlpcommon.AnyValue_builder{
		ArrayValue: &otlpcommon.ArrayValue{},
	}.Build(), value)
}

func TestReadKvlistValueInvalidArrayValue(t *testing.T) {
	jsonStr := `{"kvlistValue": {"extra":""}}`
	iter := jsoniter.ConfigFastest.BorrowIterator([]byte(jsonStr))
	defer jsoniter.ConfigFastest.ReturnIterator(iter)

	value := &otlpcommon.AnyValue{}
	ReadValue(iter, value)
	require.NoError(t, iter.Error)
	assert.EqualValues(t, otlpcommon.AnyValue_builder{
		KvlistValue: &otlpcommon.KeyValueList{},
	}.Build(), value)
}
