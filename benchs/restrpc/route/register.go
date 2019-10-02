package route

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register ...
func Register() *gin.Engine {
	r := gin.Default()
	r.MaxMultipartMemory = 500 << 20 // 500 MB

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"msg": "pong"})
		})
		v1.POST("/add", func(c *gin.Context) {
			raw, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"sum": 0})
				return
			}
			var body map[string]float64
			if err := json.Unmarshal(raw, &body); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"sum": 0})
				return
			}
			a := body["a"]
			b := body["b"]
			c.JSON(http.StatusOK, gin.H{
				"sum": a + b,
			})
		})
	}
	return r
}
