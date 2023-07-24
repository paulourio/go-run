package urn_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/paulourio/go-urn"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Parallel()

	cases := parseCases
	cases = append(cases, rfc2141examples...)
	cases = append(cases, wikiExamples...)

	for i, c := range cases {
		// If normalized is nil, we expect to be already normalized, so
		// no changes.
		if c.Norm == nil {
			c.Norm = c.URN
		}

		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			testParse(t, c)
		})
	}
}

func testParse(t *testing.T, c *urnTestCase) {
	msg := fmt.Sprintf("Input: %#v", c.Input)

	id, err := urn.Parse(c.Input)

	if c.Err == nil {
		if assert.NoError(t, err, msg) {
			if assert.Equal(t, c.URN, id, "Parse correctly - "+msg) {
				m := "Normalize identifier - " + msg
				norm := id.Normalized()

				assert.Equal(t, c.Norm.Scheme, norm.Scheme, m)
				assert.Equal(t, c.Norm.NID, norm.NID, m)
				assert.Equal(t, c.Norm.NSS, norm.NSS, m)
				assert.Equal(t, c.Norm.Resolve, norm.Resolve, m)
				assert.Equal(t, c.Norm.Query, norm.Query, m)
				assert.Equal(t, c.Norm.Fragment, norm.Fragment, m)

				// Test single string comparison.
				testEquivalences(t, id, norm)

				// Reconstruction must match input.
				assert.Equal(t, c.Input, id.String())

				testMarshaling(t, c.Input, id)
				testAssignedName(t, c.Input, id)

				// Test copy is complete.
				urnCopy := c.URN.Copy()
				if assert.NotSame(t, c.URN, urnCopy) {
					assert.Equal(t, c.URN, urnCopy)
				}
			}
		}
	} else {
		var uerr *urn.Error

		if err != nil {
			t.Log(err.Error())
		} else {
			t.Logf("No error raised, resulted in %#v", id)
		}

		// Error may be optionally wrapped by *urn.Error
		if assert.ErrorAs(t, err, &uerr, msg+"\nExpected: "+c.Err.Error()) {
			assert.EqualError(t, uerr.Err, c.Err.Error(), msg)
		}
	}
}

func testEquivalences(t *testing.T, id *urn.URN, norm *urn.URN) {
	var v bool

	t.Logf("ID: %#v", id)
	t.Logf("Norm: %#v", id)

	// Tests with itself
	v = urn.Equal(id, id, urn.AssignedName, urn.Simple)
	assert.True(t, v, "Equal to itself")

	v = urn.Equal(id, id, urn.AllParts, urn.Simple)
	assert.True(t, v, "Equal to itself")

	v = urn.Equal(id, id, urn.AssignedName, urn.CaseNormalized)
	assert.True(t, v, "Equal to itself, assigned-name, case-normalized")

	v = urn.Equal(id, id, urn.AllParts, urn.CaseNormalized)
	assert.True(t, v, "Equal to itself, full, case-normalized")

	v = urn.Equal(id, id, urn.AssignedName, urn.EncodingNormalized)
	assert.True(t, v, "Equal to itself, assigned-name, encoding-normalized")

	v = urn.Equal(id, id, urn.AllParts, urn.EncodingNormalized)
	assert.True(t, v, "Equal to itself, full, encoding-normalized")

	// Tests against normalized and not-normalized
	v = urn.Equal(norm, id, urn.AssignedName, urn.CaseNormalized)
	assert.True(t, v, "Equal to itself, assigned-name, case-normalized")

	v = urn.Equal(norm, id, urn.AllParts, urn.CaseNormalized)
	assert.True(t, v, "Equal to itself, full, case-normalized")

	v = urn.Equal(norm, id, urn.AssignedName, urn.EncodingNormalized)
	assert.True(t, v, "Equal to itself, assigned-name, encoding-normalized")

	v = urn.Equal(norm, id, urn.AllParts, urn.EncodingNormalized)
	assert.True(t, v, "Equal to itself, full, encoding-normalized")
}

