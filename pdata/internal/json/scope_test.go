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

func TestReadScope(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    *otlpcommon.InstrumentationScope
	}{
		{
			name: "scope",
			jsonStr: `{
	"name": "name_value",
	"version": "version_value"
}`,
			want: otlpcommon.InstrumentationScope_builder{
				Name:    "name_value",
				Version: "version_value",
			}.Build(),
		},
		{
			name: "with attributes",
			jsonStr: `{
	"name": "my_name",
	"version": "my_version",
	"attributes": [
		{
			"key":"string_key",
			"value":{ "stringValue": "value" }
		},
		{
			"key":"bool_key",
			"value":{ "boolValue": true }
		},
		{
			"key":"int_key",
			"value":{ "intValue": 314 }
		},
		{
			"key":"double_key",
			"value":{ "doubleValue": 3.14 }
		}
	],
	"dropped_attributes_count": 1
}`,
			want: otlpcommon.InstrumentationScope_builder{
				Name:    "my_name",
				Version: "my_version",
				Attributes: []*otlpcommon.KeyValue{
					otlpcommon.KeyValue_builder{
						Key: "string_key",
						Value: otlpcommon.AnyValue_builder{
							StringValue: ref("value"),
						}.Build(),
					}.Build(),
					otlpcommon.KeyValue_builder{
						Key: "bool_key",
						Value: otlpcommon.AnyValue_builder{
							BoolValue: ref(true),
						}.Build(),
					}.Build(),
					otlpcommon.KeyValue_builder{
						Key: "int_key",
						Value: otlpcommon.AnyValue_builder{
							IntValue: ref(int64(314)),
						}.Build(),
					}.Build(),
					otlpcommon.KeyValue_builder{
						Key: "double_key",
						Value: otlpcommon.AnyValue_builder{
							DoubleValue: ref(3.14),
						}.Build(),
					}.Build(),
				},
				DroppedAttributesCount: 1,
			}.Build(),
		},
		{
			name: "unknown field",
			jsonStr: `{
	"name": "name_value",
	"unknown": "version"
}`,
			want: otlpcommon.InstrumentationScope_builder{
				Name: "name_value",
			}.Build(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iter := jsoniter.ConfigFastest.BorrowIterator([]byte(tt.jsonStr))
			defer jsoniter.ConfigFastest.ReturnIterator(iter)
			got := &otlpcommon.InstrumentationScope{}
			ReadScope(iter, got)
			require.NoError(t, iter.Error)
			assert.Equal(t, tt.want, got)
		})
	}
}
