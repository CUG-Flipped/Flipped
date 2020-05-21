// @Title  mysqlDBConnector.go
// @Description  To provide a map between structs and database tables
// @Author  郑康
// @Update  郑康 2020.5.18
package dataBase

// 用户信息信息表，各个字段对应数据库的各个字段
type UserInfoTable struct {
	Pid        string `sql:"pid"`
	Username   string `sql:"username"`
	Password   string `sql:"password"`
	UserType   byte   `sql:"user_type"`
	Email      string `sql:"email"`
	Photo      string `sql:"photo"`
	RealName   string `sql:"realName"`
	Profession string `sql:"profession"`
	Age        byte   `sql:"age"`
	Region     string `sql:"region"`
	Hobby      string `sql:"hobby"`
}




