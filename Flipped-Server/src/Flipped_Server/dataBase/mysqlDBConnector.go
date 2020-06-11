// @Title  mysqlDBConnector.go
// @Description  To provide a database interface of mysql to the Server
// @Author  郑康
// @Update  郑康 2020.5.25
package dataBase

import (
	"Flipped_Server/initialSetting"
	"Flipped_Server/logger"
	"Flipped_Server/utils"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	_ "github.com/widuu/gojson"
)

// 包内全局变量，存放于数据库的指针
var mysqlDB *sql.DB

var (
	engine   string
	username string
	password string
	ip       string
	port     string
	dbName   string
)

// 数据库结构体用于连接时使用
type DBInfo struct {
	Engine   string //引擎名称
	UserName string //登录数据库的用户名
	PassWord string //登录数据库的密码
	IP       string //数据库的公网IP地址
	Port     string //数据库的端口号
	DBName   string //数据库名称
}

func initialSettingsMysql() {
	mysqlSettings := initialSetting.DataBaseConfig["mysql"].(map[string]interface{})
	engine = utils.AesDecrypt(mysqlSettings["engine"].(string), initialSetting.AESKey)
	username = utils.AesDecrypt(mysqlSettings["userName"].(string), initialSetting.AESKey)
	password = utils.AesDecrypt(mysqlSettings["pwd"].(string), initialSetting.AESKey)
	ip = utils.AesDecrypt(mysqlSettings["host"].(string), initialSetting.AESKey)
	port = utils.AesDecrypt(mysqlSettings["port"].(string), initialSetting.AESKey)
	dbName = utils.AesDecrypt(mysqlSettings["dbName"].(string), initialSetting.AESKey)
}

// @title    Init
// @description   数据库初始化函数，用户对数据库进行初始化、连接等操作
// @auth      郑康             2020.5.17
// @param     void
// @return    void
func Init() {
	initialSettingsMysql()
	//db := DBInfo{
	//	Engine:   "mysql",
	//	UserName: "admin",
	//	PassWord: "mountain",
	//	IP:       "47.94.134.159",
	//	Port:     "3306",
	//	DBName:   "im",
	//}
	db := DBInfo{
		Engine:   engine,
		UserName: username,
		PassWord: password,
		IP:       ip,
		Port:     port,
		DBName:   dbName,
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
		userInfo, _ = MakeUserInfoStruct(&rowMap)
		res = append(res, userInfo)
	}
	return res
}

func SelectSimilarUser(username string) (*UserInfoTable, error) {
	SQL := "Select * From im.userinfo Order By Rand() Limit 20"
	currentUser, err := FindUserInfo(username, "")
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "SelectSimilarUser", "to find user by username: "+username, err.Error())
		return nil, err
	}
	userList := ExecSelectSQL(SQL)
	userSimilarMap := make(map[int]float32)
	maxSimilarityIndex, maxSimilarValue := 0, float32(0.0)
	curFriendList, err := GetFriendListByUserName(username)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "SelectSimilarUser", "error to get friend list of "+username, err.Error())
		return nil, err
	}
	for index := range userList {
		if userList[index].Username == username || utils.Contains(curFriendList, userList[index].Username) {
			continue
		}
		res := CalculateSimilarity(currentUser, userList[index])
		if res > maxSimilarValue {
			maxSimilarValue = res
			maxSimilarityIndex = index
		}
		userSimilarMap[index] = res
	}
	return userList[maxSimilarityIndex], nil
}

// @title    	FindUserInfo
// @description   								通过用户名和密码在数据库中查找完整的信息
// @auth      	郑康           					2020.5.25
// @param     	string, string					用户名, 密码
// @return    	*dataBase.UserInfoTable, error	用户信息结构体指针, 错误信息
func FindUserInfo(username string, pwd string) (*UserInfoTable, error) {
	SQL := "SELECT * FROM im.userinfo \nWHERE username = '" + username + "'"
	if pwd != "" {
		SQL += "And password = '" + pwd + "';"
	} else {
		SQL += ";"
	}
	fmt.Println(SQL)
	data := ExecSelectSQL(SQL)
	if data == nil || len(data) != 1 {
		return nil, errors.New("the data you select is nil or has repetitive")
	}
	return data[0], nil
}

// @title    CLoseMySqlClient
// @description   			关闭mysql连接
// @auth      郑康       	2020.5.26
// @param     void
// @return    void
func CloseMySqlClient() {
	if mysqlDB != nil {
		defer mysqlDB.Close()
	}
	logger.Logger.WithFields(logrus.Fields{
		"function": "CLoseMySqlClient",
		"cause":    "close mysql connection",
	})
	mysqlDB = nil
}

// @title    	DoesUserExist
// @description   						通过用户名判断用户是否存在
// @auth      	郑康           			2020.6.10
// @param     	string					用户名
// @return    	bool					是否存在
func DoesUserExist(username string) bool {
	if username == "" {
		return false
	}
	userinfo, err := FindUserInfo(username, "")
	if err != nil || userinfo == nil {
		return false
	}
	return true
}
