package collector_test

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"

	mbrigade "github.com/slok/brigade-exporter/mocks/service/brigade"
	"github.com/slok/brigade-exporter/pkg/collector"
	"github.com/slok/brigade-exporter/pkg/log"
	"github.com/slok/brigade-exporter/pkg/service/brigade"
)

func TestProjectSubcollector(t *testing.T) {
	tests := []struct {
		name        string
		exporterCfg collector.Config
		projects    []*brigade.Project
		expMetrics  []metricResult
	}{
		{
			name: "With multiple projects the collected metrics should be of all the projects.",
			projects: []*brigade.Project{
				&brigade.Project{ID: "id1", Name: "Name1", Repository: "repo1", Namespace: "ns1", Worker: "worker1"},
				&brigade.Project{ID: "id2", Name: "Name2", Repository: "repo2", Namespace: "ns2", Worker: "worker2"},
				&brigade.Project{ID: "id3", Name: "Name3", Repository: "repo3", Namespace: "ns3", Worker: "worker3"},
			},
			expMetrics: []metricResult{
				metricResult{labels: labelMap{"id": "id1", "name": "Name1", "repository": "repo1", "namespace": "ns1", "worker": "worker1"}, value: 1, metricType: dto.MetricType_GAUGE},
				metricResult{labels: labelMap{"id": "id2", "name": "Name2", "repository": "repo2", "namespace": "ns2", "worker": "worker2"}, value: 1, metricType: dto.MetricType_GAUGE},
				metricResult{labels: labelMap{"id": "id3", "name": "Name3", "repository": "repo3", "namespace": "ns3", "worker": "worker3"}, value: 1, metricType: dto.MetricType_GAUGE},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mbsvc := &mbrigade.Interface{}
			mbsvc.On("GetProjects").Once().Return(test.projects, nil)

			clr := collector.NewProject(mbsvc, log.Dummy)

			ch := make(chan prometheus.Metric)

			go func() {
				clr.Collect(ch)
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
