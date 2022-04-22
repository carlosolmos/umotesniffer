package comms

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
)

const (
	RX = "RX"
	TX = "TX"
	BUFFER_MAX_SIZE = 1024 * 64
)

type Client struct {
	Alias    string
	Address  string
	socket   net.Conn
	DataChan chan []byte
}

func ConnectClient(targetAddress string, alias string, dataOut chan []byte) (*Client, error) {
	log.Info(fmt.Sprintf("connecting to %s", targetAddress))
	connection, err := net.Dial("tcp", targetAddress)
	if err != nil {
		return nil, err
	}
	log.Info(fmt.Sprintf("connected to %s", targetAddress))

	atakClient := &Client{socket: connection, Alias: alias, Address: targetAddress}
	atakClient.DataChan = dataOut
	return atakClient, err
}

func (client *Client) Receive() {
	for {
		message := make([]byte, BUFFER_MAX_SIZE)
		length, err := client.socket.Read(message)
		if err != nil {
			log.Error(err)
			client.socket.Close()
			os.Exit(1)
		}
		if length > 0 {
			log.Debug(fmt.Sprintf("%s-%s %d bytes", RX, client.Alias, length))
			client.DataChan <- message[:length]
		}
	}
}

func (client *Client) Send(message *string) {
	_, err := client.socket.Write([]byte(*message))
	if err != nil {
		log.Error(err)
	}
}

func (client *Client) Close() {
	_ = client.socket.Close()
}

