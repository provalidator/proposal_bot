package models

var Config config

type config struct {
	Telegram struct {
		BotName     string `yaml:"BOT_NAME"`
		BotToken    string `yaml:"BOT_TOKEN"`
		ChatID      int    `yaml:"CHAT_ID"`
		ChatIDAdmin int    `yaml:"CHAT_ID_ADMIN"`
	} `yaml:"TELEGRAM"`
	Database struct {
		MysqlServerURL    string `yaml:"MYSQL_SERVER_URL"`
		MysqlServerPort   string `yaml:"MYSQL_SERVER_PORT"`
		MysqlUserID       string `yaml:"MYSQL_USER_ID"`
		MysqlUserPW       string `yaml:"MYSQL_USER_PW"`
		MysqlSelectDBName string `yaml:"MYSQL_SELECT_DB_NAME"`
	} `yaml:"DATABASE"`
	MainnetLCD []struct {
		ChainName string `yaml:"CHAIN_NAME"`
		URL       string `yaml:"URL"`
		EX        string `yaml:"EX"`
	} `yaml:"MAINNET_LCD"`
}
