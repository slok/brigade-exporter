package collector

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/slok/brigade-exporter/pkg/log"
	"github.com/slok/brigade-exporter/pkg/service/brigade"
)

const (
	jobSubSystem = "job"
)

// job is the Brigade Job subcollector. this colletor will collect
// the metrics regarding brigade jobs.
// Satisfies internfal collector interface.
type job struct {
	brigadeSVC brigade.Interface
	logger     log.Logger

	// Metrics.
	jobInfoDesc     *prometheus.Desc
	jobStatusDesc   *prometheus.Desc
	jobDurationDesc *prometheus.Desc
	jobCreationDesc *prometheus.Desc
	jobStartDesc    *prometheus.Desc
}

// NewJob returns a new job subcollector.
func NewJob(brigadeSVC brigade.Interface, logger log.Logger) subcollector {
	return &job{
		brigadeSVC: brigadeSVC,
		logger:     logger,

		jobInfoDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, jobSubSystem, "info"),
			"Brigade job information.",
			[]string{"id", "build_id", "name", "image"}, nil,
		),
		jobStatusDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, jobSubSystem, "status"),
			"Brigade job status.",
			[]string{"id", "status"}, nil,
		),
		jobDurationDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, jobSubSystem, "duration_seconds"),
			"Brigade job duration in seconds.",
			[]string{"id"}, nil,
		),
		jobCreationDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, jobSubSystem, "create_time_seconds"),
			"Brigade job creation time in unix timestamp.",
			[]string{"id"}, nil,
		),
		jobStartDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, jobSubSystem, "start_time_seconds"),
			"Brigade job start time in unix timestamp.",
			[]string{"id"}, nil,
		),
	}
}

// Collect satisfies subcollector.
func (j *job) Collect(ctx context.Context, ch chan<- prometheus.Metric) error {
	jobs, err := j.brigadeSVC.GetJobs()
	if err != nil {
		return err
	}

	for _, job := range jobs {
		// Info metric.
		err := sendMetric(ctx, ch, prometheus.MustNewConstMetric(
			j.jobInfoDesc,
			prometheus.GaugeValue,
			1,
			job.ID, job.BuildID, job.Name, job.Image))

		if err != nil {
			return err
		}

		// Status metric.
		err = sendMetric(ctx, ch, prometheus.MustNewConstMetric(
			j.jobStatusDesc,
			prometheus.GaugeValue,
			1,
			job.ID, job.Status))

		if err != nil {
			return err
		}

		// Duration metric.
		// TODO: Think if it's 0 we should send the metric or not.
		err = sendMetric(ctx, ch, prometheus.MustNewConstMetric(
			j.jobDurationDesc,
			prometheus.GaugeValue,
			job.Duration.Seconds(),
			job.ID))

		if err != nil {
			return err
		}

		// creation and start metrics.
		// TODO: Think if it's `time.IsZero`` we should send the metric or not.
		err = sendMetric(ctx, ch, prometheus.MustNewConstMetric(
			j.jobCreationDesc,
			prometheus.GaugeValue,
			j.getUnix(job.Creation),
			job.ID))

		if err != nil {
			return err
		}

		err = sendMetric(ctx, ch, prometheus.MustNewConstMetric(
			j.jobStartDesc,
			prometheus.GaugeValue,
			j.getUnix(job.Start),
			job.ID))
		if err != nil {
			return err
		}
	}

	return nil
}

func (*job) getUnix(t time.Time) float64 {
	if t.IsZero() {
		return 0
	}

	return float64(t.Unix())
}
