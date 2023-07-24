package urn

import "fmt"

const MaxLenNID = 32

// Parse parses a raw URN into a URN identifier structure.
func Parse(s string) (*URN, error) {
	p := &parser{
		source: s,
		n:      len(s),
		id:     &URN{},
		offset: 0,
	}

	return p.Parse()
}

type parser struct {
	// Input
	source string
	n      int

	// Working data
	id     *URN
	offset int
}

func (p *parser) Parse() (*URN, error) {
	p.id = &URN{}

	if err := p.consumeScheme(); err != nil {
		return nil, err
	}

	if err := p.consumeNID(); err != nil {
		return nil, err
	}

	if err := p.consumeNSS(); err != nil {
		return nil, err
	}

	if err := p.consumeResolve(); err != nil {
		return nil, err
	}

	if err := p.consumeQuery(); err != nil {
		return nil, err
	}
	if err := p.consumeFragment(); err != nil {
		return nil, err
	}

	// Trailing characters, which could be due to failing parsing either
	// NSS, q-component, or r-component.  Return the appropriate error.
	if p.offset < p.n {
		if p.source[p.offset] == '?' {
			return nil, p.newErr(ErrInvalidIdentifier, "invalid sequence after '?'")
		}

		return nil, p.newErr(ErrInvalidNSS, "trailing characters")
	}

	return p.id, nil
}

func (p *parser) newErr(err error, msg string) error {
	return &Error{"parse", p.source, err, msg}
}

func (p *parser) eol() bool {
	return p.offset >= p.n
}

func (p *parser) consumeScheme() error {
	if p.n < 4 {
		return p.newErr(ErrInvalidScheme, "too short")
	}

	c1, c2, c3, c4 := p.source[0], p.source[1], p.source[2], p.source[3]
	if !(c1 == 'u' || c1 == 'U') || !(c2 == 'r' || c2 == 'R') || !(c3 == 'n' || c3 == 'N') || !(c4 == ':') {
		return p.newErr(ErrInvalidScheme, fmt.Sprintf("unknown scheme %q", p.source[:3]))
	}

	p.id.Scheme = p.source[:3]
	p.offset = 4

	return nil
}

func (p *parser) consumeNID() error {
	i := p.offset

	if p.eol() || !isAlphaNum(p.source[i]) {
		return p.newErr(ErrInvalidNID, "unexpected eol")
	}

	max := i + MaxLenNID
	if p.n < max {
		max = p.n
	}

	i++
	for i < max {
		c := p.source[i]
		if c == ':' {
			break
		}

		if !(isAlphaNum(c) || c == '-') {
			return p.newErr(ErrInvalidNID, fmt.Sprintf("invalid byte at pos %d", i))
		}

		i++
	}

	// NID must end with a pchar, not a hyphen.  Because we have
	// already validated that the last one is a pchar or a hyphen, we
	// just need to check that the last one is not a hyphen.
	if i >= p.n || !(p.source[i-1] != '-' && p.source[i] == ':') {
		return p.newErr(ErrInvalidNID, "invalid final byte")
	}

	// NID cannot contain fewer than two pchar.
	if i-p.offset < 2 {
		return p.newErr(ErrInvalidNID, "too short")
	}

	p.id.NID = p.source[p.offset:i]
	p.offset = i + 1

	return nil
}

func (p *parser) consumeNSS() error {
	if p.eol() {
		return p.newErr(ErrInvalidNSS, "unexpected eol")
	}

	i := p.offset

	if isPCharSingle(p.source[i]) {
		i++
	} else if p.maybePercentEncoded(i) {
		i += 3
	} else {
		return p.newErr(ErrInvalidNSS, "invalid initial byte")
	}

	for i < p.n {
		c := p.source[i]
		if isPCharSingle(c) || c == '/' {
			i++

			continue
		}

		if p.maybePercentEncoded(i) {
			i += 3

			continue
		}

		break
	}

	if i-p.offset < 1 {
		return p.newErr(ErrInvalidNSS, "too short")
	}

	p.id.NSS = p.source[p.offset:i]
	p.offset = i

	return nil
}

