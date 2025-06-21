package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	WebServerPort     string `mapstructure:"WEB_SERVER_PORT"`
	GRPCServerPort    string `mapstructure:"GRPC_SERVER_PORT"`
	GraphQLServerPort string `mapstructure:"GRAPHQL_SERVER_PORT"`
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv() // Pega variáveis de ambiente

	config := &Config{
		DBDriver:          viper.GetString("DB_DRIVER"),
		DBHost:            viper.GetString("DB_HOST"),
		DBPort:            viper.GetString("DB_PORT"),
		DBUser:            viper.GetString("DB_USER"),
		DBPassword:        viper.GetString("DB_PASSWORD"),
		DBName:            viper.GetString("DB_NAME"),
		WebServerPort:     viper.GetString("WEB_SERVER_PORT"),
		GRPCServerPort:    viper.GetString("GRPC_SERVER_PORT"),
		GraphQLServerPort: viper.GetString("GRAPHQL_SERVER_PORT"),
	}

	fmt.Printf("Config loaded: %+v\n", config)

	return config, nil
}