func testAssignedName(t *testing.T, s string, id *urn.URN) {
	pos := len(s)

	if qpos := strings.Index(s, "?"); qpos != -1 && qpos < pos {
		pos = qpos
	}

	if fpos := strings.Index(s, "#"); fpos != -1 && fpos < pos {
		pos = fpos
	}

	assert.Equal(t, s[:pos], id.AssignedName())
}

func testMarshaling(t *testing.T, s string, id *urn.URN) {
	d, err := json.Marshal(id)

	if assert.NoError(t, err, "marshaling "+s) {
		expected := `"` + strings.ReplaceAll(s, "&", `\u0026`) + `"`

		if assert.Equal(t, expected, string(d)) {
			var u urn.URN

			uerr := json.Unmarshal(d, &u)
			if assert.NoError(t, uerr) {
				assert.Equal(t, s, u.String())
			}
		}
	}
}

type urnTestCase struct {
	Input string
	Err   error
	URN   *urn.URN
	Norm  *urn.URN
}

var rfc2141examples = []*urnTestCase{
	{
		Input: "URN:foo:a123,456",
		URN:   &urn.URN{Scheme: "URN", NID: "foo", NSS: "a123,456"},
		Norm:  &urn.URN{Scheme: "urn", NID: "foo", NSS: "a123,456"},
	},
	{
		Input: "urn:foo:a123,456",
		URN:   &urn.URN{Scheme: "urn", NID: "foo", NSS: "a123,456"},
	},
	{
		Input: "urn:FOO:a123,456",
		URN:   &urn.URN{Scheme: "urn", NID: "FOO", NSS: "a123,456"},
		Norm:  &urn.URN{Scheme: "urn", NID: "foo", NSS: "a123,456"},
	},
	{
		Input: "urn:foo:A123,456",
		URN:   &urn.URN{Scheme: "urn", NID: "foo", NSS: "A123,456"},
	},
	{
		Input: "urn:foo:a123%2C456",
		URN:   &urn.URN{Scheme: "urn", NID: "foo", NSS: "a123%2C456"},
	},
	{
		Input: "URN:FOO:a123%2c456",
		URN:   &urn.URN{Scheme: "URN", NID: "FOO", NSS: "a123%2c456"},
		Norm:  &urn.URN{Scheme: "urn", NID: "foo", NSS: "a123%2C456"},
	},
}

