package etc

import "github.com/spf13/viper"

type Config struct {
	ConsumerKey       string
	ConsumerSecret    string
	ScreenName        string
	MaxId             string
	AccessToken       string
	AccessTokenSecret string
}

func LoadConfigFromFile() (*Config, error) {
	var config Config
	if err := viper.Sub("twitter").Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
