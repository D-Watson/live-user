// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package log

import (
	"encoding/base64"
	"math"
	"sync"
	"time"
	"unicode/utf8"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// For JSON-escaping; see jsonEncoder.safeAddString below.
const _hex = "0123456789abcdef"

var _jsonPool = sync.Pool{New: func() interface{} {
	return &CustomEncoder{}
}}

func getCustomEncoder() *CustomEncoder {
	return _jsonPool.Get().(*CustomEncoder)
}

func putCustomEncoder(enc *CustomEncoder) {
	//if enc.reflectBuf != nil {
	//	enc.reflectBuf.Free()
	//}
	enc.EncoderConfig = nil
	enc.buf = nil
	enc.spaced = false
	enc.openNamespaces = 0
	//enc.reflectBuf = nil
	//enc.reflectEnc = nil
	_jsonPool.Put(enc)
}

type CustomEncoder struct {
	*zapcore.EncoderConfig
	buf            *buffer.Buffer
	spaced         bool // include spaces after colons and commas
	openNamespaces int

	// for encoding generic values by reflection
	//reflectBuf *buffer.Buffer
	//reflectEnc ReflectedEncoder
}

func NewCustomEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	if cfg.ConsoleSeparator == "" {
		// Use a default delimiter of '\t' for backwards compatibility
		cfg.ConsoleSeparator = "|"
	}
	return newCustomEncoder(cfg, true)
}

func newCustomEncoder(cfg zapcore.EncoderConfig, spaced bool) *CustomEncoder {
	if cfg.SkipLineEnding {
		cfg.LineEnding = ""
	} else if cfg.LineEnding == "" {
		cfg.LineEnding = zapcore.DefaultLineEnding
	}

	// If no EncoderConfig.NewReflectedEncoder is provided by the user, then use default
	//if cfg.NewReflectedEncoder == nil {
	//	cfg.NewReflectedEncoder = defaultReflectedEncoder
	//}

	return &CustomEncoder{
		EncoderConfig: &cfg,
		buf:           _pool.Get(),
		spaced:        spaced,
	}
}

func (enc *CustomEncoder) AddArray(key string, arr zapcore.ArrayMarshaler) error {
	enc.addKey(key)
	return enc.AppendArray(arr)
}

func (enc *CustomEncoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	enc.addKey(key)
	return enc.AppendObject(obj)
}

