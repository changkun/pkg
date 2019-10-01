// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package rt_test

import (
	"reflect"
	"testing"

	"github.com/changkun/gobase/rt"
)

func TestGetFuncName(t *testing.T) {
	got := rt.GetFuncName(rt.GetFuncName)
	want := "github.com/changkun/gobase/rt.GetFuncName"
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v, got %v", want, got)
	}
}
