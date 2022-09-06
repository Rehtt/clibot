package clibot

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	MsgTypeGroup     = "Group"
	MsgTypePrivate   = "Private"
	MsgTypeGroupTemp = "GroupTemp"
)

type Msg struct {
	Id     int32
	Uin    int64
	Name   string
	Sender *Sender
	Msg    string

	MsgType string
	client  *client.QQClient
	//groupMsg   *message.GroupMessage
	//privateMsg *message.PrivateMessage
	//tempMsg    *message.TempMessage
	Original any
}

func (m *Msg) SendGroupMsg(groupCode int64, msg *message.SendingMessage) {
	m.client.SendGroupMessage(groupCode, msg)
}
func (m *Msg) SendPrivateMsg(target int64, msg *message.SendingMessage) {
	m.client.SendPrivateMessage(target, msg)
}
func (m *Msg) SendGoupTempMsg(groupCode, target int64, msg *message.SendingMessage) {
	m.client.SendGroupTempMessage(groupCode, target, msg)
}
func (m *Msg) Reply(msg *message.SendingMessage) {
	switch m.MsgType {
	case MsgTypeGroupTemp:
		m.SendGoupTempMsg(m.Uin, m.Sender.Uin, msg)
	case MsgTypeGroup:
		m.SendGroupMsg(m.Uin, msg)
	case MsgTypePrivate:
		m.SendPrivateMsg(m.Uin, msg)
	}
}

type Sender struct {
	Uin  int64
	Name string
	Role string
}

var (
	CliRoot = new(CMD)
)

type CMD struct {
	Use            string
	Instruction    string
	Func           func(data []string, cmd *CMD, client *client.QQClient, msg *Msg) error
	Permission     []string
	Ignore         bool // 忽略指令，时刻处于激活状态，只适用于头指令
	ActivationFunc func(cmd *CMD, qqClient *client.QQClient, msg *Msg) error
	OnClose        func()
	floor          int
	cmds           []*CMD
}

func (c *CMD) AddCommand(cs ...*CMD) {
	c.cmds = append(c.cmds, cs...)
}
func (c *CMD) Help() string {
	return c.Instruction
}
func (c *CMD) findCMD(data []string, uin int64) *CMD {
	if len(c.cmds) == 0 || len(data) == 0 {
		return c
	}
	for _, v := range c.cmds {
		if v.Use == data[0] {
			if HasPermission(uin, v.Permission) {
				v.floor = c.floor + 1
				return v.findCMD(data[1:], uin)
			}
			return nil
		}
	}
	if c != CliRoot {
		return c
	}
	return nil
}

func (m *Msg) parseCMD() {

	switch m.MsgType {
	case MsgTypeGroupTemp:
		tempMsg := m.Original.(*message.TempMessage)
		m.Msg = tempMsg.ToString()
		m.Sender = &Sender{
			Uin:  tempMsg.Sender.Uin,
			Name: tempMsg.Sender.DisplayName(),
		}
		m.Uin = tempMsg.GroupCode
		m.Name = tempMsg.GroupName

	case MsgTypePrivate:
		privateMsg := m.Original.(*message.PrivateMessage)
		m.Msg = privateMsg.ToString()
		m.Sender = &Sender{
			Uin:  privateMsg.Sender.Uin,
			Name: privateMsg.Sender.DisplayName(),
		}
		m.Uin = privateMsg.Sender.Uin
		m.Name = privateMsg.Sender.DisplayName()
	case MsgTypeGroup:
		groupMsg := m.Original.(*message.GroupMessage)
		m.Msg = groupMsg.ToString()
		m.Sender = &Sender{
			Uin:  groupMsg.Sender.Uin,
			Name: groupMsg.Sender.DisplayName(),
		}
		m.Uin = groupMsg.GroupCode
		m.Name = groupMsg.GroupName
	}

	c := strings.Split(m.Msg, " ")
	if len(c) < 1 || strings.Index(c[0], "/") != 0 {
		for _, c := range CliRoot.cmds {
			if c.Ignore {
				c.ActivationFunc(c, m.client, m)
			}
		}
		return
	}
	arg := make([]string, 0, len(c))
	c[0] = c[0][1:]

	for i := range c {
		if c[i] == "" {
			continue
		}
		arg = append(arg, c[i])
	}

	var err error
	defer func() {
		if err != nil {
			logger.Errorf("%s: %s", arg, err.Error())
		}
	}()

	cmdF := CliRoot.findCMD(arg, m.Sender.Uin)
	if cmdF == nil {
		CliRoot.Help()
		return
	}
	if cmdF.Func != nil {
		err = cmdF.Func(arg[cmdF.floor:], cmdF, m.client, m)
	}

}

func (m *Msg) Logger() logrus.FieldLogger {
	return logger
}

func onClose(c *CMD) {
	for _, c := range c.cmds {
		onClose(c)
	}
	if c.OnClose != nil {
		c.OnClose()
	}
}
