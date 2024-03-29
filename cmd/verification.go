package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"umotesniffer/services"
)
import log "github.com/sirupsen/logrus"

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Run verification",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("launching verification")
		backend, err := services.NewBackend(debugMode,
			viper.GetString("TopHost"),
			viper.GetString("TopAlias"),
			viper.GetString("BottomHost"),
			viper.GetString("BottomAlias"),
		)
		if err != nil {
			log.Fatal("error initializing backend ", err.Error())
			return
		}
		backend.Run()
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "debug mode")
}
