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