func (enc *CustomEncoder) AddBinary(key string, val []byte) {
	enc.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (enc *CustomEncoder) AddByteString(key string, val []byte) {
	enc.addKey(key)
	enc.AppendByteString(val)
}

func (enc *CustomEncoder) AddBool(key string, val bool) {
	enc.addKey(key)
	enc.AppendBool(val)
}

func (enc *CustomEncoder) AddComplex128(key string, val complex128) {
	enc.addKey(key)
	enc.AppendComplex128(val)
}

func (enc *CustomEncoder) AddComplex64(key string, val complex64) {
	enc.addKey(key)
	enc.AppendComplex64(val)
}

func (enc *CustomEncoder) AddDuration(key string, val time.Duration) {
	enc.addKey(key)
	enc.AppendDuration(val)
}

func (enc *CustomEncoder) AddFloat64(key string, val float64) {
	enc.addKey(key)
	enc.AppendFloat64(val)
}

func (enc *CustomEncoder) AddFloat32(key string, val float32) {
	enc.addKey(key)
	enc.AppendFloat32(val)
}

func (enc *CustomEncoder) AddInt64(key string, val int64) {
	enc.addKey(key)
	enc.AppendInt64(val)
}

func (enc *CustomEncoder) resetReflectBuf() {

}

var nullLiteralBytes = []byte("null")

// Only invoke the standard JSON encoder if there is actually something to
// encode; otherwise write JSON null literal directly.
func (enc *CustomEncoder) encodeReflected(obj interface{}) ([]byte, error) {

	return nil, nil
}

func (enc *CustomEncoder) AddReflected(key string, obj interface{}) error {
	valueBytes, err := enc.encodeReflected(obj)
	if err != nil {
		return err
	}
	enc.addKey(key)
	_, err = enc.buf.Write(valueBytes)
	return err
}

func (enc *CustomEncoder) OpenNamespace(key string) {
	enc.addKey(key)
	enc.buf.AppendByte('{')
	enc.openNamespaces++
}

func (enc *CustomEncoder) AddString(key, val string) {
	enc.addKey(key)
	enc.AppendString(val)
}

func (enc *CustomEncoder) AddTime(key string, val time.Time) {
	enc.addKey(key)
	enc.AppendTime(val)
}

func (enc *CustomEncoder) AddUint64(key string, val uint64) {
	enc.addKey(key)
	enc.AppendUint64(val)
}

func (enc *CustomEncoder) AppendArray(arr zapcore.ArrayMarshaler) error {
	//enc.addElementSeparator()
	enc.buf.AppendByte('[')
	err := arr.MarshalLogArray(enc)
	enc.buf.AppendByte(']')
	return err
}

func (enc *CustomEncoder) AppendObject(obj zapcore.ObjectMarshaler) error {
	// Close ONLY new openNamespaces that are created during
	// AppendObject().
	old := enc.openNamespaces
	enc.openNamespaces = 0
	//enc.addElementSeparator()
	enc.buf.AppendByte('{')
	err := obj.MarshalLogObject(enc)
	enc.buf.AppendByte('}')
	enc.closeOpenNamespaces()
	enc.openNamespaces = old
	return err
}

func (enc *CustomEncoder) AppendBool(val bool) {
	enc.addElementSeparator()
	enc.buf.AppendBool(val)
}

func (enc *CustomEncoder) AppendByteString(val []byte) {
	enc.addElementSeparator()
	//enc.buf.AppendByte('"')
	enc.safeAddByteString(val)
	//enc.buf.AppendByte('"')
}

// appendComplex appends the encoded form of the provided complex128 value.
// precision specifies the encoding precision for the real and imaginary
// components of the complex number.
func (enc *CustomEncoder) appendComplex(val complex128, precision int) {
	enc.addElementSeparator()
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(val)), float64(imag(val))
	enc.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendFloat(r, precision)
	// If imaginary part is less than 0, minus (-) sign is added by default
	// by AppendFloat.
	if i >= 0 {
		enc.buf.AppendByte('+')
	}
	enc.buf.AppendFloat(i, precision)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}

func (enc *CustomEncoder) AppendDuration(val time.Duration) {
	cur := enc.buf.Len()
	if e := enc.EncodeDuration; e != nil {
		e(val, enc)
	}
	if cur == enc.buf.Len() {
		// User-supplied EncodeDuration is a no-op. Fall back to nanoseconds to keep
		// JSON valid.
		enc.AppendInt64(int64(val))
	}
}

func (enc *CustomEncoder) AppendInt64(val int64) {
	enc.addElementSeparator()
	enc.buf.AppendInt(val)
}

func (enc *CustomEncoder) AppendReflected(val interface{}) error {
	valueBytes, err := enc.encodeReflected(val)
	if err != nil {
		return err
	}
	enc.addElementSeparator()
	_, err = enc.buf.Write(valueBytes)
	return err
}

func (enc *CustomEncoder) AppendString(val string) {
	enc.addElementSeparator()
	//enc.buf.AppendByte('"')
	enc.safeAddString(val)
	//enc.buf.AppendByte('"')
}

func (enc *CustomEncoder) AppendTimeLayout(time time.Time, layout string) {
	enc.addElementSeparator()
	//enc.buf.AppendByte('"')
	enc.buf.AppendTime(time, layout)
	//enc.buf.AppendByte('"')
}

