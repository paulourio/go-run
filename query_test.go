package urn_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/paulourio/go-urn"
	"github.com/stretchr/testify/assert"
)

func TestParseQuery(t *testing.T) {
	t.Parallel()

	for i, c := range parseQueryCases {
		msg := fmt.Sprintf("case %d: %q", i+1, c.Input)

		u, err := urn.Parse(c.Input)
		if assert.NoError(t, err, msg) {
			kv, perr := url.ParseQuery(u.Query)

			if assert.NoError(t, perr) {
				assert.Equal(t, c.Values, kv)
			}
		}
	}
}

type parseQueryTestCase struct {
	Input  string
	Values url.Values
}

var parseQueryCases = []*parseQueryTestCase{
	{
		Input: "urn:example:weather?=op=map&lat=39.56&lon=-104.85&datetime=1969-07-21T02:56:15Z",
		Values: url.Values{
			"op":       []string{"map"},
			"lat":      []string{"39.56"},
			"lon":      []string{"-104.85"},
			"datetime": []string{"1969-07-21T02:56:15Z"},
		},
	},
	{
		Input: "urn:foo:barr?=a=2&a=1",
		Values: url.Values{
			"a": []string{"2", "1"},
		},
	},
	{
		Input: "urn:foo:barr?=a=2&a=1%26a=2",
		Values: url.Values{
			"a": []string{"2", "1&a=2"},
		},
	},
}
