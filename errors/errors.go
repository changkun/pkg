// Package errors implements transactional error handling
package errors

import "errors"

// New creates a new erorr with given string.
func New(s string) error {
	return errors.New(s)
}

// Try catch
func Try(e func() error) catch {
	return catch{e()}
}

type catch struct {
	err error
}

func (c catch) Catch(handler func(error)) {
	handler(c.err)
}
