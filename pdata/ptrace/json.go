// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ptrace // import "go.opentelemetry.io/collector/pdata/ptrace"

import (
	"bytes"

	jsoniter "github.com/json-iterator/go"

	"go.opentelemetry.io/collector/pdata/internal"
	otlptrace "go.opentelemetry.io/collector/pdata/internal/data/protogen/trace/v1"
	"go.opentelemetry.io/collector/pdata/internal/json"
	"go.opentelemetry.io/collector/pdata/internal/otlp"
)

// JSONMarshaler marshals pdata.Traces to JSON bytes using the OTLP/JSON format.
type JSONMarshaler struct{}

// MarshalTraces to the OTLP/JSON format.
func (*JSONMarshaler) MarshalTraces(td Traces) ([]byte, error) {
	buf := bytes.Buffer{}
	pb := internal.TracesToProto(internal.Traces(td))
	err := json.Marshal(&buf, pb)
	return buf.Bytes(), err
}

// JSONUnmarshaler unmarshals OTLP/JSON formatted-bytes to pdata.Traces.
type JSONUnmarshaler struct{}

// UnmarshalTraces from OTLP/JSON format into pdata.Traces.
func (*JSONUnmarshaler) UnmarshalTraces(buf []byte) (Traces, error) {
	iter := jsoniter.ConfigFastest.BorrowIterator(buf)
	defer jsoniter.ConfigFastest.ReturnIterator(iter)
	td := NewTraces()
	td.unmarshalJsoniter(iter)
	if iter.Error != nil {
		return Traces{}, iter.Error
	}
	otlp.MigrateTraces(td.getOrig().GetResourceSpans())
	return td, nil
}

func (ms Traces) unmarshalJsoniter(iter *jsoniter.Iterator) {
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "resourceSpans", "resource_spans":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				ms.ResourceSpans().AppendEmpty().unmarshalJsoniter(iter)
				return true
			})
		default:
			iter.Skip()
		}
		return true
	})
}

func (ms ResourceSpans) unmarshalJsoniter(iter *jsoniter.Iterator) {
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "resource":
			json.ReadResource(iter, internal.GetOrigResource(internal.Resource(ms.Resource())))
		case "scopeSpans", "scope_spans":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				ms.ScopeSpans().AppendEmpty().unmarshalJsoniter(iter)
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

func (ms ScopeSpans) unmarshalJsoniter(iter *jsoniter.Iterator) {
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "scope":
			json.ReadScope(iter, ms.orig.GetScope())
		case "spans":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				ms.Spans().AppendEmpty().unmarshalJsoniter(iter)
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

func (dest Span) unmarshalJsoniter(iter *jsoniter.Iterator) {
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "traceId", "trace_id":
			dest.orig.SetTraceId(iter.ReadStringAsSlice())
		case "spanId", "span_id":
			dest.orig.SetSpanId(iter.ReadStringAsSlice())
		case "traceState", "trace_state":
			dest.TraceState().FromRaw(iter.ReadString())
		case "parentSpanId", "parent_span_id":
			dest.orig.SetParentSpanId(iter.ReadStringAsSlice())
		case "flags":
			dest.orig.SetFlags(json.ReadUint32(iter))
		case "name":
			dest.orig.SetName(iter.ReadString())
		case "kind":
			dest.orig.SetKind(otlptrace.Span_SpanKind(json.ReadEnumValue(iter, otlptrace.Span_SpanKind_value)))
		case "startTimeUnixNano", "start_time_unix_nano":
			dest.orig.SetStartTimeUnixNano(json.ReadUint64(iter))
		case "endTimeUnixNano", "end_time_unix_nano":
			dest.orig.SetEndTimeUnixNano(json.ReadUint64(iter))
		case "attributes":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				dest.orig.SetAttributes(append(dest.orig.GetAttributes(), json.ReadAttribute(iter)))
				return true
			})
		case "droppedAttributesCount", "dropped_attributes_count":
			dest.orig.SetDroppedAttributesCount(json.ReadUint32(iter))
		case "events":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				dest.Events().AppendEmpty().unmarshalJsoniter(iter)
				return true
			})
		case "droppedEventsCount", "dropped_events_count":
			dest.orig.SetDroppedEventsCount(json.ReadUint32(iter))
		case "links":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				dest.Links().AppendEmpty().unmarshalJsoniter(iter)
				return true
			})
		case "droppedLinksCount", "dropped_links_count":
			dest.orig.SetDroppedLinksCount(json.ReadUint32(iter))
		case "status":
			dest.Status().unmarshalJsoniter(iter)
		default:
			iter.Skip()
		}
		return true
	})
}

func (dest Status) unmarshalJsoniter(iter *jsoniter.Iterator) {
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "message":
			dest.orig.SetMessage(iter.ReadString())
		case "code":
			dest.orig.SetCode(otlptrace.Status_StatusCode(json.ReadEnumValue(iter, otlptrace.Status_StatusCode_value)))
		default:
			iter.Skip()
		}
		return true
	})
}

func (dest SpanLink) unmarshalJsoniter(iter *jsoniter.Iterator) {
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "traceId", "trace_id":
			dest.orig.SetTraceId(iter.ReadStringAsSlice())
		case "spanId", "span_id":
			dest.orig.SetSpanId(iter.ReadStringAsSlice())
		case "traceState", "trace_state":
			dest.orig.SetTraceState(iter.ReadString())
		case "attributes":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				dest.orig.SetAttributes(append(dest.orig.GetAttributes(), json.ReadAttribute(iter)))
				return true
			})
		case "droppedAttributesCount", "dropped_attributes_count":
			dest.orig.SetDroppedAttributesCount(json.ReadUint32(iter))
		case "flags":
			dest.orig.SetFlags(json.ReadUint32(iter))
		default:
			iter.Skip()
		}
		return true
	})
}

func (dest SpanEvent) unmarshalJsoniter(iter *jsoniter.Iterator) {
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, f string) bool {
		switch f {
		case "timeUnixNano", "time_unix_nano":
			dest.orig.SetTimeUnixNano(json.ReadUint64(iter))
		case "name":
			dest.orig.SetName(iter.ReadString())
		case "attributes":
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				dest.orig.SetAttributes(append(dest.orig.GetAttributes(), json.ReadAttribute(iter)))
				return true
			})
		case "droppedAttributesCount", "dropped_attributes_count":
			dest.orig.SetDroppedAttributesCount(json.ReadUint32(iter))
		default:
			iter.Skip()
		}
		return true
	})
}
