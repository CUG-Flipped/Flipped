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
	server.HttpServerApp.Handle("POST", "/register", server.registerHandler)
	server.HttpServerApp.Handle("POST", "/login", server.loginHandler)
}

func (server *HttpServer)registerHandler(context iris.Context) {

	userType := context.URLParam("user_type")
	name := context.URLParam("name")
	email := context.URLParam("email")
	photo := context.URLParam("photo")
	password := context.URLParam("password")

	fmt.Printf("type: %s, name: %s, email: %s, photo: %s, password: %s", userType, name,email, photo, password)


	_, _ = context.GzipResponseWriter().WriteString("You're going to Register")
}

func (server *HttpServer)loginHandler(context iris.Context) {

	_, _ = context.GzipResponseWriter().WriteString("You're going to login")
}




