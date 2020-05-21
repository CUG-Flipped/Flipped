// @Title  mysqlDBConnector.go
// @Description  To provide a database interface of mysql to the Server
// @Author  郑康
// @Update  郑康 2020.5.18
package dataBase

import (
	"Flipped_Server/logger"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// 包内全局变量，存放于数据库的指针
var mysqlDB *sql.DB

// 数据库结构体用于连接时使用
type DBInfo struct {
	Engine string //引擎名称
	UserName string //登录数据库的用户名
	PassWord string //登录数据库的密码
	IP string //数据库的公网IP地址
	Port string //数据库的端口号
	DBName string //数据库名称
}

// @title    Init
// @description   数据库初始化函数，用户对数据库进行初始化、连接等操作
// @auth      郑康             2020.5.17
// @param     void
// @return    void
func Init() {
	db := DBInfo{
		Engine: "mysql",
		UserName: "admin",
		PassWord: "mountain",
		IP: "47.94.134.159",
		Port: "3306",
		DBName: "im",
	}

	database, err := sql.Open(db.Engine, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db.UserName, db.PassWord, db.IP, db.Port, db.DBName))

	if err != nil {
		fmt.Printf("mysql open err: %s\n", err)

		logger.Logger.WithFields(logrus.Fields {
			"function" : "Init",
			"cause": "fail to connect to remote database",
		}).Error(err.Error())

		return
	} else {
		logger.Logger.WithFields(logrus.Fields {
			"function" : "Init",
			"cause": "succeed to connect to database",
		}).Info("mysql open Successfully")
	}
	mysqlDB = database

	//defer database.Close()
	//mysqlDB = nil
}

// @title    ExecSQL
// @description   			执行给定的SQL语句
// @auth      郑康       	2020.5.17
// @param     string		sql语句字符串
// @return    string；error	处理结果字符串；错误信息
func ExecSQL(sql string) (string, error){
	if mysqlDB != nil{
		result, err := mysqlDB.Exec(sql)
		if err != nil {
			return "", err
		} else {
			lastInsertID, _:= result.LastInsertId()
			affectRowCount, _ := result.RowsAffected()
			return fmt.Sprintf("LastInsertID: %d, affected row count: %d", lastInsertID, affectRowCount), err
		}
	} else {
		logger.Logger.WithFields(logrus.Fields {
			"function" : "ExecSQL",
			"cause": "DataBase does't initialise",
		}).Fatal("DataBase does't initialise, pointer is nil")
		return "DataBase does't initialise, pointer is nil", errors.New("DataBase does't initialise, pointer is nil")
	}
}

