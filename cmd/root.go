/*
Copyright © 2019 iot.ai

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var debugMode bool
var cfgFile string
var Version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "umotesniffer",
	Short: "uMote Diagnostics",
	Long:  `Tool to probe the umote network.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/umotesniffer.properties)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name "umotesniffer.properties" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("umotesniffer.properties")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	log.SetFormatter(&log.JSONFormatter{})
	if viper.Get("log_out") == "file" {
		file, err := os.OpenFile("umotesniffer.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(file)
		} else {
			log.SetOutput(os.Stderr)
			log.Info("Failed to log to file, using default stderr")
		}
	}
	switch viper.Get("log_level") {
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}
	log.Debug("initialization complete")

}
