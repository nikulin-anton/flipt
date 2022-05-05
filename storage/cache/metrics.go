package cache

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Stats struct {
	ItemCount     int64
	MissTotal     int64
	HitTotal      int64
	EvictionTotal int64
	ErrorTotal    int64
}

// statsGetter is an interface that gets cache.Stats.
type statsGetter interface {
	Stats() Stats
}

const (
	namespace = "flipt"
	subsystem = "cache"
)

//nolint
func registerMetrics(c Cacher) {
	labels := prometheus.Labels{"cache": c.String()}

	collector := &metricsCollector{
		sg: c,
		hitTotalDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "hit_total"),
			"The number of cache hits",
			nil,
			labels,
		),
		missTotalDec: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "miss_total"),
			"The number of cache misses",
			nil,
			labels,
		),
		itemCountDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "item_count"),
			"The number of items currently in the cache",
			nil,
			labels,
		),
		evictionTotalDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "eviction_total"),
			"The number of times an item is evicted from the cache",
			nil,
			labels,
		),
		errorTotalDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "error_total"),
			"The number of times an error occurred reading or writing to the cache",
			nil,
			labels,
		),
	}

	prometheus.MustRegister(collector)
}

type metricsCollector struct {
	sg statsGetter

	hitTotalDesc      *prometheus.Desc
	missTotalDec      *prometheus.Desc
	itemCountDesc     *prometheus.Desc
	evictionTotalDesc *prometheus.Desc
	errorTotalDesc    *prometheus.Desc
}

func (c *metricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.hitTotalDesc
	ch <- c.missTotalDec
	ch <- c.itemCountDesc
	ch <- c.evictionTotalDesc
	ch <- c.errorTotalDesc
}

func (c *metricsCollector) Collect(ch chan<- prometheus.Metric) {
	stats := c.sg.Stats()

	ch <- prometheus.MustNewConstMetric(
		c.hitTotalDesc,
		prometheus.CounterValue,
		float64(stats.HitTotal),
	)
	ch <- prometheus.MustNewConstMetric(
		c.missTotalDec,
		prometheus.CounterValue,
		float64(stats.MissTotal),
	)
	ch <- prometheus.MustNewConstMetric(
		c.itemCountDesc,
		prometheus.GaugeValue,
		float64(stats.ItemCount),
	)
	ch <- prometheus.MustNewConstMetric(
		c.evictionTotalDesc,
		prometheus.CounterValue,
		float64(stats.EvictionTotal),
	)
	ch <- prometheus.MustNewConstMetric(
		c.errorTotalDesc,
		prometheus.CounterValue,
		float64(stats.ErrorTotal),
	)
}
