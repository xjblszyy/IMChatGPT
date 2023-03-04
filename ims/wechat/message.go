package wechat

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/xjblszyy/im-chatgpt/openai"
	"github.com/xjblszyy/im-chatgpt/utils"
	"go.uber.org/zap"
)

var GlobalHandler *MsgHandler

func Handler(msg *openwechat.Message) {
	GlobalHandler.handle(msg)
}

type MsgHandler struct {
	logger    *zap.Logger
	keyword   string
	failedMsg string
	blackList []string
}

func NewMsgHandler(keyword string, blackList []string) MsgHandler {
	res := MsgHandler{
		logger:    zap.L(),
		keyword:   keyword,
		failedMsg: "您的问题ChatGPT回答失败，具体请见日志",
		blackList: blackList,
	}
	GlobalHandler = &res
	return res
}

func (h MsgHandler) handle(msg *openwechat.Message) {
	if !msg.IsText() {
		return
	}
	if err := h.SendText(msg); err != nil {
		h.logger.Error("send text failed", zap.Error(err))
	}
}

func (h *MsgHandler) SendText(msg *openwechat.Message) error {
	sender, err := msg.Sender()
	if err != nil {
		h.logger.Error("sender failed", zap.Error(err))
		return err
	}
	group := openwechat.Group{User: sender}
	h.logger.Info("收到微信消息", zap.String("sender", group.NickName), zap.String("content", msg.Content))
	
	content := msg.Content
	if !utils.ContainKeyword(h.keyword, content) {
		h.logger.Info("不包含关键字，不向chatGPT发送信息", zap.String("keyword", h.keyword), zap.String("question", msg.Content))
		return nil
	}
	if utils.ContainBlackList(group.NickName, h.blackList) {
		h.logger.Info("在黑名单中，不向chatGPT发送信息", zap.Strings("blacklist", h.blackList), zap.String("nickname", group.NickName))
		return nil
	}
	
	question := utils.GetQuestionFromMsg(h.keyword, content)
	h.logger.Info("收到的问题是", zap.String("question", question))
	
	answer, err := openai.ChatGPT.Ask(question)
	if err != nil {
		h.logger.Error("ask chat gpt failed", zap.String("question", question), zap.Error(err))
		
		text, err := msg.ReplyText("您的问题ChatGPT回答失败，具体请见日志")
		h.logger.Error("ChatGPT回答失败", zap.String("text", text.Content), zap.Error(err), zap.String("failed_msg", h.failedMsg))
		return err
	}
	if _, err = msg.ReplyText(answer); err != nil {
		h.logger.Error("reply text failed", zap.String("answer", answer), zap.Error(err))
		return err
	}
	
	return nil
}