func (p *parser) consumeResolve() error {
	if p.offset+1 >= p.n {
		return nil
	}

	i := p.offset
	c1, c2 := p.source[i], p.source[i+1]

	if c1 != '?' || c2 != '+' {
		// Not r-component.
		return nil
	}

	p.offset += 2
	i += 2

	if p.eol() {
		return p.newErr(ErrInvalidResolve, "unexpected eol")
	}

	// First item must be PChar only.
	if isPCharSingle(p.source[i]) {
		i++
	} else if p.maybePercentEncoded(i) {
		i += 3
	} else {
		return p.newErr(ErrInvalidResolve, "invalid first byte")
	}

	for i < p.n {
		c := p.source[i]
		if isPCharSingle(c) || c == '/' {
			i++

			continue
		}

		if c == '?' {
			if p.maybeQuerySequence(i) {
				break
			}

			i++

			continue
		}

		if p.maybePercentEncoded(i) {
			i += 3

			continue
		}

		break
	}

	p.id.Resolve = p.source[p.offset:i]
	p.offset = i

	return nil
}

func (p *parser) consumeQuery() error {
	if p.offset+1 >= p.n {
		return nil
	}

	i := p.offset
	c1, c2 := p.source[i], p.source[i+1]

	if c1 != '?' || c2 != '=' {
		return nil
	}

	// Move offset and start validating the first byte of q-component.
	p.offset += 2
	i += 2

	if p.eol() {
		return p.newErr(ErrInvalidQuery, "unexpected eol")
	}

	// First item must be PChar only.
	if isPCharSingle(p.source[i]) {
		i++
	} else if p.maybePercentEncoded(i) {
		i += 3
	} else {
		return p.newErr(ErrInvalidQuery, "invalid first byte")
	}

	for i < p.n {
		c := p.source[i]
		if isPCharSingle(c) || c == '/' || c == '?' {
			i++

			continue
		}

		if p.maybePercentEncoded(i) {
			i += 3

			continue
		}

		break
	}

	p.id.Query = p.source[p.offset:i]
	p.offset = i

	return nil
}

func (p *parser) consumeFragment() error {
	if p.eol() {
		return nil
	}

	if p.source[p.offset] != '#' {
		return nil
	}

	// Move offset and start validating the first byte of fragment.
	p.offset += 1
	i := p.offset

	for i < p.n {
		c := p.source[i]
		if isPCharSingle(c) || c == '/' || c == '?' {
			i++

			continue
		}

		if p.maybePercentEncoded(i) {
			i += 3

			continue
		}

		break
	}

	if i-p.offset < 1 {
		p.id.ForceFragment = true
	}

	p.id.Fragment = p.source[p.offset:i]
	p.offset = i

	return nil
}

func (p *parser) maybePercentEncoded(i int) bool {
	if p.n-i < 3 {
		return false
	}

	c1, c2, c3 := p.source[i], p.source[i+1], p.source[i+2]

	return c1 == '%' && isHex(c2) && isHex(c3)
}

func (p *parser) maybeQuerySequence(i int) bool {
	if p.n-i < 2 {
		return false
	}

	c1, c2 := p.source[i], p.source[i+1]
	return c1 == '?' && c2 == '='
}

func isPCharSingle(c byte) bool {
	return isUnreserved(c) || isSubDelim(c) || c == ':' || c == '@'
}

func isUnreserved(c byte) bool {
	if isAlpha(c) || isDigit(c) {
		return true
	}

	switch c {
	case '-', '.', '_', '~':
		return true
	}

	return false
}

func isSubDelim(c byte) bool {
	switch c {
	case '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=':
		return true
	}

	return false
}

func isAlphaNum(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

func isAlpha(c byte) bool {
	switch {
	case 'A' <= c && c <= 'Z':
		return true
	case 'a' <= c && c <= 'z':
		return true
	}

	return false
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isHex(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'A' <= c && c <= 'F':
		return true
	case 'a' <= c && c <= 'f':
		return true
	}

	return false
}
