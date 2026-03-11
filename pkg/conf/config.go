package conf

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var Config *Configuration

type Configuration struct {
	Service  Service  `mapstructure:"service"`
	Database Database `mapstructure:"database"`
	Redis    Redis    `mapstructure:"redis"`
	Email    Email    `mapstructure:"email"`
	LLM      LLMConfig `mapstructure:"llm"`
}

type Service struct {
	AppMode  string `mapstructure:"app_mode"`
	HttpPort string `mapstructure:"http_port"`
}

type Database struct {
	DbType    string `mapstructure:"db_type"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Host      string `mapstructure:"host"`
	DbName    string `mapstructure:"db_name"`
	Charset   string `mapstructure:"charset"`
	ParseTime bool   `mapstructure:"parse_time"`
	Loc       string `mapstructure:"loc"`
}

type Redis struct {
	RedisDebug  bool   `mapstructure:"redis_debug"`
	RedisAddr   string `mapstructure:"redis_addr"`
	RedisPw     string `mapstructure:"redis_pw"`
	RedisDb 	int    `mapstructure:"redis_db_name"`
}

type Email struct {
	Host     string `mapstructure:"host"`     
	Port     int    `mapstructure:"port"`     
	User     string `mapstructure:"user"`     
	Password string `mapstructure:"password"` 
}

type RedisKeyConfig struct {
	CaptchaPrefix string
}

var DefaultRedisKeyConfig = RedisKeyConfig{
	CaptchaPrefix: "captcha:%s",
}

type LLMConfig struct {
	Type      string `mapstructure:"type"`       // 服务商类型：如 "openai", "ollama", "deepseek"
	APIKey    string `mapstructure:"api_key"`    // 你的密钥 (Ollama 通常不需要，但 OpenAI/DeepSeek 必须有)
	BaseURL   string `mapstructure:"base_url"`   // 请求网关地址
	ModelName string `mapstructure:"model_name"` // 具体使用的模型名，比如 "gpt-4o", "qwen-turbo"
	
	// 以下是一些可选的高级配置项，根据需要启用
	// MaxTokens int    `mapstructure:"max_tokens"` // 限制 AI 最多吐多少个字，防止扣费破产
	// Timeout   int    `mapstructure:"timeout"`    // 超时时间(秒)，防止 AI 卡死导致服务器 Goroutine 泄漏

}

func Init() {

	workDir, _ := os.Getwd()

	env := os.Getenv("APP_ENV")
	configName := "config"
	// docker compose 
	if env == "prod" {
		configName = "config.prod"
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/config")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	if err := viper.Unmarshal(&Config); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	fmt.Println("config file:", viper.ConfigFileUsed())
}