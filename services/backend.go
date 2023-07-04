package services

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"umotesniffer/comms"
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
	PromMetrics = NewMetrics()
	go b.StartPrometheus()
	var err error
	if len(TopHost) > 0 {
		b.TopUmote, err = comms.ConnectClient(TopHost, RAlias, b.TopChan)
		if err != nil {
			return nil, err
		}
		log.Println("Connected TOP: %s", TopHost)
		go b.TopUmote.Receive()
	}
	if len(BottomHost) > 0 {
		b.BottomUmote, err = comms.ConnectClient(BottomHost, LAlias, b.BottomChan)
		if err != nil {
			return nil, err
		}
		log.Println("Connected Bottom: %s", BottomHost)
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

func (b *Backend) StartPrometheus() {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":2112", nil)
	if err != nil {
		log.Errorf("error starting prometheus endpoint: %s", err.Error())
		return
	}
}
