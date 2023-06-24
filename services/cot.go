package services

import (
	"fmt"
	"github.com/sbabiv/xml2map"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

const SYSTEM_STATUS = "STAT"
const TRANSPORT = "TCP"
const MAVLINK = "MAV"
const UNKNOWN = "UNKNOWN"

type CotMessageInfo struct {
	Uid           string
	Timestamp     string
	UnixTimestamp int64
	Type          string
	Origin        string
	Source        string
	Destination   string
	Size          int
	Remarks       map[string]string
	FlowTags      []string
	Sequence      string
}

func DecodeCot2Map(data string) (map[string]interface{}, error) {
	decoder := xml2map.NewDecoder(strings.NewReader(data))
	result, err := decoder.Decode()
	return result, err
}

func DecodeCotMessage(data string) *CotMessageInfo {
	msg := &CotMessageInfo{}
	cotMap, err := DecodeCot2Map(data)
	if err != nil {
		log.Errorf("error decoding cot %v", err)
		return nil
	}

	v, okay := cotMap["event"].(map[string]interface{})
	if !okay {
		log.Warning("can't parse ", data)
		return nil
	}
	if v == nil {
		return nil
	}
	msg.Size = len(data)
	msg.Uid = v["@uid"].(string)
	msg.Timestamp = v["@time"].(string)
	msg.Timestamp = msg.Timestamp[:19]

	ts, err := time.Parse("2006-01-02T15:04:05", msg.Timestamp)
	if err != nil {
		msg.UnixTimestamp = ts.Unix()
	}

	if det, ok := v["detail"].(map[string]interface{}); ok {
		if det != nil {
			if sniffer, ok := det["sniffer"].(string); ok {
				//type=tx_ack|addr=%s|id=%d|seq=%d|base=%d
				msg.Type = TRANSPORT
				msg.Remarks = make(map[string]string)
				toks := strings.Split(sniffer, "|")
				for _, v := range toks {
					if strings.Contains(v, "=") {
						kvs := strings.Split(v, "=")
						if len(kvs) == 2 {
							msg.Remarks[kvs[0]] = kvs[1]
						}
					}
				}
				if typeSniff, ok := msg.Remarks["type"]; ok && len(typeSniff) > 0 {
					msg.Type = typeSniff
					if msg.Type == "RXACK" {
						msg.Origin = msg.Remarks["addr"]
					} else {
						msg.Destination = msg.Remarks["addr"]
					}

					msg.Sequence = fmt.Sprintf("id=%s s=%s b=%s", msg.Remarks["id"], msg.Remarks["seq"], msg.Remarks["base"])
				}
			} else if remarks, ok := det["remarks"].(string); ok {
				msg.Type = SYSTEM_STATUS
				msg.Remarks = make(map[string]string)
				toks := strings.Split(remarks, ",")
				for _, v := range toks {
					if strings.Contains(v, "=") {
						kvs := strings.Split(v, "=")
						if len(kvs) == 2 {
							msg.Remarks[kvs[0]] = kvs[1]
						}
					}
				}
				if orig, ok := msg.Remarks["XLR_ORI"]; ok {
					msg.Origin = orig
				}
				if dst, ok := msg.Remarks["XLR_GW"]; ok {
					msg.Destination = dst
				}
				if seq, ok := msg.Remarks["SEQ"]; ok {
					msg.Sequence = seq
				}
			} else if payload, ok := det["payload"].(string); ok {
				msg.Type = "DUMMY"
				toks := strings.Split(payload, "|")
				if len(toks) > 0 {
					msg.Sequence = toks[0]
				}
			}

			if flowt, ok := det["_flow-tags_"].(map[string]interface{}); ok {
				msg.FlowTags = make([]string, 0)
				for k, _ := range flowt {
					if len(k) > 0 && k != "#text" {
						msg.FlowTags = append(msg.FlowTags, k)
					}
				}
			}
		}
	}
	if msg.Type == "" {
		msg.Type = UNKNOWN
	}

	return msg
}

/*

v := result["container"].
        (map[string]interface{})["cats"].
            (map[string]interface{})["cat"].
                ([]map[string]interface{})[0]["items"].
                    (map[string]interface{})["n"].([]string)[1]
*/
