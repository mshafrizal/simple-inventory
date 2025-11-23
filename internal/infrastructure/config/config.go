package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

func LoadConfig(path string) (*Config, error) {
	viper.SetDefault("APP_PORT", 8080)
	viper.SetDefault("APP_ENV", "production")
	viper.SetDefault("APP_NAME", "simple-inventory")
	viper.SetDefault("DB_HOST", "9qasp5v56q8ckkf5dc.apn.leapcellpool.com")
	viper.SetDefault("DB_PORT", 6438)
	viper.SetDefault("DB_USER", "itowpiunrzibboonxyqu")
	viper.SetDefault("DB_PASSWORD", "qsxjreudpuoufqtiuyyibdgigrregr")
	viper.SetDefault("DB_NAME", "kpipmpgdtiahyjkllzrl")
	viper.SetDefault("DB_SSLMODE", false)
	viper.SetDefault("JWT_SECRET", "2f2262daf4795e4157771dc4188a33ea541caa90032719f6241b1661e1901457")
	viper.SetDefault("JWT_EXPIRATION_HOURS", 24)

	// Try to load .env file for local development
	// In production (Leapcell), environment variables will be set directly
	if path != "" {
		viper.SetConfigFile(path)
		viper.SetConfigType("env")
		if err := viper.ReadInConfig(); err != nil {
			// Don't fail if .env file is missing - fall back to environment variables
			fmt.Printf("No .env file found at %s, using environment variables\n", path)
		}
	}
	viper.AutomaticEnv()

	config := &Config{
		App: AppConfig{
			Name: viper.GetString("APP_NAME"),
			Env:  viper.GetString("APP_ENV"),
			Port: viper.GetString("APP_PORT"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		JWT: JWTConfig{
			Secret:          viper.GetString("JWT_SECRET"),
			ExpirationHours: viper.GetInt("JWT_EXPIRATION_HOURS"),
		},
	}

	if config.JWT.ExpirationHours == 0 {
		config.JWT.ExpirationHours = 24
	}

	return config, nil
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}
