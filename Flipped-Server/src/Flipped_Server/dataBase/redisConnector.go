// @Title  httpServer.go
// @Description  To provide a database interface of Redis to the Server
// @Author  郑康
// @Update  郑康 2020.5.26

package dataBase

import (
	"Flipped_Server/initialSetting"
	"Flipped_Server/logger"
	"errors"
	"fmt"
	red "github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

//文件内全局变量，Redis连接指针
var (
	client           *red.Conn
	timeout          = 3 * 3600
	timeoutHeartBeat = 60
	redisIP          string
	redisPort        string
	redis            *Redis
)

type Redis struct {
	pool *red.Pool
}

func initialSettingsRedis() {
	redisSetting := initialSetting.DataBaseConfig["redis"].(map[string]interface{})
	redisIP = redisSetting["host"].(string)
	redisPort = redisSetting["port"].(string)
}

// @title    RedisClient_Init
// @description   			初始化Redis数据库
// @auth      郑康       	2020.5.26
// @param     void
// @return    void
func RedisClientInit() {
	initialSettingsRedis()
	redis = new(Redis)
	host := redisIP + ":" + redisPort
	redis.pool = &red.Pool{
		MaxIdle:     256,
		MaxActive:   0,
		IdleTimeout: time.Duration(120),
		Dial: func() (red.Conn, error) {
			return red.Dial(
				"tcp",
				host,
				red.DialReadTimeout(time.Duration(1000)*time.Millisecond),
				red.DialWriteTimeout(time.Duration(1000)*time.Millisecond),
				red.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
				red.DialDatabase(0),
			)
		},
	}
	//c, err := red.Dial("tcp", redisIP+":"+redisPort)
	////c, err := redis.Dial("tcp", "39.99.190.67:6379")
	//if err != nil {
	//	logger.Logger.WithFields(logrus.Fields{
	//		"function": "RedisClient_Init",
	//		"cause":    "connect to redis server",
	//	}).Error(err.Error())
	//	return
	//}
	//client = &c
}

// @title    CloseRedisClient
// @description   			关闭redis连接
// @auth      郑康       	2020.5.26
// @param     void
// @return    void
func CloseRedisClient() {
	if client != nil {
		defer (*client).Close()
	}
	logger.Logger.WithFields(logrus.Fields{
		"function": "CloseRedisClient",
		"cause":    "close redis connection",
	})
	client = nil
}

// @title    WriteToRedis
// @description   				向Redis db1数据库写入string键值对，并设置过期时间3小时
// @auth      郑康       		2020.5.26
// @param     string, string	键、值
// @return    error				错误信息
func WriteToRedis(key string, value string) error {
	conn := redis.pool.Get()
	if err:= conn.Err(); err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "WriteToRedis", "error to get instance from redis pool", err.Error())
		return err
	}
	defer conn.Close()
	res1, err1 := conn.Do("SET", key, value)
	if err1 != nil {
		logger.SetToLogger(logrus.ErrorLevel, "WriteToRedis", "set key to redis failed", err1.Error())
		return err1

	} else {
		logger.SetToLogger(logrus.InfoLevel, "WriteToRedis", "set key to redis successfully", res1.(string))
	}
	res2, err2 := conn.Do("EXPIRE", key, timeout)
	if err2 != nil {
		logger.SetToLogger(logrus.ErrorLevel, "WriteToRedis", "set expire to redis failed", err2.Error())
		return err2
	} else {
		logger.SetToLogger(logrus.InfoLevel, "WriteToRedis", "set expire to redis successfully", string(res2.(int64)))
	}
	return nil
}

// @title    ReadFromRedis
// @description   			向Redis db1数据库读取值，如果读取错误或者键对应的值不存在就会报错
// @auth      郑康       	2020.5.26
// @param     string		键
// @return    error			值、错误信息
func ReadFromRedis(key string) (string, error) {
	conn := redis.pool.Get()
	if err:= conn.Err(); err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "ReadFromRedis", "error to get instance from redis pool", err.Error())
		return "", err
	}
	defer conn.Close()

	reply, err := conn.Do("GET", key)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "ReadFromRedis", "read from redis", err.Error())
		return "", err
	}
	if reply == nil {
		logger.SetToLogger(logrus.ErrorLevel, "ReadFromRedis", "get reply", "")
		return "", errors.New("nil value of the key")
	}
	res, err1 := conn.Do("EXPIRE", key, timeout)
	if err1 != nil {
		logger.SetToLogger(logrus.ErrorLevel, "ReadFromRedis", "update the expire time failed", err1.Error())
		return "", err1
	} else {
		logger.SetToLogger(logrus.InfoLevel, "ReadFromRedis", "update the expire time successfully", string(res.(int64)))
	}
	value := string(reply.([]uint8))
	return value, nil
}

// @title    ReadFromRedis
// @description   			判断Redis数据库中是否存在某个键
// @auth      郑康       	2020.5.26
// @param     string		键
// @return    error			值、错误信息
func KeyExists(key string, dbNum int) bool {
	conn := redis.pool.Get()
	if err:= conn.Err(); err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "KeyExists", "error to get instance from redis pool", err.Error())
		return false
	}
	defer conn.Close()
	reply, err := conn.Do("Select", dbNum)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "KeyExists", "Select DB: "+strconv.Itoa(dbNum), err.Error())
		return false
	} else {
		logger.SetToLogger(logrus.InfoLevel, "KeyExists", fmt.Sprintf("reply: %v", reply), "")
	}
	exists, _ := red.Bool(conn.Do("EXISTS", key))
	return exists
}

func DeleteKey(key string, dbNum int) error {
	conn := redis.pool.Get()
	if err:= conn.Err(); err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "KeyExists", "error to get instance from redis pool", err.Error())
		return err
	}
	defer conn.Close()
	_, err := conn.Do("Select", dbNum)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "DeleteKey", "Error to switch db", err.Error())
		return err
	}
	err = conn.Send("Del", key)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "DeleteKey", "Error to delete key: "+key, err.Error())
		return err
	}
	return nil
}

func CountOnlineUsers() (int, error) {
	conn := redis.pool.Get()
	if err:= conn.Err(); err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "CountOnlineUsers", "error to get instance from redis pool", err.Error())
		return 0, err
	}
	defer conn.Close()

	reply, err := conn.Do("Select", 2)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "CountOnlineUsers", "Switch to redis db 2", err.Error())
		return -1, err
	}
	_reply, _ := red.String(reply, err)
	logger.SetToLogger(logrus.InfoLevel, "CountOnlineUsers", "Switch to redis db 2", _reply)
	reply, err = conn.Do("KEYS", "*")
	_, _ = conn.Do("Select", 0)
	return len(reply.([]interface{})), nil
}

// 实时更新用户状态， 如果用户存在于redis db2中，则更新其过期时间60s， 不存在则建立该key
func UpdateUserStatus(username string, lock *sync.Mutex) {
	conn := redis.pool.Get()
	if err:= conn.Err(); err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "UpdateUserStatus", "error to get instance from redis pool", err.Error())
		return
	}
	defer conn.Close()

	lock.Lock()
	isExist := KeyExists(username, 2)
	if isExist {
		_, _ = conn.Do("Expire", username, timeoutHeartBeat)
	} else {
		_, err := conn.Do("Set", username, "alive")
		if err != nil {
			logger.SetToLogger(logrus.ErrorLevel, "UpdateUserStatus", "Error to Set Key-Value", err.Error())
			return
		}
		_, _ = conn.Do("Expire", username, timeoutHeartBeat)
	}
	lock.Unlock()
}
