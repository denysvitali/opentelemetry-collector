// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package pmetric // import "go.opentelemetry.io/collector/pdata/pmetric"

import (
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/collector/pdata/internal"
	otlpmetrics "go.opentelemetry.io/collector/pdata/internal/data/protogen/metrics/v1"
)

var _ MarshalSizer = (*ProtoMarshaler)(nil)

type ProtoMarshaler struct{}

func (e *ProtoMarshaler) MarshalMetrics(md Metrics) ([]byte, error) {
	return proto.Marshal(internal.MetricsToProto(internal.Metrics(md)))
}

func (e *ProtoMarshaler) MetricsSize(md Metrics) int {
	return proto.Size(internal.MetricsToProto(internal.Metrics(md)))
}

type ProtoUnmarshaler struct{}

func (d *ProtoUnmarshaler) UnmarshalMetrics(buf []byte) (Metrics, error) {
	pb := otlpmetrics.MetricsData{}
	err := proto.Unmarshal(buf, &pb)
	return Metrics(internal.MetricsFromProto(&pb)), err
}
