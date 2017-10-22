package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/mpppk/unravel-twitter/etc"
)

func CreateClient(config *etc.Config) *anaconda.TwitterApi {
	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)

	api := anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)
	api.SetLogger(anaconda.BasicLogger) // logger を設定
	return api
}
