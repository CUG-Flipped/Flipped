// @Title  sqlMapper.go
// @Description  To provide an extensive sql interface to the Server
// @Author  郑康
// @Update  郑康 2020.5.18
package sqlmapper

import (
	"Flipped_Server/dataBase"
	"Flipped_Server/logger"
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"strconv"
)

//包内全局字符串列表，将tag分类，在stringAttr内的tag表示其字段为字符串，在intAttr内的tag表示器字段为整型
var (
	stringAttr = [...]string{"pid", "username", "photo","password", "email", "realName", "profession", "region", "hobby"}
	intAttr = [...]string{"user_type", "age"}
	stringDefault = ""
	intDefault = 1000
)

// @title    isStringAttr
// @description   给定一个字符串判断该字符串是否在stringAttr列表中
// @auth      郑康           2020.5.17
// @param     string		tag字符串
// @return    bool			判断结果
func isStringAttr(str string) bool {
	for i := 0; i < len(stringAttr); i++ {
		if str == stringAttr[i] {
			return true
		}
	}
	return false
}

// @title    isIntAttr
// @description   给定一个字符串判断该字符串是否在intAttr列表中
// @auth      郑康           2020.5.17
// @param     string		tag字符串
// @return    bool			判断结果
func isIntAttr(str string) bool {
	for i := 0; i < len(intAttr); i++ {
		if str == intAttr[i] {
			return true
		}
	}
	return false
}

// @title    splitDataAndStruct
// @description   使用反射的方法将给定的接口的tag字段与其值进行分离，如果该接口不是tables.go中定义的某一个struct则返回错误
// @auth      郑康           					2020.5.17
// @param     interface{}						接口变量
// @return    []string, []string, error			tag列表；值列表；错误信息
func splitDataAndStruct(metaData interface{}) ([]string, []string, error){
	dataType := reflect.TypeOf(metaData)
	dataVal := reflect.ValueOf(metaData)
	dataKind := dataVal.Kind()
	if dataKind != reflect.Struct {
		return nil, nil, errors.New("expect an struct")
	}
	fieldNum := dataVal.NumField()

	logger.Logger.WithFields(logrus.Fields {
		"function": "splitDataAndStruct",
		"cause": "count member number of the interface",
	}).Infof(fmt.Sprintf("该结构体有%d个字段\n", fieldNum))

	tagArr := make([]string, fieldNum)
	dataArr := make([]string, fieldNum)

	for i := 0; i < fieldNum; i++ {
		tagVal := dataType.Field(i).Tag.Get("sql")

		if tagVal == "" {
			return nil, nil, errors.New("tag is empty")
		}

		tagArr[i] = tagVal
		dataArr[i] = dataVal.Field(i).String()
	}

	return tagArr, dataArr, nil
}

// @title    Insert
// @description   Insert接口函数, 通过给定的结构体和表名来进行Insert操作
// @auth      郑康           					2020.5.17
// @param     interface{}, string				接口变量, 数据表名称
// @return    error								错误信息
func Insert(data interface {}, tableName string) error{
	tagArr, dataArr, err := splitDataAndStruct(data)
	tagLen := len(tagArr)
	dataLen := len(dataArr)
	if err != nil ||tagLen  == 0 ||dataLen  == 0{
		return err
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("INSERT INTO im.%s VALUES(\n", tableName))

	for i := 0; i < tagLen; i++{
		if isStringAttr(tagArr[i]) {
			buffer.WriteString("'")
			buffer.WriteString(dataArr[i])
			buffer.WriteString("'")
		} else if isIntAttr(tagArr[i]) {
			integer, _ := strconv.Atoi(dataArr[i])
			buffer.WriteString(strconv.Itoa(integer))
			buffer.WriteString("")
		} else {
			return errors.New("unexpected value: " + tagArr[i])
		}

		if i != tagLen - 1 {
			buffer.WriteString(",\n")
		}

	}
	buffer.WriteString(");")
	sql := buffer.String()
	logger.Logger.WithFields(logrus.Fields {
		"function": "Insert",
		"cause": "display sql",
	}).Info(sql)

	res, err := dataBase.ExecSQL(sql)
	if err != nil {
		return err
	}
	logger.Logger.WithFields(logrus.Fields {
		"function": "Insert",
		"cause": "succeed to insert data into database",
	}).Info(res)

	return nil
}

