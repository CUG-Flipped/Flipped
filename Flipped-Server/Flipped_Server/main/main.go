package main

import (
	"Flipped_Server/network"
	"github.com/kataras/iris/v12"
)

func main() {
	httpServer := network.HttpServer{"", 8080, iris.New()}
	httpServer.Run()
}

