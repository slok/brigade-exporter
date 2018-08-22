package collector_test

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type labelMap map[string]string

type metricResult struct {
	// Use string for decription object.
	// Comparing string is not good in test but this data is not dynamic, the dynamic data is tested
	// in object style using `labels` and `value` fields.
	desc       string
	labels     labelMap
	value      float64
	metricType dto.MetricType
}

func readMetric(m prometheus.Metric) metricResult {
	// Write metric on  model object so we can get the value afterwards.
	pb := &dto.Metric{}
	m.Write(pb)

	// Get labels
	labels := make(labelMap, len(pb.Label))
	for _, v := range pb.Label {
		labels[v.GetName()] = v.GetValue()
	}

	// Description string object.
	desc := m.Desc().String()

	// Metric type.
	if pb.Gauge != nil {
		return metricResult{desc: desc, labels: labels, value: pb.GetGauge().GetValue(), metricType: dto.MetricType_GAUGE}
	}
	if pb.Counter != nil {
		return metricResult{desc: desc, labels: labels, value: pb.GetCounter().GetValue(), metricType: dto.MetricType_COUNTER}
	}
	if pb.Untyped != nil {
		return metricResult{desc: desc, labels: labels, value: pb.GetUntyped().GetValue(), metricType: dto.MetricType_UNTYPED}
	}
	panic("Unsupported metric type")
}
