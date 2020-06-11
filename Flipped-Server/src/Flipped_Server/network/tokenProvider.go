// @Title  httpServer.go
// @Description  To provide a token provider interface
// @Author  郑康
// @Update  郑康 2020.5.25
package network

import (
	"Flipped_Server/dataBase"
	"Flipped_Server/logger"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"time"
)

//自定义的Claim
type MyClaims struct {
	UserName string `json:"username"`
	jwt.StandardClaims
}

//设置Token过期时间为1小时
const TokenExpireDuration = 3 * 3600

//自定义Secret
var MySecret = []byte("First Blood")

// @title    GenerateToken
// @description   			通过用户名生成token字符串,并将token字符串与用户名组成键值对存入redis中
// @auth      郑康       	2020.5.25
// @param     string		用户名
// @return    string；error	token字符串；错误信息
func GenerateToken(username string) (string, error) {
	//创建自定义声明
	claim := MyClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "MrSecond",
		},
	}
	//使用指定的方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString(MySecret)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"function": "GenerateToken",
			"cause":    "create new Claims",
		}).Error(err.Error())
		return "", err
	}
	err1 := dataBase.WriteToRedis(username, tokenStr)

	if err1 != nil {
		logger.Logger.WithFields(logrus.Fields{
			"function": "GenerateToken",
			"cause":    "Write To Redis",
		}).Error(err1.Error())
		return "", err1
	}
	return tokenStr, nil
}

// @title    ParseToken
// @description   			通过token字符串解析用户名, 同时判断token是否合法
// @auth      郑康       	2020.5.25
// @param     string		token字符串
// @return    string；error	用户名字符串；错误信息
func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		logger.SetToLogger(logrus.InfoLevel, "ParseToken", "Parse tokenStr", err.Error())
	}
	if token == nil {
		return "", err
	}

	if claims, ok := token.Claims.(*MyClaims); ok && dataBase.KeyExists(claims.UserName, 0) {
		return claims.UserName, nil
	}
	return "", errors.New("invalid token")
}
