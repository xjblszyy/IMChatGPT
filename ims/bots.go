package ims

import (
	"os"
	"os/signal"
	"syscall"
	
	"github.com/xjblszyy/im-chatgpt/config"
	"github.com/xjblszyy/im-chatgpt/ims/telegram"
	"github.com/xjblszyy/im-chatgpt/ims/wechat"
	"go.uber.org/zap"
)

type Bots struct {
	logger *zap.Logger
	cfg    config.Config
}

func NewBots(cfg config.Config) Bots {
	return Bots{
		logger: zap.L(),
		cfg:    cfg,
	}
}

func (b Bots) Start() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGINT)
	if b.cfg.Wechat.Enabled {
		go func() {
			wechatBot := wechat.NewBot(b.cfg.Wechat, b.cfg.Keyword)
			if err := wechatBot.Start(); err != nil {
				b.logger.Error("start wechat bot failed", zap.Error(err))
			}
		}()
	}
	if b.cfg.Telegram.Enabled {
		go func() {
			tgBot := telegram.NewBot(b.cfg.Telegram, b.cfg.Keyword)
			if err := tgBot.Start(); err != nil {
				b.logger.Error("start tg bot failed", zap.Error(err))
			}
		}()
	}
	data := <-s
	b.logger.Info("收到退出新型号", zap.String("signal", data.String()))
}
