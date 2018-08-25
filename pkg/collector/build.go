package collector

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/slok/brigade-exporter/pkg/log"
	"github.com/slok/brigade-exporter/pkg/service/brigade"
)

const (
	buildSubSystem = "build"
)

// build is the Brigade build subcollector. this colletor will collect
// the metrics regarding brigade builds.
// Satisfies internfal collector interface.
type build struct {
	brigadeSVC brigade.Interface
	logger     log.Logger

	// Metrics.
	buildInfoDesc     *prometheus.Desc
	buildStatusDesc   *prometheus.Desc
	buildDurationDesc *prometheus.Desc
}

// NewBuild returns a new build subcollector.
func NewBuild(brigadeSVC brigade.Interface, logger log.Logger) subcollector {
	return &build{
		brigadeSVC: brigadeSVC,
		logger:     logger,

		buildInfoDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, buildSubSystem, "info"),
			"Brigade build information.",
			[]string{"id", "project_id", "event_type", "provider", "version"}, nil,
		),
		buildStatusDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, buildSubSystem, "status"),
			"Brigade build status.",
			[]string{"id", "status"}, nil,
		),
		buildDurationDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, buildSubSystem, "duration_seconds"),
			"Brigade build duration in seconds.",
			[]string{"id"}, nil,
		),
	}
}

// Collect satisfies subcollector.
func (b *build) Collect(ctx context.Context, ch chan<- prometheus.Metric) error {
	blds, err := b.brigadeSVC.GetBuilds()
	if err != nil {
		return err
	}

	for _, bld := range blds {
		// Info metric.
		err := sendMetric(ctx, ch, prometheus.MustNewConstMetric(
			b.buildInfoDesc,
			prometheus.GaugeValue,
			1,
			bld.ID, bld.ProjectID, bld.Type, bld.Provider, bld.Version))

		if err != nil {
			return err
		}

		// Status metric.
		err = sendMetric(ctx, ch, prometheus.MustNewConstMetric(
			b.buildStatusDesc,
			prometheus.GaugeValue,
			1,
			bld.ID, bld.Status))

		if err != nil {
			return err
		}

		// Duration metric.
		// TODO: Think if it's 0 we should send the metric or not.
		err = sendMetric(ctx, ch, prometheus.MustNewConstMetric(
			b.buildDurationDesc,
			prometheus.GaugeValue,
			bld.Duration.Seconds(),
			bld.ID))

		if err != nil {
			return err
		}
	}

	return nil
}