var parseCases = []*urnTestCase{
	// Invalid scheme.
	{Input: "", Err: urn.ErrInvalidScheme},
	{Input: ":", Err: urn.ErrInvalidScheme},
	{Input: "urn", Err: urn.ErrInvalidScheme},
	{Input: "urn!", Err: urn.ErrInvalidScheme},
	{Input: ":urn:", Err: urn.ErrInvalidScheme},
	{Input: " urn:", Err: urn.ErrInvalidScheme},
	// Invalid NID.
	{Input: "urn:", Err: urn.ErrInvalidNID},
	{Input: "urn::", Err: urn.ErrInvalidNID},
	{Input: "urn:a:", Err: urn.ErrInvalidNID},
	{Input: "urn:a:b", Err: urn.ErrInvalidNID},
	{Input: "urn:a-:b", Err: urn.ErrInvalidNID},
	{Input: "urn:-a:b", Err: urn.ErrInvalidNID},
	{Input: "urn:--:b", Err: urn.ErrInvalidNID},
	{Input: "urn:-a-:b", Err: urn.ErrInvalidNID},
	{Input: "urn:aa", Err: urn.ErrInvalidNID},
	{Input: "urn:a a:b", Err: urn.ErrInvalidNID},
	{Input: "urn:a~a:b", Err: urn.ErrInvalidNID},
	{Input: "urn:a%20c:x", Err: urn.ErrInvalidNID},
	{Input: "urn:my.path:b", Err: urn.ErrInvalidNID},
	{Input: `urn:bb:"abc"`, Err: urn.ErrInvalidNSS},
	{Input: fmt.Sprintf("urn:a%sa:b", strings.Repeat("-", 31)), Err: urn.ErrInvalidNID},
	// Valid NID.
	{
		Input: "urn:aa:b",
		URN:   &urn.URN{Scheme: "urn", NID: "aa", NSS: "b"},
	},
	{
		Input: "urn:a-a:b",
		URN:   &urn.URN{Scheme: "urn", NID: "a-a", NSS: "b"},
	},
	{
		Input: fmt.Sprintf("urn:a%sa:b", strings.Repeat("-", 30)),
		URN: &urn.URN{
			Scheme: "urn",
			NID:    fmt.Sprintf("a%sa", strings.Repeat("-", 30)),
			NSS:    "b",
		},
	},
	{
		Input: "Urn:abcdefghilmnopqrstuvzabcdefghilm:x",
		URN: &urn.URN{
			Scheme: "Urn",
			NID:    "abcdefghilmnopqrstuvzabcdefghilm",
			NSS:    "x",
		},
		Norm: &urn.URN{
			Scheme: "urn",
			NID:    "abcdefghilmnopqrstuvzabcdefghilm",
			NSS:    "x",
		},
	},
	{
		Input: "urn:123:x",
		URN:   &urn.URN{Scheme: "urn", NID: "123", NSS: "x"},
	},
	{
		Input: "urn:1ab:x",
		URN:   &urn.URN{Scheme: "urn", NID: "1ab", NSS: "x"},
	},
	{
		Input: "urn:cd1:x",
		URN:   &urn.URN{Scheme: "urn", NID: "cd1", NSS: "x"},
	},
	{
		Input: "urn:colon::::;",
		URN:   &urn.URN{Scheme: "urn", NID: "colon", NSS: ":::;"},
	},
	{
		Input: "urn:foo:my.path",
		URN:   &urn.URN{Scheme: "urn", NID: "foo", NSS: "my.path"},
	},
	{
		Input: "urn:foo:=@",
		URN:   &urn.URN{Scheme: "urn", NID: "foo", NSS: "=@"},
	},
	{
		Input: "urn:foo:@!=%2C(xyz)+a,b.*@g=$_'",
		URN:   &urn.URN{Scheme: "urn", NID: "foo", NSS: "@!=%2C(xyz)+a,b.*@g=$_'"},
	},
	{
		Input: "Urn:Xx:abc%1Dz%2F%3az",
		URN:   &urn.URN{Scheme: "Urn", NID: "Xx", NSS: "abc%1Dz%2F%3az"},
		Norm:  &urn.URN{Scheme: "urn", NID: "xx", NSS: "abc%1Dz%2F%3Az"},
	},
	// NSS Cannot start with a slash.
	{Input: "urn:foo:/", Err: urn.ErrInvalidNSS},
	{Input: "urn:foo:/p", Err: urn.ErrInvalidNSS},
	// NSS with invalid percent escape.
	{Input: "urn:foo:a%P", Err: urn.ErrInvalidNSS},
	{Input: "urn:foo:a%P2", Err: urn.ErrInvalidNSS},
	{
		Input: "urn:foo:p/",
		URN:   &urn.URN{Scheme: "urn", NID: "foo", NSS: "p/"},
	},
	{
		Input: "urn:foo:%2c",
		URN:   &urn.URN{Scheme: "urn", NID: "foo", NSS: "%2c"},
		Norm:  &urn.URN{Scheme: "urn", NID: "foo", NSS: "%2C"},
	},
	{
		Input: "urn:example:apple:pear:plum:cherry",
		URN:   &urn.URN{Scheme: "urn", NID: "example", NSS: "apple:pear:plum:cherry"},
	},
	// Invalid r-component
	{Input: "urn:abc:def?+?", Err: urn.ErrInvalidResolve},
	{Input: "urn:abc:def?+", Err: urn.ErrInvalidResolve},
	{Input: "urn:abc:def?+/abc", Err: urn.ErrInvalidResolve},
	// Valid r-component
	{
		Input: "urn:example:%2a?+%2b",
		URN:   &urn.URN{Scheme: "urn", NID: "example", NSS: "%2a", Resolve: "%2b"},
		Norm:  &urn.URN{Scheme: "urn", NID: "example", NSS: "%2A", Resolve: "%2B"},
	},
	{
		Input: "urn:example:a?+%1A%2B",
		URN:   &urn.URN{Scheme: "urn", NID: "example", NSS: "a", Resolve: "%1A%2B"},
	},
	{
		Input: "urn:ab:a?+%2b????",
		URN:   &urn.URN{Scheme: "urn", NID: "ab", NSS: "a", Resolve: "%2b????"},
		Norm:  &urn.URN{Scheme: "urn", NID: "ab", NSS: "a", Resolve: "%2B????"},
	},
	{
		Input: "urn:ab:a?+@?",
		URN:   &urn.URN{Scheme: "urn", NID: "ab", NSS: "a", Resolve: "@?"},
	},
	{
		Input: "urn:example:foo-bar-baz-qux?+CCResolve:cc=uk",
		URN: &urn.URN{
			Scheme:  "urn",
			NID:     "example",
			NSS:     "foo-bar-baz-qux",
			Resolve: "CCResolve:cc=uk",
		},
	},
	{
		Input: "urn:lex:it:ministero.giustizia:decreto:1992-07-24;358~art5",
		URN: &urn.URN{
			Scheme: "urn",
			NID:    "lex",
			NSS:    "it:ministero.giustizia:decreto:1992-07-24;358~art5",
		},
	},
	// Invalid q-component
	{Input: "urn:abc:def?=?", Err: urn.ErrInvalidQuery},
	{Input: "urn:abc:def?=", Err: urn.ErrInvalidQuery},
	{Input: "urn:abc:def?=/abc", Err: urn.ErrInvalidQuery},
	// Valid q-component
	{
		Input: "urn:aa:b?=c",
		URN:   &urn.URN{Scheme: "urn", NID: "aa", NSS: "b", Query: "c"},
	},
	{
		Input: "urn:aa:b?=%B3",
		URN:   &urn.URN{Scheme: "urn", NID: "aa", NSS: "b", Query: "%B3"},
	},
	{
		Input: "urn:aa:b?=a%b3",
		URN:   &urn.URN{Scheme: "urn", NID: "aa", NSS: "b", Query: "a%b3"},
		Norm:  &urn.URN{Scheme: "urn", NID: "aa", NSS: "b", Query: "a%B3"},
	},
	{
		Input: "urn:example:weather?=op=map&lat=39.56&lon=-104.85&datetime=1969-07-21T02:56:15Z",
		URN: &urn.URN{
			Scheme: "urn",
			NID:    "example",
			NSS:    "weather",
			Query:  "op=map&lat=39.56&lon=-104.85&datetime=1969-07-21T02:56:15Z",
		},
	},
	// Valid fragment
	{
		Input: "urn:aa:b#c",
		URN:   &urn.URN{Scheme: "urn", NID: "aa", NSS: "b", Fragment: "c"},
	},
	{
		Input: "urn:aa:b#%B3",
		URN:   &urn.URN{Scheme: "urn", NID: "aa", NSS: "b", Fragment: "%B3"},
	},
	{
		Input: "urn:aa:b#a%B3",
		URN:   &urn.URN{Scheme: "urn", NID: "aa", NSS: "b", Fragment: "a%B3"},
	},
	// Fragment may be empty, and we test that rebuilding from parsed
	// URN will result in the same input.
	{
		Input: "urn:aa:b#",
		URN:   &urn.URN{Scheme: "urn", NID: "aa", NSS: "b", ForceFragment: true},
	},
	{
		Input: "urn:Aa:b#",
		URN:   &urn.URN{Scheme: "urn", NID: "Aa", NSS: "b", ForceFragment: true},
		Norm:  &urn.URN{Scheme: "urn", NID: "aa", NSS: "b", ForceFragment: true},
	},
	{
		Input: "urn:example:foo-bar-baz-qux#somepart",
		URN: &urn.URN{
			Scheme:   "urn",
			NID:      "example",
			NSS:      "foo-bar-baz-qux",
			Fragment: "somepart",
		},
	},
	// [RFC 8141 §2](urn:ietf:rfc:8141#section-2)
	//
	//   The question mark character "?" can be used without percent-encoding
	//   inside r-components, q-components, and f-components.  Other than
	//   inside those components, a "?" that is not immediately followed by
	//   "=" or "+" is not defined for URNs and SHOULD be treated as a syntax
	//   error by URN-specific parsers and other processors.
	//
	// This means that we can have "?+" inside r-components, and
	// both "?=" and "?+" inside q-components and fragments.
	// The parser should not raise an error.
	{
		Input: "urn:foo:bar?++?+?+?==?=?=?+?#+?+?+?==?=?=?",
		URN: &urn.URN{
			Scheme:   "urn",
			NID:      "foo",
			NSS:      "bar",
			Resolve:  "+?+?+",
			Query:    "=?=?=?+?",
			Fragment: "+?+?+?==?=?=?",
		},
	},
	{
		Input: "urn:abc:def?=a?+b",
		URN: &urn.URN{
			Scheme: "urn",
			NID:    "abc",
			NSS:    "def",
			Query:  "a?+b",
		},
	},
	// Invalid sequence of "?" outside components.
	{Input: "urn:abc:def?#?+Resolve", Err: urn.ErrInvalidIdentifier},
	// Valid fragment.
	{
		Input: "urn:abc:def#fragment?=a?+b?",
		URN: &urn.URN{
			Scheme:   "urn",
			NID:      "abc",
			NSS:      "def",
			Fragment: "fragment?=a?+b?",
		},
	},
}

