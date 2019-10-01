// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package gen

import (
	"math/rand"
	"time"
)

// RandomString generates a random string
func RandomString(n int) string {
	const alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var str string
	length := len(alphanum)
	for i := 0; i < n; i++ {
		a := alphanum[r.Intn(len(alphanum))%length]
		str += string(a)
	}
	return str
}
