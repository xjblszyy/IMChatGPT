package config

import (
	"github.com/jinzhu/configor"
	"github.com/xjblszyy/im-chatgpt/ims/telegram"
	"github.com/xjblszyy/im-chatgpt/ims/wechat"
	"go.uber.org/zap"
)

type Config struct {
	Debug    bool            `yaml:"debug,omitempty" default:"false"`
	Proxy    string          `yaml:"proxy" default:""`
	Logger   LoggerConfig    `yaml:"logger,omitempty"`
	ApiKey   string          `yaml:"api_key" default:""`
	Keyword  string          `yaml:"keyword" default:"im"`
	Wechat   wechat.Config   `yaml:"wechat"`
	Telegram telegram.Config `yaml:"telegram"`
	Wecon    WeconConfig     `yaml:"wecon"`
}

type LoggerConfig struct {
	Level string `yaml:"level,omitempty" default:"debug"`
	// json or text
	Format string `yaml:"format,omitempty" default:"json"`
	// file
	Output string `yaml:"output,omitempty" default:""`
}

type TelegramConfig struct {
}

type WeconConfig struct {
}

var C *Config

func Init(cfgFile string) {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	
	if cfgFile != "" {
		if err := configor.New(&configor.Config{AutoReload: true}).Load(C, cfgFile); err != nil {
			zap.L().Panic("init config fail", zap.Error(err))
		}
	} else {
		if err := configor.New(&configor.Config{AutoReload: true}).Load(C); err != nil {
			zap.L().Panic("init config fail", zap.Error(err))
		}
	}
	
	zap.L().Debug("loaded config")
}

func init() {
	C = &Config{}
}
