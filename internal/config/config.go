package config

import "github.com/spf13/viper"

type Config struct {
	Server
	DB
	RabbitMQ
	Redis
	JWT
}

type Server struct {
	AppPort int
	AppHost string
}

type DB struct {
	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string
}

type RabbitMQ struct {
	RabbitMQHost string
	RabbitMQPort int
	RabbitMQUser string
	RabbitMQPass string
}

type Redis struct {
	RedisAddr     string
	RedisDB       int
	RedisPassword string
}

type JWT struct {
	JWTSecret     string
	JWTExpiration int
}

var config Config

func Load() (*Config, error) {
	return LoadPath(".")
}

func LoadPath(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func Get() *Config {
	return &config
}
