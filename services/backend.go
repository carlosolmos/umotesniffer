package services

import (
	"github.com/carlosolmos/umotesniffer/comms"
	log "github.com/sirupsen/logrus"
)

type Backend struct {
	TopUmote    *comms.Client
	TopChan     chan []byte
	BottomUmote *comms.Client
	BottomChan  chan []byte
	debugMode   bool
}

func NewBackend(debug bool, TopHost, RAlias, BottomHost, LAlias string) (*Backend, error) {
	log.Debug("New Backend")
	b := &Backend{
		debugMode: debug,
	}
	b.TopChan = make(chan []byte)
	b.BottomChan = make(chan []byte)

	var err error
	if len(TopHost) > 0 {
		b.TopUmote, err = comms.ConnectClient(TopHost, RAlias, b.TopChan)
		if err != nil {
			return nil, err
		}
		go b.TopUmote.Receive()
	}
	if len(BottomHost) > 0 {
		b.BottomUmote, err = comms.ConnectClient(BottomHost, LAlias, b.BottomChan)
		if err != nil {
			return nil, err
		}
		go b.BottomUmote.Receive()
	}
	return b, nil
}

func (b *Backend) Run() {
	for {
		select {
		case buffer := <-b.TopChan:
			log.Info("TopHost", buffer)
		case buffer := <-b.BottomChan:
			log.Info("BottomHost", buffer)
		}
	}
}

func (b *Backend) Shutdown() {
	b.TopUmote.Close()
	b.BottomUmote.Close()
}