var wikiExamples = []*urnTestCase{
	// The 1968 book The Last Unicorn, identified by its International
	// Standard Book Number.
	{
		Input: "urn:isbn:0451450523",
		URN:   &urn.URN{Scheme: "urn", NID: "isbn", NSS: "0451450523"},
	},
	// The 2002 film Spider-Man, identified by its International
	// Standard Audiovisual Number.
	{
		Input: "urn:isan:0000-0000-2CEA-0000-1-0000-0000-Y",
		URN: &urn.URN{
			Scheme: "urn",
			NID:    "isan",
			NSS:    "0000-0000-2CEA-0000-1-0000-0000-Y",
		},
	},
	// The scientific journal Science of Computer Programming,
	// identified by its International Standard Serial Number.
	{
		Input: "urn:ISSN:0167-6423",
		URN:   &urn.URN{Scheme: "urn", NID: "ISSN", NSS: "0167-6423"},
		Norm:  &urn.URN{Scheme: "urn", NID: "issn", NSS: "0167-6423"},
	},
	// The IETF's RFC 2648.
	{
		Input: "urn:ietf:rfc:2648",
		URN:   &urn.URN{Scheme: "urn", NID: "ietf", NSS: "rfc:2648"},
	},
	// The default namespace rules for MPEG-7 video metadata.
	{
		Input: "urn:mpeg:mpeg7:schema:2001",
		URN:   &urn.URN{Scheme: "urn", NID: "mpeg", NSS: "mpeg7:schema:2001"},
	},
	// The OID for the United States.
	{
		Input: "urn:oid:2.16.840",
		URN:   &urn.URN{Scheme: "urn", NID: "oid", NSS: "2.16.840"},
	},
	// A version 1 UUID.
	{
		Input: "urn:uuid:6e8bc430-9c3a-11d9-9669-0800200c9a66",
		URN: &urn.URN{
			Scheme: "urn",
			NID:    "uuid",
			NSS:    "6e8bc430-9c3a-11d9-9669-0800200c9a66",
		},
	},
	// A National Bibliography Number for a document, indicating
	// country (de), regional network (bvb = Bibliotheksverbund Bayern),
	// library number (19) and document number.
	{
		Input: "urn:nbn:de:bvb:19-146642",
		URN:   &urn.URN{Scheme: "urn", NID: "nbn", NSS: "de:bvb:19-146642"},
	},
	// A directive of the European Union, using the proposed Lex URN
	// namespace.
	{
		Input: "urn:lex:eu:council:directive:2010-03-09;2010-19-UE",
		URN: &urn.URN{
			Scheme: "urn",
			NID:    "lex",
			NSS:    "eu:council:directive:2010-03-09;2010-19-UE",
		},
	},
	// A Life Science Identifiers that may be resolved to
	// http://zoobank.org/urn:lsid:zoobank.org:pub:CDC8D258-8F57-41DC-B560-247E17D3DC8C
	{
		Input: "urn:lsid:zoobank.org:pub:CDC8D258-8F57-41DC-B560-247E17D3DC8C",
		URN: &urn.URN{
			Scheme: "urn",
			NID:    "lsid",
			NSS:    "zoobank.org:pub:CDC8D258-8F57-41DC-B560-247E17D3DC8C",
		},
	},
	// Global Trade Item Number with lot/batch number. As defined by Tag
	// Data Standard[11] (TDS). See more examples at EPC Identification Keys.
	{
		Input: "urn:epc:class:lgtin:4012345.012345.998877",
		URN: &urn.URN{
			Scheme: "urn",
			NID:    "epc",
			NSS:    "class:lgtin:4012345.012345.998877",
		},
	},
	// Global Trade Item Number with an individual serial number.
	{
		Input: "urn:epc:id:sgtin:0614141.112345.400",
		URN: &urn.URN{
			Scheme: "urn",
			NID:    "epc", NSS: "id:sgtin:0614141.112345.400"},
	},
	// Serial Shipping Container Code.
	{
		Input: "urn:epc:id:sscc:0614141.1234567890",
		URN: &urn.URN{
			Scheme: "urn",
			NID:    "epc", NSS: "id:sscc:0614141.1234567890"},
	},
	// Global Location Number with extension.
	{
		Input: "urn:epc:id:sgln:0614141.12345.400",
		URN: &urn.URN{
			Scheme: "urn",
			NID:    "epc",
			NSS:    "id:sgln:0614141.12345.400",
		},
	},
	// BIC Intermodal container Code as per ISO 6346.
	{
		Input: "urn:epc:id:bic:CSQU3054383",
		URN:   &urn.URN{Scheme: "urn", NID: "epc", NSS: "id:bic:CSQU3054383"},
	},
	// IMO Vessel Number of marine vessels
	{
		Input: "urn:epc:id:imovn:9176187",
		URN:   &urn.URN{Scheme: "urn", NID: "epc", NSS: "id:imovn:9176187"},
	},
	// Global Document Type Identifier of a document instance
	{
		Input: "urn:epc:id:gdti:0614141.12345.400",
		URN: &urn.URN{
			Scheme: "urn",
			NID:    "epc",
			NSS:    "id:gdti:0614141.12345.400",
		},
	},
	// Identifier for Marine Aids to Navigation
	{
		Input: "urn:mrn:iala:aton:us:1234.5",
		URN:   &urn.URN{Scheme: "urn", NID: "mrn", NSS: "iala:aton:us:1234.5"},
	},
	// Identifier for Vessel Traffic Services
	{
		Input: "urn:mrn:iala:vts:ca:ecareg",
		URN:   &urn.URN{Scheme: "urn", NID: "mrn", NSS: "iala:vts:ca:ecareg"},
	},
	// Identifier for Waterways
	{
		Input: "urn:mrn:iala:wwy:us:atl:chba:potri",
		URN: &urn.URN{
			Scheme: "urn",
			NID:    "mrn",
			NSS:    "iala:wwy:us:atl:chba:potri",
		},
	},
	// Identifier for IALA publications
	{
		Input: "urn:mrn:iala:pub:g1143",
		URN:   &urn.URN{Scheme: "urn", NID: "mrn", NSS: "iala:pub:g1143"},
	},
	// Identifier for federated identity; this example is from Claims X-Ray
	{
		Input: "urn:microsoft:adfs:claimsxray",
		URN:   &urn.URN{Scheme: "urn", NID: "microsoft", NSS: "adfs:claimsxray"},
	},
	// European Network of Transmission System Operators for Electricity
	// (ENTSO-E), identified by its Energy Identification Code.
	{
		Input: "urn:eic:10X1001A1001A450",
		URN:   &urn.URN{Scheme: "urn", NID: "eic", NSS: "10X1001A1001A450"},
	},
}
