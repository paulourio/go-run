package urn

type ComparisonPart int

const (
	AssignedName ComparisonPart = iota // "scheme:nid:nss"
	AllParts                           // entire URN
)

type ComparisonMethod int

const (
	// Simple is the simple string comparison from RFC 3986.
	//
	// [RFC 3986 §6.2.1](urn:ietf:rfc:3986#section-6.2.1):
	//
	//   6.2.1.  Simple String Comparison
	//   If two URIs, when considered as character strings, are identical,
	//   then it is safe to conclude that they are equivalent.  This type of
	//   equivalence test has very low computational cost and is in wide use
	//   in a variety of applications, particularly in the domain of parsing.
	Simple ComparisonMethod = iota
	// Normalized is the case-normalization procedure from RFC 3986 and
	// also RFC 8141.
	//
	// [RFC 3986 §6.2.2.1](urn:ietf:rfc:3986#section-6.2.2.1):
	//
	//   For all URIs, the hexadecimal digits within a percent-encoding
	//   triplet (e.g., "%3a" versus "%3A") are case-insensitive and therefore
	//   should be normalized to use uppercase letters for the digits A-F.
	//   When a URI uses components of the generic syntax, the component
	//   syntax equivalence rules always apply; namely, that the scheme and
	//   host are case-insensitive and therefore should be normalized to
	//   lowercase.  For example, the URI <HTTP://www.EXAMPLE.com/> is
	//   equivalent to <http://www.example.com/>.  The other generic syntax
	//   components are assumed to be case-sensitive unless specifically
	//   defined otherwise by the scheme (see Section 6.2.3).
	//
	// [RFC 8141 §3.1](urn:ietf:rfc:8141#section-3.1)
	//
	//   Two URNs are URN-equivalent if their assigned-name portions are
	//   octet-by-octet equal after applying case normalization (as specified
	//   in Section 6.2.2.1 of [RFC3986]) to the following constructs:
	//
	//   1.  the URI scheme "urn", by conversion to lower case
	//
	//   2.  the NID, by conversion to lower case
	//
	//   3.  any percent-encoded characters in the NSS (that is, all character
	//       triplets that match the <pct-encoding> production found in
	//       Section 2.1 of the base URI specification [RFC3986]), by
	//       conversion to upper case for the digits A-F.
	CaseNormalized
	// EncodingNormalized is the percent-encoding normalization from RFC 3986.
	//
	// [RFC 3986 §6.2.2.2](urn:ietf:rfc:3986#section-6.2.2.2):
	//
	//   The percent-encoding mechanism (Section 2.1) is a frequent source of
	//   variance among otherwise identical URIs.  In addition to the case
	//   normalization issue noted above, some URI producers percent-encode
	//   octets that do not require percent-encoding, resulting in URIs that
	//   are equivalent to their non-encoded counterparts.  These URIs should
	//   be normalized by decoding any percent-encoded octet that corresponds
	//   to an unreserved character, as described in Section 2.3.
	EncodingNormalized
)

// EqualString returns whether the string representations of the two URNs
// match.
func EqualAssignedName(a *URN, b *URN) bool {
	if a == b {
		return true
	}

	return a.AssignedName() == b.AssignedName()
}

// Equal compares two URNs parts according to the specified method.
// For RFC 8141's equivalence test:
//
//	urn.Equal(a, b, urn.AssignedName, urn.CaseNormalized)
func Equal(a *URN, b *URN, part ComparisonPart, method ComparisonMethod) bool {
	if a == b {
		return true
	}

	switch method {
	case Simple:
		// No transformation required.
	case CaseNormalized:
		a = a.Normalized()
		b = b.Normalized()
	case EncodingNormalized:
		a = a.EncodingNormalized()
		b = b.EncodingNormalized()
	}

	switch part {
	case AssignedName:
		return a.AssignedName() == b.AssignedName()
	case AllParts:
		return a.String() == b.String()
	}

	panic("unexpected equal params")
}
