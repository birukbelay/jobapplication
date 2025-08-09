package config

import conf "github.com/birukbelay/gocmn/src/config"

type EnvConfig struct {
	ServerPort string `koanf:"SERVER_PORT"`
	ServerHost string `koanf:"SERVER_HOST"`

	// SqlDbConfig SqlDbConfig `koanf:",squash"`
	JwtVar

	conf.SqlDbConfig

	SmtpCred
	conf.CloudinaryConfig
	conf.S3Config
	conf.KeyValConfig
}

type JwtVar struct {
	AccessSecret  string `koanf:"ACCESS_SECRET"`
	RefreshSecret string `koanf:"REFRESH_SECRET"`

	AccessExpireMin  int `koanf:"ACCESS_SECRET_EXPIRE_MIN"`
	RefreshExpireMin int `koanf:"REFRESH_SECRET_EXPIRES_MIN"`
}
type SmtpCred struct {
	SmtpHost     string `koanf:"SMTP_HOST"`
	SmtpPort     string `koanf:"SMTP_PORT"`
	SmtpPwd      string `koanf:"SMTP_PWD"`
	SmtpUsername string `koanf:"SMTP_USERNAME"`
}
