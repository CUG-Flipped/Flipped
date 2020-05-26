// @Title  mysqlDBConnector.go
// @Description  To provide a database interface of mysql to the Server
// @Author  郑康
// @Update  郑康 2020.5.25
package dataBase

import (
	"Flipped_Server/logger"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	_ "github.com/widuu/gojson"
)

// 包内全局变量，存放于数据库的指针
var mysqlDB *sql.DB

// 数据库结构体用于连接时使用
type DBInfo struct {
	Engine   string //引擎名称
	UserName string //登录数据库的用户名
	PassWord string //登录数据库的密码
	IP       string //数据库的公网IP地址
	Port     string //数据库的端口号
	DBName   string //数据库名称
}

// @title    Init
// @description   数据库初始化函数，用户对数据库进行初始化、连接等操作
// @auth      郑康             2020.5.17
// @param     void
// @return    void
func Init() {
	db := DBInfo{
		Engine:   "mysql",
		UserName: "admin",
		PassWord: "mountain",
		IP:       "47.94.134.159",
		Port:     "3306",
		DBName:   "im",
	}

	database, err := sql.Open(db.Engine, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db.UserName, db.PassWord, db.IP, db.Port, db.DBName))

	if err != nil {
		fmt.Printf("mysql open err: %s\n", err)

		logger.Logger.WithFields(logrus.Fields{
			"function": "Init",
			"cause":    "fail to connect to remote database",
		}).Error(err.Error())

		return
	} else {
		logger.Logger.WithFields(logrus.Fields{
			"function": "Init",
			"cause":    "succeed to connect to database",
		}).Info("mysql open Successfully")
	}
	mysqlDB = database
}

// @title    ExecSQL
// @description   			执行给定(insert、 update、delete)的SQL语句
// @auth      郑康       	2020.5.17
// @param     string		sql语句字符串
// @return    string；error	处理结果字符串；错误信息
func ExecSQL(sql string) (string, error) {
	if mysqlDB != nil {
		result, err := mysqlDB.Exec(sql)
		if err != nil {
			return "", err
		} else {
			lastInsertID, _ := result.LastInsertId()
			affectRowCount, _ := result.RowsAffected()
			return fmt.Sprintf("LastInsertID: %d, affected row count: %d", lastInsertID, affectRowCount), err
		}
	} else {
		logger.Logger.WithFields(logrus.Fields{
			"function": "ExecSQL",
			"cause":    "DataBase does't initialise",
		}).Fatal("DataBase does't initialise, pointer is nil")
		return "DataBase does't initialise, pointer is nil", errors.New("DataBase does't initialise, pointer is nil")
	}
}

// @title    ExecSelectSQL
// @description   			执行给定(select)的SQL语句
// @auth      郑康       	2020.5.25
// @param     string		sql语句字符串
// @return    []*UserInfoTable 根据sql语句选择出的用户信息表
func ExecSelectSQL(sql string) []*UserInfoTable {
	rows, err := mysqlDB.Query(sql)

	if err != nil {
		fmt.Println(err.Error())
	}
	if rows == nil {
		return nil
	}
	userInfoList := rowsMapper(rows)
	return userInfoList
}

// @title    rowsMapper
// @description   			解析select结果
// @auth      郑康       	2020.5.25
// @param     string		select结果
// @return    []*UserInfoTable 根据解析select结果构成的用户信息表
func rowsMapper(rows *sql.Rows) []*UserInfoTable {
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err.Error())
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	var res []*UserInfoTable

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			fmt.Println(err.Error())
		}
		rowMap := make(map[string]string)
		var value string

		for i, col := range values {
			if col != nil {
				value = string(col)
				rowMap[columns[i]] = value
			}
		}
		var userInfo *UserInfoTable
		userInfo, err = MakeUserInfoStruct(&rowMap)
		res = append(res, userInfo)
	}
	return res
}

// @title    CLoseMySqlClient
// @description   			关闭mysql连接
// @auth      郑康       	2020.5.26
// @param     void
// @return    void
func CLoseMySqlClient()  {
	if mysqlDB != nil {
		defer mysqlDB.Close()
	}
	logger.Logger.WithFields(logrus.Fields{
		"function": "CLoseMySqlClient",
		"cause": "close mysql connection",
	})
	mysqlDB = nil
}