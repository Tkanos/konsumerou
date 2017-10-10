package config

import (
	"os"

	"github.com/spf13/viper"
)

// AppConfig represents the application-side configurations.
type AppConfig struct {
	Port                int    `mapstructure:"APP_PORT"`
	KafkaBrokers        string `mapstructure:"KAFKA_BROKERS"`
	MyServiceKafkaTopic string `mapstructure:"MY_SERVICE_TOPIC"`
	ZipkinURI           string `mapstructure:"ZIPKIN_URI"`
	Verbose             bool   `mapstructure:"VERBOSE"`
}

// Config stores application configuration
var Config AppConfig

// initViper inits viper, reads data from config file if 'development' environement variable is defined or from environement
func initViper() (*viper.Viper, error) {
	vp := viper.New()

	if os.Getenv("ENVIRONMENT") == "DEV" {

		// Set configuration file path and read it if it is present
		vp.SetConfigName("config")
		vp.SetConfigType("toml")
		vp.AddConfigPath(".")

		// Read configuration file if provided.
		if err := vp.ReadInConfig(); err == nil {
			if err := vp.Unmarshal(&Config); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		vp.AutomaticEnv()
	}

	vp.SetDefault("APP_PORT", 8080)
	vp.SetDefault("VERBOSE", false)

	return vp, nil
}

// InitConfig reads data from config file or environment variables and stores it to Config
func InitConfig() (err error) {
	// Configurations
	Config = AppConfig{}
	// set input source for viper either config file or environement variables
	vp, err := initViper()
	if err != nil {
		return err
	}

	Config.Port = vp.GetInt("APP_PORT")
	Config.KafkaBrokers = vp.GetString("KAFKA_BROKERS")
	Config.MyServiceKafkaTopic = vp.GetString("MY_SERVICE_TOPIC")
	Config.ZipkinURI = vp.GetString("ZIPKIN_URI")
	Config.Verbose = vp.GetBool("VERBOSE")

	return
}
