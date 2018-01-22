package aggregate

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/iyacontrol/telegraf-proxy/discovery"
	"github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

// Aggregator defines aggregate telegraf metrics
type Aggregator struct {
	HTTP *http.Client
}

type Result struct {
	URL          string
	SecondsTaken float64
	MetricFamily map[string]*io_prometheus_client.MetricFamily
	Error        error
}

func (a *Aggregator) Aggregate(reg *discovery.Registery, output io.Writer) {
	resultChan := make(chan *Result, 100)

	targets := reg.Translate(":9273/metrics")

	for _, target := range targets {
		go a.fetch(target, resultChan)
	}

	func(numTargets int, resultChan chan *Result) {

		numResuts := 0

		allFamilies := make(map[string]*io_prometheus_client.MetricFamily)

		for {
			if numTargets == numResuts {
				break
			}
			select {
			case result := <-resultChan:
				numResuts++

				if result.Error != nil {
					log.Printf("Fetch error: %s", result.Error.Error())
					continue
				}

				for mfName, mf := range result.MetricFamily {
					if existingMf, ok := allFamilies[mfName]; ok {
						for _, m := range mf.Metric {
							existingMf.Metric = append(existingMf.Metric, m)
						}
					} else {
						allFamilies[*mf.Name] = mf
					}
				}

				log.Printf("OK: %s was refreshed in %.3f seconds", result.URL, result.SecondsTaken)

			}
		}

		encoder := expfmt.NewEncoder(output, expfmt.FmtText)
		for _, f := range allFamilies {
			encoder.Encode(f)
		}

	}(len(targets), resultChan)

}

func (a *Aggregator) fetch(target string, resultChan chan *Result) {

	startTime := time.Now()
	res, err := a.HTTP.Get(target)

	result := &Result{URL: target, SecondsTaken: time.Since(startTime).Seconds(), Error: nil}
	if res != nil {
		result.MetricFamily, err = getMetricFamilies(res.Body)
		if err != nil {
			result.Error = fmt.Errorf("failed to add labels to target %s metrics: %s", target, err.Error())
			resultChan <- result
			return
		}
	}
	if err != nil {
		result.Error = fmt.Errorf("failed to fetch URL %s due to error: %s", target, err.Error())
	}
	resultChan <- result
}

func getMetricFamilies(sourceData io.Reader) (map[string]*io_prometheus_client.MetricFamily, error) {
	parser := expfmt.TextParser{}
	metricFamiles, err := parser.TextToMetricFamilies(sourceData)
	if err != nil {
		return nil, err
	}
	return metricFamiles, nil
}
