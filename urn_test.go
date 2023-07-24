package urn_test

import (
	"fmt"
	"testing"

	"github.com/paulourio/go-urn"
	"github.com/stretchr/testify/assert"
)

func TestEquivalenceRFC2141(t *testing.T) {
	t.Parallel()

	for i, c := range rfc2141equivExamples {
		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			testEquivalence(t, c)
		})
	}
}

func TestEquivalenceRFC8141(t *testing.T) {
	t.Parallel()

	cases := make([]*urnEquivalenceTestCase, 0, 120)

	// First and second groups are all equal to each other.
	for i := 1; i <= 6; i++ {
		for j := 1; j <= 6; j++ {
			cases = append(cases, equivTestRFC8141(i, j, true))
		}
	}

	// Third group is distinct from each other and previous groups.
	for i := 7; i <= 9; i++ {
		for j := 1; j <= 9; j++ {
			cases = append(cases, equivTestRFC8141(i, j, i == j))
		}
	}

	// Four group is equiv from each other but different from all others.
	for i := 10; i <= 11; i++ {
		for j := 1; j <= 11; j++ {
			cases = append(cases, equivTestRFC8141(i, j, j >= 10))
		}
	}

	// Fifth group is not equiv with the first, second, and third group.
	for i := 12; i <= 13; i++ {
		for j := 1; j <= 9; j++ {
			cases = append(cases, equivTestRFC8141(i, j, j >= 10))
		}
	}

	// Sixth group is not equiv no other group.
	for j := 1; j <= 14; j++ {
		cases = append(cases, equivTestRFC8141(14, j, j == 14))
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			testEquivalence(t, c)
		})
	}
}

func testEquivalence(t *testing.T, c *urnEquivalenceTestCase) {
	a, aerr := urn.Parse(c.InputA)
	b, berr := urn.Parse(c.InputB)

	if !assert.NoError(t, aerr) && !assert.NoError(t, berr) {
		eqAB := urn.Equal(a, b, urn.AssignedName, urn.CaseNormalized)
		eqBA := urn.Equal(b, a, urn.AssignedName, urn.CaseNormalized)

		if c.Equal {
			assert.True(t, eqAB)
			assert.True(t, eqBA)
		} else {
			assert.False(t, eqAB)
			assert.False(t, eqBA)
		}
	}
}

type urnEquivalenceTestCase struct {
	InputA string
	InputB string
	Equal  bool
}

func equivTestRFC2141(i, j int, eq bool) *urnEquivalenceTestCase {
	return &urnEquivalenceTestCase{
		InputA: rfc2141[i-1],
		InputB: rfc2141[j-1],
		Equal:  eq,
	}
}

func equivTestRFC8141(i, j int, eq bool) *urnEquivalenceTestCase {
	return &urnEquivalenceTestCase{
		InputA: rfc8141[i-1],
		InputB: rfc8141[j-1],
		Equal:  eq,
	}
}

var rfc2141 = []string{
	// URNs 1, 2, and 3 are all lexically equivalent.
	"URN:foo:a123,456", // 1
	"urn:foo:a123,456", // 2
	"urn:FOO:a123,456", // 3
	// URN 4 is not lexically equivalent any of the other URNs.
	"urn:foo:A123,456", // 4
	// URNs 5 and 6 are only lexically equivalent to each other.
	"urn:foo:a123%2C456", // 5
	"URN:FOO:a123%2c456", // 6
}

var rfc2141equivExamples = []*urnEquivalenceTestCase{
	// URNs 1, 2, and 3 are all lexically equivalent.
	equivTestRFC2141(1, 1, true),
	equivTestRFC2141(1, 2, true),
	equivTestRFC2141(1, 3, true),
	equivTestRFC2141(2, 2, true),
	equivTestRFC2141(2, 3, true),
	equivTestRFC2141(2, 5, false),
	equivTestRFC2141(2, 6, false),
	equivTestRFC2141(3, 3, true),
	equivTestRFC2141(3, 5, false),
	equivTestRFC2141(3, 6, false),
	// URN 4 is not lexically equivalent any of the other URNs.
	equivTestRFC2141(4, 1, false),
	equivTestRFC2141(4, 2, false),
	equivTestRFC2141(4, 3, false),
	equivTestRFC2141(4, 4, false),
	equivTestRFC2141(4, 5, false),
	equivTestRFC2141(4, 6, false),
	// URNs 5 and 6 are only lexically equivalent to each other.
	equivTestRFC2141(5, 1, false),
	equivTestRFC2141(5, 5, true),
	equivTestRFC2141(5, 6, true),
	equivTestRFC2141(6, 1, false),
	equivTestRFC2141(6, 5, true),
	equivTestRFC2141(6, 6, true),
}

var rfc8141 = []string{
	// First, because the scheme and NID are case insensitive, the
	// following three URNs are URN-equivalent to each other:
	"urn:example:a123,z456", // 1
	"URN:example:a123,z456", // 2
	"urn:EXAMPLE:a123,z456", // 3
	// Second, because the r-component, q-component, and f-component are
	// not taken into account for purposes of testing URN-equivalence,
	// the following three URNs are URN-equivalent to the first three
	// examples above:
	"urn:example:a123,z456?+abc", // 4
	"urn:example:a123,z456?=xyz", // 5
	"urn:example:a123,z456#789",  // 6
	//  Third, because the "/" character (and anything that follows it)
	// in the NSS is taken into account for purposes of URN-equivalence,
	// the following URNs are not URN-equivalent to each other or to the
	// six preceding URNs:
	"urn:example:a123,z456/foo", // 7
	"urn:example:a123,z456/bar", // 8
	"urn:example:a123,z456/baz", // 9
	//   Fourth, because of percent-encoding, the following URNs are
	// URN-equivalent only to each other and not to any of those above
	// (note that, although %2C is the percent-encoded transformation
	// of "," from the previous examples, such sequences are not decoded
	// for purposes of testing URN-equivalence):
	"urn:example:a123%2Cz456", // 10
	"URN:EXAMPLE:a123%2cz456", // 11
	// Fifth, because characters in the NSS other than percent-encoded
	// sequences are treated in a case-sensitive manner (unless otherwise
	// specified for the URN namespace in question), the following URNs
	// are not URN-equivalent to the first three URNs:
	"urn:example:A123,z456", // 12
	"urn:example:a123,Z456", // 13
	// Sixth, on casual visual inspection of a URN presented in a
	// human-oriented interface, the following URN might appear the same
	// as the first three URNs (because U+0430 CYRILLIC SMALL LETTER A
	// can be confused with U+0061 LATIN SMALL LETTER A), but it is not
	// URN-equivalent to the first three URNs:
	"urn:example:%D0%B0123,z456", // 14
}
