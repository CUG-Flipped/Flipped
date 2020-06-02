// @Title  main.go
// @Description  The main function to run the application
// @Author  郑康
// @Update  郑康 2020.5.17
package main

import (
	"Flipped_Server/initialSetting"
	"Flipped_Server/network"
	"Flipped_Server/utils"
	"runtime"
)



// @title    main
// @description   main函数用于启动服务器
// @auth      郑康             2020.5.17
// @param     void
// @return    void
func main() {
	runtime.GOMAXPROCS(4)
	utils.ExitFlag = make(chan bool)

	initialSetting.InitSettings()
	socketServer := network.SocketServer{IPAddr: "", Port:8081}
	go socketServer.Run()

	httpServer := network.HttpServer{IPAddr: "", Port: 8080}
	httpServer.Run()
	<-utils.ExitFlag
}
