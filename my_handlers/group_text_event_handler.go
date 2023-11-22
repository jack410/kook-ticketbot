package my_handlers

import (
	"errors"
	"fmt"
	kook_CardBuild "github.com/Quinlivanner/kook-CardBuild"
	"github.com/bytedance/sonic"
	"github.com/gookit/event"
	"github.com/kaiheila/golang-bot/api/base"
	event2 "github.com/kaiheila/golang-bot/api/base/event"
	"github.com/kaiheila/golang-bot/api/helper"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

const prefix string = "!bot"

type GroupTextEventHandler struct {
	Token   string
	BaseUrl string
}

type Answers struct {
	TargetChannelId string //ID of the channel that the message was sent back to(new group channel or origin channel)
	AuthorName      string
	ReportContent   string // content of 1st question
	ReportSecond    string // content of 2nd question
}

var responses map[string]Answers = map[string]Answers{}

func (gteh *GroupTextEventHandler) Handle(e event.Event) error {
	log.WithField("event", fmt.Sprintf("%+v", e.Data())).Info("收到频道内的文字消息.")
	var err error
	err = func() error {
		log.Infof("匿名函数开始————————————————————————————")
		if _, ok := e.Data()[base.EventDataFrameKey]; !ok {
			return errors.New("data has no frame field")
		}
		frame := e.Data()[base.EventDataFrameKey].(*event2.FrameMap)
		data, err := sonic.Marshal(frame.Data)
		if err != nil {
			return err
		}
		msgEvent := &event2.MessageKMarkdownEvent{}
		err = sonic.Unmarshal(data, msgEvent)
		log.Infof("Received json event:%+v", msgEvent)
		if err != nil {
			return err
		}
		client := helper.NewApiHelper("/v3/message/create", gteh.Token, gteh.BaseUrl, "", "")
		if msgEvent.Author.Bot {
			log.Info("bot message")
			return nil
		}

		//dm logic
		log.Infof("guilid为：%+v", msgEvent.GuildID)
		log.Infof("codeid为：%+v", msgEvent.Code)
		// guildid是空字符串，则表示是direct message
		if msgEvent.Code != "" {
			//每个direct message 都有一个唯一的chat code
			answers, ok := responses[msgEvent.Code]
			if !ok {
				return nil
			}

			if answers.ReportContent == "" {
				answers.ReportContent = msgEvent.Content
				gteh.DirectMessageSend(msgEvent.Code, "#2 请提交下截图或者视频信息，如果没有请回复：没有图像信息。")

				responses[msgEvent.Code] = answers
				return nil
			} else {
				answers.ReportSecond = msgEvent.Content

				currentTime := time.Now().Format("2006-01-02 15:04:05")

				log.Infof("%v answers: %v, %v", answers.AuthorName, answers.ReportContent, answers.ReportSecond)

				userReportCards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
				if err != nil {
					return err
				}

				userReportCards.AddKmarkdown("**提交玩家:** " + answers.AuthorName)
				userReportCards.AddDivider()
				userReportCards.AddKmarkdown("**提交时间:** " + currentTime)
				userReportCards.AddKmarkdown("**提交内容:** " + answers.ReportContent)

				//判断第二个问题用户提交的是视频还是图片还是文字信息
				if isImageLink(answers.ReportSecond) {
					log.Infof("用户输入的是图片链接：%v", answers.ReportSecond)
					userReportCards.AddImage(answers.ReportSecond)
				} else if isVideoLink(answers.ReportSecond) {
					log.Infof("用户输入的是视频链接：%v", answers.ReportSecond)
					//linkUrl输出[url](url),在markdown中显示为link
					userReportCards.AddKmarkdown("**用户提交的视频信息:** " + linkUrl(answers.ReportSecond))
				} else {
					log.Infof("用户输入的是文字信息：%v", answers.ReportSecond)
					userReportCards.AddKmarkdown("**备注:** " + answers.ReportSecond)
				}

				userReportCardsContent, err := kook_CardBuild.GenerateCardMessageContent(userReportCards)
				if err != nil {
					return err
				}

				SendGroupCardessage(userReportCardsContent, answers.TargetChannelId, client)
				//卡片信息发送至对应频道后，回复玩家提交成功信息。
				gteh.DirectMessageSend(msgEvent.Code, "举报信息提交成功。我们会进行跟踪处理。感谢您的反馈！")

				delete(responses, msgEvent.Code)
			}

		}

		// server logic
		args := strings.Split(msgEvent.KMarkdown.RawContent, " ")

		if args[0] != prefix {
			return nil
		}

		if len(args) == 1 && msgEvent.GuildID != "" {
			SendGroupTextMessage(msgEvent.TargetId, "请输入!bot help查看本机器人支持哪些命令。", client)

			return nil
		}

		if args[1] == "help" && msgEvent.GuildID != "" {
			helpCards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			helpCards.AddHeader("本机器人支持的命令如下：")
			helpCards.AddDivider()
			helpCards.AddKmarkdown("(font)!bot help(font)[purple]    " + "查看机器人命令")
			helpCards.AddKmarkdown("(font)!bot report或者!bot 举报(font)[purple]    " + "触发机器人私信，并通过与机器人私信提交举报信息")
			helpCards.AddKmarkdown("(font)!bot npc(font)[purple]    " + "随机发送魔兽世界npc经典台词")

			helpCardsContent, err := kook_CardBuild.GenerateCardMessageContent(helpCards)
			if err != nil {
				return err
			}

			SendGroupCardessage(helpCardsContent, msgEvent.TargetId, client)

		}

		if args[1] == "npc" && msgEvent.GuildID != "" {
			npcquote, err := getRandomNpcquote()
			if err != nil {
				return err
			}

			cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			//cards.AddColorText(kook_CardBuild.TextColorPink, "@"+msgEvent.Author.Username+" "+npc[selection])
			cards.AddKmarkdown("(met)" + msgEvent.Author.ID + "(met)" + "  " + npcquote)

			cardMessageContent, err := kook_CardBuild.GenerateCardMessageContent(cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(cardMessageContent, msgEvent.TargetId, client)
		}

		if args[1] == "report" || args[1] == "举报" && msgEvent.GuildID != "" {
			gteh.UserReportHandler(msgEvent.Author.ID)
		}
		return nil
	}()
	if err != nil {
		log.WithError(err).Error("GroupTextEventHandler err")
	}
	return nil
}

func (gteh *GroupTextEventHandler) UserReportHandler(targetId string) {
	// user channel
	chatInfo := gteh.UserChatCreate(targetId)
	if chatInfo.Code != 0 {
		log.Infof("请求出错", chatInfo.Message)
	} else {
		log.Println("创建私聊成功，chat_cod为：", chatInfo.Data.Code)
	}
	//if the user is already answer question, ignore it, otherwise ask question
	//true means answer exists in the map
	if _, ok := responses[chatInfo.Data.Code]; !ok {
		responses[chatInfo.Data.Code] = Answers{
			TargetChannelId: os.Getenv("ReportChannelId"),
			AuthorName:      chatInfo.Data.TargetInfo.Username,
			ReportContent:   "",
			ReportSecond:    "",
		}
		gteh.DirectMessageSend(chatInfo.Data.Code, "您好！欢迎使用举报机器人。请按照接下来的两个步骤来完成举报!")
		gteh.DirectMessageSend(chatInfo.Data.Code, "#1：请描述您要举报的内容:ID、时间、地点、发生事件。 如：玩家Abc于2023/11/19 20:30左右在灼热峡谷跟随双开。")

	} else { // if answer exists
		gteh.DirectMessageSend(chatInfo.Data.Code, "我们还在等待等待您回复之前的问题。")
	}
}

func isVideoLink(input string) bool {
	// 判断是否以 "https" 或以指定视频文件扩展名开头
	videoExtensions := []string{".mp4", ".avi", ".mov"}
	for _, ext := range videoExtensions {
		if strings.HasPrefix(input, "[https") || strings.HasPrefix(input, "https") || strings.HasPrefix(input, ext) {
			return true
		}
	}
	return false
}

func isImageLink(input string) bool {
	// 使用正则表达式判断是否匹配图片链接的模式
	imagePattern := `^https?://.*\.(?:png|jpg|jpeg|gif)$`
	matched, _ := regexp.MatchString(imagePattern, input)
	return matched
}

func linkUrl(input string) string {
	if strings.HasPrefix(input, "[") {
		return input
	}
	return "[" + input + "](" + input + ")"
}
func getRandomNpcquote() (string, error) {
	content, err := ioutil.ReadFile("npc_quotes.txt")
	if err != nil {
		return "", err
	}

	proverbs := strings.Split(string(content), "\n")

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(proverbs))

	return proverbs[randomIndex], nil
}
