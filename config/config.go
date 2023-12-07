package config

import (
	"github.com/spf13/viper"
)

type Config struct{}

func NewConfig() *Config {
	viper.AutomaticEnv()

	return &Config{}
}

func (c *Config) GetGrpcPort() string {
	return viper.GetString("SELF_GRPC_PORT")
}

func (c *Config) GetDbHost() string {
	return viper.GetString("DB_HOST")
}

func (c *Config) GetDbUsername() string {
	return viper.GetString("DB_USER")
}

func (c *Config) GetDbName() string {
	return viper.GetString("DB_NAME")
}

func (c *Config) GetDbPort() string {
	return viper.GetString("DB_PORT")
}

func (c *Config) GetDbPassword() string {
	return viper.GetString("DB_PASSWORD")
}

func (c *Config) GetCodeSalt() string {
	return viper.GetString("SALT")
}

func (c *Config) GetEncryptKey() string {
	return viper.GetString("ENCRYPT_KEY")
}

func (c *Config) GetGlobalCacheAddress() string {
	return viper.GetString("CACHE_ADDRESS")
}

func (c *Config) GetGlobalCachePassword() string {
	return viper.GetString("CACHE_PASSWORD")
}

func (c *Config) GetClientSecretKey() string {
	return viper.GetString("CLIENT_SECRET_KEY")
}

func (c *Config) GetApiSecret() string {
	return viper.GetString("API_KEY_SECRET")
}

func (c *Config) GetSecretKey() string {
	return viper.GetString("SECRET_KEY")
}

func (c *Config) GetApiPort() string {
	return viper.GetString("SELF_API_PORT")
}

func (c *Config) GetApiServerMode() string {
	return viper.GetString("SERVER_MODE")
}

func (c *Config) GetTgToken() string {
	return viper.GetString("TELEGRAM_TOKEN")
}

func (c *Config) GetTestnet() bool {
	return viper.GetBool("TESTNET")
}
