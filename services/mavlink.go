package services

import (
	"bytes"
	"fmt"
	"github.com/aler9/gomavlib/pkg/dialect"
	"github.com/aler9/gomavlib/pkg/dialects/standard"
	"github.com/aler9/gomavlib/pkg/frame"
	"time"

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
var mavlinkCounter uint64

func NewMavlinkDecoder() *MavlinkDecoder {
	mavDecoder := MavlinkDecoder{
		inBuf:  bytes.NewBuffer([]byte{}),
		outBuf: bytes.NewBuffer([]byte{}),
	}
	dialectRW, err := dialect.NewReadWriter((*dialect.Dialect)(standard.Dialect))
	if err != nil {
		panic(err)
	}

	mavDecoder.rw, err = frame.NewReadWriter(frame.ReadWriterConf{
		ReadWriter: &readWriter{
			Reader: mavDecoder.inBuf,
			Writer: mavDecoder.outBuf,
		},
		DialectRW:   dialectRW,
		OutVersion:  frame.V1, // change to V1 if you're unable to communicate with the target
		OutSystemID: 10,
	})
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return &mavDecoder
}

func (mavDecoder *MavlinkDecoder) mavlinkDecode(data []byte) (frame.Frame, error) {
	// read a frame, that contains a message
	mavDecoder.inBuf.Reset()
	mavDecoder.inBuf.Write(data)

	mavFrame, err := mavDecoder.rw.Read()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	//log.Printf("MAVLINK: id=%d, %+v\n", mavFrame.GetMessage().GetID(), mavFrame.GetMessage())
	return mavFrame, nil
}

func DecodeBinaryMessage(data []byte) *CotMessageInfo {
	msg := &CotMessageInfo{
		Type:          MAVLINK,
		Size:          len(data),
		Uid:           "UAV",
		Timestamp:     time.Now().Format("15:04:05"),
		UnixTimestamp: time.Now().Unix(),
		Level:         "",
		Ppower:        "",
	}
	if MavDecoder == nil {
		MavDecoder = NewMavlinkDecoder()
		mavlinkCounter = 0
	}
	mavlinkCounter++
	mavframe, err := MavDecoder.mavlinkDecode(data)
	if err != nil || mavframe == nil {
		return msg
	}
	msg.Remarks = map[string]string{}

	msg.Remarks["m"] = fmt.Sprintf("%+v", mavframe.GetMessage())
	msg.Sequence = fmt.Sprintf("%d", mavlinkCounter)
	return msg
}