func (enc *CustomEncoder) AppendTime(val time.Time) {
	cur := enc.buf.Len()
	if e := enc.EncodeTime; e != nil {
		e(val, enc)
	}
	if cur == enc.buf.Len() {
		// User-supplied EncodeTime is a no-op. Fall back to nanos since epoch to keep
		// output JSON valid.
		enc.AppendInt64(val.UnixNano())
	}
}

func (enc *CustomEncoder) AppendUint64(val uint64) {
	enc.addElementSeparator()
	enc.buf.AppendUint(val)
}

func (enc *CustomEncoder) AddInt(k string, v int)         { enc.AddInt64(k, int64(v)) }
func (enc *CustomEncoder) AddInt32(k string, v int32)     { enc.AddInt64(k, int64(v)) }
func (enc *CustomEncoder) AddInt16(k string, v int16)     { enc.AddInt64(k, int64(v)) }
func (enc *CustomEncoder) AddInt8(k string, v int8)       { enc.AddInt64(k, int64(v)) }
func (enc *CustomEncoder) AddUint(k string, v uint)       { enc.AddUint64(k, uint64(v)) }
func (enc *CustomEncoder) AddUint32(k string, v uint32)   { enc.AddUint64(k, uint64(v)) }
func (enc *CustomEncoder) AddUint16(k string, v uint16)   { enc.AddUint64(k, uint64(v)) }
func (enc *CustomEncoder) AddUint8(k string, v uint8)     { enc.AddUint64(k, uint64(v)) }
func (enc *CustomEncoder) AddUintptr(k string, v uintptr) { enc.AddUint64(k, uint64(v)) }
func (enc *CustomEncoder) AppendComplex64(v complex64)    { enc.appendComplex(complex128(v), 32) }
func (enc *CustomEncoder) AppendComplex128(v complex128)  { enc.appendComplex(complex128(v), 64) }
func (enc *CustomEncoder) AppendFloat64(v float64)        { enc.appendFloat(v, 64) }
func (enc *CustomEncoder) AppendFloat32(v float32)        { enc.appendFloat(float64(v), 32) }
func (enc *CustomEncoder) AppendInt(v int)                { enc.AppendInt64(int64(v)) }
func (enc *CustomEncoder) AppendInt32(v int32)            { enc.AppendInt64(int64(v)) }
func (enc *CustomEncoder) AppendInt16(v int16)            { enc.AppendInt64(int64(v)) }
func (enc *CustomEncoder) AppendInt8(v int8)              { enc.AppendInt64(int64(v)) }
func (enc *CustomEncoder) AppendUint(v uint)              { enc.AppendUint64(uint64(v)) }
func (enc *CustomEncoder) AppendUint32(v uint32)          { enc.AppendUint64(uint64(v)) }
func (enc *CustomEncoder) AppendUint16(v uint16)          { enc.AppendUint64(uint64(v)) }
func (enc *CustomEncoder) AppendUint8(v uint8)            { enc.AppendUint64(uint64(v)) }
func (enc *CustomEncoder) AppendUintptr(v uintptr)        { enc.AppendUint64(uint64(v)) }

func (enc *CustomEncoder) Clone() zapcore.Encoder {
	clone := enc.clone()
	clone.buf.Write(enc.buf.Bytes())
	return clone
}

func (enc *CustomEncoder) clone() *CustomEncoder {
	clone := getCustomEncoder()
	clone.EncoderConfig = enc.EncoderConfig
	clone.spaced = enc.spaced
	clone.openNamespaces = enc.openNamespaces
	clone.buf = _pool.Get()
	return clone
}

