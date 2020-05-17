package dataBase

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var mysqlDB *sql.DB

type DBInfo struct {
	Engine string
	UserName string
	PassWord string
	IP string
	Port string
	Table string
}

func Init() {
	db := DBInfo{
		Engine: "mysql",
		UserName: "admin",
		PassWord: "mountain",
		IP: "47.94.134.159",
		Port: "3306",
		Table: "userinfo",
	}

	database, err := sql.Open(db.Engine, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db.UserName, db.PassWord, db.IP, db.Port, db.Table))

	if err != nil {
		fmt.Printf("mysql open err: %s\n", err)
		return
	} else {
		fmt.Println("mysql open Successfully")
	}
	mysqlDB = database

	defer database.Close()
	mysqlDB = nil
}

