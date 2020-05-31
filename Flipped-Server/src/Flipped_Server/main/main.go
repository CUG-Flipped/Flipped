// @Title  main.go
// @Description  The main function to run the application
// @Author  郑康
// @Update  郑康 2020.5.17
package main

import (
	"Flipped_Server/network"
	"Flipped_Server/utils"
)



// @title    main
// @description   main函数用于启动服务器
// @auth      郑康             2020.5.17
// @param     void
// @return    void
func main() {
	//runtime.GOMAXPROCS(4)
	utils.ExitFlag = make(chan bool)
	httpServer := network.HttpServer{IPAddr: "", Port: 8080}
	go httpServer.Run()
	socketServer := network.SocketServer{IPAddr: "", Port:8081}
	socketServer.Run()
	<-utils.ExitFlag
}
