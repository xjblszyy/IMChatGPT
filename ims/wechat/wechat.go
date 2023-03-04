package wechat

import (
	"os"
	
	"github.com/eatmoreapple/openwechat"
	"go.uber.org/zap"
)

type Bot struct {
	logger       *zap.Logger
	tokenStorage string
	cfg          Config
}

func NewBot(cfg Config, globalKeyword string) Bot {
	if cfg.Keyword == "" {
		cfg.Keyword = globalKeyword
	}
	return Bot{
		logger:       zap.L(),
		tokenStorage: "token.wechat",
		cfg:          cfg,
	}
}

func (b Bot) Start() error {
	NewMsgHandler(b.cfg.Keyword, b.cfg.BlackList)
	bot := openwechat.DefaultBot(openwechat.Desktop)
	bot.MessageHandler = Handler
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			b.logger.Error("remove token file failed", zap.Error(err))
		}
	}(b.tokenStorage)
	storage := openwechat.NewFileHotReloadStorage(b.tokenStorage)
	if err := bot.HotLogin(storage, openwechat.NewRetryLoginOption()); err != nil {
		b.logger.Error("hot login failed", zap.Error(err))
		return err
	}
	
	user, err := bot.GetCurrentUser()
	if err != nil {
		b.logger.Error("get current user failed", zap.Error(err))
		return err
	}
	b.logger.Info("登陆成功", zap.String("user", user.String()))
	
	err = bot.Block()
	if err != nil {
		b.logger.Error("bot block failed", zap.Error(err))
		return err
	}
	return nil
}
