package dataBase

import (
	"Flipped_Server/logger"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var mysqlDB *sql.DB

type DBInfo struct {
	Engine string
	UserName string
	PassWord string
	IP string
	Port string
	DBName string
}

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

