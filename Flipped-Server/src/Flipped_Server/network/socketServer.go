package network

import (
	"Flipped_Server/dataBase"
	"Flipped_Server/logger"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
	"sync"
	"time"
)

var(
	UserConnectionMap map[string] net.Conn
	rwLock = &sync.RWMutex{}
)

type SocketServer struct {
	IPAddr string
	Port int
}

type fromClientMsg struct {
	MsgType int `json:"MsgType"`
	MsgFrom string `json:"MsgFrom"`
	MsgTo string `json:"MsgTo"`
	MsgContent string `json:"MsgContent"`
}

type toClientMsg struct {
	MsgFrom string `json:"MsgFrom"`
	MsgContent string `json:"MsgContent"`
}

type replyMsg struct {
	ResultCode int `json:"ResultCode"`
	MsgContent string `json:"MsgContent"`
}

func (ss *SocketServer) Run() {
	server, err := net.Listen("tcp", ss.IPAddr + ":" + strconv.Itoa(ss.Port))
	UserConnectionMap = make(map[string] net.Conn)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "Run", "error to start socket server", err.Error())
		return
	}
	for  {
		conn, err := server.Accept()
		go loopFindOnlineUserAndSendMessage()
		if err != nil {
			logger.SetToLogger(logrus.ErrorLevel, "Run", "error to Accept socket", err.Error())
		} else {
			logger.SetToLogger(logrus.ErrorLevel, "Run", "succeed to Accept socket", "")
			go connectionHandler(conn)
		}
	}
}

//ToDo: 持续监听该客户端的消息
func connectionHandler(conn net.Conn){
	data := json.NewDecoder(conn)
	var msg fromClientMsg
	err := data.Decode(&msg)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "connectionHandler", "error to decode message sent from client", err.Error())
		return
	}
	switch msg.MsgType {
		case 1:
			err = communicationRequestHandler(&msg, conn)
			if err != nil {
				reply := replyMsg{
					ResultCode: 500,
					MsgContent: "Some error occur in the server, please try it latter",
				}
				buf, _ := json.Marshal(reply)
				_, _ = conn.Write(buf)
			}
		case 2:
			err = connectionRequestHandler(&msg, conn)
			if err != nil {
				reply := replyMsg{
					ResultCode: 500,
					MsgContent: "Some error occur in the server, please try it latter",
				}
				buf, _ := json.Marshal(reply)
				_, _ = conn.Write(buf)
			}
		default:
			logger.SetToLogger(logrus.ErrorLevel, "connectionHandler", "error type in communicationMsg struct", "")
	}
}
//交流请求
func communicationRequestHandler(msg *fromClientMsg, conn net.Conn) error{
	realMsg := *msg
	sourceUserToken := realMsg.MsgFrom
	sourceUser, err := ParseToken(sourceUserToken)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "communicationRequestHandler", "error to parse token from client", err.Error())
		return err
	}

	rwLock.Lock()
	UserConnectionMap[sourceUser] = conn
	rwLock.Unlock()

	targetUser := realMsg.MsgTo
	msgContent := realMsg.MsgContent
	//如果目标用户处于与服务器TCP连接状态
	if isUserConnected(targetUser) {
		//获取目标用户连接
		if targetConn := getUserConnection(targetUser); targetConn != nil {
			communicationMsg := toClientMsg{
				MsgFrom: targetUser,
				MsgContent: msgContent,
			}
			buf, err := json.Marshal(communicationMsg)
			if err != nil {
				logger.SetToLogger(logrus.ErrorLevel, "communicationRequestHandler", "error to parse struct to json", err.Error())
				return err
			}
			count, err := targetConn.Write(buf)
			if err != nil {
				logger.SetToLogger(logrus.ErrorLevel, "communicationRequestHandler", "error to send to client with bytes", err.Error())
				return err
			} else {
				logger.SetToLogger(logrus.InfoLevel,"communicationRequestHandler", "Succeed to send message from " + sourceUser + " to " + targetUser, "bytes num: " + strconv.Itoa(count))
			}
		}
	} else {
		err := dataBase.WriteMessage(sourceUser, targetUser, msgContent)
		if err != nil {
			logger.SetToLogger(logrus.ErrorLevel, "communicationRequestHandler", "error to wirte message to mongodb", err.Error())
			return err
		}
	}
	return nil
}

//连接请求
func connectionRequestHandler(msg *fromClientMsg, conn net.Conn) error {
	realMsg := *msg
	sourceUserToken := realMsg.MsgFrom
	sourceUser, err := ParseToken(sourceUserToken)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "communicationRequestHandler", "error to parse token from client", err.Error())
		return err
	}

	rwLock.Lock()
	UserConnectionMap[sourceUser] = conn
	rwLock.Unlock()

	return nil
}

func loopFindOnlineUserAndSendMessage(){
	for {
		if len(UserConnectionMap) == 0 {
			time.Sleep(30 * time.Second)
			continue
		}
		for k, v := range UserConnectionMap {
			logger.SetToLogger(logrus.InfoLevel, "loopFindOnlineUserAndSendMessage", "online user: " + k, "")
			go responseToClient(k, v)
		}
		time.Sleep(5 * time.Second)
	}
}

func responseToClient(username string, conn net.Conn){
	recorder := dataBase.ReadMessageOfUser(username)
	if recorder != nil {
		msg := toClientMsg{
			MsgFrom: (*recorder).SourceUser,
			MsgContent: (*recorder).Content,
		}
		buf, err := json.Marshal(msg)
		if err != nil {
			logger.SetToLogger(logrus.ErrorLevel, "responseToClient", "error to marshal Recorder to json bytes",err.Error())
			return
		}
		count, err := conn.Write(buf)
		if err != nil {
			logger.SetToLogger(logrus.ErrorLevel, "responseToClient", "error to write bytes to client",err.Error())
			return
		} else {
			logger.SetToLogger(logrus.ErrorLevel, "responseToClient", "succeed send message to client","bytes number: "+ strconv.Itoa(count))
		}
		_ = dataBase.DeleteMessageOfUser(username)
	}
}

func isUserConnected(username string) bool {
	rwLock.RLock()
	_, ok := UserConnectionMap[username]
	rwLock.RUnlock()
	return ok
}

func getUserConnection(username string) net.Conn {
	rwLock.RLock()
	res, ok := UserConnectionMap[username]
	rwLock.RUnlock()
	if !ok {
		return nil
	} else {
		return res
	}
}