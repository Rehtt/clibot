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
			var cmds = CliRoot.cmds
			if len(args) != 0 {
				cmds = CliRoot.findCMD(args, msg.Sender.Uin).cmds
			}
			for _, v := range cmds {
				if !HasPermission(msg.Sender.Uin, v.Permission) {
					continue
				}
				cmdHelp.WriteByte('\n')
				if len(args) == 0 {
					cmdHelp.WriteByte('/')
				}
				cmdHelp.WriteString(v.Use)
				cmdHelp.WriteString(" ")
				cmdHelp.WriteString(v.Instruction)
			}
			sendMsg.Append(message.NewText(cmdHelp.String()))
			msg.Reply(sendMsg)
			return nil
		},
	}
	CliRoot.AddCommand(help)
}
