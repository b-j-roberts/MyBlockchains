package metrics

import (
  "github.com/prometheus/client_golang/prometheus"
)

// Create a Prometheus metric to track the number of requests
var RequestCount = prometheus.NewCounter(prometheus.CounterOpts{
  Name: "requests_total",
  Help: "Total number of requests",
})
// Create a Prometheus metric to track the duration of requests
var RequestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
  Name:    "request_duration_seconds",
  Help:    "Duration of requests",
  Buckets: prometheus.LinearBuckets(0.01, 0.05, 20),
})
// Create a Prometheus metric with block height
var Blockheight = prometheus.NewGauge(prometheus.GaugeOpts{
  Name: "block_height",
  Help: "Blockchain current height",
})
// Create a Prometheus metric to track number of blocks mined
var BlocksMined = prometheus.NewCounter(prometheus.CounterOpts{
  Name: "blocks_mined",
  Help: "Total number of blocks mined",
})
// Create a Prometheus metric with status of loading snapshot
var SnapshotLoading = prometheus.NewGauge(prometheus.GaugeOpts{
  Name: "snapshot_loading",
  Help: "Snapshot loading status : 0 - Pre / No Snapshot, 1 - Loading Snapshot, 2 - Verifying Snapshot, 3 - Done",
})
// Create a Prometheus metric to track the number of inbound peers
var InboundPeers = prometheus.NewCounter(prometheus.CounterOpts{
  Name: "inbound_peers",
  Help: "Inbound Peers to node",
})
// Total value transfered? TransactionCount?

func PromSetup() {
  // Register the metrics with Prometheus
  prometheus.MustRegister(RequestCount)
  prometheus.MustRegister(RequestDuration)
  prometheus.MustRegister(Blockheight)
  prometheus.MustRegister(BlocksMined)
  prometheus.MustRegister(SnapshotLoading)
  prometheus.MustRegister(InboundPeers)
}
