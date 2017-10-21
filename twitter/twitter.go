package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/jinzhu/gorm"
	"github.com/mpppk/unravel-twitter/etc"
)

type Image struct {
	gorm.Model
	Url         string
	Description string
	Labels      []Label `gorm:"many2many:imageLabel;"`
}

type Label struct {
	gorm.Model
	Name   string
	Value  string
	Images []Image `gorm:"many2many:imageLabel;"`
}

func CreateClient(config *etc.Config) *anaconda.TwitterApi {
	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)

	api := anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)
	api.SetLogger(anaconda.BasicLogger) // logger を設定
	return api
}
