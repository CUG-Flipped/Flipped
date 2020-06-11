package network

import (
	"Flipped_Server/initialSetting"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var Exit chan int

func before() {
	Exit = make(chan int)
	initialSetting.InitSettings("F:\\Go_WorkSpace\\Projects\\src\\Flipped_Server\\defaultSettings.json")
	server := HttpServer{IPAddr: "127.0.0.1", Port: 8081}
	server.Run()
	<-Exit
}

func TestLogin(t *testing.T) {
	initialSetting.InitSettings("F:\\Go_WorkSpace\\Projects\\src\\Flipped_Server\\defaultSettings.json")
	server := HttpServer{IPAddr: "127.0.0.1", Port: 8081}
	router := server.SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8081/login?username=mrfirst&password=123456789", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	jsonData := make(map[string]interface{})
	var data = w.Body.Bytes()
	_ = json.Unmarshal(data, &jsonData)
	assert.Equal(t, "succeed to login", jsonData["message"].(string))
}

func TestGetFriendList(t *testing.T) {
	go before()
	req, _ := http.NewRequest("GET", "http://127.0.0.1:8081/friendList", nil)
	req.Header.Add("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1yZmlyc3QiLCJleHAiOjE1OTEzNDE4NDEsImlzcyI6Ik1yU2Vjb25kIn0.yZ1gDjcYGj_u7iM30uaf-b8ymqqnZHyUeApcXnKVu7o")
	var resp *http.Response
	resp, _ = http.DefaultClient.Do(req)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()
}

func TestAcquireRecommendUser(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://39.99.190.67:8081/recommendUser", nil)
	req.Header.Add("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1yZmlyc3QiLCJleHAiOjE1OTEzNDE4NDEsImlzcyI6Ik1yU2Vjb25kIn0.yZ1gDjcYGj_u7iM30uaf-b8ymqqnZHyUeApcXnKVu7o")
	var resp *http.Response
	resp, _ = http.DefaultClient.Do(req)
	res, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(res))
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()
}

func TestHeartBeat(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://39.99.190.67:8081/heartBeat", nil)
	req.Header.Add("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1yZmlyc3QiLCJleHAiOjE1OTEzNDE4NDEsImlzcyI6Ik1yU2Vjb25kIn0.yZ1gDjcYGj_u7iM30uaf-b8ymqqnZHyUeApcXnKVu7o")
	var resp *http.Response
	resp, _ = http.DefaultClient.Do(req)
	res, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(res))
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()
}

func TestGetOnlineUserNumber(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://39.99.190.67:8081/onlineUserNumber", nil)
	req.Header.Add("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1yZmlyc3QiLCJleHAiOjE1OTEzNDE4NDEsImlzcyI6Ik1yU2Vjb25kIn0.yZ1gDjcYGj_u7iM30uaf-b8ymqqnZHyUeApcXnKVu7o")
	var resp *http.Response
	resp, _ = http.DefaultClient.Do(req)
	res, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(res))
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()
}

func TestJudgeUserAlive(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://39.99.190.67:8081/isAlive", nil)
	req.Header.Add("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1yZmlyc3QiLCJleHAiOjE1OTEzNDE4NDEsImlzcyI6Ik1yU2Vjb25kIn0.yZ1gDjcYGj_u7iM30uaf-b8ymqqnZHyUeApcXnKVu7o")
	req.URL.Query().Add("username", "909")
	var resp *http.Response
	resp, _ = http.DefaultClient.Do(req)
	res, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(res))
	assert.Equal(t, 404, resp.StatusCode)
	defer resp.Body.Close()
}

func TestAddFriend(t *testing.T) {
	req, _ := http.NewRequest("POST", "http://39.99.190.67:8081/addFriend?friend=1", nil)
	req.Header.Add("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1yZmlyc3QiLCJleHAiOjE1OTE3NTI5NzUsImlzcyI6Ik1yU2Vjb25kIn0.aSgcN9K8zEFSZVcOO9uexZOlDbpxTOL9KaFc1ZcHI_A")
	var resp *http.Response
	resp, _ = http.DefaultClient.Do(req)
	res, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(res))
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()
}
