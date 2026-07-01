package middlewares

import (
	"math/rand/v2"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var mu sync.Mutex

func RateLimit(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	d := 3 + rand.IntN(8)
	time.Sleep(time.Duration(d) * time.Second)
	c.Next()
}