func (enc *CustomEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	final := enc.clone()
	//final.buf.AppendByte('{')

	if final.TimeKey != "" {
		final.AppendTime(ent.Time)
		final.buf.AppendString(enc.ConsoleSeparator)
	}

	if final.LevelKey != "" && final.EncodeLevel != nil {
		final.buf.AppendString(ent.Level.String())

		final.buf.AppendString(enc.ConsoleSeparator)
	}

	if ent.Caller.Defined {
		if final.CallerKey != "" && final.EncodeCaller != nil {
			final.EncodeCaller(ent.Caller, final)
		}
		if final.FunctionKey != "" {
			final.buf.AppendString(ent.Caller.Function)
		}
		final.buf.AppendString(enc.ConsoleSeparator)
	}

	fields = final.addTraceInfo(fields)

	if final.MessageKey != "" {
		final.buf.AppendString(enc.MessageKey)
		final.buf.AppendByte('=')
		final.buf.AppendString(ent.Message)
	}

	if enc.buf.Len() > 0 {
		final.addElementSeparator()
		final.buf.Write(enc.buf.Bytes())
	}

	if len(fields) != 0 {
		final.buf.AppendString(enc.ConsoleSeparator)
		addFields(final, fields)
	}

	final.closeOpenNamespaces()
	if ent.Stack != "" && final.StacktraceKey != "" {
		final.AddString(final.StacktraceKey, ent.Stack)
	}

	final.buf.AppendString(final.LineEnding)

	ret := final.buf
	putCustomEncoder(final)
	return ret, nil
}

func (enc *CustomEncoder) truncate() {
	enc.buf.Reset()
}

func (enc *CustomEncoder) closeOpenNamespaces() {
	for i := 0; i < enc.openNamespaces; i++ {
		enc.buf.AppendByte('}')
	}
	enc.openNamespaces = 0
}

func (enc *CustomEncoder) addKey(key string) {
	enc.addElementSeparator()
	//enc.buf.AppendByte('')
	enc.safeAddString(key)
	//enc.buf.AppendByte('"')
	enc.buf.AppendByte('=')
	//if enc.spaced {
	//	enc.buf.AppendByte(' ')
	//}
}

func (enc *CustomEncoder) addElementSeparator() {
	last := enc.buf.Len() - 1
	if last < 0 {
		return
	}

	switch enc.buf.Bytes()[last] {
	case '{', '[', ':', ',', ' ', '|', '=':
		return
	default:
		enc.buf.AppendByte(',')
		if enc.spaced {
			enc.buf.AppendByte(' ')
		}
	}
}

func (enc *CustomEncoder) appendFloat(val float64, bitSize int) {
	enc.addElementSeparator()
	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(val, bitSize)
	}
}

// safeAddString JSON-escapes a string and appends it to the internal buffer.
// Unlike the standard library's encoder, it doesn't attempt to protect the
// user from browser vulnerabilities or JSONP-related problems.
func (enc *CustomEncoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.AppendString(s[i : i+size])
		i += size
	}
}

// safeAddByteString is no-alloc equivalent of safeAddString(string(s)) for s []byte.
func (enc *CustomEncoder) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.Write(s[i : i+size])
		i += size
	}
}

// tryAddRuneSelf appends b if it is valid UTF-8 character represented in a single byte.
func (enc *CustomEncoder) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		enc.buf.AppendByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte(b)
	case '\n':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('n')
	case '\r':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('r')
	case '\t':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		enc.buf.AppendString(`\u00`)
		enc.buf.AppendByte(_hex[b>>4])
		enc.buf.AppendByte(_hex[b&0xF])
	}
	return true
}

func (enc *CustomEncoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		enc.buf.AppendString(`\ufffd`)
		return true
	}
	return false
}

func (enc *CustomEncoder) addTraceInfo(fields []zapcore.Field) []zapcore.Field {
	retField := make([]zapcore.Field, 0, len(fields))
	for _, field := range fields {
		if field.Key == TraceIDFlag || field.Key == SpanIDFlag {
			enc.buf.AppendString(field.Key)
			enc.buf.AppendByte('=')
			enc.buf.AppendString(field.String)
			enc.buf.AppendString(enc.ConsoleSeparator)
		} else {
			retField = append(retField, field)
		}
	}

	return retField
}

func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}
