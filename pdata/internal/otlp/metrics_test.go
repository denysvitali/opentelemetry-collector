// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlp

import (
	"testing"

	"github.com/stretchr/testify/assert"

	otlpmetrics "go.opentelemetry.io/collector/pdata/internal/data/protogen/metrics/v1"
)

func TestDeprecatedScopeMetrics(t *testing.T) {
	sm := new(otlpmetrics.ScopeMetrics)
	sm1 := &otlpmetrics.ResourceMetrics{}
	sm1.SetScopeMetrics([]*otlpmetrics.ScopeMetrics{sm})
	sm1.SetDeprecatedScopeMetrics([]*otlpmetrics.ScopeMetrics{sm})
	sm2 := &otlpmetrics.ResourceMetrics{}
	sm2.SetDeprecatedScopeMetrics([]*otlpmetrics.ScopeMetrics{sm})
	rms := []*otlpmetrics.ResourceMetrics{sm1, sm2}

	MigrateMetrics(rms)
	assert.Same(t, sm, rms[0].GetScopeMetrics()[0])
	assert.Same(t, sm, rms[1].GetScopeMetrics()[0])
	assert.Nil(t, rms[0].GetDeprecatedScopeMetrics())
	assert.Nil(t, rms[0].GetDeprecatedScopeMetrics())
}
