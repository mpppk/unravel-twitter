package twitter

import (
	"net/url"

	"fmt"

	"github.com/ChimeraCoder/anaconda"
	"github.com/mpppk/unravel-twitter/adapter"
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
	client         *anaconda.TwitterApi
	screenName     string
	unravelAdapter *adapter.Adapter
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

func (c *Crawler) SaveTweet(tweet anaconda.Tweet) error {
	for _, media := range tweet.Entities.Media {
		image := &adapter.Image{
			Url:         media.Media_url,
			Description: tweet.Text,
		}

		err := c.unravelAdapter.AddLabelsToImage(image, []adapter.NewLabel{
			{Name: "twitter"},
			{Name: "twitterid", Value: fmt.Sprint(tweet.Id)},
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Crawler) FetchAndSave(maxId int64) error {
	tweets, err := c.Fetch(maxId)
	if err != nil {
		return err
	}

	for _, tweet := range tweets {
		err := c.SaveTweet(tweet)
		return err
	}
	return nil
}

func (c *Crawler) Close() {
	c.unravelAdapter.Close()
}

func NewCrawler(config *Config) (*Crawler, error) {
	client := CreateClient(config)
	adpt, err := adapter.New(false)
	return &Crawler{
		client:         client,
		screenName:     config.ScreenName,
		unravelAdapter: adpt,
	}, err
}
