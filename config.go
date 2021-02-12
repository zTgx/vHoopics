package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// 配置信息数据结构
type Config struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Db       string `json:"db"`
}

// 读取配置信息
func ReadConfig(conf *Config) {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	// conf := Config{}
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("Error:", err)
	}
	// fmt.Println("数据库配置信息： ", conf.User)
}
