package main

import (
	"github.com/joho/godotenv"
	"github.com/kaiheila/golang-bot/api/base"
	"github.com/kaiheila/golang-bot/example/handler"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"kookbot01/my_handlers"
	"os"
	"time"
)

func main() {
	godotenv.Load()

	log.SetReportCaller(true)

	//log.SetFormatter(&log.TextFormatter{})
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	session := base.NewWebSocketSession(os.Getenv("Token"), os.Getenv("BaseUrl"), "./session.pid", "", 1)

	session.On(base.EventReceiveFrame, &handler.ReceiveFrameHandler{})
	//session.On("GROUP*", &handler.GroupEventHandler{})
	//session.On(base.EventReceiveFrame, &my_handlers.PersonTextEventHandler{os.Getenv("Token"), os.Getenv("BaseUrl")})
	//GROUP*代表监听群里的所有信息，GROUP_9代表监听文字信息
	session.On(base.EventReceiveFrame, &my_handlers.GroupTextEventHandler{os.Getenv("Token"), os.Getenv("BaseUrl")})
	session.On(base.EventReceiveFrame, &my_handlers.FaqEventHandler{os.Getenv("Token"), os.Getenv("BaseUrl")})
	session.On(base.EventReceiveFrame, &my_handlers.MessageDelHandler{os.Getenv("Token"), os.Getenv("BaseUrl")})
	session.Start()
}

func init() {

	path := "message.log"

	/* 日志轮转相关函数

	   `WithLinkName` 为最新的日志建立软连接

	   `WithRotationTime` 设置日志分割的时间，隔多久分割一次

	   WithMaxAge 和 WithRotationCount二者只能设置一个

	    `WithMaxAge` 设置文件清理前的最长保存时间

	    `WithRotationCount` 设置文件清理前最多保存的个数

	*/

	// 下面配置日志每隔 1 分钟轮转一个新文件，保留最近 3 分钟的日志文件，多余的自动清理掉。

	writer, _ := rotatelogs.New(

		path+".%Y%m%d%H",

		rotatelogs.WithLinkName(path),

		rotatelogs.WithMaxAge(time.Duration(8760)*time.Hour),

		rotatelogs.WithRotationTime(time.Duration(60)*time.Minute),
	)

	log.SetOutput(writer)

	//log.SetFormatter(&log.JSONFormatter{})

}
