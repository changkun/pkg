package errors_test

import (
	"testing"

	"github.com/changkun/gobase/errors"
)

func TestTryCatch(t *testing.T) {
	errors.Try(func() error {
		return nil
	}).Catch(func(err error) {
		if err != nil {
			println("on error")
		}
	})

	errors.Try(func() error {
		return errors.New("e")
	}).Catch(func(err error) {
		if err == nil {
			println("no errors")
		}
	})
}
