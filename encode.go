package urn

import (
	"strings"
)

// Decode unescapes a string and returns a byte slice.
func Decode(d string) []byte {
	count := strings.Count(d, "%")
	data := make([]byte, len(d)-count*2)
	n := len(d)
	i := 0
	pos := 0

	for i < n {
		c := d[i]

		if i+2 < n && c == '%' && isHex(d[i+1]) && isHex(d[i+2]) {
			data[pos] = unhex(d[i+1])<<4 + unhex(d[i+2])
			i += 3
			pos++
		} else {
			data[pos] = c
			i++
			pos++
		}
	}

	return data
}

// EncodeStringComponent escapes a string so that it is suitable for use
// as a Resolve, Query, or Fragment component of a URN.
func EncodeStringComponent(d string) string {
	return EncodeComponent([]byte(d))
}

// EncodeStringComponent escapes a byte slice so that it is suitable for
// use as a Resolve, Query, or Fragment component of a URN.
func EncodeComponent(d []byte) string {
	e := NewEncoder(WithKeepUnescaped('/', '?'))
	return e.Encode(d)
}

// EncodeStringNSS escapes a string so that it is suitable for use
// as the NSS part of the URN identifier.
func EncodeStringNSS(d string) string {
	return EncodeNSS([]byte(d))
}

// EncodeNSS escapes a byte slice so that it is suitable for use
// as a Resolve, Query, or Fragment component of a URN.
func EncodeNSS(d []byte) string {
	e := NewEncoder(WithKeepUnescaped('/'))
	return e.Encode(d)
}

// RecodeStringComponent will decode and encode a percent-encoded string
// to guarantee that only the necessary characters are escaped.  This is
// particularly useful when performing Percent-Encoding Normalized
// Comparison.
func RecodeStringComponent(d string) string {
	return EncodeComponent(Decode(d))
}

// RecodeStringNSS will decode and encode a percent-encoded string
// to guarantee that only the necessary characters are escaped.  This is
// particularly useful when performing Percent-Encoding Normalized
// Comparison.
func RecodeStringNSS(d string) string {
	return EncodeNSS(Decode(d))
}

// Encoder is a type for escaping sequence of data that are inside
// one component of a URN.
type Encoder struct {
	Allowed []byte // additional single bytes that should not escape
}

// Encode escapes bytes that are reserved for URI syntax, according to
// the rules from RFC 3986.
func (e *Encoder) Encode(d []byte) string {
	if len(d) == 0 {
		return ""
	}

	impl := &escaper{
		input: d,
		n:     len(d),
		allow: e.Allowed,
	}

	size := impl.computeOutputSize()
	if size == len(d) {
		// No escaping is necessary.
		return string(d)
	}

	impl.data = make([]byte, size)
	impl.Encode()

	return string(impl.data)
}

// NewEncoder returns a new encoder for the given options.
func NewEncoder(opts ...func(*Encoder)) *Encoder {
	e := &Encoder{}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

// WithKeepUnescaped defines a list of single bytes that should not be
// escaped by an encoder.
func WithKeepUnescaped(b ...byte) func(*Encoder) {
	return func(e *Encoder) {
		e.Allowed = append(e.Allowed, b...)
	}
}

type escaper struct {
	// Input
	input []byte
	n     int
	allow []byte

	// Working data
	data []byte
}

func (e *escaper) Encode() {
	pos := 0

	for i := 0; i < e.n; i++ {
		c := e.input[i]
		if e.shouldEscape(c) {
			e.data[pos] = '%'
			e.data[pos+1] = upperHex[c>>4]
			e.data[pos+2] = upperHex[c&0xF]
			pos += 3
		} else {
			e.data[pos] = c
			pos++
		}
	}
}

func (e *escaper) computeOutputSize() (sz int) {
	sz = len(e.input)

	for _, c := range e.input {
		if e.shouldEscape(c) {
			sz += 2
		}
	}

	return
}

func (e *escaper) shouldEscape(c byte) bool {
	if isPCharSingle(c) {
		return false
	}

	for _, b := range e.allow {
		if c == b {
			return false
		}
	}

	return true
}

func unhex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}

const upperHex = "0123456789ABCDEF"
