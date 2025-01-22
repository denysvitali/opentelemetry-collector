// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal // import "go.opentelemetry.io/collector/pdata/internal"

import (
	otlpcommon "go.opentelemetry.io/collector/pdata/internal/data/protogen/common/v1"
)

type Map struct {
	orig  *[]*otlpcommon.KeyValue
	state *State
}

func GetOrigMap(ms Map) *[]*otlpcommon.KeyValue {
	return ms.orig
}

func GetMapState(ms Map) *State {
	return ms.state
}

func NewMap(orig *[]*otlpcommon.KeyValue, state *State) Map {
	return Map{orig: orig, state: state}
}

func GenerateTestMap() Map {
	var orig []*otlpcommon.KeyValue
	state := StateMutable
	ms := NewMap(&orig, &state)
	FillTestMap(ms)
	return ms
}

func FillTestMap(dest Map) {
	*dest.orig = nil
	kv := otlpcommon.KeyValue{}
	kv.SetKey("k")
	v := &otlpcommon.AnyValue{}
	v.SetStringValue("v")
	kv.SetValue(v)
	*dest.orig = append(*dest.orig, &kv)
}
