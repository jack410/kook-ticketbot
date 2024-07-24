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
	"strings"
)

type FaqEventHandler struct {
	Token   string
	BaseUrl string
}

//func (gteh *FaqEventHandler) Handle(e event.Event) error {
//	fmt.Println("----------------faq handler--------------")
//	log.WithField("event", fmt.Sprintf("%+v", e.Data())).Info("收到频道内的文字消息.")
//	err := func() error {
//		if _, ok := e.Data()[base.EventDataFrameKey]; !ok {
//			return errors.New("data has no frame field")
//		}
//		frame := e.Data()[base.EventDataFrameKey].(*event2.FrameMap)
//		data, err := sonic.Marshal(frame.Data)
//		if err != nil {
//			return err
//		}
//		msgEvent := &event2.MessageKMarkdownEvent{}
//		err = sonic.Unmarshal(data, msgEvent)
//		log.Infof("Received json event:%+v", msgEvent)
//		if err != nil {
//			return err
//		}
//		client := helper.NewApiHelper("/v3/message/create", gteh.Token, gteh.BaseUrl, "", "")
//		if msgEvent.Author.Bot {
//			log.Info("bot message")
//			return nil
//		}
//
//		if containsWord(msgEvent.KMarkdown.RawContent, "英文") {
//			fmt.Println("-------------------------------------------")
//			cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
//			if err != nil {
//				return err
//			}
//			cards.AddKmarkdown("(met)" + msgEvent.Author.ID + "(met)" + "  " + "客户端英文改成中文方法：找到\\winterspring-data\\WoW 1.14.2 Everlook Asia\\_classic_era_\\WTF文件夹" +
//				"里面的Config.wtf，鼠标右键该文件，打开方式选择记事本，将第2第3行的“enUS”那改为 \"zhCN\" ，保存退出，重登游戏")
//
//			cardMessageContent, err := kook_CardBuild.GenerateCardMessageContent(cards)
//			if err != nil {
//				return err
//			}
//
//			SendGroupCardessage(cardMessageContent, msgEvent.TargetId, client)
//			fmt.Println("-------------------------------------------")
//		}
//
//		return nil
//	}()
//	if err != nil {
//		log.WithError(err).Error("FaqEventHandler err")
//	}
//
//	return nil
//}

