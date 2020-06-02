// @Title  httpServer.go
// @Description  To provide a network interface of Http、TCP and UDP to the Server
// @Author  郑康
// @Update  郑康 2020.5.17
package network

import (
	"Flipped_Server/dataBase"
	"Flipped_Server/logger"
	"Flipped_Server/sqlmapper"
	"Flipped_Server/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"sync"
)

// IFunction接口包含了http路由处理函数
type IFunction interface {
	registerHandler(context *gin.Context)
	loginHandler(context *gin.Context)
	friendsListHandler(context *gin.Context)
	recommendedFriendsListHandler(context *gin.Context)
	heartBeatHandler(context *gin.Context)
	countOnlineUserNumber(context *gin.Context)
	judgeUserAlive(context *gin.Context)
	closeServer(context *gin.Context)
}

// HttpServer结构体包含了Http服务器绑定的IP地址和端口号
type HttpServer struct {
	IPAddr string
	Port   int
}

// 全局变量，gin实例
var (
	Router = gin.Default()
	lock = &sync.Mutex{}
)

// @title    Run
// @description   绑定路由处理函数、初始化mysql数据库、日志模块、初始化Redis数据库、初始化MongoDB数据库、启动Http服务器
// @auth      郑康             2020.5.17
// @param     void
// @return    void
func (server *HttpServer) Run() {
	server.bindRouteAndHandler()
	dataBase.Init()
	logger.InitLog()
	dataBase.RedisClientInit()
	dataBase.InitializeMongoDB()
	_ = Router.Run(server.IPAddr + ":" + strconv.Itoa(server.Port))
}

// @title    bindRouteAndHandler
// @description   将路由与处理函数绑定
// @auth      郑康             2020.5.17
// @param     void
// @return    void
func (server *HttpServer) bindRouteAndHandler() {
	Router.POST("/login", server.loginHandler)
	Router.POST("/register", server.registerHandler)
	Router.GET("/friendList", server.friendsListHandler)
	Router.GET("/recommendUser", server.recommendedFriendsListHandler)
	Router.GET("heartBeat", server.heartBeatHandler)
	Router.GET("/onlineUserNumber", server.countOnlineUserNumber)
	Router.POST("/isAlive", server.judgeUserAlive)
	Router.GET("/closeServer", server.closeServer)
}

// @title    registerHandler
// @description   注册路由的处理函数
// @auth      郑康             2020.5.17
// @param     *gin.Context	  gin的上下文指针
// @return    void
func (server *HttpServer) registerHandler(context *gin.Context) {
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

		logger.Logger.WithFields(logrus.Fields{
			"function": "registerHandler",
			"cause":    "convert user_type",
		}).Error(err1.Error())

		status = http.StatusBadRequest
		responseStr = "wrong data type of 'userType'"
	}

	if err != nil {
		fmt.Println(err.Error())
		logger.Logger.WithFields(logrus.Fields{
			"function": "registerHandler",
			"cause":    "open photoInfo",
		}).Error(err.Error())

		status = http.StatusBadRequest
		responseStr = "upload file is unacceptable"

	} else {
		count, _ := file.Read(photo)
		res.WriteString(fmt.Sprintf("Photo total %d bytes\n", count))
		defer file.Close()
	}

	res.WriteString(fmt.Sprintf("type: %d, name: %s, email: %s, password: %s\n", userType, name, email, password))

	registerTable := dataBase.UserInfoTable{
		Pid:        utils.GeneratorUUID(),
		Username:   name,
		Password:   password,
		UserType:   userType,
		Email:      email,
		Photo:      utils.GetImageURL(photoName, photo),
		RealName:   "",
		Profession: "",
		Age:        0,
		Region:     "",
		Hobby:      "",
	}

	err2 := sqlmapper.Insert(registerTable, "userinfo")

	if err2 != nil {
		status = http.StatusInternalServerError
		responseStr = "Get an error when insert into DataBase"

		logger.Logger.WithFields(logrus.Fields{
			"function": "registerHandler",
			"cause":    "Insert into database using 'sqlmapper.Insert'",
		}).Error(err2.Error())
	}

	if status == http.StatusOK {
		context.String(status, res.String())
	} else {
		context.String(status, responseStr)
	}

	_ = dataBase.InitUserFriendList(name)

	logger.Logger.WithFields(logrus.Fields{
		"function": "registerHandler",
		"cause":    "receive Request from client",
	}).Info("response: " + responseStr + ", Status: " + strconv.Itoa(status))
}

// @title    loginHandler
// @description   登录路由的处理函数
// @auth      郑康             2020.5.17
// @param     *gin.Context	  gin的上下文指针
// @return    void
func (server *HttpServer) loginHandler(context *gin.Context) {
	username := context.DefaultQuery("username", "")
	pwd := context.DefaultQuery("password", "")

	var status int
	var msg string
	var data interface{}

	userInfo, err := dataBase.FindUserInfo(username, pwd)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"Function": "loginHandler",
			"cause":    "execute function of FindUserInfo",
		}).Error(err.Error())
		status = http.StatusInternalServerError
		msg = "there is something going wrong with the server, please try it again"
		data = ""
	}
	if userInfo == nil {
		logger.Logger.WithFields(logrus.Fields{
			"Function": "loginHandler",
			"cause":    "the username or password of the request is incorrect",
		}).Info(username, pwd)
		//context.String(http.StatusUnauthorized, "user does't exist or wrong username or wrong password")
		status = http.StatusUnauthorized
		msg = "account does't exist or wrong username or wrong password"
		data = ""
	}
	var tokenStr string
	if dataBase.KeyExists(username, 0) {
		tokenStr, err = dataBase.ReadFromRedis(username)
		if err != nil {
			status = http.StatusInternalServerError
			msg = "there is something going wrong with the server, please try it again"
			data = ""
		} else {
			status = http.StatusOK
			msg = "succeed to login"
			data = gin.H{
				"token": tokenStr,
			}
		}
	} else {
		tokenStr, err = GenerateToken(username)
		if err != nil {
			status = http.StatusInternalServerError
			msg = "there is something going wrong with the server, please try it again"
			data = ""
		} else {
			status = http.StatusOK
			msg = "succeed to login"
			data = gin.H{
				"token": tokenStr,
			}
		}
	}
	context.JSON(status, gin.H{
		"message": msg,
		"data":    data,
	})
}

