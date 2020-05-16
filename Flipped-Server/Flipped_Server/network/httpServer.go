package network

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"strconv"
)

type HttpServer struct {
	IPAddr string
	Port int
	HttpServerApp *iris.Application
}


func (server *HttpServer) Run() {
	host := iris.Addr(server.IPAddr + ":" + strconv.Itoa(server.Port))
	fmt.Println(host)
	server.bindRouteAndHandler()
	server.HttpServerApp.Run(host)
}

func (server *HttpServer) bindRouteAndHandler() {
	server.HttpServerApp.Handle("GET", "/register", server.registerHandler)
	server.HttpServerApp.Handle("GET", "/login", server.loginHandler)
}

func (server *HttpServer)registerHandler(context iris.Context) {
	fmt.Println(context.Path())
	context.GzipResponseWriter().WriteString("You're going to Register")
}

func (server *HttpServer)loginHandler(context iris.Context) {
	fmt.Println(context.Path())
	context.GzipResponseWriter().WriteString("You're going to login")
}




