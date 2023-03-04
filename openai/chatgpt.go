package openai

import (
	"encoding/json"
	"fmt"
	"strings"
	
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type Client struct {
	apiKey string
	debug  bool
	c      *resty.Client
	url    string
	logger *zap.Logger
}

var ChatGPT *Client

func NewClient(apiKey string, debug bool, proxy string) Client {
	c := resty.New()
	if proxy != "" {
		c.SetProxy(proxy)
	}
	res := Client{
		apiKey: apiKey,
		c:      c,
		url:    "https://api.openai.com/v1/completions",
		logger: zap.L(),
		debug:  debug,
	}
	ChatGPT = &res
	return res
}

func (c Client) Ask(question string) (string, error) {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", c.apiKey),
	}
	body := c.buildBody(question)
	resp, err := c.c.SetDebug(c.debug).R().SetHeaders(headers).SetBody(body).Post(c.url)
	if err != nil {
		c.logger.Error("post failed", zap.Error(err), zap.Any("body", body))
		return "", err
	}
	return c.analyseResp(resp.Body())
}

func (c Client) buildBody(msg string) ChatGPTRequest {
	return ChatGPTRequest{
		Model:            "text-davinci-003",
		Prompt:           msg,
		MaxTokens:        4000,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
}

func (c Client) analyseResp(data []byte) (string, error) {
	res := ChatGPTResponse{}
	if err := json.Unmarshal(data, &res); err != nil {
		c.logger.Error("json unmarshal failed", zap.Error(err))
		return "", err
	}
	var answer string
	if len(res.Choices) > 0 {
		for _, v := range res.Choices {
			answer = v["text"].(string)
			break
		}
	}
	
	errRes := &ChatGPTError{}
	if err := json.Unmarshal(data, errRes); err != nil {
		c.logger.Error("json unmarshal failed", zap.Error(err))
		return "", err
	}
	
	if len(answer) == 0 {
		answer = errRes.Error["message"].(string)
	}
	c.logger.Info("get answer from gpt", zap.String("answer", answer))
	result := strings.TrimSpace(answer)
	
	// 如果在提问的时候没有包含？,AI会自动在开头补充个？看起来很奇怪
	if strings.HasPrefix(result, "?") {
		answer = strings.Replace(result, "?", "", -1)
	}
	if strings.HasPrefix(answer, "？") {
		answer = strings.Replace(answer, "？", "", -1)
	}
	// 微信不支持markdown格式，所以把反引号直接去掉
	answer = strings.Replace(answer, "`", "", -1)
	
	return answer, nil
}
