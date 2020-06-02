package network

import (
	"Flipped_Server/logger"
	"Flipped_Server/messageQueue"
	"Flipped_Server/utils"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
)

type SocketServer struct {
	IPAddr string
	Port   int
}

func (ss *SocketServer) Run() {
	server, err := net.Listen("tcp", ss.IPAddr+":"+strconv.Itoa(ss.Port))
	utils.UserConnectionMap = make(map[string]net.Conn)

	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "Run", "error to start socket server", err.Error())
		return
	}

	go messageQueue.InitMQ()

	for {
		conn, err := server.Accept()
		if err != nil {
			logger.SetToLogger(logrus.ErrorLevel, "Run", "error to Accept socket", err.Error())
		} else {
			logger.SetToLogger(logrus.InfoLevel, "Run", "succeed to Accept socket", "")
			connectionHandler(conn)
		}
	}
}

//ToDo: 持续监听该客户端的消息
func connectionHandler(conn net.Conn) {

	data := json.NewDecoder(conn)
	var msg utils.FromClientMsg
	err := data.Decode(&msg)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "connectionHandler", "error to decode message sent from client", err.Error())
		return
	}
	switch msg.MsgType {
	case 1:
		err = communicationRequestHandler(&msg, conn)
		if err != nil {
			reply := utils.ReplyMsg{
				ResultCode: 500,
				MsgContent: "Some error occur in the server, please try it latter",
			}
			buf, _ := json.Marshal(reply)
			_, _ = conn.Write(buf)
		}
	case 2:
		err = connectionRequestHandler(&msg, conn)
		if err != nil {
			replyToClient(conn, 500, "Some error occur in the server, please try it latter")
		}
	default:
		logger.SetToLogger(logrus.ErrorLevel, "connectionHandler", "error type in communicationMsg struct", "")
	}
}

//交流请求
func communicationRequestHandler(msg *utils.FromClientMsg, conn net.Conn) error {
	sourceUserToken := msg.MsgFrom
	sourceUser, err := ParseToken(sourceUserToken)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "communicationRequestHandler", "error to parse token from client", err.Error())
		return err
	}

	utils.RWLock.Lock()
	utils.UserConnectionMap[sourceUser] = conn
	utils.RWLock.Unlock()

	targetUser := msg.MsgTo
	msgContent := msg.MsgContent
	//如果目标用户处于与服务器TCP连接状态,并且能够正常获得目标用户的连接
	if targetConn := utils.GetUserConnection(targetUser); utils.IsUserConnected(targetUser) && targetConn != nil {
		resultCode := 200
		replyContent := "your friend is not online currently, server will send the message latter to your friend"
		//获取目标用户连接
		communicationMsg := utils.ToClientMsg{
			MsgFrom:    targetUser,
			MsgContent: msgContent,
		}
		buf, err := json.Marshal(communicationMsg)
		if err != nil {
			logger.SetToLogger(logrus.ErrorLevel, "communicationRequestHandler", "error to parse struct to json", err.Error())
			resultCode = 400
			replyContent = "your request is not formatted in required form, please check your source code"
			replyToClient(conn, resultCode, replyContent)
			return err
		}
		count, err := targetConn.Write(buf)
		if err != nil {
			logger.SetToLogger(logrus.ErrorLevel, "communicationRequestHandler", "error to send to client with bytes", err.Error())
			resultCode = 500
			replyContent = "error to send message to your friend, server will push the message to the messageQueue"
			err = messageQueue.PublishMsg("im", &utils.Recorder{TargetUser: targetUser, SourceUser: sourceUser, Content: msgContent})
			//无法向目标客户端写入数据说明与该客户端的TCP连接已经已经失效
			utils.RWLock.Lock()
			defer targetConn.Close()
			delete(utils.UserConnectionMap, targetUser)
			utils.RWLock.Unlock()
		} else {
			logger.SetToLogger(logrus.InfoLevel, "communicationRequestHandler", "Succeed to send message from "+sourceUser+" to "+targetUser, "bytes num: "+strconv.Itoa(count))
		}
		replyToClient(conn, resultCode, replyContent)
		return err
	} else {
		newRecorder := &utils.Recorder{SourceUser: sourceUser, TargetUser: targetUser, Content: msgContent}
		err := messageQueue.PublishMsg("im", newRecorder)
		resultCode := 200
		replyContent := "your friend is not online currently, server will send the message latter to your friend"
		if err != nil {
			logger.SetToLogger(logrus.ErrorLevel, "communicationRequestHandler", "error to publish message to NSQ", err.Error())
			resultCode = 500
			replyContent = "there is something wrong in server, please send your message again"
		} else {
			logger.SetToLogger(logrus.InfoLevel, "communicationRequestHandler", "succeed to publish message to NSQ", "")
		}
		replyToClient(conn, resultCode, replyContent)
		return err
	}
}

//连接请求
func connectionRequestHandler(msg *utils.FromClientMsg, conn net.Conn) error {
	realMsg := *msg
	sourceUserToken := realMsg.MsgFrom
	sourceUser, err := ParseToken(sourceUserToken)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "communicationRequestHandler", "error to parse token from client", err.Error())
		return err
	}

	utils.RWLock.Lock()
	utils.UserConnectionMap[sourceUser] = conn
	utils.RWLock.Unlock()
	replyToClient(conn, 200, "succeed to connect to socketServer")
	return nil
}

func replyToClient(conn net.Conn, resultCode int, msgContent string) {
	reply := utils.ReplyMsg{
		ResultCode: resultCode,
		MsgContent: msgContent,
	}
	buf, _ := json.Marshal(reply)
	_, _ = conn.Write(buf)
}
