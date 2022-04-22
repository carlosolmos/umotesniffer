package cmd

import (
	"github.com/carlosolmos/umotesniffer/gui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Launch in GUI mode",
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("launching GUI")
		gui.RenderFWgui(debugMode)
	},
}

func init() {
	rootCmd.AddCommand(guiCmd)
	guiCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "debug mode")
}
