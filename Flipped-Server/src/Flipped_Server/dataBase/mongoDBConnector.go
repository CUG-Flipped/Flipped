package dataBase

import (
	"Flipped_Server/logger"
	"Flipped_Server/utils"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	mongoHost = "47.94.134.159"
	mongoPort = "27017"
	mongoUser = "root"
	mongoPassword = "mountain"
	collectionName = "friendMap"
	msgCollectionName = "messageMap"
)

var (
	currentDB *mgo.Database
	currentCollection *mgo.Collection
)

type friendMap struct {
	SourceUser string `bson:"sourceUser"`
	FriendList []string `bson:"friendList"`
}

type Recorder struct {
	TargetUser string `bson:"targetUser"`
	SourceUser string `bson:"sourceUser"`
	Content string `bson:"content"`
}

func (fl *friendMap)String() string {
	res := fmt.Sprintf("Source User: %s, Friends: ", fl.SourceUser)
	for index := range fl.FriendList{
		res = res + fl.FriendList[index] + ", "
	}
	return res
}

func InitializeMongoDB() {
	session, err := mgo.Dial(mongoHost + ":" + mongoPort)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"function": "InitializeMongoDB",
			"cause": "connect to mongodb",
		}).Error(err.Error())
		return
	}
	session.SetMode(mgo.Eventual, true)
	tmpDB := session.DB("admin")
	err = tmpDB.Login(mongoUser, mongoPassword)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"function": "InitializeMongoDB",
			"cause": "login mongodb",
		}).Error(err.Error())
		return
	}
	currentDB = session.DB("im")
	currentCollection = currentDB.C(collectionName)
	session.SetPoolLimit(5)
}

func GetFriendListByUserName(username string) ([]string, error) {
	resFriendList := friendMap{}
	err := currentCollection.Find(bson.M {"sourceUser": username}).One(&resFriendList)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"function": "GetFriendListByUserName",
			"cause": "find data in MongoDB",
		}).Error(err.Error())
		return []string{}, err
	} else {
		return resFriendList.FriendList, nil
	}
}

func InitUserFriendList(username string) error {
	newFriendPair := friendMap{SourceUser: username, FriendList: []string{}}
	err := currentCollection.Insert(newFriendPair)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "InitUserFriendList", "insert into mongodb", err.Error())
		return err
	}
	return nil
}

func AddFriend(sourceUser string, targetUser string) error {
	friendList, err := GetFriendListByUserName(sourceUser)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "AddFriend", "find friends of User: " + sourceUser, err.Error())
	}
	selector := bson.M{"sourceUser": sourceUser}
	if len(friendList) != 0 {
		data := bson.M{"$push": bson.M{"friendList":targetUser}}
		err := currentCollection.Update(selector, data)
		if err != nil {
			logger.SetToLogger(logrus.ErrorLevel, "AddFriend", "update friends of User: " + sourceUser, err.Error())
			return err
		}
	} else {
		newFriendPair := friendMap{SourceUser: sourceUser, FriendList: []string{targetUser}}
		err := currentCollection.Insert(newFriendPair)
		if err != nil {
			logger.SetToLogger(logrus.ErrorLevel, "AddFriend", "insert into mongodb", err.Error())
			return err
		}
	}
	return nil
}

func DeleteFriend(sourceUser string, targetUser string) error{
	friendList, err := GetFriendListByUserName(sourceUser)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "DeleteFriend", "find friend list of User: " + sourceUser, err.Error())
	}
	if friendList == nil || len(friendList) == 0{
		logger.SetToLogger(logrus.InfoLevel, "DeleteFriend", "User: "+ sourceUser + " does't have a friend","")
		return nil
	} else {
		if utils.Contains(friendList, targetUser) {
			selector := bson.M{"sourceUser": sourceUser}
			data := bson.M{"$pull": bson.M{"friendList": targetUser}}
			err := currentCollection.Update(selector, data)
			if err != nil {
				logger.SetToLogger(logrus.ErrorLevel, "DeleteFriend", "delete friend of User: "+ sourceUser, err.Error())
				return err
			}
		}
	}
	logger.SetToLogger(logrus.InfoLevel, "DeleteFriend", "Succeed to Delete friend", "")
	return nil
}

func WriteMessage(sourceUser string, targetUser string, content string) error {
	curRecorder := Recorder{SourceUser: sourceUser, Content: content, TargetUser: targetUser}
	currentCollection = currentDB.C(msgCollectionName)
	err := currentCollection.Insert(curRecorder)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "WriteMessage", "insert into mongodb", err.Error())
		return err
	}
	currentCollection = currentDB.C(collectionName)
	return nil
}

func ReadMessageOfUser(username string) *Recorder {
	currentCollection = currentDB.C(msgCollectionName)
	recorder := Recorder{}
	err := currentCollection.Find(bson.M{"targetUser": username}).One(&recorder)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "ReadMessageOfUser", "error to read recorder from mongodb", err.Error())
		return nil
	}
	currentCollection = currentDB.C(collectionName)
	return &recorder
}

func DeleteMessageOfUser(username string) error {
	currentCollection = currentDB.C(msgCollectionName)
	selector := bson.M{"targetUser": username}
	err := currentCollection.Remove(selector)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "DeleteMessageOfUser", "error to remove message", err.Error())
		return err
	}
	currentCollection = currentDB.C(collectionName)
	return nil
}