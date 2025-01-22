// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package pcommon // import "go.opentelemetry.io/collector/pdata/pcommon"

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"go.opentelemetry.io/collector/pdata/internal"
	otlpcommon "go.opentelemetry.io/collector/pdata/internal/data/protogen/common/v1"
)

// ValueType specifies the type of Value.
type ValueType int32

const (
	ValueTypeEmpty ValueType = iota
	ValueTypeStr
	ValueTypeInt
	ValueTypeDouble
	ValueTypeBool
	ValueTypeMap
	ValueTypeSlice
	ValueTypeBytes
)

// String returns the string representation of the ValueType.
func (avt ValueType) String() string {
	switch avt {
	case ValueTypeEmpty:
		return "Empty"
	case ValueTypeStr:
		return "Str"
	case ValueTypeBool:
		return "Bool"
	case ValueTypeInt:
		return "Int"
	case ValueTypeDouble:
		return "Double"
	case ValueTypeMap:
		return "Map"
	case ValueTypeSlice:
		return "Slice"
	case ValueTypeBytes:
		return "Bytes"
	}
	return ""
}

// Value is a mutable cell containing any value. Typically used as an element of Map or Slice.
// Must use one of NewValue+ functions below to create new instances.
//
// Intended to be passed by value since internally it is just a pointer to actual
// value representation. For the same reason passing by value and calling setters
// will modify the original, e.g.:
//
//	func f1(val Value) { val.SetInt(234) }
//	func f2() {
//	    v := NewValueStr("a string")
//	    f1(v)
//	    _ := v.Type() // this will return ValueTypeInt
//	}
//
// Important: zero-initialized instance is not valid for use. All Value functions below must
// be called only on instances that are created via NewValue+ functions.
type Value internal.Value

// NewValueEmpty creates a new Value with an empty value.
func NewValueEmpty() Value {
	state := internal.StateMutable
	return newValue(&otlpcommon.AnyValue{}, &state)
}

// NewValueStr creates a new Value with the given string value.
func NewValueStr(v string) Value {
	state := internal.StateMutable
	return newValue(otlpcommon.AnyValue_builder{StringValue: &v}.Build(), &state)
}

// NewValueInt creates a new Value with the given int64 value.
func NewValueInt(v int64) Value {
	state := internal.StateMutable
	return newValue(otlpcommon.AnyValue_builder{IntValue: &v}.Build(), &state)
}

// NewValueDouble creates a new Value with the given float64 value.
func NewValueDouble(v float64) Value {
	state := internal.StateMutable
	return newValue(otlpcommon.AnyValue_builder{DoubleValue: &v}.Build(), &state)
}

// NewValueBool creates a new Value with the given bool value.
func NewValueBool(v bool) Value {
	state := internal.StateMutable
	return newValue(otlpcommon.AnyValue_builder{BoolValue: &v}.Build(), &state)
}

// NewValueMap creates a new Value of map type.
func NewValueMap() Value {
	state := internal.StateMutable
	return newValue(otlpcommon.AnyValue_builder{KvlistValue: otlpcommon.KeyValueList_builder{}.Build()}.Build(), &state)
}

// NewValueSlice creates a new Value of array type.
func NewValueSlice() Value {
	state := internal.StateMutable
	return newValue(otlpcommon.AnyValue_builder{ArrayValue: otlpcommon.ArrayValue_builder{}.Build()}.Build(), &state)
}

// NewValueBytes creates a new empty Value of byte type.
func NewValueBytes() Value {
	state := internal.StateMutable
	return newValue(otlpcommon.AnyValue_builder{BytesValue: nil}.Build(), &state)
}

func newValue(orig *otlpcommon.AnyValue, state *internal.State) Value {
	return Value(internal.NewValue(orig, state))
}

func (v Value) getOrig() *otlpcommon.AnyValue {
	return internal.GetOrigValue(internal.Value(v))
}

func (v Value) getState() *internal.State {
	return internal.GetValueState(internal.Value(v))
}

