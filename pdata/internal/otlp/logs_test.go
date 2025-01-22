// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlp

import (
	"testing"

	"github.com/stretchr/testify/assert"

	otlplogs "go.opentelemetry.io/collector/pdata/internal/data/protogen/logs/v1"
)

func TestDeprecatedScopeLogs(t *testing.T) {
	sl := new(otlplogs.ScopeLogs)
	rl1 := &otlplogs.ResourceLogs{}
	rl1.SetScopeLogs([]*otlplogs.ScopeLogs{sl})
	rl1.SetDeprecatedScopeLogs([]*otlplogs.ScopeLogs{sl})

	rl2 := &otlplogs.ResourceLogs{}
	rl2.SetDeprecatedScopeLogs([]*otlplogs.ScopeLogs{sl})
	rls := []*otlplogs.ResourceLogs{rl1, rl2}

	MigrateLogs(rls)
	assert.Same(t, sl, rls[0].GetScopeLogs()[0])
	assert.Same(t, sl, rls[1].GetScopeLogs()[0])
	assert.Nil(t, rls[0].GetDeprecatedScopeLogs())
	assert.Nil(t, rls[0].GetDeprecatedScopeLogs())
}
