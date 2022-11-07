package pkg

import (
	"time"

	"github.com/app-sre/aws-resource-exporter/pkg/awsclient"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type TagExporter struct {
	awsAccountId string

	logger   log.Logger
	interval time.Duration
	exporter *ArnTagExporter
}

type ArnTagExporter struct {
	tagMap map[string][]awsclient.Tag
	desc   *prometheus.Desc
}

func (ate ArnTagExporter) Collect(ch chan<- prometheus.Metric) {
	if ate.desc != nil {
		ch <- prometheus.MustNewConstMetric(ate.desc, prometheus.GaugeValue, 1.0)
	}
}

func (ate ArnTagExporter) Describe(ch chan<- *prometheus.Desc) {
	// labelMap := make(map[string]string, 0)

	// for arn, tags := range ate.tagMap {
	// 	ate.desc = prometheus.NewDesc("aws_resource_exporter_aws_tags", "", labelNames, nil)
	// 	ch <- ate.desc
	// }
}

func NewTagExporter(logger log.Logger, config EC2Config, awsAccountId string) *TagExporter {
	return &TagExporter{
		awsAccountId: awsAccountId,
		logger:       logger,
		interval:     *config.Interval,
	}
}

func (e *TagExporter) CollectLoop() {
	for {

		// x := awsclient.GetTagCache()

		// if e.exporter == nil {
		// 	fmt.Printf("Registering tag exporter \n")
		// 	e.exporter = &ArnTagExporter{
		// 		tagMap: x.GetTagsPerARN(),
		// 	}

		// 	prometheus.MustRegister(e.exporter)
		// }

		x := awsclient.GetTagCache()

		for _, tags := range x.GetTagsPerARN() {

			labelMap := make(prometheus.Labels)

			for _, tag := range tags {
				labelMap[tag.Key] = tag.Value
			}

			promauto.NewGauge(prometheus.GaugeOpts{
				Namespace:   "aws_resource_exporter",
				Subsystem:   "tags",
				Name:        "rds",
				ConstLabels: labelMap,
			})
		}

		time.Sleep(1 * time.Second)
	}
}
