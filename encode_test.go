package urn_test

import (
	"fmt"
	"testing"

	"github.com/paulourio/go-urn"
	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	t.Parallel()

	for i, c := range encodingCases {
		var (
			repr    string
			strmode bool
		)

		if c.String != "" {
			strmode = true
			repr = fmt.Sprintf("%q", c.String)
		} else {
			repr = fmt.Sprintf("%#v", c.Bytes)
		}

		msg := fmt.Sprintf("case %d: %v %s", i+1, c.Op, repr)

		switch c.Op {
		case Decode:
			assert.Empty(t, c.Bytes) // Decode enabled only for strings.
			assert.Equal(t, c.Decode, urn.Decode(c.String), msg)
		case EncodeComp:
			if strmode {
				r := urn.EncodeStringComponent(c.String)
				assert.Equal(t, c.Encode, r, msg)
			} else {
				r := urn.EncodeComponent(c.Bytes)
				assert.Equal(t, c.Encode, r, msg)
			}
		case EncodeNSS:
			if strmode {
				r := urn.EncodeStringNSS(c.String)
				assert.Equal(t, c.Encode, r, msg)
			} else {
				r := urn.EncodeNSS(c.Bytes)
				assert.Equal(t, c.Encode, r, msg)
			}
		case RecodeComp:
			assert.Empty(t, c.Bytes) // Recode enabled only for strings.
			assert.Equal(t, c.Encode, urn.RecodeStringComponent(c.String), msg)
		case RecodeNSS:
			assert.Empty(t, c.Bytes) // Recode enabled only for strings.
			assert.Equal(t, c.Encode, urn.RecodeStringNSS(c.String), msg)
		}
	}
}

type EncodingOp int

const (
	Decode EncodingOp = iota + 1
	EncodeComp
	EncodeNSS
	RecodeComp
	RecodeNSS
)

func (e EncodingOp) String() string {
	switch e {
	case Decode:
		return "Decode"
	case EncodeComp:
		return "EncodeComponent"
	case EncodeNSS:
		return "EncodeNSS"
	case RecodeComp:
		return "RecodeComponent"
	case RecodeNSS:
		return "RecodeNSS"
	}

	panic("unknown op")
}

type encodingCase struct {
	Op EncodingOp

	// Only one of the three must be filled.
	Bytes  []byte
	String string

	// Only one of the three must be filled.
	Decode []byte
	Encode string
	Err    error
}

var encodingCases = []*encodingCase{
	// Decode
	{Op: Decode, String: "", Decode: []byte{}},
	{Op: Decode, String: "abc", Decode: []byte("abc")},
	{Op: Decode, String: "@!=%2c(xyz)+a,b.*@g=$_", Decode: []byte("@!=,(xyz)+a,b.*@g=$_")},
	{Op: Decode, String: "@!=%2C(xyz)+a,b.*@g=$_", Decode: []byte("@!=,(xyz)+a,b.*@g=$_")},
	{Op: Decode, String: "%20", Decode: []byte(" ")},
	{Op: Decode, String: "%41%00%1A", Decode: []byte{0x41, 0x0, 0x1a}},
	// Encode
	{Op: EncodeComp, String: "", Encode: ""},
	{Op: EncodeNSS, String: "", Encode: ""},
	{Op: EncodeComp, Bytes: []byte("abc"), Encode: "abc"},
	{Op: EncodeNSS, Bytes: []byte("abc"), Encode: "abc"},
	{Op: EncodeComp, Bytes: []byte(" "), Encode: "%20"},
	{Op: EncodeNSS, Bytes: []byte(" "), Encode: "%20"},
	{Op: EncodeComp, Bytes: []byte("?a&b=1/"), Encode: "?a&b=1/"},
	{Op: EncodeNSS, Bytes: []byte("?a&b=1/"), Encode: "%3Fa&b=1/"},
	{Op: EncodeNSS, Bytes: []byte("?a&b=1/"), Encode: "%3Fa&b=1/"},
	{Op: EncodeNSS, String: "?a&b=1/", Encode: "%3Fa&b=1/"},
	{Op: EncodeNSS, String: "?a&b=1/", Encode: "%3Fa&b=1/"},
	{Op: EncodeNSS, Bytes: []byte{0x0, 0x10, 0x20, 0x61}, Encode: "%00%10%20a"},
	{Op: EncodeNSS, Bytes: []byte{0x0, 0x10, 0x20, 0x61}, Encode: "%00%10%20a"},
	// Recode
	{Op: RecodeComp, String: "abc", Encode: "abc"},
	{Op: RecodeComp, String: "ab%32c%61", Encode: "ab2ca"},
	{Op: RecodeComp, String: "%32%33%34", Encode: "234"},
	{Op: RecodeComp, String: "%32%33%34%3f", Encode: "234?"},
	{Op: RecodeNSS, String: "abc", Encode: "abc"},
	{Op: RecodeNSS, String: "ab%32c%61", Encode: "ab2ca"},
	{Op: RecodeNSS, String: "%32%33%34", Encode: "234"},
	{Op: RecodeNSS, String: "%32%33%34%3f", Encode: "234%3F"},
}
