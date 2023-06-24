package services

import (
	"bytes"
	"github.com/aler9/gomavlib/pkg/frame"
	log "github.com/sirupsen/logrus"
	"io"
)

type readWriter struct {
	io.Reader
	io.Writer
}

type MavlinkDecoder struct {
	inBuf  *bytes.Buffer
	outBuf *bytes.Buffer
	rw     *frame.ReadWriter
}

var MavDecoder *MavlinkDecoder

func NewMavlinkDecoder() *MavlinkDecoder {
	mavDecoder := MavlinkDecoder{
		inBuf:  bytes.NewBuffer([]byte{}),
		outBuf: bytes.NewBuffer([]byte{}),
	}
	var err error
	mavDecoder.rw, err = frame.NewReadWriter(frame.ReadWriterConf{
		ReadWriter: &readWriter{
			Reader: mavDecoder.inBuf,
			Writer: mavDecoder.outBuf,
		},
		DialectRW:   nil,
		OutVersion:  frame.V2, // change to V1 if you're unable to communicate with the target
		OutSystemID: 10,
	})
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return &mavDecoder
}

func (mavDecoder *MavlinkDecoder) mavlinkDecode(data []byte) error {
	// read a frame, that contains a message
	mavDecoder.inBuf.Reset()
	mavDecoder.inBuf.Write(data)

	mavFrame, err := mavDecoder.rw.Read()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.Printf("MAVLINK: id=%d, %+v\n", mavFrame.GetMessage().GetID(), mavFrame.GetMessage())
	return nil
}

func DecodeBinaryMessage(data []byte) *CotMessageInfo {
	log.Println("binary message")
	msg := &CotMessageInfo{}
	/*if MavDecoder == nil {
		MavDecoder = NewMavlinkDecoder()
	}
	err := MavDecoder.mavlinkDecode(data)
	if err != nil {
		return msg
	}*/
	return msg
}
