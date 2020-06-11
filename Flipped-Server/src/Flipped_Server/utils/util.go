package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"github.com/gofrs/uuid"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var ExitFlag chan bool

func GeneratorUUID() string {
	u2, err := uuid.NewV4()
	if err != nil {
		return err.Error()
	}
	return u2.String()
}

func GetImageURL(imageName string, imageData []byte) string {
	imageDir := filepath.Dir(GetCurrentPath()) + "/imageContainer/"
	if !Exists(imageDir) {
		panic(errors.New("image Dir does't exist"))
	}
	imageName = imageDir + imageName
	imageName = strings.Replace(imageName, "\\", "/", -1)
	f, err := os.OpenFile(imageName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	_, _ = f.Write(imageData)
	defer f.Close()

	return imageName
}

func GetCurrentPath() string {
	_, filename, _, ok := runtime.Caller(1)
	var cwdPath string
	if ok {
		cwdPath = path.Join(path.Dir(filename), "")
	} else {
		cwdPath = "./"
	}
	return cwdPath
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func Contains(container []string, item string) bool {
	if len(container) == 0 {
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

func Image2Base64(imagePath string) (string, error) {
	if !Exists(imagePath) {
		return "", errors.New("image doesn't exist")
	}
	_, imageName := filepath.Split(imagePath)
	imageType := path.Ext(imageName)[1:]
	imageBytes, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return "", err
	}
	imageStr := "data:image/" + imageType + ";base64," + base64.StdEncoding.EncodeToString(imageBytes)
	return imageStr, nil
}

func VerifyEmail(email string) bool {
	if email == "" || len(email) == 0 {
		return false
	} else {
		pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
		reg := regexp.MustCompile(pattern)
		return reg.MatchString(email)
	}
}

func AesEncrypt(orig string, key string) string {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)

	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		panic(err)
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)

	return base64.StdEncoding.EncodeToString(cryted)

}

func AesDecrypt(cryted string, key string) string {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := []byte(key)

	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		panic(err)
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	return string(orig)
}

//补码
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//去码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
