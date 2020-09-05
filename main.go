package main

import (
	"chat-test/app"
	"chat-test/chat"
	"chat-test/conf"
	"chat-test/gate"
	"chat-test/log"
	"chat-test/wordfilter"
)

func main() {

	log.Info("launching ...")
	conf.Init()

	gateModule := gate.NewModule()
	wordFilterModule := wordfilter.NewModule()
	chatModule := chat.NewModule()
	app.Instance().Run(gateModule, wordFilterModule, chatModule)
}
