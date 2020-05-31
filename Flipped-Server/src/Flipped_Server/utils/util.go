package utils

import (
	"Flipped_Server/logger"
	"encoding/base64"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func GeneratorUUID() string{
	u2, err := uuid.NewV4()
	if err != nil {
		return  err.Error()
	}
	return u2.String()
}

func GetImageURL(imageName string, imageData []byte) string{
	imageDir := filepath.Dir(GetCurrentPath()) + "/imageContainer/"
	if !Exists(imageDir){
		panic(errors.New("image Dir does't exist"))
	}
	imageName = imageDir + imageName
	imageName = strings.Replace(imageName, "\\", "/", -1)
	f, err := os.OpenFile(imageName, os.O_WRONLY | os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	_, _ = f.Write(imageData)
	defer f.Close()

	return imageName
}

func GetCurrentPath() string{
	_, filename, _,ok := runtime.Caller(1)
	var cwdPath string
	if ok {
		cwdPath = path.Join(path.Dir(filename), "")
	} else {
		cwdPath = "./"
	}
	return cwdPath
}

func Exists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func Contains(container []string, item string) bool {
	if container == nil || len(container) == 0 {
		return false
	} else {
		for i := range container {
			if container[i] == item {
				return true
			}
		}
		return false
	}
}

func Image2Base64(imagePath string) (string, error){
	if !Exists(imagePath){
		logger.SetToLogger(logrus.ErrorLevel, "Image2Base64", "image doesn't exist", "")
		return "", errors.New("image doesn't exist")
	}
	_, imageName := filepath.Split(imagePath)
	imageType := path.Ext(imageName)[1:]
	imageBytes, err := ioutil.ReadFile(imagePath)
	if err != nil {
		logger.SetToLogger(logrus.ErrorLevel, "Image2Base64", "read from path: " + imagePath, err.Error())
		return "", err
	}
	imageStr := "data:image/"+imageType+";base64," + base64.StdEncoding.EncodeToString(imageBytes)
	return imageStr, nil
}