package agent

import (
	"github.com/go-resty/resty/v2"
)

func SendMetric(requestURL, metricType, metricName, metricValue string, req *resty.Request) (*resty.Response, error) {
	return req.Post(req.URL + "/update/{metricType}/{metricName}/{metricValue}")
}









