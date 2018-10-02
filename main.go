package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	namespace = "ethminer"
)

var (
	totalhashrate = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "totalhashrate"),
		"Was the last query of pids check successful.",
		nil, nil,
	)
)

type EthminerAPI struct {
	Id      int    `json:"id"`
	JsonRPC string `json:"jsonrpc"`
	Result  Result `json:"result"`
	Error   Error  `json:"error"`
}

type Result struct {
	PoolAddress    string    `json:"pooladdrs"`
	PoolSw         int       `json:"ethpoolsw"`
	TotalHashRate  int       `json:"ethhashrate"`
	HashRates      []int     `json:"ethhashrates"`
	Shares         int       `json:"ethshares"`
	Invalid        int       `json:"ethinvalid"`
	Rejected       int       `json:"ethrejected"`
	PowerUsages    []float32 `json:"powerusages"`
	Temperatures   []int     `json:"temperatures"`
	FanPercentages []int     `json:"fanpercentages"`
	Runtime        string    `json:"runtime"`
	Version        string    `json:"version"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ethminerExporter struct {
	target        string
	totalhashrate prometheus.Gauge
	hashrates     []prometheus.Gauge
}

func newEthminerExporter(target string) (*ethminerExporter, error) {
	return &ethminerExporter{
		target: target,
		totalhashrate: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "totalhashrate",
				Help:      "Hashrate [H/s]",
			}),
	}, nil
}

func (e *ethminerExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.totalhashrate.Desc()
}

func (e *ethminerExporter) Collect(ch chan<- prometheus.Metric) {
	conn, err := net.Dial("tcp", e.target)
	if err != nil {
		log.Errorln(err)
		ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("connection_error", "Error connecting to target", nil, nil), err)
		return
	}
	defer conn.Close()

	message := "{\"id\":0, \"jsonrpc\": \"2.0\", \"method\":\"miner_getstathr\"}\n"
	conn.Write([]byte(message))
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}

	ethstats := new(EthminerAPI)
	if err := json.Unmarshal(buf[:n], ethstats); err != nil {
		log.Errorln(err)
	}
	if ethstats.Error.Code != 0 {
		log.Errorln(ethstats.Error.Message)
	}

	ch <- prometheus.MustNewConstMetric(
		e.totalhashrate.Desc(),
		prometheus.GaugeValue,
		float64(ethstats.Result.TotalHashRate),
	)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	if target == "" {
		http.Error(w, fmt.Sprintf("target is not specified."), http.StatusBadRequest)
		return
	}
	e, _ := newEthminerExporter(target)

	registry := prometheus.NewRegistry()
	registry.MustRegister(e)

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func main() {
	listenAddress := kingpin.Flag("listen", "Address to listen on for web interface and telemetry.").Default("0.0.0.0:8555").String()
	kingpin.Version(version.Print("ethminer_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Infoln("Starting ethminer exporter")
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metricsHandler(w, r)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
        <head><title>ethminer Exporter</title></head>
        <body>
        <h1>ethminer Exporter</h1>
        <p><a href='/metrics'>Metrics</a></p>
        </body>
        </html>`))
	})

	log.Infoln("Listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
