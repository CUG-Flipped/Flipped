package main

import "Flipped_Server/network"

func main() {
	httpServer := network.HttpServer{"", 8080}
	httpServer.Run()
}


