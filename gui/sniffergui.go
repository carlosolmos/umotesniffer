package gui

import (
	"github.com/carlosolmos/umotesniffer/services"
	ui "github.com/gizak/termui/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

const UPDATE_INTERVAL = 1

var umotesBackend *services.Backend

// field and gw panels
func RenderFWgui(debugMode bool) {
	var err error
	umotesBackend, err = services.NewBackend(debugMode,
		viper.GetString("TopHost"),
		viper.GetString("TopAlias"),
		viper.GetString("BottomHost"),
		viper.GetString("BottomAlias"),
	)
	if err != nil {
		log.Fatalf("error initializing backend  %v", err)
		return
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	UmoteTop := NewUmoteTable(viper.GetString("TopAlias"), 0, 0)
	UmoteBottom := NewUmoteTable(viper.GetString("BottomAlias"), 0, UMTABLE_H+1)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second * UPDATE_INTERVAL).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				umotesBackend.Shutdown()
				return
			}
		case <-ticker:
			ui.Render(UmoteTop.UmTable, UmoteBottom.UmTable)
		case buffer := <-umotesBackend.TopChan:
			UmoteTop.UpdateUmoteTable(buffer)
		case buffer := <-umotesBackend.BottomChan:
			UmoteBottom.UpdateUmoteTable(buffer)
		}
	}
}