// FromRaw sets the value from the given raw value.
// Calling this function on zero-initialized Value will cause a panic.
func (v Value) FromRaw(iv any) error {
	switch tv := iv.(type) {
	case nil:
		v.getOrig().Reset()
	case string:
		v.SetStr(tv)
	case int:
		v.SetInt(int64(tv))
	case int8:
		v.SetInt(int64(tv))
	case int16:
		v.SetInt(int64(tv))
	case int32:
		v.SetInt(int64(tv))
	case int64:
		v.SetInt(tv)
	case uint:
		// nolint:gosec
		v.SetInt(int64(tv))
	case uint8:
		v.SetInt(int64(tv))
	case uint16:
		v.SetInt(int64(tv))
	case uint32:
		v.SetInt(int64(tv))
	case uint64:
		// nolint:gosec
		v.SetInt(int64(tv))
	case float32:
		v.SetDouble(float64(tv))
	case float64:
		v.SetDouble(tv)
	case bool:
		v.SetBool(tv)
	case []byte:
		v.SetEmptyBytes().FromRaw(tv)
	case map[string]any:
		return v.SetEmptyMap().FromRaw(tv)
	case []any:
		return v.SetEmptySlice().FromRaw(tv)
	default:
		return fmt.Errorf("<Invalid value type %T>", tv)
	}
	return nil
}

// Type returns the type of the value for this Value.
// Calling this function on zero-initialized Value will cause a panic.
func (v Value) Type() ValueType {
	switch v.getOrig().WhichValue() {
	case otlpcommon.AnyValue_StringValue_case:
		return ValueTypeStr
	case otlpcommon.AnyValue_BoolValue_case:
		return ValueTypeBool
	case otlpcommon.AnyValue_IntValue_case:
		return ValueTypeInt
	case otlpcommon.AnyValue_DoubleValue_case:
		return ValueTypeDouble
	case otlpcommon.AnyValue_KvlistValue_case:
		return ValueTypeMap
	case otlpcommon.AnyValue_ArrayValue_case:
		return ValueTypeSlice
	case otlpcommon.AnyValue_BytesValue_case:
		return ValueTypeBytes
	}
	return ValueTypeEmpty
}

// Str returns the string value associated with this Value.
// The shorter name is used instead of String to avoid implementing fmt.Stringer interface.
// If the Type() is not ValueTypeStr then returns empty string.
func (v Value) Str() string {
	return v.getOrig().GetStringValue()
}

// Int returns the int64 value associated with this Value.
// If the Type() is not ValueTypeInt then returns int64(0).
func (v Value) Int() int64 {
	return v.getOrig().GetIntValue()
}

// Double returns the float64 value associated with this Value.
// If the Type() is not ValueTypeDouble then returns float64(0).
func (v Value) Double() float64 {
	return v.getOrig().GetDoubleValue()
}

// Bool returns the bool value associated with this Value.
// If the Type() is not ValueTypeBool then returns false.
func (v Value) Bool() bool {
	return v.getOrig().GetBoolValue()
}

// Map returns the map value associated with this Value.
// If the function is called on zero-initialized Value or if the Type() is not ValueTypeMap
// then it returns an invalid map. Note that using such map can cause panic.
func (v Value) Map() Map {
	kvlist := v.getOrig().GetKvlistValue()
	if kvlist == nil {
		return Map{}
	}
	return newMap(kvlist.GetValues(), internal.GetValueState(internal.Value(v)))
}

// Slice returns the slice value associated with this Value.
// If the function is called on zero-initialized Value or if the Type() is not ValueTypeSlice
// then returns an invalid slice. Note that using such slice can cause panic.
func (v Value) Slice() Slice {
	arr := v.getOrig().GetArrayValue()
	if arr == nil {
		return Slice{}
	}
	values := arr.GetValues()
	return newSlice(&values, internal.GetValueState(internal.Value(v)))
}

// Bytes returns the ByteSlice value associated with this Value.
// If the function is called on zero-initialized Value or if the Type() is not ValueTypeBytes
// then returns an invalid ByteSlice object. Note that using such slice can cause panic.
func (v Value) Bytes() ByteSlice {
	bv := v.getOrig().GetBytesValue()
	if bv == nil {
		return ByteSlice{}
	}
	return ByteSlice(internal.NewByteSlice(&bv, internal.GetValueState(internal.Value(v))))
}

