// Package errors implements transactional error handling
package errors

import "errors"

// New creates a new erorr with given string.
func New(s string) error {
	return errors.New(s)
}

// Try tries a given function and see if it throw any error
func Try(e func() error) catch {
	return catch{err: e(), result: nil}
}

type catch struct {
	err    error
	result interface{}
}

// Catch catches the error that throwed in Try call.
func (c catch) Catch(handler func(error) interface{}) catch {
	if c.err != nil {
		c.result = handler(c.err)
	}
	return c
}

// Final runs in any case if it is called.
func (c catch) Final(handler func(interface{})) {
	handler(c.result)
}
