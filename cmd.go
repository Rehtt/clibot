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
	Uin    int64
	Name   string
	Sender *Sender
	Msg    string

	MsgType    string
	client     *client.QQClient
	groupMsg   *message.GroupMessage
	privateMsg *message.PrivateMessage
	tempMsg    *message.TempMessage
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
	return nil
}

func (m *Msg) parseCMD() {

	switch m.MsgType {
	case MsgTypeGroupTemp:
		m.Msg = m.tempMsg.ToString()
		m.Sender = &Sender{
			Uin:  m.tempMsg.Sender.Uin,
			Name: m.tempMsg.Sender.DisplayName(),
		}
		m.Uin = m.tempMsg.GroupCode
		m.Name = m.tempMsg.GroupName

	case MsgTypePrivate:
		m.Msg = m.privateMsg.ToString()
		m.Sender = &Sender{
			Uin:  m.privateMsg.Sender.Uin,
			Name: m.privateMsg.Sender.DisplayName(),
		}
		m.Uin = m.privateMsg.Sender.Uin
		m.Name = m.privateMsg.Sender.DisplayName()
	case MsgTypeGroup:
		m.Msg = m.groupMsg.ToString()
		m.Sender = &Sender{
			Uin:  m.groupMsg.Sender.Uin,
			Name: m.groupMsg.Sender.DisplayName(),
		}
		m.Uin = m.groupMsg.GroupCode
		m.Name = m.groupMsg.GroupName
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
	c[0] = c[0][1:]

	var err error
	defer func() {
		if err != nil {
			logger.Errorf("%s: %s", c, err.Error())
		}
	}()

	cmdF := CliRoot.findCMD(c, m.Sender.Uin)
	if cmdF == nil {
		CliRoot.Help()
		return
	}
	if cmdF.Func != nil {
		err = cmdF.Func(c[cmdF.floor:], cmdF, m.client, m)
	}

}

func (m *Msg) Logger() logrus.FieldLogger {
	return logger
}
