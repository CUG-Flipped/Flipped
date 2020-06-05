package initialSetting

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var(
	configData map[string] interface {}
	DataBaseConfig map[string] interface {}
	LoggerConfig map[string] interface{}
)

func InitSettings(path string){
	var buf []byte
	var err error
	if path != "" {
		buf, err = ioutil.ReadFile(path)
	} else {
		buf, err = ioutil.ReadFile("./defaultSettings.json")
	}
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(buf, &configData)
	if err != nil {
		fmt.Println("Fail to unmarshal json file")
		panic(err)
	}
	DataBaseConfig = configData["dataBase"].(map[string] interface {})
	LoggerConfig = configData["logger"].(map[string] interface {})
}
