package cmd

import (
	"fmt"
	"os"

	"io/ioutil"
	"net/url"
	"path"

	"github.com/ChimeraCoder/anaconda"
	"github.com/mpppk/unravel-twitter/etc"
	"github.com/mpppk/unravel-twitter/twitter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var cfgFile string
var consumerKey string
var consumerSecret string
var screenName string
var maxId string
var accessToken string
var accessTokenSecret string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "unravel-twitter",
	Short: "A brief description of your application",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := etc.LoadConfigFromFile()
		if err != nil {
			fmt.Println(err)
		}

		anaconda.SetConsumerKey(config.ConsumerKey)
		anaconda.SetConsumerSecret(config.ConsumerSecret)

		api := anaconda.NewTwitterApi(config.AccessToken, config.AccessTokenSecret)
		api.SetLogger(anaconda.BasicLogger) // logger を設定

		tweets, err := api.GetUserTimeline(url.Values{
			"screen_name":     []string{config.ScreenName},
			"count":           []string{"200"},
			"exclude_replies": []string{"true"},
			"trim_user":       []string{"true"},
			"include_rts":     []string{"false"},
			"max_id":          []string{config.MaxId}})

		if err != nil {
			fmt.Println("GetUserTimeline error")
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
			Source:    config.ScreenName,
			List:      tweetImages,
		}

		d, err := yaml.Marshal(&imageMetaData)
		if err != nil {
			panic(err)
		}

		ioutil.WriteFile("test.yaml", d, 0777)
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.unravel-worker.yaml)")
	RootCmd.PersistentFlags().StringVar(&consumerKey, "consumer-key", "", "twitter consumer key")
	viper.BindPFlag("twitter.consumerKey", RootCmd.PersistentFlags().Lookup("consumer-key"))
	RootCmd.PersistentFlags().StringVar(&consumerSecret, "consumer-secret", "", "twitter consumer key")
	viper.BindPFlag("twitter.consumerSecret", RootCmd.PersistentFlags().Lookup("consumer-secret"))
	RootCmd.PersistentFlags().StringVar(&screenName, "screen-name", "", "")
	viper.BindPFlag("twitter.screenName", RootCmd.PersistentFlags().Lookup("screen-name"))
	RootCmd.PersistentFlags().StringVar(&maxId, "max-id", "", "")
	viper.BindPFlag("twitter.maxId", RootCmd.PersistentFlags().Lookup("max-id"))
	RootCmd.PersistentFlags().StringVar(&accessToken, "access-token", "", "")
	viper.BindPFlag("twitter.accessToken", RootCmd.PersistentFlags().Lookup("access-token"))
	RootCmd.PersistentFlags().StringVar(&accessTokenSecret, "access-token-secret", "", "")
	viper.BindPFlag("twitter.accessTokenSecret", RootCmd.PersistentFlags().Lookup("access-token-secret"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".worker")                                          // name of config file (without extension)
	viper.AddConfigPath(path.Join(os.Getenv("HOME"), ".config", "unravel")) // adding home directory as first search path
	viper.AutomaticEnv()                                                    // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
