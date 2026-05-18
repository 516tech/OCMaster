package config

import (
	"os"
	"path/filepath"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string
	DBType     string // "sqlite" or "mysql"
	DBPath     string // sqlite file path
	DBDSN      string // mysql DSN
	JWTSecret  string
}

func Load() *Config {
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("db.type", "sqlite")
	viper.SetDefault("db.path", "./data/ocmaster.db")
	viper.SetDefault("db.dsn", "root:ocmaster123@tcp(127.0.0.1:3306)/ocmaster?charset=utf8mb4&parseTime=True&loc=Local")
	viper.SetDefault("jwt.secret", "ocmaster-jwt-secret-dev")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	return &Config{
		ServerPort: viper.GetString("server.port"),
		DBType:     viper.GetString("db.type"),
		DBPath:     viper.GetString("db.path"),
		DBDSN:      viper.GetString("db.dsn"),
		JWTSecret:  viper.GetString("jwt.secret"),
	}
}

func (c *Config) OpenDB() (*gorm.DB, error) {
	if c.DBType == "mysql" {
		return gorm.Open(mysql.Open(c.DBDSN), &gorm.Config{})
	}
	// sqlite
	dir := filepath.Dir(c.DBPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return gorm.Open(sqlite.Open(c.DBPath), &gorm.Config{})
}