// SetStr replaces the string value associated with this Value,
// it also changes the type to be ValueTypeStr.
// The shorter name is used instead of SetString to avoid implementing
// fmt.Stringer interface by the corresponding getter method.
// Calling this function on zero-initialized Value will cause a panic.
func (v Value) SetStr(sv string) {
	v.getState().AssertMutable()
	if v.getOrig() != nil {
		v.getOrig().SetStringValue(sv)
	} else {
		panic("calling SetStr on invalid Value")
	}
}

// SetInt replaces the int64 value associated with this Value,
// it also changes the type to be ValueTypeInt.
// Calling this function on zero-initialized Value will cause a panic.
func (v Value) SetInt(iv int64) {
	v.getState().AssertMutable()
	v.getOrig().SetIntValue(iv)
}

// SetDouble replaces the float64 value associated with this Value,
// it also changes the type to be ValueTypeDouble.
// Calling this function on zero-initialized Value will cause a panic.
func (v Value) SetDouble(dv float64) {
	v.getState().AssertMutable()
	v.getOrig().SetDoubleValue(dv)
}

// SetBool replaces the bool value associated with this Value,
// it also changes the type to be ValueTypeBool.
// Calling this function on zero-initialized Value will cause a panic.
func (v Value) SetBool(bv bool) {
	v.getState().AssertMutable()
	v.getOrig().SetBoolValue(bv)
}

// SetEmptyBytes sets value to an empty byte slice and returns it.
// Calling this function on zero-initialized Value will cause a panic.
func (v Value) SetEmptyBytes() ByteSlice {
	v.getState().AssertMutable()
	var bv []byte
	v.getOrig().SetBytesValue(bv)
	return ByteSlice(internal.NewByteSlice(&bv, v.getState()))
}

// SetEmptyMap sets value to an empty map and returns it.
// Calling this function on zero-initialized Value will cause a panic.
func (v Value) SetEmptyMap() Map {
	v.getState().AssertMutable()
	kvList := []*otlpcommon.KeyValue{}
	kv := otlpcommon.KeyValueList_builder{Values: kvList}.Build()
	v.getOrig().SetKvlistValue(kv)
	return newMap(kvList, v.getState())
}

// SetEmptySlice sets value to an empty slice and returns it.
// Calling this function on zero-initialized Value will cause a panic.
func (v Value) SetEmptySlice() Slice {
	v.getState().AssertMutable()
	values := []*otlpcommon.AnyValue{}
	v.getOrig().SetArrayValue(otlpcommon.ArrayValue_builder{Values: values}.Build())
	return newSlice(&values, v.getState())
}

// CopyTo copies the Value instance overriding the destination.
// Calling this function on zero-initialized Value will cause a panic.
func (v Value) CopyTo(dest Value) {
	dest.getState().AssertMutable()
	destOrig := dest.getOrig()
	switch v.getOrig().WhichValue() {
	case otlpcommon.AnyValue_KvlistValue_case:
		destKv := destOrig.GetKvlistValue()
		if destKv == nil {
			destKv = otlpcommon.KeyValueList_builder{Values: []*otlpcommon.KeyValue{}}.Build()
			destOrig.SetKvlistValue(destKv)
		}
		if v.getOrig().GetKvlistValue() == nil {
			destKv.Reset()
			return
		}
		// Deep copy to dest.
		newMap(v.getOrig().GetKvlistValue().GetValues(), v.getState()).CopyTo(newMap(destKv.GetValues(), dest.getState()))
	case otlpcommon.AnyValue_ArrayValue_case:
		av := destOrig.GetArrayValue()
		if av == nil {
			av = otlpcommon.ArrayValue_builder{Values: []*otlpcommon.AnyValue{}}.Build()
			destOrig.SetArrayValue(av)
		}
		if v.getOrig().GetArrayValue() == nil {
			av.Reset()
			return
		}
		// Deep copy to dest.
		vals := v.getOrig().GetArrayValue().GetValues()
		avVals := av.GetValues()
		newSlice(&vals, v.getState()).CopyTo(newSlice(&avVals, dest.getState()))
	case otlpcommon.AnyValue_BytesValue_case:
		destOrig.SetBytesValue(make([]byte, len(v.getOrig().GetBytesValue())))
		destBv := destOrig.GetBytesValue()
		copy(destBv, v.getOrig().GetBytesValue())
	default:
		// Primitive immutable type, no need for deep copy.
		destOrig = v.getOrig()
	}
}

