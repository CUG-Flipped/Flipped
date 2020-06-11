package utils

import (
	"net"
	"sync"
)

var (
	UserConnectionMap map[string]net.Conn
	RWLock            = &sync.RWMutex{}
)

type FromClientMsg struct {
	MsgType    int    `json:"MsgType"`
	MsgFrom    string `json:"MsgFrom"`
	MsgTo      string `json:"MsgTo"`
	MsgContent string `json:"MsgContent"`
}

type ToClientMsg struct {
	MsgFrom    string `json:"MsgFrom"`
	MsgContent string `json:"MsgContent"`
}

type ReplyMsg struct {
	ResultCode int    `json:"ResultCode"`
	MsgContent string `json:"MsgContent"`
}

type Recorder struct {
	TargetUser string `bson:"targetUser"`
	SourceUser string `bson:"sourceUser"`
	Content    string `bson:"content"`
}

func IsUserConnected(username string) bool {
	RWLock.RLock()
	_, ok := UserConnectionMap[username]
	RWLock.RUnlock()
	return ok
}

func GetUserConnection(username string) net.Conn {
	RWLock.RLock()
	res, ok := UserConnectionMap[username]
	RWLock.RUnlock()
	if !ok {
		return nil
	} else {
		return res
	}
}

func DeleteConnection(conn net.Conn) {
	for username, curConn := range UserConnectionMap {
		if curConn == conn {
			delete(UserConnectionMap, username)
		}
	}
}
