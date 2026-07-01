package initialize

import (
	"aurora/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRouter() *gin.Engine {
	handler := NewHandle(
		checkProxy(),
		readAccessToken(),
	)

	// 初始化基础前置参数
	handler.InitBasicConfigForChatGPT()

	router := gin.Default()
	router.Use(middlewares.Cors)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/auth/session", handler.session)
	router.POST("/auth/refresh", handler.refresh)
	router.OPTIONS("/v1/chat/completions", optionsHandler)
	router.OPTIONS("/v1/models", optionsHandler)
	router.OPTIONS("/v1/responses", optionsHandler)
	router.OPTIONS("/v1/images/generations", optionsHandler)
	router.OPTIONS("/v1/images/edits", optionsHandler)
	router.OPTIONS("/v1/images/variations", optionsHandler)
	router.OPTIONS("/v1/files", optionsHandler)

	authGroup := router.Group("").Use(middlewares.Authorization)
	authGroup.POST("/v1/files", handler.files)
	authGroup.GET("/v1/models", handler.engines)
	authGroup.POST("/backend-api/conversation", handler.chatgptConversation)

	// rate-limited endpoints to avoid triggering ChatGPT web rate limits (~5s between requests)
	limitedGroup := router.Group("").Use(middlewares.Authorization, middlewares.RateLimit)
	limitedGroup.POST("/v1/chat/completions", handler.nightmare)
	limitedGroup.POST("/v1/responses", handler.responses)
	limitedGroup.POST("/v1/images/generations", handler.imageGenerations)
	// 改图 + 图生图(变体)统一入口:
	//   - 传 prompt     → 按 prompt 改图
	//   - 不传 prompt   → 自动注入默认指令,生成图像变体(图生图)
	limitedGroup.POST("/v1/images/edits", handler.imageEdits)
	limitedGroup.POST("/v1/images/variations", handler.imageVariations)
	limitedGroup.OPTIONS("/v1/audio/speech", optionsHandler)
	limitedGroup.POST("/v1/audio/speech", handler.tts)
	limitedGroup.OPTIONS("/v1/audio/transcriptions", optionsHandler)
	limitedGroup.POST("/v1/audio/transcriptions", handler.transcriptions)
	limitedGroup.OPTIONS("/v1/audio/translations", optionsHandler)
	limitedGroup.POST("/v1/audio/translations", handler.translations)

	return router
}
