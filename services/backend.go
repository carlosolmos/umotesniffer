package services

import (
	"github.com/carlosolmos/umotesniffer/comms"
	log "github.com/sirupsen/logrus"
)

type Backend struct {
	RUmote    *comms.Client
	RChan     chan []byte
	LUmote    *comms.Client
	LChan     chan []byte
	debugMode bool
}

func NewBackend(debug bool, RHost, RAlias, LHost, LAlias string) (*Backend, error) {
	log.Debug("New Backend")
	b := &Backend{
		debugMode: debug,
	}
	b.RChan = make(chan []byte)
	b.LChan = make(chan []byte)

	var err error
	b.RUmote, err = comms.ConnectClient(RHost, RAlias, b.RChan)
	if err != nil {
		return nil, err
	}
	go b.RUmote.Receive()

	b.LUmote, err = comms.ConnectClient(LHost, LAlias, b.LChan)
	if err != nil {
		return nil, err
	}
	go b.LUmote.Receive()

	return b, nil
}

func (b *Backend) Run() {
	for {
		select {
		case buffer := <-b.RChan:
			log.Info("RHost", buffer)
		case buffer := <-b.LChan:
			log.Info("LHost", buffer)
		}
	}
}
