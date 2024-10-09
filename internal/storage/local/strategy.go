package local

import (
	"fmt"
	"sync"

	"github.com/igortoigildin/go-metrics-altering/internal/models"
)

type Strategy interface {
	Update(metricType string, metricName string, metricValue any) error
	Get(metricType string, metricName string) (models.Metrics, error)
}

type count struct {
	Counter map[string]int64
	rm      sync.RWMutex
}

func (c *count) Update(metricType string, metricName string, metricValue any) error {
	if c.Counter == nil {
		c.Counter = make(map[string]int64)
	}
	c.rm.Lock()
	v, ok := metricValue.(*int64)
	if !ok {
		fmt.Println("FALSE")
	}
	var temp int64
	if v == nil {
		temp = int64(0)
	} else {
		temp = *v
	}
	fmt.Println("RESULT")
	c.Counter[metricName] += temp
	c.rm.Unlock()
	return nil
}

func (c *count) Get(metricType string, metricName string) (models.Metrics, error) {
	var metric models.Metrics

	c.rm.RLock()
	d := c.Counter[metricName]
	metric.Delta = &d
	c.rm.RUnlock()
	metric.MType = metricType

	return metric, nil
}

type gauge struct {
	Gauge map[string]float64
	rm    sync.RWMutex
}

func (g *gauge) Update(metricType string, metricName string, metricValue any) error {
	if g.Gauge == nil {
		g.Gauge = make(map[string]float64)
	}
	g.rm.Lock()
	v, _ := metricValue.(*float64)
	var temp float64
	if v == nil {
		temp = float64(0)
	} else {
		temp = *v
	}
	g.Gauge[metricName] = temp
	g.rm.Unlock()
	return nil
}
func (g *gauge) Get(metricType string, metricName string) (models.Metrics, error) {
	var metric models.Metrics

	g.rm.RLock()
	v := g.Gauge[metricName]
	metric.Value = &v
	g.rm.RUnlock()
	metric.MType = metricType

	return metric, nil
}
