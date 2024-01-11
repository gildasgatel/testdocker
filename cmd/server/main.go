package main

import (
	"testdocker/internal/handler"
)

func main() {
	// init handler
	var handle handler.Handler = handler.New()

	// setup routes
	router := handler.Routes(handle)

	//start server port 8080
	router.Run(":8080")
}
