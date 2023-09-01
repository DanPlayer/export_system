package feishu

import (
	"encoding/json"
	"export_system/utils"
)

// ApiRobotLink API告警机器人webhook链接
const ApiRobotLink = ""

type MessageResult struct {
	Extra         interface{} `json:"Extra"`
	StatusCode    int         `json:"StatusCode"`
	StatusMessage string      `json:"StatusMessage"`
}

// TextMessage 文本
type TextMessage struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

// SendTextMsg 发送文本消息
func SendTextMsg(url string, content string) error {
	msg := TextMessage{
		MsgType: "text",
		Content: struct {
			Text string `json:"text"`
		}{
			Text: content,
		},
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = utils.HttpPostJson(url, bytes)
	if err != nil {
		return err
	}
	return nil
}

// InteractiveMessage 卡片
type InteractiveMessage struct {
	MsgType string `json:"msg_type"`
	Card    struct {
		Config struct {
			WideScreenMode bool `json:"wide_screen_mode"`
			EnableForward  bool `json:"enable_forward"`
		} `json:"config"`
		Elements []struct {
			Tag  string `json:"tag"`
			Text struct {
				Content string `json:"content"`
				Tag     string `json:"tag"`
			} `json:"text,omitempty"`
			Actions []struct {
				Tag  string `json:"tag"`
				Text struct {
					Content string `json:"content"`
					Tag     string `json:"tag"`
				} `json:"text"`
				Url   string `json:"url"`
				Type  string `json:"type"`
				Value struct {
				} `json:"value"`
			} `json:"actions,omitempty"`
		} `json:"elements"`
		Header struct {
			Title struct {
				Content string `json:"content"`
				Tag     string `json:"tag"`
			} `json:"title"`
		} `json:"header"`
	} `json:"card"`
}

// SendInteractiveMsg 发送卡片消息
func SendInteractiveMsg(url string, msg InteractiveMessage) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = utils.HttpPostJson(url, bytes)
	if err != nil {
		return err
	}
	return nil
}
