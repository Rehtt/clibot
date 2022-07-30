/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/7/30 上午 08:41
 */

package clibot

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"strings"
)

var help *CMD

func init() {
	help = &CMD{
		Use:         "help",
		Instruction: "显示所有命令及其帮助内容",
		Func: func(args []string, cmd *CMD, client *client.QQClient, msg *Msg) error {
			sendMsg := message.NewSendingMessage()
			var cmdHelp strings.Builder

			if msg.MsgType == MsgTypeGroup {
				sendMsg.Append(message.NewAt(msg.Sender.Uin))
			}

			for _, v := range root.cmds {
				if !HasPermission(msg.Sender.Uin, v.Permission) {
					continue
				}
				cmdHelp.WriteByte('\n')
				cmdHelp.WriteByte('/')
				cmdHelp.WriteString(v.Use)
				cmdHelp.WriteString(" ")
				cmdHelp.WriteString(v.Instruction)
			}
			sendMsg.Append(message.NewText(cmdHelp.String()))
			msg.Reply(sendMsg)
			return nil
		},
	}
	root.AddCommand(help)
}

//type help struct {
//}
//
//func (h *help) Key() string {
//	return "help"
//}
//
//func (h *help) Help() string {
//	return "输出所有命令列表，及其帮助"
//}
//
//func (h *help) Run(client *client.QQClient, msg *Msg) error {
//	sendMsg := message.NewSendingMessage()
//	var cmdHelp strings.Builder
//
//	if msg.MsgType == MsgTypeGroup {
//		sendMsg.Append(message.NewAt(msg.Sender.Uin))
//		cmdHelp.WriteByte('\n')
//	}
//
//	for k, v := range cmdMap {
//		cmdHelp.WriteByte('/')
//		cmdHelp.WriteString(k)
//		cmdHelp.WriteString(" ")
//		cmdHelp.WriteString(v.Help())
//		cmdHelp.WriteByte('\n')
//	}
//	sendMsg.Append(message.NewText(cmdHelp.String()))
//	msg.Reply(sendMsg)
//	return nil
//}
//
//var (
//	h = new(help)
//)
//
//func init() {
//	Register(h)
//}
