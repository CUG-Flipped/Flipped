// @Title  mysqlDBConnector.go
// @Description  To provide a map between struct and database tables
// @Author  郑康
// @Update  郑康 2020.5.18
package dataBase

import (
	"fmt"
	"strconv"
)

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


