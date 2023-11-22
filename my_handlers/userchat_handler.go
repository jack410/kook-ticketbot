package my_handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/bytedance/sonic"
	"github.com/gookit/event"
	"github.com/kaiheila/golang-bot/api/base"
	event2 "github.com/kaiheila/golang-bot/api/base/event"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type UserChatCreateResponse struct {
	Code    int                        `json:"code"`
	Message string                     `json:"message"`
	Data    UserChatCreateResponseData `json:"data"`
}

type UserChatCreateResponseData struct {
	Code            string     `json:"code"`
	LastReadTime    int64      `json:"last_read_time"`
	LatestMsgTime   int64      `json:"latest_msg_time"`
	UnreadCount     int        `json:"unread_count"`
	IsFriend        bool       `json:"is_friend"`
	IsBlocked       bool       `json:"is_blocked"`
	IsTargetBlocked bool       `json:"is_target_blocked"`
	TargetInfo      TargetInfo `json:"target_info"`
}

type TargetInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Online   bool   `json:"online"`
	Avatar   string `json:"avatar"`
}

type DirectMessageResponse struct {
	Code    int                       `json:"code"`
	Message string                    `json:"message"`
	Data    DirectMessageResponseData `json:"data"`
}

type DirectMessageResponseData struct {
	MsgID        string `json:"msg_id"`
	MsgTimestamp int    `json:"msg_timestamp"`
	Nonce        string `json:"nonce"`
}

type DirectMessagePayload struct {
	Type     int    `json:"type,omitempty"`
	TargetID string `json:"target_id,omitempty"`
	ChatCode string `json:"chat_code,omitempty"`
	Content  string `json:"content,omitempty"`
	Quote    string `json:"quote,omitempty"`
	Nonce    string `json:"nonce,omitempty"`
}

// 创建一个私信聊天会话，返回chat_code https://developer.kookapp.cn/doc/http/user-chat#%E5%88%9B%E5%BB%BA%E7%A7%81%E4%BF%A1%E8%81%8A%E5%A4%A9%E4%BC%9A%E8%AF%9D
func (gteh *GroupTextEventHandler) UserChatCreate(targetId string) UserChatCreateResponse {
	// 构建请求数据
	requestData := map[string]interface{}{
		"target_id": targetId,
	}
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		log.Println(err)
	}

	// 发送 HTTP POST 请求
	url := "https://www.kookapp.cn/api/v3/user-chat/create"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println(err)
	}
	request.Header.Set("Authorization", "Bot "+gteh.Token)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}

	defer response.Body.Close()

	// 解析返回的 JSON 消息
	var responseObject UserChatCreateResponse
	err = json.NewDecoder(response.Body).Decode(&responseObject)
	if err != nil {
		log.Println(err)
	}

	return responseObject
}

// send direct message to user
func (gteh *GroupTextEventHandler) DirectMessageSend(chatCode string, s string) {
	//构建请求数据
	payload := DirectMessagePayload{
		Type:     1, //send text message, 1 text, 9 kmarkdown, 10 card.
		ChatCode: chatCode,
		Content:  s,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshal direct message payload:", err)
		return
	}

	// 发送 HTTP POST 请求
	url := "https://www.kookapp.cn/api/v3/direct-message/create"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Println("Failed to create request:", err)
		return
	}
	request.Header.Set("Authorization", "Bot "+gteh.Token)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Failed to send request:", err)
		return
	}
	defer response.Body.Close()

	// 解析返回的 JSON 响应
	var responseObject DirectMessageResponse
	err = json.NewDecoder(response.Body).Decode(&responseObject)
	if err != nil {
		log.Println("Failed to parse JSON:", err)
		return
	}

	if responseObject.Code != 0 {
		log.Println("返回错误码：", responseObject.Code)
		log.Println("返回错误消息：", responseObject.Message)
	} else {
		log.Println(responseObject.Message, "msg_id:", responseObject.Data.MsgID)
	}
}

type DirectMessageFrameHandler struct {
}

func (dm *DirectMessageFrameHandler) Handle(e event.Event) error {
	log.WithField("event", e).WithField("data", e.Data()).Info("ReceiveDirectMessageHandler receive DirectMessage.")
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
	log.Infof("Received directmessage json event:%+v", msgEvent)
	if err != nil {
		return err
	}
	if msgEvent.Author.Bot {
		log.Info("bot message")
		return nil
	}
	return nil
}
