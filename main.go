package main

import (
	"net/url"

	"io/ioutil"

	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type TweetImageMetaData struct {
	Id   int64
	Url  string
	Text string
}

func (t *TweetImageMetaData) GetId() int64 {
	return t.Id
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

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))

	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"))
	api.SetLogger(anaconda.BasicLogger) // logger を設定

	SCREEN_NAME := os.Getenv("TWITTER_SCREEN_NAME")
	MAX_ID := os.Getenv("TWITTER_MAX_ID")

	tweets, err := api.GetUserTimeline(url.Values{
		"screen_name":     []string{SCREEN_NAME},
		"count":           []string{"200"},
		"exclude_replies": []string{"true"},
		"trim_user":       []string{"true"},
		"include_rts":     []string{"false"},
		"max_id":          []string{MAX_ID}})

	if err != nil {
		panic(err)
	}

	tweetImages := []MetaData{}
	for _, tweet := range tweets {
		for _, media := range tweet.Entities.Media {
			tweetImages = append(tweetImages, &TweetImageMetaData{
				Id:   tweet.Id,
				Url:  media.Media_url,
				Text: tweet.Text,
			})
		}
	}

	imageMetaData := MetaDataSet{
		MediaType: "twitter",
		Source:    SCREEN_NAME,
		List:      tweetImages,
	}

	d, err := yaml.Marshal(&imageMetaData)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile("test.yaml", d, 0777)
}
