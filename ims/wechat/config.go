package wechat

type Config struct {
	BlackList []string `yaml:"black_list"`
	Enabled   bool     `yaml:"enabled"`
	Keyword   string   `yaml:"keyword"`
}
