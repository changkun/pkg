package errors_test

import (
	"testing"

	"changkun.de/x/pkg/errors"
)

func TestTryCatch(t *testing.T) {
	errors.Try(func() (interface{}, error) {
		return 0, nil
	}).Catch(func(_ interface{}, err error) interface{} {
		t.Fatalf("catch error: %v", err)
		return nil
	}).Final(func(result interface{}) {
		t.Log("errors: everything is good")
	})

	errors.Try(func() (interface{}, error) {
		return 1, errors.New("e")
	}).Catch(func(result interface{}, err error) interface{} {
		t.Logf("captured result: %v", result.(int))
		t.Logf("captured error: %v", err)
		return err
	}).Final(func(result interface{}) {
		if result == nil {
			t.Fatalf("cannot capture error")
		}
	})

	errors.Try(func() (interface{}, error) {
		return 1, nil
	}).Final(func(r interface{}) {
		if r.(int) != 1 {
			t.Fatalf("result from try block is not as expected")
		}
	})
}
