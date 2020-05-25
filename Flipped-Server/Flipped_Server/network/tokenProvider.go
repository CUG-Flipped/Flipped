// @Title  httpServer.go
// @Description  To provide a token provider interface
// @Author  郑康
// @Update  郑康 2020.5.25
package network

import (
	"Flipped_Server/logger"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"time"
)

type MyClaims struct {
	UserName string `json:"username"`
	jwt.StandardClaims
}

//设置Token过期时间为1小时
const TokenExpireDuration = time.Hour

//自定义Secret
var MySecret = []byte("First Blood")

func GenerateToken(username string) (string, error) {
	//创建自定义声明
	claim := MyClaims{
		username,
	 	jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer: "MrSecond",
		},
	}
	//使用指定的方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString(MySecret)
}

func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"function": "ParseToken",
			"cause": "Parse tokenStr",
		}).Info(err.Error())
		return "", err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims.UserName, nil
	}
	return "", errors.New("invalid token")
}