func (gteh *FaqEventHandler) Handle(e event.Event) error {
	log.WithField("event", fmt.Sprintf("%+v", e.Data())).Info("收到频道内的文字消息.")
	err := func() error {
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

		if ContainsPhrase(msgEvent.KMarkdown.RawContent, "巫妖王") || ContainsPhrase(msgEvent.KMarkdown.RawContent, "wlk") || ContainsPhrase(msgEvent.KMarkdown.RawContent, "通行证") {

			authorId := msgEvent.Author.ID
			askUser := " 您好，如果您想咨询everlook巫妖王项目的问题，\n请至这个kook服务器咨询: "
			inviteLink := linkUrl("https://kook.top/20CAGp")

			content := "(met)" + authorId + "(met)" + askUser + inviteLink

			SendGroupTextMessage(msgEvent.TargetId, content, client)
		}

		if ContainsPhrase(msgEvent.KMarkdown.RawContent, "举报") && msgEvent.GuildID != "" && DoesNotContainPhrase(msgEvent.KMarkdown.RawContent, "bot") {
			Cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			Cards.AddKmarkdown("如果你想举报，目前everlook亚服支持的举报方式如下：")
			Cards.AddKmarkdown("1. 在" + "(chn)" + "7124698304119930" + "(chn)" + "发送 !bot 举报 就会有机器人小E私信你，按提示和机器人对话2次提交举报信息。")
			Cards.AddKmarkdown("2. 在1.12客户端或者定制1.142客户端内点❓提交举报信息。")
			Cards.AddKmarkdown("更多问题请查看左侧" + "(chn)" + "5573736130652654" + "(chn)" + "频道")

			CardsContent, err := kook_CardBuild.GenerateCardMessageContent(Cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(CardsContent, msgEvent.TargetId, client)

		}

		if ContainsPhrase(msgEvent.KMarkdown.RawContent, "暗影之眼") || ContainsPhrase(msgEvent.KMarkdown.RawContent, "肌腱") {
			Cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			Cards.AddKmarkdown("(met)" + msgEvent.Author.ID + "(met)" + "  " + "成年蓝龙的肌腱和暗影之眼现在会从恰当等级的恶魔和蓝龙怪物身上掉落。值得注意的是卡扎克和艾索雷葛斯肯定会掉落这两样东西，但是从别的生物上获得这些东西的几率要低很多。")
			Cards.AddKmarkdown("更多问题请查看左侧" + "(chn)" + "5573736130652654" + "(chn)" + "频道")

			CardsContent, err := kook_CardBuild.GenerateCardMessageContent(Cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(CardsContent, msgEvent.TargetId, client)
		}

		if ContainsPhrase(msgEvent.KMarkdown.RawContent, "官网是") || ContainsPhrase(msgEvent.KMarkdown.RawContent, "官网地址") || ContainsPhrase(msgEvent.KMarkdown.RawContent, "亚服官网") {
			Cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			Cards.AddKmarkdown("(met)" + msgEvent.Author.ID + "(met)")
			Cards.AddKmarkdown("everlook60级香草时代官网是: " + linkUrl("https://cn.everlook-wow.net"))
			Cards.AddKmarkdown("everlook80级巫妖王官网是: " + linkUrl("https://wotlk.everlook-wow.net"))
			Cards.AddKmarkdown("更多问题请查看左侧" + "(chn)" + "5573736130652654" + "(chn)" + "频道")

			CardsContent, err := kook_CardBuild.GenerateCardMessageContent(Cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(CardsContent, msgEvent.TargetId, client)
		}

		if ContainsPhrase(msgEvent.KMarkdown.RawContent, "密语") {
			Cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			Cards.AddKmarkdown("(met)" + msgEvent.Author.ID + "(met)" + "  " + "密语问题请到FAQ里面查看相应处理办法或者下载聊天修复插件。")
			Cards.AddKmarkdown("更多问题请查看左侧" + "(chn)" + "5573736130652654" + "(chn)" + "频道")

			CardsContent, err := kook_CardBuild.GenerateCardMessageContent(Cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(CardsContent, msgEvent.TargetId, client)
		}

		if (ContainsPhrase(msgEvent.KMarkdown.RawContent, "辛特兰") || ContainsPhrase(msgEvent.KMarkdown.RawContent, "灼热")) && (ContainsPhrase(msgEvent.KMarkdown.RawContent, "npc") || ContainsPhrase(msgEvent.KMarkdown.RawContent, "营地")) {
			Cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			Cards.AddKmarkdown("(met)" + msgEvent.Author.ID + "(met)" + "  " + "辛特兰营地在1.5版本增加。")
			Cards.AddKmarkdown("(met)" + msgEvent.Author.ID + "(met)" + "  " + "灼热峡谷营地在1.5版本增加。")
			Cards.AddKmarkdown("更多问题请查看左侧" + "(chn)" + "5573736130652654" + "(chn)" + "频道")

			CardsContent, err := kook_CardBuild.GenerateCardMessageContent(Cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(CardsContent, msgEvent.TargetId, client)
		}

		if (ContainsPhrase(msgEvent.KMarkdown.RawContent, "鸟点") || ContainsPhrase(msgEvent.KMarkdown.RawContent, "飞行点")) && ContainsPhrase(msgEvent.KMarkdown.RawContent, "没有") {
			Cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			Cards.AddKmarkdown("(met)" + msgEvent.Author.ID + "(met)" + "  " + "陶拉祖营地飞行点在1.6版本增加。")
			Cards.AddKmarkdown("(met)" + msgEvent.Author.ID + "(met)" + "  " + "棘齿城飞行点在1.11版本增加。")
			Cards.AddKmarkdown("更多问题请查看左侧" + "(chn)" + "5573736130652654" + "(chn)" + "频道")

			CardsContent, err := kook_CardBuild.GenerateCardMessageContent(Cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(CardsContent, msgEvent.TargetId, client)
		}

		if ContainsPhrase(msgEvent.KMarkdown.RawContent, "T1") && ContainsPhrase(msgEvent.KMarkdown.RawContent, "属性") {
			Cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			Cards.AddKmarkdown("(met)" + msgEvent.Author.ID + "(met)" + "  " + "T1装备属性将在1.5版本改变。")
			Cards.AddKmarkdown("更多问题请查看左侧" + "(chn)" + "5573736130652654" + "(chn)" + "频道")

			CardsContent, err := kook_CardBuild.GenerateCardMessageContent(Cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(CardsContent, msgEvent.TargetId, client)
		}

		if ContainsPhrase(msgEvent.KMarkdown.RawContent, "T2") && ContainsPhrase(msgEvent.KMarkdown.RawContent, "属性") {
			Cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			Cards.AddKmarkdown("(met)" + msgEvent.Author.ID + "(met)" + "  " + "T2装备属性将在1.5版本改变。")
			Cards.AddKmarkdown("更多问题请查看左侧" + "(chn)" + "5573736130652654" + "(chn)" + "频道")

			CardsContent, err := kook_CardBuild.GenerateCardMessageContent(Cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(CardsContent, msgEvent.TargetId, client)
		}

		if ContainsPhrase(msgEvent.KMarkdown.RawContent, "商城工单") {
			Cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			Cards.AddKmarkdown("这里的GM都没有权限去处理商城相关问题。如果遇到金额错误等问题。请到" + linkUrl("https://everlook.zendesk.com/hc/en-us/requests/new") + "提交表单，表单内上传充值记录截图 和你所遇到的问题。越详细越好。提交后请耐心等待。因为有时差。处理时间会慢点。但放心，有问题的都会处理好。感谢支持和理解。")
			Cards.AddKmarkdown("更多问题请查看左侧" + "(chn)" + "5573736130652654" + "(chn)" + "频道")

			CardsContent, err := kook_CardBuild.GenerateCardMessageContent(Cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(CardsContent, msgEvent.TargetId, client)
		}

		if ContainsPhrase(msgEvent.KMarkdown.RawContent, "工单") && ContainsPhrase(msgEvent.KMarkdown.RawContent, "游戏") {
			Cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			Cards.AddKmarkdown("游戏内问题，请在1.12和1.142定制版点右下角红色问号❓或者esc支持内填写工单，选项可任意填写，中文填写提交，提交后可以下线可以换号，一般1天内处理，请耐心等待。")
			Cards.AddKmarkdown("更多问题请查看左侧" + "(chn)" + "5573736130652654" + "(chn)" + "频道")

			CardsContent, err := kook_CardBuild.GenerateCardMessageContent(Cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(CardsContent, msgEvent.TargetId, client)
		}

		if ContainsPhrase(msgEvent.KMarkdown.RawContent, "工程传送器") {
			Cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			Cards.AddKmarkdown("(met)" + msgEvent.Author.ID + "(met)" + "  " + "工程传送器将在1.5版本增加。")
			Cards.AddKmarkdown("更多问题请查看左侧" + "(chn)" + "5573736130652654" + "(chn)" + "频道")

			CardsContent, err := kook_CardBuild.GenerateCardMessageContent(Cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(CardsContent, msgEvent.TargetId, client)
		}

		if ContainsPhrase(msgEvent.KMarkdown.RawContent, "野蛮角斗士胸甲") && ContainsPhrase(msgEvent.KMarkdown.RawContent, "属性") {
			Cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			Cards.AddKmarkdown("(met)" + msgEvent.Author.ID + "(met)" + "  " + "野蛮角斗士胸甲属性将在1.11版本变更")
			Cards.AddKmarkdown("更多问题请查看左侧" + "(chn)" + "5573736130652654" + "(chn)" + "频道")

			CardsContent, err := kook_CardBuild.GenerateCardMessageContent(Cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(CardsContent, msgEvent.TargetId, client)
		}

		if ContainsPhrase(msgEvent.KMarkdown.RawContent, "猎人") && ContainsPhrase(msgEvent.KMarkdown.RawContent, "技能") {
			Cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			Cards.AddKmarkdown("(met)" + msgEvent.Author.ID + "(met)" + "  " + "请登录1.12或者 1.142新定制客户端的 红色问号❓联系在线GM处理.")
			Cards.AddKmarkdown("更多问题请查看左侧" + "(chn)" + "5573736130652654" + "(chn)" + "频道")

			CardsContent, err := kook_CardBuild.GenerateCardMessageContent(Cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(CardsContent, msgEvent.TargetId, client)
		}

		if ContainsPhrase(msgEvent.KMarkdown.RawContent, "改密码") || ContainsPhrase(msgEvent.KMarkdown.RawContent, "改账号密码") {
			Cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			Cards.AddKmarkdown("(met)" + msgEvent.Author.ID + "(met)" + "  " + "请登录官网" + linkUrl("https://cn.everlook-wow.net") + "在“个人资料”中进行修改")
			Cards.AddKmarkdown("更多问题请查看左侧" + "(chn)" + "5573736130652654" + "(chn)" + "频道")

			CardsContent, err := kook_CardBuild.GenerateCardMessageContent(Cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(CardsContent, msgEvent.TargetId, client)
		}

		if ContainsPhrase(msgEvent.KMarkdown.RawContent, "下载") || ContainsPhrase(msgEvent.KMarkdown.RawContent, "客户端") {
			Cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			Cards.AddKmarkdown("请前往" + "(chn)" + "4736496800586949" + "(chn)" + "频道下载")

			CardsContent, err := kook_CardBuild.GenerateCardMessageContent(Cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(CardsContent, msgEvent.TargetId, client)
		}

		if ContainsPhrase(msgEvent.KMarkdown.RawContent, "启动器") {
			Cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			Cards.AddKmarkdown("(met)" + msgEvent.Author.ID + "(met)" + "  " + "请前往" + linkUrl("https://cn.everlook-wow.net/launcher/") + "下载最新官方启动器")

			CardsContent, err := kook_CardBuild.GenerateCardMessageContent(Cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(CardsContent, msgEvent.TargetId, client)
		}

		if ContainsPhrase(msgEvent.KMarkdown.RawContent, "充值") || ContainsPhrase(msgEvent.KMarkdown.RawContent, "捐助") {
			Cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			Cards.AddKmarkdown("请查看" + "(chn)" + "2777878836706338" + "(chn)" + "里面对应教程")

			CardsContent, err := kook_CardBuild.GenerateCardMessageContent(Cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(CardsContent, msgEvent.TargetId, client)
		}

		if ContainsPhrase(msgEvent.KMarkdown.RawContent, "幻化") || ContainsPhrase(msgEvent.KMarkdown.RawContent, "改名") {
			Cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			Cards.AddKmarkdown("请查看" + "(chn)" + "2041871976617459" + "(chn)" + "里面对应教程")

			CardsContent, err := kook_CardBuild.GenerateCardMessageContent(Cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(CardsContent, msgEvent.TargetId, client)
		}

		if ContainsPhrase(msgEvent.KMarkdown.RawContent, "处理") || ContainsPhrase(msgEvent.KMarkdown.RawContent, "解决") {
			Cards, err := kook_CardBuild.NewCardWithOption(kook_CardBuild.CardThemeNone, "lg", "#000000")
			if err != nil {
				return err
			}

			Cards.AddKmarkdown("请查看" + "(chn)" + "5573736130652654" + "(chn)" + "频道内容")

			CardsContent, err := kook_CardBuild.GenerateCardMessageContent(Cards)
			if err != nil {
				return err
			}

			SendGroupCardessage(CardsContent, msgEvent.TargetId, client)
		}

		return nil
	}()
	if err != nil {
		log.WithError(err).Error("GroupTextEventHandler err")
	}

	return nil
}

func ContainsPhrase(text, phrase string) bool {
	// 将文字和词组都转换为小写，然后检查是否包含词组
	return strings.Contains(strings.ToLower(text), strings.ToLower(phrase))
}

func DoesNotContainPhrase(text, phrase string) bool {
	// 将文字和词组都转换为小写，然后检查是否不包含词组
	return !strings.Contains(strings.ToLower(text), strings.ToLower(phrase))
}
