package handlers

import (
	"github.com/eatmoreapple/openwechat"
	"log"
	"wechatbot/config"
	"wechatbot/gtp"
)

// MessageHandlerInterface 消息处理接口
type MessageHandlerInterface interface {
	handle(*openwechat.Message) error
	ReplyText(*openwechat.Message) error
}

type HandlerType string

const (
	GroupHandler = "group"
	UserHandler  = "user"
)

// handlers 所有消息类型类型的处理器
var handlers map[HandlerType]MessageHandlerInterface

func init() {
	handlers = make(map[HandlerType]MessageHandlerInterface)
	handlers[GroupHandler] = NewGroupMessageHandler()
	handlers[UserHandler] = NewUserMessageHandler()
}

// Handler 全局处理入口
func Handler(msg *openwechat.Message) {
	log.Printf("hadler Received msg : %v", msg.Content)
	// 处理群消息
	if msg.IsSendByGroup() {
		handlers[GroupHandler].handle(msg)
		return
	}

	// 好友申请
	if msg.IsFriendAdd() {
		if config.LoadConfig().AutoPass {
			_, err := msg.Agree(config.LoadConfig().Greet)
			if err != nil {
				log.Fatalf("add friend agree error : %v", err)
				return
			}
			if gtp.NewIntentAnalyzer().SendGroupAddMsg(msg) {
				return
			}
		}
	}

	// 私聊
	handlers[UserHandler].handle(msg)
}
