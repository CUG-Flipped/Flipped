package messageQueue

import (
	"Flipped_Server/logger"
	"Flipped_Server/utils"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"strconv"
)

var(
	Client *nats.EncodedConn
)

func InitMQ(){
	natsConnection, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(natsConnection, nats.JSON_ENCODER)
	defer c.Close()
	Client = c
	_, err := Client.Subscribe("im", SendMessageToClient)
	if err != nil {
		panic(err)
	}
	<- utils.ExitFlag
}

func PublishMsg(subject string, value interface{}) error{
	err := Client.Publish(subject, value)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "PublishMsg", "error to publish message to mq", err.Error())
		return err
	}
	logger.SetToLogger(logrus.InfoLevel,"PublishMsg", "Succeed to publish message to mq", "")
	return nil
}

func SendMessageToClient(recorder *utils.Recorder){
	conn := utils.GetUserConnection(recorder.TargetUser)
	if conn == nil{
		logger.SetToLogger(logrus.ErrorLevel, "SendMessageToClient", "error to get user connection: " + recorder.TargetUser, "republish message to mq")
		_ = PublishMsg("im", recorder)
		return
	}
	msg := utils.ToClientMsg{
		MsgFrom: (*recorder).SourceUser,
		MsgContent: (*recorder).Content,
	}
	buf, err := json.Marshal(msg)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "SendMessageToClient", "error to marshal Recorder to json bytes",err.Error())
		return
	}
	count, err := conn.Write(buf)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "SendMessageToClient", "error to write bytes to client", err.Error())
		return
	} else {
		logger.SetToLogger(logrus.ErrorLevel, "SendMessageToClient", "succeed send message to client","bytes number: "+ strconv.Itoa(count))
	}
}

