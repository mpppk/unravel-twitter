package twitter

import (
	"net/url"

	"fmt"

	"time"

	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/mpppk/unravel-twitter-worker/adapter"
)

type Config struct {
	ScreenName        string
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
	SinceDate         time.Time
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
	config         *Config
	unravelAdapter *adapter.Adapter
	beforeMaxId    int64
}

func (c *Crawler) Fetch() ([]anaconda.Tweet, bool, error) {
	values := url.Values{
		"screen_name":     []string{c.config.ScreenName},
		"count":           []string{"200"},
		"exclude_replies": []string{"true"},
		"trim_user":       []string{"true"},
		"include_rts":     []string{"false"},
	}

	if c.beforeMaxId != -1 {
		values["max_id"] = []string{fmt.Sprint(c.beforeMaxId)}
	}

	id, ok, err := c.GetLatestSavedTweetId()
	if err != nil {
		return nil, false, err
	}

	if ok {
		values["since_id"] = []string{fmt.Sprint(id)}
	}

	tweets, err := c.client.GetUserTimeline(values)
	if err != nil {
		return nil, false, err
	}

	layout := "Mon Jan 2 15:04:05 -0700 2006"
	var retTweets []anaconda.Tweet

	for _, tweet := range tweets {
		date, err := time.Parse(layout, tweet.CreatedAt)
		if err != nil {
			return nil, false, err
		}
		timeDiff := date.Sub(c.config.SinceDate).Nanoseconds()
		if timeDiff < 0 {
			return retTweets, false, nil
		}
		retTweets = append(retTweets, tweet)
		if c.beforeMaxId > tweet.Id || c.beforeMaxId == -1 {
			c.beforeMaxId = tweet.Id - 1
		}
	}
	return retTweets, true, err
}

func (c *Crawler) SaveTweet(tweet anaconda.Tweet) error {
	for _, media := range tweet.Entities.Media {
		image := &adapter.Image{
			Url:         media.Media_url,
			Description: tweet.Text,
		}

		err := c.unravelAdapter.AddLabelsToImage(image, []adapter.NewLabel{
			{Name: "twitter"},
			{Name: "twitter_id", Value: fmt.Sprint(tweet.Id)},
			{Name: "post_date", Value: fmt.Sprint(tweet.CreatedAt)},
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Crawler) FetchAndSave() error {
	for {
		tweets, hasNext, err := c.Fetch()
		if err != nil {
			return err
		}

		for _, tweet := range tweets {
			err := c.SaveTweet(tweet)
			if err != nil {
				return err
			}
		}

		if len(tweets) <= 0 || !hasNext {
			return nil
		}
	}
	return nil
}

func (c *Crawler) GetLatestSavedTweetId() (uint, bool, error) {
	image, err := c.unravelAdapter.FindByMaxLabelValue("twitter_id")
	if err != nil {
		return 0, false, err
	}

	for _, label := range image.Labels {
		if label.Name == "twitter_id" {
			v, err := strconv.ParseUint(label.Value, 10, 64)
			return uint(v), true, err
		}
	}

	return 0, false, nil
}

func (c *Crawler) Close() {
	c.unravelAdapter.Close()
}

func NewCrawler(config *Config) (*Crawler, error) {
	client := CreateClient(config)
	adpt, err := adapter.New(false)

	if err != nil {
		panic(err)
	}

	return &Crawler{
		client:         client,
		config:         config,
		unravelAdapter: adpt,
		beforeMaxId:    -1,
	}, err
}
