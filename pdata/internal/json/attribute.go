// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package json // import "go.opentelemetry.io/collector/pdata/internal/json"

import (
	"encoding/base64"
	"fmt"

	jsoniter "github.com/json-iterator/go"

	otlpcommon "go.opentelemetry.io/collector/pdata/internal/data/protogen/common/v1"
)

// ReadAttribute Unmarshal JSON data and return otlpcommon.KeyValue
func ReadAttribute(iter *jsoniter.Iterator) *otlpcommon.KeyValue {
	kv := otlpcommon.KeyValue{}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "key":
			kv.SetKey(iter.ReadString())
		case "value":
			ReadValue(iter, kv.GetValue())
		default:
			iter.Skip()
		}
		return true
	})
	return &kv
}

// ReadValue Unmarshal JSON data and return otlpcommon.AnyValue
func ReadValue(iter *jsoniter.Iterator, val *otlpcommon.AnyValue) {
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "stringValue", "string_value":
			val.SetStringValue(iter.ReadString())

		case "boolValue", "bool_value":
			val.SetBoolValue(iter.ReadBool())
		case "intValue", "int_value":
			val.SetIntValue(ReadInt64(iter))
		case "doubleValue", "double_value":
			val.SetDoubleValue(ReadFloat64(iter))
		case "bytesValue", "bytes_value":
			v, err := base64.StdEncoding.DecodeString(iter.ReadString())
			if err != nil {
				iter.ReportError("bytesValue", fmt.Sprintf("base64 decode:%v", err))
				break
			}
			val.SetBytesValue(v)
		case "arrayValue", "array_value":
			val.SetArrayValue(readArray(iter))
		case "kvlistValue", "kvlist_value":
			val.SetKvlistValue(readKvlistValue(iter))
		default:
			iter.Skip()
		}
		return true
	})
}

func readArray(iter *jsoniter.Iterator) *otlpcommon.ArrayValue {
	v := &otlpcommon.ArrayValue{}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "values":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				v.SetValues(append(v.GetValues(), &otlpcommon.AnyValue{}))
				ReadValue(iter, v.GetValues()[len(v.GetValues())-1])
				return true
			})
		default:
			iter.Skip()
		}
		return true
	})
	return v
}

func readKvlistValue(iter *jsoniter.Iterator) *otlpcommon.KeyValueList {
	v := &otlpcommon.KeyValueList{}
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "values":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				v.SetValues(append(v.GetValues(), ReadAttribute(iter)))
				return true
			})
		default:
			iter.Skip()
		}
		return true
	})
	return v
}
