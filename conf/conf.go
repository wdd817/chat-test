package conf

import (
	"chat-test/log"
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"time"
)

const confFile = "conf/conf.yaml"

const (
	GateModule       = "gate"
	ChatModule       = "chat"
	WordFilterModule = "word_filter"
)

var (
	LenStackBuf = 4096

	// log
	LogLevel log.LogLevel

	// gateway
	GatewayAddr            = ":8080"
	PendingWriteNum        = 2000
	MaxMsgLen       uint32 = 4096
	HTTPTimeout            = 10 * time.Second
	LenMsgLen              = 2
	MaxConnNum             = 100

	// skeleton conf
	AsynCallLen = 1000
	ChanRPCLen  = 1000

	// chat
	MaxHistory = 50
	PopularSecond = 5
)

func Init() {
	config.AddDriver(yaml.Driver)
	err := config.LoadFiles(confFile)
	if err != nil {
		log.Fatal("conf init err: %v", err)
	}

	log.Info("loadData: %v", config.Data())
	initLogLevel()

	gatePort := config.Int("gateway.port", 8080)
	GatewayAddr = fmt.Sprintf(":%d", gatePort)
}

func initLogLevel() {
	level := config.String("log.level", "info")
	if level == "debug" {
		LogLevel = log.DebugLevel
	} else if level == "info" {
		LogLevel = log.InfoLevel
	} else if level == "warn" {
		LogLevel = log.WarnLevel
	} else if level == "error" {
		LogLevel = log.ErrorLevel
	} else if level == "fatal" {
		LogLevel = log.FatalLevel
	} else {
		LogLevel = log.DebugLevel
	}
}
