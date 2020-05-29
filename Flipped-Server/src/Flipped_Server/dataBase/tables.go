// @Title  mysqlDBConnector.go
// @Description  To provide a map between struct and database tables
// @Author  郑康
// @Update  郑康 2020.5.18
package dataBase

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type UserInfo UserInfoTable

// 用户信息信息表，各个字段对应数据库的各个字段
type UserInfoTable struct {
	Pid        string `sql:"pid"`
	Username   string `sql:"username"`
	Password   string `sql:"password"`
	UserType   int   `sql:"user_type"`
	Email      string `sql:"email"`
	Photo      string `sql:"photo"`
	RealName   string `sql:"realName"`
	Profession string `sql:"profession"`
	Age        int   `sql:"age"`
	Region     string `sql:"region"`
	Hobby      string `sql:"hobby"`
}

func (userInfo *UserInfoTable)String() string {
	res := fmt.Sprintf("pid = %s, username = %s, password = %s, user_type = %d, email = %s, photo = %s, RealName = %s, Profession = %s, Age = %d, Region = %s, Hobby = %s", userInfo.Pid, userInfo.Username, userInfo.Password, userInfo.UserType, userInfo.Email, userInfo.Photo, userInfo.RealName, userInfo.Profession, userInfo.Age, userInfo.Region, userInfo.Hobby)
	return res
}

func MakeUserInfoStruct(data *map[string]string) (*UserInfoTable, error) {
	var userInfo UserInfoTable
	userInfo.Pid = (*data)["pid"]
	userInfo.Username = (*data)["username"]
	userInfo.Password = (*data)["password"]
	integer, err := strconv.Atoi((*data)["user_type"])
	if err != nil {
		fmt.Println(err.Error())
		return new(UserInfoTable) , err
	}
	userInfo.UserType = integer
	userInfo.Email = (*data)["email"]
	userInfo.Photo = (*data)["photo"]
	userInfo.RealName = (*data)["realName"]
	userInfo.Profession = (*data)["profession"]
	integer, err = strconv.Atoi((*data)["age"])
	if err != nil {
		fmt.Println(err.Error())
		return new(UserInfoTable) , err
	}
	userInfo.Age = integer
	userInfo.Region = (*data)["region"]
	userInfo.Hobby = (*data)["hobby"]

	return &userInfo, nil
}

func CalculateSimilarity(user1 *UserInfoTable, user2 *UserInfoTable) float32{
	var similarity float32 = 1.0
	if user1.UserType == user2.UserType {
		similarity *= 0.6
	} else {
		similarity *= 0.4
	}
	user1EmailSuffix := strings.Split(user1.Email,"@")[1]
	user2EmailSuffix := strings.Split(user2.Email, "@")[1]
	if user1EmailSuffix == user2EmailSuffix {
		similarity *= 0.6
	} else {
		similarity *= 0.4
	}

	if int(math.Abs(float64(user1.Age - user2.Age))) < 5 {
		similarity *= 0.7
	} else {
		similarity *= 0.3
	}

	if user1.Region == user2.Region {
		similarity *= 0.9
	} else {
		similarity *= 0.1
	}

	if user1.Profession == user2.Profession {
		similarity *= 0.9
	} else {
		similarity *= 0.1
	}

	if strings.Contains(user1.Hobby, user2.Hobby) || strings.Contains(user2.Hobby, user1.Hobby) {
		similarity *= 0.8
	} else {
		similarity *= 0.2
	}
	return similarity
}

