// Copyright 2020 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// maxMultipartMemory caps the request body gin buffers in memory.
const maxMultipartMemory = 500 << 20 // 500 MB

// addRequest is the body of POST /api/v1/add. Binding it lets gin reject a
// malformed body before the handler runs.
type addRequest struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

// Register builds the REST half of the REST versus gRPC comparison. Both
// halves answer the same addition, so the difference the benchmark reports is
// transport and encoding, not work.
func Register() *gin.Engine {
	r := gin.Default()
	r.MaxMultipartMemory = maxMultipartMemory

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"msg": "pong"})
		})
		v1.POST("/add", func(c *gin.Context) {
			var body addRequest
			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"sum": 0})
				return
			}
			c.JSON(http.StatusOK, gin.H{"sum": body.A + body.B})
		})
	}
	return r
}
