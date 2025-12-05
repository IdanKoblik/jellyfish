package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := ReadConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.POST("/webhook", func(context *gin.Context) {
		HandleWebhook(cfg, context)
	})

	err = router.Run(fmt.Sprintf(":%s", os.Getenv("API_PORT")))
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
