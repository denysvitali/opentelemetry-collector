// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlp

import (
	"testing"

	"github.com/stretchr/testify/assert"

	otlptrace "go.opentelemetry.io/collector/pdata/internal/data/protogen/trace/v1"
)

func TestDeprecatedScopeSpans(t *testing.T) {
	ss := new(otlptrace.ScopeSpans)
	ss1 := &otlptrace.ResourceSpans{}
	ss1.SetScopeSpans([]*otlptrace.ScopeSpans{ss})
	ss1.SetDeprecatedScopeSpans([]*otlptrace.ScopeSpans{ss})
	ss2 := &otlptrace.ResourceSpans{}
	ss2.SetDeprecatedScopeSpans([]*otlptrace.ScopeSpans{ss})
	rss := []*otlptrace.ResourceSpans{ss1, ss2}

	MigrateTraces(rss)
	assert.Same(t, ss, rss[0].GetScopeSpans()[0])
	assert.Same(t, ss, rss[1].GetScopeSpans()[0])
	assert.Nil(t, rss[0].GetDeprecatedScopeSpans())
	assert.Nil(t, rss[0].GetDeprecatedScopeSpans())
}
