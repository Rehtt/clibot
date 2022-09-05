/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/7/30 上午 09:11
 */

package clibot

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"strconv"
)

func init() {
	admin := &CMD{
		Use:         "admin",
		Instruction: "管理员配置",
		Func: func(args []string, cmd *CMD, client *client.QQClient, msg *Msg) error {
			return nil
		},
		Permission: []string{"admin"},
	}
	admin.AddCommand(
		&CMD{
			Use:         "is",
			Instruction: "判断是否是管理员",
			Func: func(data []string, cmd *CMD, client *client.QQClient, msg *Msg) (err error) {
				uin := msg.Sender.Uin
				if len(data) != 0 {
					uin, err = strconv.ParseInt(data[0], 10, 64)
					if err != nil {
						return err
					}
				}
				sendMsg := message.NewSendingMessage()
				if GetRole(uin) == RoleAdmin {
					sendMsg.Append(message.NewText("是管理员"))
				} else {
					sendMsg.Append(message.NewText("不是管理员"))
				}
				msg.Reply(sendMsg)
				return
			},
		},
		&CMD{
			Use:         "set",
			Instruction: "设置管理员",
			Func: func(data []string, cmd *CMD, client *client.QQClient, msg *Msg) error {
				if len(data) == 0 {
					return fmt.Errorf("设置失败")
				}
				uin, err := strconv.ParseInt(data[0], 10, 64)
				if err != nil {
					return err
				}
				SetRole(uin, RoleAdmin)
				msg.Reply(message.NewSendingMessage().Append(message.NewText("设置成功")))
				return nil
			},
		},
	)
	CliRoot.AddCommand(admin)
}
