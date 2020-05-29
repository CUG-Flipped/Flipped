// @Title  httpServer.go
// @Description  To provide a database interface of Redis to the Server
// @Author  郑康
// @Update  郑康 2020.5.26

package dataBase

import (
	"Flipped_Server/logger"
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
)

//文件内全局变量，Redis连接指针
var (
	client *redis.Conn
	timeout = 3*3600
)

// @title    RedisClient_Init
// @description   			初始化Redis数据库
// @auth      郑康       	2020.5.26
// @param     void
// @return    void
func RedisClientInit(){
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	//Client, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		logger.Logger.WithFields(logrus.Fields {
			"function": "RedisClient_Init",
			"cause": "connect to redis server",
		}).Error(err.Error())
		return
	}
	client = &c
}

// @title    CloseRedisClient
// @description   			关闭redis连接
// @auth      郑康       	2020.5.26
// @param     void
// @return    void
func CloseRedisClient()  {
	if client != nil {
		defer (*client).Close()
	}
	logger.Logger.WithFields(logrus.Fields{
		"function": "CloseRedisClient",
		"cause": "close redis connection",
	})
	client = nil
}

// @title    WriteToRedis
// @description   				向Redis数据库写入string键值对，并设置过期时间3小时
// @auth      郑康       		2020.5.26
// @param     string, string	键、值
// @return    error				错误信息
func WriteToRedis(key string, value string) error{
	res1, err1 := (*client).Do("SET", key, value)
	if err1 != nil {
		logger.Logger.WithFields(logrus.Fields {
			"function": "WriteToRedis",
			"cause": "set key to redis failed",
		}).Error(err1.Error())
		return err1

	} else {
		logger.Logger.WithFields(logrus.Fields {
			"function": "WriteToRedis",
			"cause": "set key to redis successfully",
		}).Info(res1)
	}

	res2, err2 := (*client).Do("EXPIRE", key, timeout)
	if err2 != nil {
		logger.Logger.WithFields(logrus.Fields {
			"function": "WriteToRedis",
			"cause": "set expire to redis failed",
		}).Error(err2.Error())
		return err2

	} else {
		logger.Logger.WithFields(logrus.Fields {
			"function": "WriteToRedis",
			"cause": "set expire to redis successfully",
		}).Info(res2)
	}
	return nil
}

// @title    ReadFromRedis
// @description   			向Redis数据库读取值，如果读取错误或者键对应的值不存在就会报错
// @auth      郑康       	2020.5.26
// @param     string		键
// @return    error			值、错误信息
func ReadFromRedis(key string) (string, error) {
	reply, err := (*client).Do("GET", key)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields {
			"function": "ReadFromRedis",
			"cause": "read from redis",
		}).Error(err.Error())
		return "", err
	}
	if reply == nil {
		logger.Logger.WithFields(logrus.Fields {
			"function": "ReadFromRedis",
			"cause": "get reply",
		}).Error("value is Empty")
		return "", errors.New("nil value of the key")
	}
	res, err1 := (*client).Do("EXPIRE", key, timeout)
	if err1 != nil {
		logger.Logger.WithFields(logrus.Fields {
			"function": "ReadFromRedis",
			"cause": "update the expire time failed",
		}).Error(err1.Error())
		return "", err1
	} else {
		logger.Logger.WithFields(logrus.Fields {
			"function": "ReadFromRedis",
			"cause": "update the expire time successfully",
		}).Info(res)
	}
	value := string(reply.([]uint8))
	return value, nil
}

// @title    ReadFromRedis
// @description   			判断Redis数据库中是否存在某个键
// @auth      郑康       	2020.5.26
// @param     string		键
// @return    error			值、错误信息
func KeyExists(key string) bool{
	exists, _ := redis.Bool((*client).Do("EXISTS", key))
	return exists
}
