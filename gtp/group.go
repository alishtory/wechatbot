package gtp

import (
	"github.com/eatmoreapple/openwechat"
	"log"
	"regexp"
	"sync"
	"wechatbot/config"
)

// IntentAnalyzer 是用户意图分析的工具类
type IntentAnalyzer struct {
	pattern *regexp.Regexp
	once    sync.Once
}

var (
	instance *IntentAnalyzer
	once     sync.Once
)

// NewIntentAnalyzer 返回 IntentAnalyzer 的单例实例
func NewIntentAnalyzer() *IntentAnalyzer {
	once.Do(func() {
		instance = &IntentAnalyzer{
			pattern: regexp.MustCompile(config.LoadConfig().GroupIntentPattern),
		}
	})
	return instance
}

// HasGroupIntent 检查用户输入是否包含加群意愿
func (analyzer *IntentAnalyzer) HasGroupIntent(userInput string) bool {
	// 不需要再次初始化 pattern
	return analyzer.pattern.MatchString(userInput)
}

func (analyzer *IntentAnalyzer) FriendAddSendGroupAddMsg(msg *openwechat.Message) bool {
	self, err := msg.Bot().GetCurrentUser()
	if err != nil {
		log.Fatalf("GetCurrentUser error : %v", err)
		return false
	}
	groups, err := self.Groups()
	if err != nil {
		log.Fatalf("get Groups error : %v", err)
		return false
	}
	searchGroups := groups.SearchByID(config.LoadConfig().GroupId)
	if g := searchGroups.First(); g != nil {
		friends, err := self.Friends()
		if err != nil {
			log.Fatalf("get Friends error : %v", err)
			return false
		}
		friends = friends.SearchByUserName(1, msg.RecommendInfo.UserName)
		if err := g.AddFriendsIn(friends...); err != nil {
			log.Println(err)
		} else {
			return true
		}
	} else {
		log.Fatalf("can not find group: %v", config.LoadConfig().GroupId)
	}
	return false
}

func (analyzer *IntentAnalyzer) SendGroupAddMsg(msg *openwechat.Message) bool {
	content := msg.Content
	if !analyzer.HasGroupIntent(content) {
		return false
	}
	sender, err := msg.Sender()
	if err != nil {
		log.Fatalf("get Sender error : %v", err)
		return false
	}
	self, err := msg.Bot().GetCurrentUser()
	if err != nil {
		log.Fatalf("GetCurrentUser error : %v", err)
		return false
	}
	groups, err := self.Groups()
	if err != nil {
		log.Fatalf("get Groups error : %v", err)
		return false
	}
	searchGroups := groups.SearchByID(config.LoadConfig().GroupId)
	if g := searchGroups.First(); g != nil {
		friend, _ := sender.AsFriend()
		if err := g.AddFriendsIn(friend); err != nil {
			log.Println(err)
		} else {
			return true
		}
	} else {
		log.Fatalf("can not find group: %v", config.LoadConfig().GroupId)
	}
	return false
}
