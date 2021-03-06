package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type VideoInfo struct {
	Id     int    `json:"index"`
	Author string `json:"author"`
	Url    string `json:"url"`
}

// 配置信息数据结构
type Config struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Db       string `json:"db"`
}

var cfg = Config{}

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

func getUrl(c echo.Context) error {
	info := VideoInfo{}

	query(&info)

	// fmt.Println("id: ", info.Id, "\nauthor: ", info.Author, "\nurl: ", info.Url)

	jsonBytes, err := json.Marshal(&info)
	if err != nil {
		fmt.Println(err)
	}

	return c.String(http.StatusOK, string(jsonBytes))
}

func query(info *VideoInfo) {
	mysql_user := cfg.User
	mysql_passwd := cfg.Password
	mysql_ip := cfg.Ip
	mysql_port := cfg.Port
	mysql_qyDB := cfg.Db

	//打开数据库
	db, errOpen := sql.Open("mysql", mysql_user+":"+mysql_passwd+"@tcp("+mysql_ip+":"+mysql_port+")/"+mysql_qyDB+"?charset=utf8")
	if errOpen != nil {
		//TODO，这里只是打印了一下，并没有做异常处理
		fmt.Println("query Open is error")
	}

	var totals int
	err := db.QueryRow("select count(*) as totals from hvideos").Scan(&totals)
	if err != nil {
		fmt.Println("查询记录总数失败。")
	}
	fmt.Println("数据totals: ", totals)

	rand.Seed(time.Now().Unix())
	index := rand.Intn(totals) + 1

	var id int
	var url, author string
	errTables := db.QueryRow("SELECT id, url, author FROM hvideos WHERE id=?", index).Scan(&id, &url, &author)
	if errTables != nil {
		//TODO，这里只是打印了一下，并没有做异常处理
		fmt.Println("funReadSql SELECT t_knowledge_tree is error", errTables)
	}

	//关闭数据库
	db.Close()

	info.Id = id
	info.Author = author
	info.Url = url

	// fmt.Println("id: ", info.id, "\nauthor: ", author, "\nurl: ", url)
}

func main() {
	e := echo.New()
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })

	// conf := Config{}
	ReadConfig(&cfg)

	e.GET("/query", getUrl)

	e.Logger.Fatal(e.Start(":9279"))
}
