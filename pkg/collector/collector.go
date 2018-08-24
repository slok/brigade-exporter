package collector

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/slok/brigade-exporter/pkg/log"
	"github.com/slok/brigade-exporter/pkg/service/brigade"
)

const (
	// namespace identifies all the metrics by the exporter.
	namespace = "brigade"

	// Defaults.
	collectTimeoutDef = 60 * time.Second
)

// Config is the Exporter configuration.
type Config struct {
	// CollectTimeout is the timeout to collect the metrics.
	CollectTimeout time.Duration
	// DisableProjects will disable the project metrics subcollector.
	DisableProjects bool
	// DisableBuilds will disable the builds metrics subcollector.
	DisableBuilds bool
	// DisableJobs will disable the Jobs metrics subcollector.
	DisableJobs bool
}

// defaults sets the required defaults.
func (c *Config) defaults() {
	if c.CollectTimeout == 0 {
		c.CollectTimeout = collectTimeoutDef
	}
}

// Exporter is the main exporter that implements the prometheus.Collector interface
// and executes the other collectors
type Exporter struct {
	scrapeDurationDesc *prometheus.Desc
	scrapeSuccessDesc  *prometheus.Desc

	// Subcollectors.
	subcolls map[string]subcollector

	cfg    Config
	logger log.Logger
}

// NewExporter returns a new exporter.
func NewExporter(cfg Config, brigadeSVC brigade.Interface, logger log.Logger) prometheus.Collector {
	// Fill the required defaults.
	cfg.defaults()

	exporter := &Exporter{
		scrapeDurationDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "exporter", "collector_duration_seconds"),
			"Collector time duration.",
			[]string{"collector"},
			nil,
		),

		scrapeSuccessDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "exporter", "collector_success"),
			"Whether a collector succeeded.",
			[]string{"collector"},
			nil,
		),
		cfg:    cfg,
		logger: logger,
	}

	exporter.initSubcollectors(brigadeSVC)
	return exporter
}

func (e *Exporter) initSubcollectors(brigadeSVC brigade.Interface) {
	e.subcolls = map[string]subcollector{}

	// Generate subcollectors.
	if !e.cfg.DisableProjects {
		e.subcolls["projects"] = NewProject(brigadeSVC, e.logger.With("collector", "projects"))
	} else {
		e.logger.Warnf("projects collector disabled")
	}

	if !e.cfg.DisableBuilds {
		e.subcolls["builds"] = NewBuild(brigadeSVC, e.logger.With("collector", "builds"))
	} else {
		e.logger.Warnf("builds collector disabled")
	}
	if !e.cfg.DisableJobs {
		e.subcolls["jobs"] = NewJob(brigadeSVC, e.logger.With("collector", "jobs"))
	} else {
		e.logger.Warnf("jobs collector disabled")
	}
}

// Describe satisfies prometheus.Collector interface.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.scrapeDurationDesc
	ch <- e.scrapeSuccessDesc

}

// Collect satisfies prometheus.Collector interface.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.logger.Debugf("starting collect")
	var wg sync.WaitGroup

	// Call all the subcollectors.
	wg.Add(len(e.subcolls))
	for scName, sc := range e.subcolls {
		go func(scName string, sc subcollector) {
			defer wg.Done()
			e.subcollect(scName, sc, ch)
		}(scName, sc)
	}

	var wgch = make(chan struct{})
	go func() {
		wg.Wait()
		wgch <- struct{}{}
	}()

	// Wait to all the subscrapes.
	select {
	case <-wgch:
		// TODO: ok.

	case <-time.After(e.cfg.CollectTimeout):
		e.logger.Errorf("timeout collecting metrics")
		// TODO timeout.
	}

	e.logger.Debugf("finished collect")
}

func (e *Exporter) subcollect(scName string, sc subcollector, ch chan<- prometheus.Metric) {
	logger := e.logger.With("collector", scName)
	logger.Debugf("starting subcollection")

	startTime := time.Now()
	err := sc.Collect(ch)

	var success float64 = 1
	if err != nil {
		logger.Errorf("subcollection failed: %s", err)
		success = 0
	}

	ch <- prometheus.MustNewConstMetric(e.scrapeDurationDesc, prometheus.GaugeValue, time.Since(startTime).Seconds(), scName)
	ch <- prometheus.MustNewConstMetric(e.scrapeSuccessDesc, prometheus.GaugeValue, success, scName)
}

// subcollector is an internal type of collector that allows us to
// collect custmizing the collection pieces and track if the collect
// process failed and.
type subcollector interface {
	// Collect will collect and return if the collection has been made successfully.
	Collect(ch chan<- prometheus.Metric) error
}
