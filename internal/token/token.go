package token

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"ffmpeg_hls_go/internal/configs"
	"ffmpeg_hls_go/internal/logger"
	"ffmpeg_hls_go/pkg/sm3"

	"github.com/sirupsen/logrus"
)

type Response struct {
	Code         string `json:"code"`
	Msg          string `json:"msg"`
	Datetime     string `json:"datetime"`
	CompanyCode  string `json:"company_code"`
	AccessToken  string `json:"access_token"`
	PublicKey    string `json:"public_key"`
	PkExpiryTime string `json:"pk_expiry_time"`
}

var (
	response Response
)

func GetResponse() Response {
	return response
}
func Init() {
	//启动一个线程 专门获取 token
	// 并每 10分钟重新取一次token
	go func() {
		config := configs.GetConfigInstance()
		log := logger.GetLogger("token.log", true)
		xcount := 0
		for {
			log.Info("#####xcount:", xcount)
			xcount += 1
			fetchToken(config, log)
			time.Sleep(time.Minute * 1)
		}
	}()
}
func fetchToken(config *configs.Config, log *logrus.Logger) {
	url := config.TokenServer.TokenApiUrl

	payload := map[string]interface{}{
		"username":         config.TokenServer.UserName,
		"password":         config.TokenServer.Password,
		"requestTimestamp": "",
		"encrypt_method":   "1",
	}

	// Password encryption using SM3, encryption method:
	// SM3(username+requestTime+allocPassword), generate a 32-bit hexadecimal string.
	payload["requestTimestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	//payload["requestTimestamp"] = "1729687611"
	xStr := fmt.Sprintf("%s%s%s", payload["username"].(string), payload["requestTimestamp"].(string), payload["password"].(string))

	dataInBytes := []byte(xStr)
	s3hash := sm3.Sm3Sum(dataInBytes)
	payload["password"] = hex.EncodeToString(s3hash)
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println("json.Marshal( Error:", err)
		return
	}
	log.Info("token payload: ", string(jsonPayload))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Println("http.Post Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ioutil.ReadAll Error:", err)
		return
	}

	log.Info("token body: ", string(body))
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("json.Unmarshal( Error:", err)
		return
	}
	log.Info("token response: ", response)
	log.Info("x", response.AccessToken)
}
