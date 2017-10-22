package twitter

import (
	"net/url"

	"fmt"

	"github.com/ChimeraCoder/anaconda"
)

type Config struct {
	ScreenName        string
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

func CreateClient(config *Config) *anaconda.TwitterApi {
	anaconda.SetConsumerKey(config.ConsumerKey)
	anaconda.SetConsumerSecret(config.ConsumerSecret)

	api := anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)
	api.SetLogger(anaconda.BasicLogger) // logger を設定
	return api
}

type Crawler struct {
	client     *anaconda.TwitterApi
	screenName string
}

func (c *Crawler) Fetch(maxId int64) ([]anaconda.Tweet, error) {
	fmt.Println(c.screenName)
	return c.client.GetUserTimeline(url.Values{
		"screen_name":     []string{c.screenName},
		"count":           []string{"200"},
		"exclude_replies": []string{"true"},
		"trim_user":       []string{"true"},
		"include_rts":     []string{"false"},
		"max_id":          []string{fmt.Sprint(maxId)}})
}

func NewCrawler(config *Config) *Crawler {
	client := CreateClient(config)
	return &Crawler{client: client, screenName: config.ScreenName}
}
