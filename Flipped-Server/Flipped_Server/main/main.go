// @Title  main.go
// @Description  The main function to run the application
// @Author  郑康
// @Update  郑康 2020.5.17
package main

import "Flipped_Server/network"

// @title    main
// @description   main函数用于启动服务器
// @auth      郑康             2020.5.17
// @param     void
// @return    void
func main() {
	httpServer := network.HttpServer{IPAddr: "", Port: 8080}
	httpServer.Run()
}
