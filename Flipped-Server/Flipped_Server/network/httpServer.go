package network

import (
	"Flipped_Server/dataBase"
	"Flipped_Server/logger"
	"Flipped_Server/sqlmapper"
	"Flipped_Server/utils"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type IFunction interface {
	registerHandler(context *gin.Context)
	loginHandler(context *gin.Context)
}

type HttpServer struct {
	IPAddr string
	Port int
}

var(
	Router = gin.Default()
)

func (server *HttpServer) Run() {
	server.bindRouteAndHandler()
	dataBase.Init()
	logger.InitLog()
	_ = Router.Run(server.IPAddr + ":" + strconv.Itoa(server.Port))
}


func (server *HttpServer) bindRouteAndHandler() {
	Router.POST("/login", server.loginHandler)
	Router.POST("/register", server.registerHandler)
}

func (server *HttpServer)registerHandler(context *gin.Context) {
	var res bytes.Buffer
	status := http.StatusOK
	responseStr := ""

	userType, err1 := strconv.Atoi(context.DefaultQuery("user_type", "-1"))
	name := context.DefaultQuery("name", "")
	email := context.DefaultQuery("email", "")
	password := context.DefaultQuery("password", "")

	form, _ := context.MultipartForm()
	photoInfo := form.File["photo"][0]

	photoName := photoInfo.Filename
	photo := make([]byte, photoInfo.Size)
	file, err := photoInfo.Open()

	if err1 != nil {

		logger.Logger.WithFields(logrus.Fields {
			"function": "registerHandler",
			"cause": "convert user_type",
		}).Error(err1.Error())

		status = http.StatusBadRequest
		responseStr = "wrong data type of 'userType'"
	}

	if err != nil {
		fmt.Println(err.Error())
		logger.Logger.WithFields(logrus.Fields {
			"function": "registerHandler",
			"cause": "open photoInfo",
		}).Error(err1.Error())

		status = http.StatusBadRequest
		responseStr = "upload file is unacceptable"
		
	} else {
		count, _ := file.Read(photo)
		res.WriteString(fmt.Sprintf("Photo total %d bytes\n", count))
		defer file.Close()
	}

	res.WriteString(fmt.Sprintf("type: %d, name: %s, email: %s, password: %s\n", userType, name,email, password))

	registerTable := dataBase.RegisterTable{
			Pid:        utils.GeneratorUUID(),
			Username:   name,
			Password:   password,
			UserType:   byte(userType),
			Email:      email,
			Photo:      utils.GetImageURL(photoName, photo),
			RealName:   "",
			Profession: "",
			Age:        0,
			Region:     ""  ,
			Hobby:      "",
	}

	err2 := sqlmapper.Insert(registerTable, "userinfo")

	if err2 != nil {
		status = http.StatusInternalServerError
		responseStr = "Get an error when insert into DataBase"

		logger.Logger.WithFields(logrus.Fields {
			"function": "registerHandler",
			"cause": "Insert into database using 'sqlmapper.Insert'",
		}).Error(err2.Error())
	}

	if status == http.StatusOK {
		context.String(status, res.String())
	} else {
		context.String(status, responseStr)
	}

	logger.Logger.WithFields(logrus.Fields{
		"function": "registerHandler",
		"cause": "receive Request from client",
	}).Info("response: " + responseStr + ", Status: "+ strconv.Itoa(status))
}

func (server *HttpServer)loginHandler(context *gin.Context) {
	context.String(http.StatusOK, "You're going to login")

}

