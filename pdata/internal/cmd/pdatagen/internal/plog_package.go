// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal // import "go.opentelemetry.io/collector/pdata/internal/cmd/pdatagen/internal"

var plog = &Package{
	info: &PackageInfo{
		name: "plog",
		path: "plog",
		imports: []string{
			`"sort"`,
			``,
			`"go.opentelemetry.io/collector/pdata/internal"`,
			`"go.opentelemetry.io/collector/pdata/internal/data"`,
			`otlpcommon "go.opentelemetry.io/collector/pdata/internal/data/protogen/common/v1"`,
			`otlplogs "go.opentelemetry.io/collector/pdata/internal/data/protogen/logs/v1"`,
			`otlpresource "go.opentelemetry.io/collector/pdata/internal/data/protogen/resource/v1"`,
			`"go.opentelemetry.io/collector/pdata/pcommon"`,
		},
		testImports: []string{
			`"testing"`,
			`"unsafe"`,
			``,
			`"github.com/stretchr/testify/assert"`,
			``,
			`"go.opentelemetry.io/collector/pdata/internal"`,
			`"go.opentelemetry.io/collector/pdata/internal/data"`,
			`otlplogs "go.opentelemetry.io/collector/pdata/internal/data/protogen/logs/v1"`,
			`"go.opentelemetry.io/collector/pdata/pcommon"`,
		},
	},
	structs: []baseStruct{
		resourceLogsSlice,
		resourceLogs,
		scopeLogsSlice,
		scopeLogs,
		logSlice,
		logRecord,
	},
}

var resourceLogsSlice = &sliceOfPtrs{
	structName: "ResourceLogsSlice",
	element:    resourceLogs,
	parentType: "otlplogs.LogsData",
}

var resourceLogs = &messageValueStruct{
	structName:     "ResourceLogs",
	description:    "// ResourceLogs is a collection of logs from a Resource.",
	originFullName: "otlplogs.ResourceLogs",
	fields: []baseField{
		resourceField,
		schemaURLField,
		&sliceField{
			fieldName:   "ScopeLogs",
			returnSlice: scopeLogsSlice,
			newType:     "otlplogs.ResourceLogs",
		},
	},
}

var scopeLogsSlice = &sliceOfPtrs{
	structName: "ScopeLogsSlice",
	element:    scopeLogs,
	parentType: "otlplogs.ResourceLogs",
}

var scopeLogs = &messageValueStruct{
	structName:     "ScopeLogs",
	description:    "// ScopeLogs is a collection of logs from a LibraryInstrumentation.",
	originFullName: "otlplogs.ScopeLogs",
	fields: []baseField{
		scopeField,
		schemaURLField,
		&sliceField{
			fieldName:   "LogRecords",
			returnSlice: logSlice,
			newType:     "otlplogs.ScopeLogs",
		},
	},
}

var logSlice = &sliceOfPtrs{
	structName: "LogRecordSlice",
	element:    logRecord,
	parentType: "otlplogs.ScopeLogs",
}

var logRecord = &messageValueStruct{
	structName:     "LogRecord",
	description:    "// LogRecord are experimental implementation of OpenTelemetry Log Data Model.\n",
	originFullName: "otlplogs.LogRecord",
	fields: []baseField{
		&primitiveTypedField{
			fieldName:       "ObservedTimestamp",
			originFieldName: "ObservedTimeUnixNano",
			returnType:      timestampType,
		},
		&primitiveTypedField{
			fieldName:       "Timestamp",
			originFieldName: "TimeUnixNano",
			returnType:      timestampType,
		},
		traceIDField,
		spanIDField,
		&primitiveTypedField{
			fieldName: "Flags",
			returnType: &primitiveType{
				structName: "LogRecordFlags",
				rawType:    "uint32",
				defaultVal: "0",
				testVal:    "1",
			},
		},
		&primitiveField{
			fieldName:  "EventName",
			returnType: "string",
			defaultVal: `""`,
			testVal:    `""`,
		},
		&primitiveField{
			fieldName:  "SeverityText",
			returnType: "string",
			defaultVal: `""`,
			testVal:    `"INFO"`,
		},
		&primitiveTypedField{
			fieldName: "SeverityNumber",
			returnType: &primitiveType{
				structName: "SeverityNumber",
				rawType:    "otlplogs.SeverityNumber",
				defaultVal: `otlplogs.SeverityNumber(0)`,
				testVal:    `otlplogs.SeverityNumber(5)`,
			},
		},
		bodyField,
		attributes,
		droppedAttributesCount,
	},
}

var bodyField = &messageValueField{
	fieldName:     "Body",
	returnMessage: anyValue,
}
