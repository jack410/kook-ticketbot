package my_handlers

import (
	"github.com/bytedance/sonic"
	"github.com/kaiheila/golang-bot/api/helper"
	log "github.com/sirupsen/logrus"
)

func SendGroupTextMessage(channel_id, content string, client *helper.ApiHelper) {
	echoData := map[string]string{
		"channel_id": channel_id,
		"content":    content,
	}
	echoDataByte, err := sonic.Marshal(echoData)
	if err != nil {
		log.Infof("echoData 序列化出错", err)
	}
	resp, err := client.SetBody(echoDataByte).Post()
	log.Info("sent post:%s", client.String())
	if err != nil {
		log.Infof("发送echoDataByte出错", err)
	}
	log.Infof("resp:%s", string(resp))

}

func SendGroupCardessage(cardMessageContent, channel_id string, client *helper.ApiHelper) {
	echoData := map[string]interface{}{
		"type":       10,
		"channel_id": channel_id,
		"content":    cardMessageContent,
	}
	//将map转化成[]byte
	echoDataByte, err := sonic.Marshal(echoData)
	if err != nil {
		log.Infof("将map转化成[]byte出错", err)
	}

	resp, err := client.SetBody(echoDataByte).Post()
	log.Info("sent post:%s", client.String())
	if err != nil {
		log.Infof("发送echoDataByte出错", err)
	}
	log.Infof("resp:%s", string(resp))
}
