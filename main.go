package main

import (
	"os"
	"strconv"

	// "fmt"
	"g2ww/config"
	"g2ww/router"

	"github.com/gofiber/compression"
	"github.com/gofiber/fiber"
)

// 读取配置文件中的config
var (
	Port = strconv.Itoa(config.Config.Port)
)

func main() {

	app := fiber.New()
	var ListenAddress string

	// 判断：默认端口号为2408
	if Port == "" {
		Port = "2408"
	}

	if os.Getenv("DOCKER") != "" {
		ListenAddress = "0.0.0.0" + ":" + Port
	} else {
		ListenAddress = "127.0.0.1" + ":" + Port
	}

	// Server Info
	app.Use(compression.New())
	app.Get("/", router.GwStat())
	// app.All("/:key", GwWorker())
	app.Post("/:key", router.GwWorker())

	v1 := app.Group("/send")
	v1.Post("/:key", router.GwWorker())

	app.Listen(ListenAddress)

}
