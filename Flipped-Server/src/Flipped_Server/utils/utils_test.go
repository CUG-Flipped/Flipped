package utils

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestGeneratorUUID(t *testing.T) {
	uuid := GeneratorUUID()
	fmt.Println(uuid)
}

func TestVerifyEmail(t *testing.T) {
	assert.Equal(t, true, VerifyEmail("906747215@qq.com"))
	assert.Equal(t, false, VerifyEmail("123456789"))
}

func TestContains(t *testing.T) {
	container := []string{"MrSecond", "MrFirst", "MrThird"}
	assert.Equal(t, true, Contains(container, "MrSecond"))
	assert.Equal(t, false, Contains(container, "1221"))
}

func TestImage2Base64(t *testing.T) {
	_, err := Image2Base64("C:\\Users\\13407\\Desktop\\Flipped.png")
	assert.Equal(t, nil, err)
}

func TestGetCurrentPath(t *testing.T) {
	path := GetCurrentPath()
	assert.Equal(t, "F:/Go_WorkSpace/Projects/src/Flipped_Server/utils", path)
}

func TestExists(t *testing.T) {
	assert.Equal(t, true, Exists("C:\\Users\\13407\\Desktop\\Flipped.png"))
}

func TestAesEncrypt(t *testing.T) {
	assert.Equal(t, "imLkEPRpDT8QHX+B2i5oFg==", AesEncrypt("mysql", "You Are Thieves "))
}

func TestAesDecrypt(t *testing.T) {
	assert.Equal(t, "mysql", AesDecrypt("imLkEPRpDT8QHX+B2i5oFg==", "You Are Thieves "))
}
