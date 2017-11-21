package cmd

import (
	"fmt"
	"os"

	"path"

	"time"

	"github.com/mpppk/unravel-twitter-worker/etc"
	"github.com/mpppk/unravel-twitter-worker/twitter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	Short: "twitter crawler for unravel",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := etc.LoadConfigFromFile()
		if err != nil {
			fmt.Println(err)
		}

		crawler, err := twitter.NewCrawler(&twitter.Config{
			ScreenName:        screenName,
			ConsumerKey:       config.ConsumerKey,
			ConsumerSecret:    config.ConsumerSecret,
			AccessToken:       config.AccessToken,
			AccessTokenSecret: config.AccessTokenSecret,
			SinceDate:         time.Now().Add(-time.Duration(24*30*12) * time.Hour),
		})
		if err != nil {
			panic(err)
		}

		defer crawler.Close()

		err = crawler.FetchAndSave()
		if err != nil {
			panic(err)
		}
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
