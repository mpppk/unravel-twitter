package main

import (
	"net/url"

	"io/ioutil"

	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
	"github.com/mpppk/unravel-twitter/twitter"
	"gopkg.in/yaml.v2"
)

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

	tweetImages := []twitter.MetaData{}
	for _, tweet := range tweets {
		for _, media := range tweet.Entities.Media {
			tweetImages = append(tweetImages, &twitter.TweetImageMetaData{
				Id:   tweet.Id,
				Url:  media.Media_url,
				Text: tweet.Text,
			})
		}
	}

	imageMetaData := twitter.MetaDataSet{
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
