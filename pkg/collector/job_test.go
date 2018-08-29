package collector_test

import (
	"context"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"

	mbrigade "github.com/slok/brigade-exporter/mocks/service/brigade"
	"github.com/slok/brigade-exporter/pkg/collector"
	"github.com/slok/brigade-exporter/pkg/log"
	"github.com/slok/brigade-exporter/pkg/service/brigade"
)

const (
	jobInfoDesc     = `Desc{fqName: "brigade_job_info", help: "Brigade job information.", constLabels: {}, variableLabels: [id build_id name image]}`
	jobStatusDesc   = `Desc{fqName: "brigade_job_status", help: "Brigade job status.", constLabels: {}, variableLabels: [id status]}`
	jobDurationDesc = `Desc{fqName: "brigade_job_duration_seconds", help: "Brigade job duration in seconds.", constLabels: {}, variableLabels: [id]}`
	jobCreationDesc = `Desc{fqName: "brigade_job_create_time_seconds", help: "Brigade job creation time in unix timestamp.", constLabels: {}, variableLabels: [id]}`
	jobStartDesc    = `Desc{fqName: "brigade_job_start_time_seconds", help: "Brigade job start time in unix timestamp.", constLabels: {}, variableLabels: [id]}`
)

func TestJobSubcollector(t *testing.T) {
	// Test times.
	t1 := time.Now()
	t2 := t1.Add(265 * time.Second)
	t3 := t2.Add(12 * time.Minute)
	t4 := t3.Add(1 * time.Hour)

	tests := []struct {
		name       string
		jobs       []*brigade.Job
		expMetrics []metricResult
	}{
		{
			name: "With multiple jobs the collected metrics should be of all the jobs.",
			jobs: []*brigade.Job{
				&brigade.Job{ID: "id1", BuildID: "bld1", Name: "id-name-1", Image: "image1", Status: "Running", Duration: 125 * time.Second, Creation: t1, Start: t2},
				&brigade.Job{ID: "id2", BuildID: "bld2", Name: "id-name-2", Image: "image2", Status: "Pending", Duration: 340 * time.Second, Creation: t3},
				&brigade.Job{ID: "id3", BuildID: "bld3", Name: "id-name-3", Image: "image3", Status: "Failed", Duration: 18 * time.Second, Start: t4},
			},
			expMetrics: []metricResult{
				metricResult{
					desc:       jobInfoDesc,
					labels:     labelMap{"id": "id1", "build_id": "bld1", "name": "id-name-1", "image": "image1"},
					value:      1,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       jobStatusDesc,
					labels:     labelMap{"id": "id1", "status": "Running"},
					value:      1,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       jobDurationDesc,
					labels:     labelMap{"id": "id1"},
					value:      125,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       jobCreationDesc,
					labels:     labelMap{"id": "id1"},
					value:      float64(t1.Unix()),
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       jobStartDesc,
					labels:     labelMap{"id": "id1"},
					value:      float64(t2.Unix()),
					metricType: dto.MetricType_GAUGE,
				},

				metricResult{
					desc:       jobInfoDesc,
					labels:     labelMap{"id": "id2", "build_id": "bld2", "name": "id-name-2", "image": "image2"},
					value:      1,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       jobStatusDesc,
					labels:     labelMap{"id": "id2", "status": "Pending"},
					value:      1,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       jobDurationDesc,
					labels:     labelMap{"id": "id2"},
					value:      340,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       jobCreationDesc,
					labels:     labelMap{"id": "id2"},
					value:      float64(t3.Unix()),
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       jobStartDesc,
					labels:     labelMap{"id": "id2"},
					value:      0,
					metricType: dto.MetricType_GAUGE,
				},

				metricResult{
					desc:       jobInfoDesc,
					labels:     labelMap{"id": "id3", "build_id": "bld3", "name": "id-name-3", "image": "image3"},
					value:      1,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       jobStatusDesc,
					labels:     labelMap{"id": "id3", "status": "Failed"},
					value:      1,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       jobDurationDesc,
					labels:     labelMap{"id": "id3"},
					value:      18,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       jobCreationDesc,
					labels:     labelMap{"id": "id3"},
					value:      0,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       jobStartDesc,
					labels:     labelMap{"id": "id3"},
					value:      float64(t4.Unix()),
					metricType: dto.MetricType_GAUGE,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mbsvc := &mbrigade.Interface{}
			mbsvc.On("GetJobs").Once().Return(test.jobs, nil)

			clr := collector.NewJob(mbsvc, log.Dummy)

			ch := make(chan prometheus.Metric)

			go func() {
				clr.Collect(context.TODO(), ch)
				close(ch)
			}()

			// Get the metrics
			var got []metricResult
			for m := range ch {
				got = append(got, readMetric(m))
			}

			// Check metrics are ok.
			assert.Equal(test.expMetrics, got)
		})
	}
}
