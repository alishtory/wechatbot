package config

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

// Configuration 项目配置
type Configuration struct {
	// gtp apikey
	ApiKey string `json:"api_key"`
	// 自动通过好友
	AutoPass bool `json:"auto_pass"`
	//代理地址
	Proxy string `json:"proxy"`
	// gpt 版本
	GptModel string `json:"gpt_model"`
	// 用户回复白名单
	ReplyUids []string `json:"reply_uids"`
	// 打招呼语句
	Greet string `json:"greet"`
	// 出错语句
	ErrorReply []string `json:"error_reply"`
}

var config *Configuration
var once sync.Once

// LoadConfig 加载配置
func LoadConfig() *Configuration {
	once.Do(func() {
		// 从文件中读取
		config = &Configuration{}
		f, err := os.Open("config.json")
		if err != nil {
			log.Fatalf("open config err: %v", err)
			return
		}
		defer f.Close()
		encoder := json.NewDecoder(f)
		err = encoder.Decode(config)
		if err != nil {
			log.Fatalf("decode config err: %v", err)
			return
		}

		// 如果环境变量有配置，读取环境变量
		ApiKey := os.Getenv("ApiKey")
		AutoPass := os.Getenv("AutoPass")
		GptModel := os.Getenv("GptModel")
		if ApiKey != "" {
			config.ApiKey = ApiKey
		}
		if AutoPass == "true" {
			config.AutoPass = true
		}
		if AutoPass != "" {
			config.GptModel = GptModel
		}
	})
	return config
}

func RandErrorReplay() string {
	rand.Seed(time.Now().UnixNano())
	errorReply := LoadConfig().ErrorReply
	return errorReply[rand.Intn(len(errorReply))]
}
