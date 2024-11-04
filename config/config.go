package config

import "github.com/caarlos0/env/v10"

type Config struct {
	ServerConfig   ServerConfig   `envPrefix:"SERVER_"`
	DBConfig       DBConfig       `envPrefix:"DB_"`
	BinanceConfig  BinanceConfig  `envPrefix:"BINANCE_"`
	CacheStorage   CacheStorage   `envPrefix:"CACHE_STORAGE_"`
	TelegramConfig TelegramConfig `envPrefix:"TELEGRAM_"`
	ServiceConfig  ServiceConfig  `envPrefix:"SERVICE_"`
}

func NewConfig() *Config {
	cfg := &Config{}

	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	return cfg
}

type DBConfig struct {
	Username string `env:"USERNAME" envDefault:"postgres"`
	Password string `env:"PASSWORD" envDefault:"postgres"`
	Dbname   string `env:"DATABASE" envDefault:"postgres"`
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"5432"`
}

func (d *DBConfig) GetDbHost() string {
	return d.Host
}

func (d *DBConfig) GetDbUsername() string {
	return d.Username
}

func (d *DBConfig) GetDbName() string {
	return d.Dbname
}

func (d *DBConfig) GetDbPort() string {
	return d.Port
}

func (d *DBConfig) GetDbPassword() string {
	return d.Password
}

type ServerConfig struct {
	Port string `env:"PORT" envDefault:"8080"`
	Mode string `env:"MODE" envDefault:"release"`
}

func (s *ServerConfig) GetPort() string {
	return s.Port
}

func (s *ServerConfig) GetMode() string {
	return s.Mode
}

type BinanceConfig struct {
	ClientSecretKey string `env:"CLIENT_SECRET_KEY" envDefault:""`
	ApiSecretKey    string `env:"API_SECRET_KEY" envDefault:""`
	Testnet         bool   `env:"TESTNET" envDefault:"true"`
}

func (b *BinanceConfig) GetTestnet() bool {
	return b.Testnet
}

func (b *BinanceConfig) GetApiSecretKey() string {
	return b.ApiSecretKey
}

func (b *BinanceConfig) GetClientSecretKey() string {
	return b.ClientSecretKey
}

type CacheStorage struct {
	Address  string `env:"ADDRESS" envDefault:"localhost:6379"`
	Password string `env:"PASSWORD" envDefault:""`
	DB       int    `env:"DB" envDefault:"0"`
}

func (c *CacheStorage) GetAddress() string {
	return c.Address
}

func (c *CacheStorage) GetPassword() string {
	return c.Password
}

func (c *CacheStorage) GetDB() int {
	return c.DB
}

type TelegramConfig struct {
	Token string `env:"TOKEN" envDefault:""`
}

func (t *TelegramConfig) GetTgToken() string {
	return t.Token
}

type ServiceConfig struct {
	CodeSalt   string `env:"CODE_SALT" envDefault:""`
	EncryptKey string `env:"ENCRYPT_KEY" envDefault:""`
	SecretKey  string `env:"SECRET_KEY" envDefault:""`
}

func (s *ServiceConfig) GetCodeSalt() string {
	return s.CodeSalt
}

func (s *ServiceConfig) GetEncryptKey() string {
	return s.EncryptKey
}

func (s *ServiceConfig) GetSecretKey() string {
	return s.SecretKey
}
