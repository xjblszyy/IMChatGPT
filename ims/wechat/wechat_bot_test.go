package wechat

import (
	"os"
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/xjblszyy/im-chatgpt/openai"
)

func TestBot_Start(t *testing.T) {
	cfg := Config{
		Keyword: "小白",
		Enabled: true,
	}
	apiKey := os.Getenv("test_api_key")
	openai.NewClient(apiKey, true, "")
	bot := NewBot(cfg, "")
	err := bot.Start()
	assert.NoError(t, err)
}
