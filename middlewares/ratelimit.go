package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
)

var limiter = time.NewTicker(5 * time.Second)

func RateLimit(c *gin.Context) {
	<-limiter.C
	c.Next()
}
