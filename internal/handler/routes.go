package handler

import "github.com/gin-gonic/gin"

func Routes(h Handler) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", h.Home)
	router.GET("/ws", h.WS)

	return router
}
