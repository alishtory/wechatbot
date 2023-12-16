package gtp

import (
	"github.com/eatmoreapple/openwechat"
	"log"
	"regexp"
	"sync"
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
			pattern: regexp.MustCompile(`[加进入].*群|新西兰|交流`),
		}
	})
	return instance
}

// HasGroupIntent 检查用户输入是否包含加群意愿
func (analyzer *IntentAnalyzer) HasGroupIntent(userInput string) bool {
	// 不需要再次初始化 pattern
	return analyzer.pattern.MatchString(userInput)
}

func (analyzer *IntentAnalyzer) SendGroupAddMsg(msg *openwechat.Message) bool {
	content := msg.Content
	if msg.IsFriendAdd() {
		content = msg.RecommendInfo.Content
	}

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
	searchGroups := groups.SearchByUserName(1, "@@64ba28a6970a9a72e86658178a60ecfceaf40ca3b76659642ad4103f3a76e382")
	if g := searchGroups.First(); g != nil {
		if msg.IsFriendAdd() {
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
			friend, _ := sender.AsFriend()
			if err := g.AddFriendsIn(friend); err != nil {
				log.Println(err)
			} else {
				return true
			}
		}

	}
	return false
}
