package errors_test

import (
	"testing"

	"github.com/changkun/gobase/errors"
)

func TestTryCatch(t *testing.T) {
	errors.Try(func() error {
		return nil
	}).Catch(func(err error) interface{} {
		t.Fatalf("catch error: %v", err)
		return nil
	}).Final(func(result interface{}) {
		t.Log("errors: everything is good")
	})

	errors.Try(func() error {
		return errors.New("e")
	}).Catch(func(err error) interface{} {
		t.Logf("capctured error: %v", err)
		return err
	}).Final(func(result interface{}) {
		if result == nil {
			t.Fatalf("cannot capture error")
		}
	})

	errors.Try(func() error {
		return nil
	}).Final(func(r interface{}) {
		if r != nil {
			t.Fatalf("result contains error")
		}
	})
}
