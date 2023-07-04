package services

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
	"strings"
	"sync"
)

type Metrics struct {
	mu           sync.Mutex
	TxRetry      *prometheus.GaugeVec
	XlrTx        *prometheus.GaugeVec
	XlrRx        *prometheus.GaugeVec
	PRssi        *prometheus.GaugeVec
	XlrTX_KBps   *prometheus.GaugeVec
	XlrRX_KBps   *prometheus.GaugeVec
	XlrTX_Pps    *prometheus.GaugeVec
	XlrRX_Pps    *prometheus.GaugeVec
	Qlen         *prometheus.GaugeVec
	P_Pow        *prometheus.GaugeVec
	Lvl          *prometheus.GaugeVec
	Errors       *prometheus.GaugeVec
	Sequence     *prometheus.GaugeVec
	MavlinkTotal *prometheus.CounterVec
}

var PromMetrics *Metrics

const (
	METRIC_LBL_UUID = "uuid"
	METRIC_LBL_FW   = "fw"
	METRIC_LBL_MAC  = "mac"
)

func NewMetrics() *Metrics {
	umoteLabels := []string{
		METRIC_LBL_UUID, METRIC_LBL_FW, METRIC_LBL_MAC,
	}
	m := Metrics{}
	m.TxRetry = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "umote_txretry_total",
		Help: "The total number of TX Retry events",
	}, umoteLabels)
	m.XlrTx = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "umote_xlrtx_total",
		Help: "The total number of TX events",
	}, umoteLabels)
	m.XlrRx = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "umote_xlrrx_total",
		Help: "The total number of RX events",
	}, umoteLabels)
	m.PRssi = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "umote_p_rssi_level",
		Help: "Parent RSSI Level",
	}, umoteLabels)
	m.XlrTX_KBps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "umote_xlrtx_kbps",
		Help: "radio tx kbps calculated at the umote (KB per second)",
	}, umoteLabels)
	m.XlrRX_KBps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "umote_xlrrx_kbps",
		Help: "radio rx kbps calculated at the umote (KB per second)",
	}, umoteLabels)
	m.XlrTX_Pps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "umote_xlrtx_pps",
		Help: "radio tx pps calculated at the umote (packets per second)",
	}, umoteLabels)
	m.XlrRX_Pps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "umote_xlrrx_pps",
		Help: "radio rx pps calculated at the umote (packets per second)",
	}, umoteLabels)
	m.Qlen = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "umote_qlen_messages",
		Help: "TX queue size at the umote",
	}, umoteLabels)
	m.P_Pow = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "umote_p_pow_messages",
		Help: "Node parent power level",
	}, umoteLabels)
	m.Lvl = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "umote_lvl_messages",
		Help: "Node mesh level",
	}, umoteLabels)
	m.Sequence = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "umote_seq_total",
		Help: "Message sequence number",
	}, umoteLabels)
	m.Errors = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "umote_errors_total",
		Help: "The total number of Error events",
	}, umoteLabels)
	m.MavlinkTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "umote_mavlink_total",
		Help: "The total number of Mavlink messages",
	}, []string{METRIC_LBL_UUID})
	return &m
}

func (m *Metrics) UpdateMessageMetrics(cotMsg *CotMessageInfo) {
	if cotMsg == nil {
		return
	}
	m.UpdateRemarksMetrics(cotMsg.Uid, cotMsg.Remarks)
}

func (m *Metrics) UpdateRemarksMetrics(uuid string, remarks map[string]string) {
	var mac string
	var fw string
	mac, ok := remarks["XLR_ORI"]
	if !ok {
		mac = "NA"
	}
	fw, ok = remarks["fw_version"]
	if !ok {
		fw = "NA"
	}
	m.mu.Lock()
	for key, value := range remarks {
		key = strings.ToLower(key)
		if n, err := strconv.ParseFloat(value, 64); err == nil {
			switch key {
			case "txretry":
				m.TxRetry.With(prometheus.Labels{
					METRIC_LBL_UUID: uuid, METRIC_LBL_FW: fw, METRIC_LBL_MAC: mac}).Set(n)
				break
			case "xlrtx":
				m.XlrTx.With(prometheus.Labels{
					METRIC_LBL_UUID: uuid, METRIC_LBL_FW: fw, METRIC_LBL_MAC: mac}).Set(n)
				break
			case "xlrrx":
				m.XlrRx.With(prometheus.Labels{
					METRIC_LBL_UUID: uuid, METRIC_LBL_FW: fw, METRIC_LBL_MAC: mac}).Set(n)
				break
			case "p_rssi":
				m.PRssi.With(prometheus.Labels{
					METRIC_LBL_UUID: uuid, METRIC_LBL_FW: fw, METRIC_LBL_MAC: mac}).Set(n)
				break
			case "xlrtx_kbps":
				m.XlrTX_KBps.With(prometheus.Labels{
					METRIC_LBL_UUID: uuid, METRIC_LBL_FW: fw, METRIC_LBL_MAC: mac}).Set(n)
				break
			case "xlrrx_kbps":
				m.XlrRX_KBps.With(prometheus.Labels{
					METRIC_LBL_UUID: uuid, METRIC_LBL_FW: fw, METRIC_LBL_MAC: mac}).Set(n)
				break
			case "xlrtx_pps":
				m.XlrTX_Pps.With(prometheus.Labels{
					METRIC_LBL_UUID: uuid, METRIC_LBL_FW: fw, METRIC_LBL_MAC: mac}).Set(n)
				break
			case "xlrrx_pps":
				m.XlrRX_Pps.With(prometheus.Labels{
					METRIC_LBL_UUID: uuid, METRIC_LBL_FW: fw, METRIC_LBL_MAC: mac}).Set(n)
				break
			case "seq":
				m.Sequence.With(prometheus.Labels{
					METRIC_LBL_UUID: uuid, METRIC_LBL_FW: fw, METRIC_LBL_MAC: mac}).Set(n)
				break
			case "qlen":
				m.Qlen.With(prometheus.Labels{
					METRIC_LBL_UUID: uuid, METRIC_LBL_FW: fw, METRIC_LBL_MAC: mac}).Set(n)
				break
			case "lvl":
				m.Lvl.With(prometheus.Labels{
					METRIC_LBL_UUID: uuid, METRIC_LBL_FW: fw, METRIC_LBL_MAC: mac}).Set(n)
				break
			case "p_pow":
				m.P_Pow.With(prometheus.Labels{
					METRIC_LBL_UUID: uuid, METRIC_LBL_FW: fw, METRIC_LBL_MAC: mac}).Set(n)
				break
			case "errors":
				m.Errors.With(prometheus.Labels{
					METRIC_LBL_UUID: uuid, METRIC_LBL_FW: fw, METRIC_LBL_MAC: mac}).Set(n)
				break
			}

		}

	}
	m.mu.Unlock()
}

func (m *Metrics) IncMavlinkMessages(uuid string) {
	m.MavlinkTotal.With(prometheus.Labels{
		METRIC_LBL_UUID: uuid}).Inc()
}
