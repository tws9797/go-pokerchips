package config

import "github.com/spf13/viper"

type Config struct {
	DBUri    string `mapstructure:"MONGO_URI"`
	RedisUri string `mapstructure:"REDIS_URI"`
	Port     string `mapstructure:"PORT"`
}

func LoadConfig(path string) (config Config, err error) {

	viper.AddConfigPath(path)  // path to look for the config file, in this case the current working dir
	viper.SetConfigType("env") // name of the config file (without extension)
	viper.SetConfigName("app") // REQUIRED if the config file does not have the extension in the name

	viper.AutomaticEnv()

	err = viper.ReadInConfig() // Find and read the config file

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
