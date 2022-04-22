package cmd

import (
	"github.com/carlosolmos/umotesniffer/services"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)
import log "github.com/sirupsen/logrus"

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Run verification",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("launching verification")
		backend, err := services.NewBackend(debugMode,
			viper.GetString("RHost"),
			"RHost",
			viper.GetString("LHost"),
			"LHost",
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
