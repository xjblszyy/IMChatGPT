package telegram

import (
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xjblszyy/im-chatgpt/openai"
	"github.com/xjblszyy/im-chatgpt/utils"
	"go.uber.org/zap"
)

type Bot struct {
	cfg    Config
	logger *zap.Logger
}

func NewBot(cfg Config, globalKeyword string) Bot {
	if cfg.Keyword == "" {
		cfg.Keyword = globalKeyword
	}
	return Bot{
		cfg:    cfg,
		logger: zap.L(),
	}
}

func (b Bot) Start() error {
	botAPI, err := tgBot.NewBotAPI(b.cfg.Token)
	if err != nil {
		b.logger.Error("new tg bot failed", zap.Error(err), zap.String("token", b.cfg.Token))
		return err
	}
	// botAPI.Debug = b.debug
	b.logger.Info("登陆telegram成功", zap.String("user_name", botAPI.Self.UserName))
	u := tgBot.NewUpdate(0)
	u.Timeout = 60
	
	updates := botAPI.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil { // If we got a message
			b.logger.Info("收到telegram消息", zap.String("user_name", update.Message.From.UserName),
				zap.String("msg", update.Message.Text))
			content := update.Message.Text
			keyword := b.cfg.Keyword
			if !utils.ContainKeyword(keyword, content) {
				b.logger.Info("不包含关键字，不向chatGPT发送信息", zap.String("keyword", keyword), zap.String("question", content))
				return nil
			}
			if utils.ContainBlackList(update.Message.From.UserName, b.cfg.BlackList) {
				b.logger.Info("在黑名单中，不向chatGPT发送信息", zap.Strings("blacklist", b.cfg.BlackList),
					zap.String("username", update.Message.From.UserName))
				return nil
			}
			
			question := utils.GetQuestionFromMsg(b.cfg.Keyword, update.Message.Text)
			b.logger.Info("收到的问题是", zap.String("question", question))
			
			var answer string
			answer, err = openai.ChatGPT.Ask(question)
			if err != nil {
				b.logger.Error("ask chat gpt failed", zap.String("question", question), zap.Error(err))
				answer = "您的问题ChatGPT回答失败，具体请见日志"
			}
			
			msg := tgBot.NewMessage(update.Message.Chat.ID, answer)
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err := botAPI.Send(msg); err != nil {
				b.logger.Error("send telegram msg failed", zap.Error(err))
			}
		}
	}
	return nil
}
