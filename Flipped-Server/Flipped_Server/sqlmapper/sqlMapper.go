package sqlmapper

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	stringAttr = [...]string{"pid", "username", "password", "email", "realName", "profession", "region", "hobby"}
	intAttr = [...]string{"user_type", "age"}
	hexAttr = [1]string{"photo"}
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

func isHexAttr(str string) bool {
	return str == hexAttr[0]
}

func byteToInt(b []byte) int{
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)

	return int(x)
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
	fmt.Printf("该结构体有%d个字段\n", fieldNum)

	tagArr := make([]string, fieldNum)
	dataArr := make([]string, fieldNum)

	for i := 0; i < fieldNum; i++ {
		tagVal := dataType.Field(i).Tag.Get("sql")

		if tagVal == "" {
			return nil, nil, errors.New("tag is empty")
		}

		//fmt.Printf("Field %d: 值=%v, tag=%v\n", i, dataVal.Field(i), tagVal)

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
			buffer.WriteString("',\n")
		} else if isIntAttr(tagArr[i]) {
			integer, _ := strconv.Atoi(dataArr[i])
			buffer.WriteString(strconv.Itoa(integer))
			buffer.WriteString(",\n")
		} else if isHexAttr(tagArr[i]) {
			buffer.WriteString("x'")

			//ToDO: 将[]byte转为hex string

			buffer.WriteString("',\n")
		} else {
			return errors.New("unexpected value: " + tagArr[i])
		}
	}
	buffer.WriteString(");")
	sql := buffer.String()
	fmt.Println(sql)

	return nil
}


