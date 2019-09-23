// Copyright 2019 Booking.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package graphite

import (
	"strconv"
	"time"

	"github.com/bookingcom/cloudsec-metrics/api"
	"github.com/marpaia/graphite-golang"
	"github.com/pkg/errors"
)

// New is creating new Graphite instance with given parameters, disabling graphite.Graphite default logging
func New(host string, port int, prefix string) (G *graphite.Graphite, err error) {
	if G, err = graphite.GraphiteFactory("tcp", host, port, prefix); err != nil {
		return nil, errors.Wrap(err, "can't connect to Graphite")
	}
	G.DisableLog = true
	return G, nil
}

// SendComplianceInfo tries to get assets compliance information for last day, thread-safe
func SendComplianceInfo(g *graphite.Graphite, prefix string, ci []api.ComplianceInfo) error {
	timeNow := time.Now().Unix()
	var metrics []graphite.Metric
	for _, x := range ci {
		metricPrefix := prefix
		for _, c := range x.Name {
			switch c {
			case ' ', '.', '{', '}', '(', ')', '/':
				metricPrefix += "_"
			default:
				metricPrefix += string(c)
			}
		}
		metrics = append(metrics, graphite.NewMetric(metricPrefix+".policies_total", strconv.Itoa(x.PoliciesCount), timeNow))
		metrics = append(metrics, graphite.NewMetric(metricPrefix+".assets_passed", strconv.Itoa(x.PassedAssetsCount), timeNow))
		metrics = append(metrics, graphite.NewMetric(metricPrefix+".assets_failed", strconv.Itoa(x.FailedAssetsCount), timeNow))
		metrics = append(metrics, graphite.NewMetric(metricPrefix+".assets_total", strconv.Itoa(x.TotalAssetsCount), timeNow))
	}
	return g.SendMetrics(metrics)
}
