package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/jinzhu/gorm"
	"github.com/mpppk/unravel-twitter/etc"
)

type TweetImageMetaData struct {
	gorm.Model
	MediaNo int64
	Url     string
	Text    string
}

func (t *TweetImageMetaData) GetId() int64 {
	return t.MediaNo
}

func (t *TweetImageMetaData) GetUrl() string {
	return t.Url
}

func (t *TweetImageMetaData) GetText() string {
	return t.Text
}

type MetaData interface {
	GetId() int64
	GetUrl() string
	GetText() string
}

type MetaDataSet struct {
	MediaType string
	Source    string
	List      []MetaData
}

func CreateClient(config *etc.Config) *anaconda.TwitterApi {
	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)

	api := anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)
	api.SetLogger(anaconda.BasicLogger) // logger を設定
	return api
}
