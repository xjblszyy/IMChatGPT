package telegram

type Config struct {
	BlackList []string `yaml:"black_list"`
	Enabled   bool     `yaml:"enabled"`
	Keyword   string   `yaml:"keyword"`
	Token     string   `yaml:"token"`
}
