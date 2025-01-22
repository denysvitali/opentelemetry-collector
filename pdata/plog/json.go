// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package plog // import "go.opentelemetry.io/collector/pdata/plog"

import (
	"bytes"

	jsoniter "github.com/json-iterator/go"

	"go.opentelemetry.io/collector/pdata/internal"
	otlplogs "go.opentelemetry.io/collector/pdata/internal/data/protogen/logs/v1"
	"go.opentelemetry.io/collector/pdata/internal/json"
	"go.opentelemetry.io/collector/pdata/internal/otlp"
)

// JSONMarshaler marshals pdata.Logs to JSON bytes using the OTLP/JSON format.
type JSONMarshaler struct{}

// MarshalLogs to the OTLP/JSON format.
func (*JSONMarshaler) MarshalLogs(ld Logs) ([]byte, error) {
	buf := bytes.Buffer{}
	pb := internal.LogsToProto(internal.Logs(ld))
	err := json.Marshal(&buf, pb)
	return buf.Bytes(), err
}

var _ Unmarshaler = (*JSONUnmarshaler)(nil)

// JSONUnmarshaler unmarshals OTLP/JSON formatted-bytes to pdata.Logs.
type JSONUnmarshaler struct{}

// UnmarshalLogs from OTLP/JSON format into pdata.Logs.
func (*JSONUnmarshaler) UnmarshalLogs(buf []byte) (Logs, error) {
	iter := jsoniter.ConfigFastest.BorrowIterator(buf)
	defer jsoniter.ConfigFastest.ReturnIterator(iter)
	ld := NewLogs()
	ld.unmarshalJsoniter(iter)
	if iter.Error != nil {
		return Logs{}, iter.Error
	}
	otlp.MigrateLogs(ld.getOrig().GetResourceLogs())
	return ld, nil
}

func (ms Logs) unmarshalJsoniter(iter *jsoniter.Iterator) {
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "resource_logs", "resourceLogs":
			iter.ReadArrayCB(func(*jsoniter.Iterator) bool {
				ms.ResourceLogs().AppendEmpty().unmarshalJsoniter(iter)
				return true
			})
		default:
			iter.Skip()
		}
		return true
	})
}

func (ms ResourceLogs) unmarshalJsoniter(iter *jsoniter.Iterator) {
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "resource":
			json.ReadResource(iter, ms.orig.GetResource())
		case "scope_logs", "scopeLogs":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				ms.ScopeLogs().AppendEmpty().unmarshalJsoniter(iter)
				return true
			})
		case "schemaUrl", "schema_url":
			ms.orig.SetSchemaUrl(iter.ReadString())
		default:
			iter.Skip()
		}
		return true
	})
}

func (ms ScopeLogs) unmarshalJsoniter(iter *jsoniter.Iterator) {
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "scope":
			json.ReadScope(iter, ms.orig.GetScope())
		case "log_records", "logRecords":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				ms.LogRecords().AppendEmpty().unmarshalJsoniter(iter)
				return true
			})
		case "schemaUrl", "schema_url":
			ms.orig.SetSchemaUrl(iter.ReadString())
		default:
			iter.Skip()
		}
		return true
	})
}

func (ms LogRecord) unmarshalJsoniter(iter *jsoniter.Iterator) {
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "timeUnixNano", "time_unix_nano":
			ms.orig.SetTimeUnixNano(json.ReadUint64(iter))
		case "observed_time_unix_nano", "observedTimeUnixNano":
			ms.orig.SetObservedTimeUnixNano(json.ReadUint64(iter))
		case "severity_number", "severityNumber":
			ms.orig.SetSeverityNumber(otlplogs.SeverityNumber(json.ReadEnumValue(iter, otlplogs.SeverityNumber_value)))
		case "severity_text", "severityText":
			ms.orig.SetSeverityText(iter.ReadString())
		case "body":
			json.ReadValue(iter, ms.orig.GetBody())
		case "attributes":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				ms.orig.SetAttributes(append(ms.orig.GetAttributes(), json.ReadAttribute(iter)))
				return true
			})
		case "droppedAttributesCount", "dropped_attributes_count":
			ms.orig.SetDroppedAttributesCount(json.ReadUint32(iter))
		case "flags":
			ms.orig.SetFlags(json.ReadUint32(iter))
		case "traceId", "trace_id":
			ms.orig.SetTraceId(iter.ReadStringAsSlice())
		case "spanId", "span_id":
			ms.orig.SetSpanId(iter.ReadStringAsSlice())
		default:
			iter.Skip()
		}
		return true
	})
}