// @title    Update
// @description   Update接口函数，通过给定旧数据和新数据来进行Update操作
// @auth      郑康           											2020.5.17
// @param     *map[string]string, *map[string]string, string			旧数据, 新数据, 表名
// @return    error														错误信息
func Update(oldData *map[string]string, newData *map[string]string, tableName string)  error {
	var buffer bytes.Buffer
	buffer.WriteString("Update im." + tableName + " set ")

	newDataLen := len(*newData)
	index := 0
	for key, value := range *newData {
		if isStringAttr(key){
			buffer.WriteString(key + "='" + value + "'")
		} else if isIntAttr(key) {
			integer, _ := strconv.Atoi(value)
			buffer.WriteString(key + "=" + strconv.Itoa(integer))
		} else {
			return errors.New("unexpected value: " + key)
		}
		if index != newDataLen - 1{
			buffer.WriteString(", ")
		}
		index++
	}

	index = 0
	buffer.WriteString(" where ")
	for key, value := range *oldData {
		if isStringAttr(key){
			buffer.WriteString(key + "='" + value + "'")
		} else if isIntAttr(key) {
			buffer.WriteString(key + "=" + value)
		} else {
			return errors.New("unexpected value: " + key)
		}
		if index != newDataLen - 1{
			buffer.WriteString(" and ")
		}
	}
	buffer.WriteString(";")
	sql := buffer.String()
	logger.Logger.WithFields(logrus.Fields {
		"function": "Update",
		"cause": "display sql",
	}).Info(sql)

	fmt.Println(sql)
	return nil
}

// @title    	Delete
// @description   							Delete接口函数，通过给定结构体和表名来进行Update操作
// @auth      	郑康           				2020.5.17
// @param     	interface {}, string		数据, 表名
// @return    	error						错误信息
func Delete(data interface {}, tableName string) error {
	tagArr, dataArr, err := splitDataAndStruct(data)
	tagLen := len(tagArr)
	dataLen := len(dataArr)
	if err != nil ||tagLen  == 0 ||dataLen  == 0{
		return err
	}

	var buffer bytes.Buffer
	buffer.WriteString("Delete from im." + tableName + " where ")
	for i := 0; i < tagLen; i++ {
		if i != 0  {
			buffer.WriteString("and ")
		}

		if isStringAttr(tagArr[i]){
			if dataArr[i] != stringDefault{
				buffer.WriteString(tagArr[i] + "='" + dataArr[i] + "' ")
			}
		} else if isIntAttr(tagArr[i]) {
			integer, _ := strconv.Atoi(dataArr[i])
			if integer != intDefault{
				buffer.WriteString(tagArr[i] + "=" + strconv.Itoa(integer) + " ")
			}
		} else {
			return errors.New("unexpected value: " + tagArr[i])
		}
	}
	buffer.WriteString(";")
	sql := buffer.String()
	logger.Logger.WithFields(logrus.Fields {
		"function": "Delete",
		"cause": "display sql",
	}).Info(sql)
	fmt.Println(sql)

	res, err := dataBase.ExecSQL(sql)
	if err != nil {
		return err
	}
	logger.Logger.WithFields(logrus.Fields {
		"function": "Delete",
		"cause": "succeed to Delete data from database",
	}).Info(res)

	return nil
}

// @title    	FindUserInfo
// @description   								通过用户名和密码在数据库中查找完整的信息
// @auth      	郑康           					2020.5.25
// @param     	string, string					用户名, 密码
// @return    	*dataBase.UserInfoTable, error	用户信息结构体指针, 错误信息
func FindUserInfo(username string, pwd string) (*dataBase.UserInfoTable, error) {
	sql := "SELECT * FROM im.userinfo \nWHERE username = '" + username + "' AND password = '" + pwd + "';"
	fmt.Println(sql)
	data := dataBase.ExecSelectSQL(sql)
	if data == nil || len(data) != 1 {
		return nil, errors.New("the data you select is nil or has repetitive")
	}
	return data[0], nil
}