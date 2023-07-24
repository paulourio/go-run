package urn

// I would rather name the URN identifier type urn.Identifier,
// but opted to use urn.URN to be similar to net/url's url.URL.

import (
	"bytes"
	"fmt"
	"strings"
)

// A URN represents a parsed URN.
//
// The general form represented is:
//
//	scheme:nid:nss[?+resolve][?=query][#fragment]
type URN struct {
	Scheme string
	NID    string // Namespace Identifier.
	NSS    string // Namespace Specific String.

	Resolve       string // encoded r-component.
	Query         string // encoded q-component.
	Fragment      string // encoded fragment hint.
	ForceFragment bool   // append symbol ('#') even if Fragment is empty

	normalized bool
}

// AssignedName returns schema:nid:nss.
func (u *URN) AssignedName() string {
	var b strings.Builder

	b.Grow(len(u.Scheme) + len(u.NID) + len(u.NSS) + 2)

	b.WriteString(u.Scheme)
	b.WriteByte(':')
	b.WriteString(u.NID)
	b.WriteByte(':')
	b.WriteString(u.NSS)

	return b.String()
}

// Copy returns a copy of the current URN.
func (u *URN) Copy() *URN {
	return &URN{
		Scheme:        u.Scheme,
		NID:           u.NID,
		NSS:           u.NSS,
		Resolve:       u.Resolve,
		Query:         u.Query,
		Fragment:      u.Fragment,
		ForceFragment: u.ForceFragment,
		normalized:    u.normalized,
	}
}

// Normalize applies case normalization used for URN-equivalence
// procedure.
func (u *URN) Normalize() {
	u.Scheme = strings.ToLower(u.Scheme)
	u.NID = strings.ToLower(u.NID)
	u.NSS = normalizePercentEncoding(u.NSS)
	// Though not strictly required for RFC's 8141 equivalence comparison,
	// we normalize percent encoding of components as well to facilitate
	// RFC 3986's case-normalization comparison  method.
	u.Resolve = normalizePercentEncoding(u.Resolve)
	u.Query = normalizePercentEncoding(u.Query)
	u.Fragment = normalizePercentEncoding(u.Fragment)
	u.normalized = true
}

// EncodingNormalize applies percent-encoding normalization specified
// for URIs. This has the effect of partially decoding unreserved
// characters while performing case-normalization of percent-encoded
// characters.
func (u *URN) EncodingNormalize() {
	u.Scheme = strings.ToLower(u.Scheme)
	u.NID = strings.ToLower(u.NID)
	u.NSS = RecodeStringNSS(u.NSS)
	u.Resolve = RecodeStringComponent(u.Resolve)
	u.Query = RecodeStringComponent(u.Query)
	u.Fragment = RecodeStringComponent(u.Fragment)

	// A percent-encoding normalized URN is also regularly normalized,
	// so we can set it as normalized.
	u.normalized = true
}

func (u *URN) IsNormalized() bool {
	return u.normalized
}

// EncodingNormalized returns a percent-encoding normalized identifier
// without changing contents of the current identifier.
func (u *URN) EncodingNormalized() *URN {
	// Encoding normalized will be performed even if IsNormalized() is
	// true, as that method refers to URN normalization rather than
	// percent-encoding normalization.

	if u == nil {
		return u
	}

	n := u.Copy()
	n.EncodingNormalize()

	return n
}

// Normalized returns a normalized identifier without changing
// contents of the current identifier. If the current identifier is
// already normalized, returns itself instead of a copy.
func (u *URN) Normalized() *URN {
	if u == nil || u.IsNormalized() {
		return u
	}

	n := u.Copy()
	n.Normalize()

	return n
}

// String returns the complete identifier, including components.
func (u *URN) String() string {
	var b strings.Builder

	b.Grow(len(u.Scheme) + len(u.NID) + len(u.NSS) + len(u.Query) +
		len(u.Fragment) + 6)

	b.WriteString(u.Scheme)
	b.WriteByte(':')
	b.WriteString(u.NID)
	b.WriteByte(':')
	b.WriteString(u.NSS)

	if u.Resolve != "" {
		b.WriteString("?+")
		b.WriteString(u.Resolve)
	}

	if u.Query != "" {
		b.WriteString("?=")
		b.WriteString(u.Query)
	}

	if u.Fragment != "" || u.ForceFragment {
		b.WriteByte('#')
		b.WriteString(u.Fragment)
	}

	return b.String()
}

func normalizePercentEncoding(s string) string {
	// If input does not contains an escaped sequence, we can just
	// return itself.
	if !strings.Contains(s, "%") {
		return s
	}

	data := []byte(s)
	i := 0
	n := len(data)

	for i < n {
		if data[i] != '%' {
			i++

			continue
		}

		// If processing an invalid escape at the end of the string,
		// just stop now.
		if i+2 >= n {
			break
		}

		up := bytes.ToUpper(data[i+1 : i+3])

		data[i+1], data[i+2] = up[0], up[1]
		i += 2
	}

	return string(data)
}

// MarshalText marshals the URN as text.
func (u *URN) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

// MarshalText parses URN text input.
func (u *URN) UnmarshalText(b []byte) error {
	n, err := Parse(string(b))
	if err != nil {
		return fmt.Errorf("urn.UnmarshalText: %w", err)
	}

	*u = *n
	return nil
}
