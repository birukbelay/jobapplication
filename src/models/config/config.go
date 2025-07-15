package config

type EnvConfig struct {
	ServerPort string `koanf:"SERVER_PORT"`
	ServerHost string `koanf:"SERVER_HOST"`

	SqlDbConfig SqlDbConfig `koanf:",squash"`
	JwtVar
	S3Config
	CloudinaryConfig
	KeyValConfig
	SmtpCred
}

type S3Config struct {
	Endpoint         string `koanf:"S3_ENDPOINT"`
	S3WebEndpoint    string `koanf:"S3_WEB_ENDPOINT"`
	Region           string `koanf:"S3_REGION"`
	S3ForcePathStyle bool   `koanf:"S3_FORCE_PATH_STYLE"` // TODO: Remove?
	AccessKey        string `koanf:"S3_ACCESS_KEY_ID"`
	SecretKey        string `koanf:"S3_SECRET_ACCESS_KEY"`
	BucketName       string `koanf:"S3_BUCKET_NAME"`
}
type CloudinaryConfig struct {
	CloudinaryCloudName string `koanf:"CLOUDINARY_CLOUD_NAME"`
	CloudinaryAPIKey    string `koanf:"CLOUDINARY_API_KEY"`
	CloudinaryAPISecret string `koanf:"CLOUDINARY_API_SECRET"`
	CloudinaryFolder    string `koanf:"CLOUDINARY_FOLDER"`
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
type KeyValConfig struct {
	KVDbName   int    `koanf:"KV_DB"`
	KVUsername string `koanf:"KV_USER"`
	KVPassword string `koanf:"KV_PASSWORD"`
	KVHost     string `koanf:"KV_HOST"`
	KVPort     string `koanf:"KV_PORT"`
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
