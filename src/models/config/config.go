package config

import conf "github.com/birukbelay/gocmn/src/config"

type EnvConfig struct {
	ServerPort string `koanf:"SERVER_PORT"`
	ServerHost string `koanf:"SERVER_HOST"`

	SqlDbConfig SqlDbConfig `koanf:",squash"`
	JwtVar
	// S3Config
	// CloudinaryConfig

	SmtpCred
	conf.CloudinaryConfig
	conf.S3Config
	conf.KeyValConfig
}

type SqlDbConfig struct {
	DbName   string `koanf:"POSTGRES_DB"`
	Username string `koanf:"POSTGRES_USER"`
	Password string `koanf:"POSTGRES_PASSWORD"`
	SqlHost  string `koanf:"POSTGRES_HOST"`
	SqlPort  string `koanf:"POSTGRES_PORT"`
	Driver   string `koanf:"SQL_DRIVER"`
	SSLMode  string `koanf:"SSL_MODE"`
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