// @title    loginHandler
// @description   获取好友列表路由的处理函数
// @auth      郑康             2020.5.28
// @param     *gin.Context	  gin的上下文指针
// @return    void
func (server *HttpServer) friendsListHandler(context *gin.Context)  {
	tokenStr := context.Request.Header.Get("token")
	username, err := ParseToken(tokenStr)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "friendsListHandler", "Parse Token which sent by client", "")
		context.JSON(http.StatusUnauthorized, gin.H {
			"message": "some error occur when parsing tokenStr, Please login again",
			"data": err.Error(),
		})
	} else {
		friendList, err := dataBase.GetFriendListByUserName(username)
		if err != nil {
			logger.SetToLogger(logrus.ErrorLevel, "friendsListHandler", "execute GetFriendListByName function", err.Error())
			context.JSON(http.StatusInternalServerError, gin.H {
				"message": "some error occur when parsing tokenStr",
				"data": err.Error(),
			})
		} else {
			context.JSON(http.StatusOK, gin.H{
				"message": "succeed to find friend list",
				"data": friendList,
			})
		}
	}
}

// @title    loginHandler
// @description   获取好友列表路由的处理函数
// @auth      郑康             2020.5.28
// @param     *gin.Context	  gin的上下文指针
// @return    void
func (server *HttpServer) recommendedFriendsListHandler(context *gin.Context) {
	tokenStr := context.Request.Header.Get("token")
	username, err := ParseToken(tokenStr)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "recommendedFriendsListHandler", "Parse Token which sent by client", "")
		context.JSON(http.StatusUnauthorized, gin.H {
			"message": "some error occur when parsing tokenStr, Please login again",
			"data": err.Error(),
		})
		return
	}
	selectedUser, err := dataBase.SelectSimilarUser(username)
	if err != nil{
		logger.SetToLogger(logrus.ErrorLevel, "recommendedFriendsListHandler", "select similar user", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H {
			"message": "some error occur in the server, Please try again",
			"data": err.Error(),
		})
		return
	}
	imageStr, err := utils.Image2Base64(selectedUser.Photo)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "recommendedFriendsListHandler", "convert image to Base64 string", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H {
			"message": "some error occur in the server, Please try again",
			"data": err.Error(),
		})
		return
	}
	selectedUser.Photo = imageStr
	dataJson, err := json.Marshal(selectedUser)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "recommendedFriendsListHandler", "convert Struct to Json", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H {
			"message": "some error occur in the server, Please try again",
			"data": err.Error(),
		})
		return
	}

	var dataMap map[string] interface{}
	err = json.Unmarshal(dataJson, &dataMap)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "recommendedFriendsListHandler", "convert Json to Map", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H {
			"message": "some error occur in the server, Please try again",
			"data": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H {
		"message": "succeed to handle the request",
		"data": dataMap,
	})
}

func (server *HttpServer) heartBeatHandler(context *gin.Context) {
	status := http.StatusOK
	msg := "succeed to handle the request"
	data := ""
	tokenStr := context.Request.Header.Get("token")
	username, err := ParseToken(tokenStr)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "heartBeatHandler", "parse token error", err.Error())
		status = http.StatusUnauthorized
		msg = "token is invalid, Please login again"
		data = err.Error()
		return
	}
	go dataBase.UpdateUserStatus(username, lock)
	context.JSON(status, gin.H{
		"message": msg,
		"data":data,
	})
}

func (server *HttpServer) countOnlineUserNumber(context *gin.Context){
	num, err := dataBase.CountOnlineUsers()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H {
			"message": "some error occur in the server, Please try again",
			"data": err.Error(),
		})
		return
	} else {
		context.JSON(http.StatusOK, gin.H {
			"message": "succeed to handle the request",
			"data": num,
		})
	}
}

func (server *HttpServer) judgeUserAlive(context *gin.Context) {
	status := http.StatusOK
	msg := "succeed to handle the request"
	data := ""
	tokenStr := context.Request.Header.Get("token")
	_, err := ParseToken(tokenStr)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "heartBeatHandler", "parse token error", err.Error())
		status = http.StatusUnauthorized
		msg = "token is invalid, Please login again"
		data = err.Error()
	} else {
		targetUser := context.DefaultQuery("username", "")
		if targetUser == ""{
			msg = "the user you want to query can't be empty"
		} else {
			msg = "succeed top handle the request"
			if dataBase.KeyExists(targetUser, 2) {
				data = "user: " + targetUser + " is alive"
			} else {
				data = "user: " + targetUser + " is not alive"
			}
		}
	}
	context.JSONP(status, gin.H {
		"message": msg,
		"data": data,
	})
}

func (server *HttpServer) closeServer(context *gin.Context){
	utils.ExitFlag <- true
}
