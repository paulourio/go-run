package urn

import (
	"errors"
	"fmt"
)

// Error represents an error that occurred during an operation, such
// as during parse or some specific sub-operation such as unescape.
type Error struct {
	Op   string
	Data string // Data associated with the error, eg URN or component
	Err  error
	Msg  string // Optional explanation
}

func (e *Error) Error() string {
	if e.Msg != "" {
		return fmt.Sprintf("%s %q: %s: %s", e.Op, e.Data, e.Err, e.Msg)
	}

	return fmt.Sprintf("%s %q: %s", e.Op, e.Data, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

var (
	ErrInvalidIdentifier = errors.New("invalid identifier")
	ErrInvalidScheme     = errors.New("invalid scheme")
	ErrInvalidNID        = errors.New("invalid NID")
	ErrInvalidNSS        = errors.New("invalid NSS")
	ErrInvalidResolve    = errors.New("invalid resolve component")
	ErrInvalidQuery      = errors.New("invalid query component")
)
