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

var (
	stringAttr = [...]string{"pid", "username", "photo","password", "email", "realName", "profession", "region", "hobby"}
	intAttr = [...]string{"user_type", "age"}
)

func isStringAttr(str string) bool {
	for i := 0; i < len(stringAttr); i++ {
		if str == stringAttr[i] {
			return true
		}
	}
	return false
}

func isIntAttr(str string) bool {
	for i := 0; i < len(intAttr); i++ {
		if str == intAttr[i] {
			return true
		}
	}
	return false
}


type SqlMapper interface {
	Insert(data interface {}) error
	Update(data interface {}) error
	Select(data interface {}) error
	Delete(data interface {}) error
}


func splitDataAndStruct(metaData interface{}) ([]string, []string, error){
	dataType := reflect.TypeOf(metaData)
	dataVal := reflect.ValueOf(metaData)
	dataKind := dataVal.Kind()
	if dataKind != reflect.Struct {
		return nil, nil, errors.New("Expect Struct")
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
	fmt.Println(sql)

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


