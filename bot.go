package clibot

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"sync"
)

type cmdBot struct{}

const (
	name = "rehtt/cli"
)

var (
	c      = cmdBot{}
	logger = utils.GetModuleLogger(name)
	pool   = sync.Pool{
		New: func() interface{} {
			return new(Msg)
		},
	}
)

func (c *cmdBot) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       name,
		Instance: c,
	}
}

func (c *cmdBot) Init() {

}

func (c *cmdBot) PostInit() {

}

func (c *cmdBot) Serve(bot *bot.Bot) {
	bot.GroupMessageEvent.Subscribe(func(client *client.QQClient, message *message.GroupMessage) {
		msg := pool.Get().(*Msg)
		defer pool.Put(msg)
		msg.client = client
		msg.groupMsg = message
		msg.MsgType = MsgTypeGroup
		msg.parseCMD()
	})
	bot.PrivateMessageEvent.Subscribe(func(qqClient *client.QQClient, privateMessage *message.PrivateMessage) {
		msg := pool.Get().(*Msg)
		defer pool.Put(msg)
		msg.client = qqClient
		msg.privateMsg = privateMessage
		msg.MsgType = MsgTypePrivate
		msg.parseCMD()
	})
	bot.TempMessageEvent.Subscribe(func(qqClient *client.QQClient, event *client.TempMessageEvent) {
		msg := pool.Get().(*Msg)
		defer pool.Put(msg)
		msg.client = qqClient
		msg.tempMsg = event.Message
		msg.MsgType = MsgTypeGroupTemp
		msg.parseCMD()
	})
}

func (c *cmdBot) Start(bot *bot.Bot) {

}

func (c *cmdBot) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	wg.Done()
}

func init() {
	bot.RegisterModule(&c)
}
