package json_test

import (
	"fmt"
	"testing"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/plog/plogotlp"
)

func BenchmarkEncoding(b *testing.B) {
	logs := plog.NewLogs()
	rl := logs.ResourceLogs().AppendEmpty()
	resAttrs := rl.Resource().Attributes()
	resAttrs.PutStr("host", "test-host")
	resAttrs.PutStr("version", "v1.0.0")
	sl := rl.ScopeLogs().AppendEmpty()

	for i := 0; i < 1000; i++ {
		logRecord := sl.LogRecords().AppendEmpty()
		logRecord.Body().SetStr(fmt.Sprintf("Log line %d", i))
		logRecord.SetSpanID([]byte{0, 1, 2, 3, 4, 5, 6, 7})
		logRecord.SetEventName(fmt.Sprintf("event_%d", i))
		logRecord.SetSeverityText("INFO")
		logRecord.SetTimestamp(
			pcommon.NewTimestampFromTime(
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			),
		)
	}
	logs.ResourceLogs().AppendEmpty()
	exp := plogotlp.NewExportRequestFromLogs(logs)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := exp.MarshalJSON()
		if err != nil {
			b.Fatal(err)
		}
	}
}