// AsString converts an OTLP Value object of any type to its equivalent string
// representation. This differs from Str which only returns a non-empty value
// if the ValueType is ValueTypeStr.
// Calling this function on zero-initialized Value will cause a panic.
func (v Value) AsString() string {
	switch v.Type() {
	case ValueTypeEmpty:
		return ""

	case ValueTypeStr:
		return v.Str()

	case ValueTypeBool:
		return strconv.FormatBool(v.Bool())

	case ValueTypeDouble:
		return float64AsString(v.Double())

	case ValueTypeInt:
		return strconv.FormatInt(v.Int(), 10)

	case ValueTypeMap:
		jsonStr, _ := json.Marshal(v.Map().AsRaw())
		return string(jsonStr)

	case ValueTypeBytes:
		return base64.StdEncoding.EncodeToString(*v.Bytes().getOrig())

	case ValueTypeSlice:
		jsonStr, _ := json.Marshal(v.Slice().AsRaw())
		return string(jsonStr)

	default:
		return fmt.Sprintf("<Unknown OpenTelemetry attribute value type %q>", v.Type())
	}
}

// See https://cs.opensource.google/go/go/+/refs/tags/go1.17.7:src/encoding/json/encode.go;l=585.
// This allows us to avoid using reflection.
func float64AsString(f float64) string {
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return "json: unsupported value: " + strconv.FormatFloat(f, 'g', -1, 64)
	}

	// Convert as if by ES6 number to string conversion.
	// This matches most other JSON generators.
	// See golang.org/issue/6384 and golang.org/issue/14135.
	// Like fmt %g, but the exponent cutoffs are different
	// and exponents themselves are not padded to two digits.
	scratch := [64]byte{}
	b := scratch[:0]
	abs := math.Abs(f)
	fmt := byte('f')
	if abs != 0 && (abs < 1e-6 || abs >= 1e21) {
		fmt = 'e'
	}
	b = strconv.AppendFloat(b, f, fmt, -1, 64)
	if fmt == 'e' {
		// clean up e-09 to e-9
		n := len(b)
		if n >= 4 && b[n-4] == 'e' && b[n-3] == '-' && b[n-2] == '0' {
			b[n-2] = b[n-1]
			b = b[:n-1]
		}
	}
	return string(b)
}

func (v Value) AsRaw() any {
	switch v.Type() {
	case ValueTypeEmpty:
		return nil
	case ValueTypeStr:
		return v.Str()
	case ValueTypeBool:
		return v.Bool()
	case ValueTypeDouble:
		return v.Double()
	case ValueTypeInt:
		return v.Int()
	case ValueTypeBytes:
		return v.Bytes().AsRaw()
	case ValueTypeMap:
		return v.Map().AsRaw()
	case ValueTypeSlice:
		return v.Slice().AsRaw()
	}
	return fmt.Sprintf("<Unknown OpenTelemetry value type %q>", v.Type())
}

func newKeyValueString(k string, v string) *otlpcommon.KeyValue {
	orig := otlpcommon.KeyValue_builder{Key: k, Value: &otlpcommon.AnyValue{}}.Build()
	state := internal.StateMutable
	akv := newValue(orig.GetValue(), &state)
	akv.SetStr(v)
	return orig
}

func newKeyValueInt(k string, v int64) *otlpcommon.KeyValue {
	orig := otlpcommon.KeyValue_builder{Key: k}.Build()
	state := internal.StateMutable
	akv := newValue(orig.GetValue(), &state)
	akv.SetInt(v)
	return orig
}

func newKeyValueDouble(k string, v float64) *otlpcommon.KeyValue {
	orig := otlpcommon.KeyValue_builder{Key: k}.Build()
	state := internal.StateMutable
	akv := newValue(orig.GetValue(), &state)
	akv.SetDouble(v)
	return orig
}

func newKeyValueBool(k string, v bool) *otlpcommon.KeyValue {
	orig := otlpcommon.KeyValue_builder{Key: k}.Build()
	state := internal.StateMutable
	akv := newValue(orig.GetValue(), &state)
	akv.SetBool(v)
	return orig
}
