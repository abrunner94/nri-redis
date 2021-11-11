package main

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
)

type Specs struct {
	SpecVersion                  string     `yaml:"specVersion"`
	OwningTeam                   string     `yaml:"owningTeam"`
	IntegrationName              string     `yaml:"integrationName"`
	HumanReadableIntegrationName string     `yaml:"humanReadableIntegrationName"`
	Entities                     []Entities `yaml:"entities"`
}

type Entities struct {
	EntityType         string               `yaml:"entityType"`
	Metrics            []Metrics            `yaml:"metrics"`
	Tags               []Tags               `yaml:"tags"`
	InternalAttributes []InternalAttributes `yaml:"internalAttributes"`
	IgnoredAttributes  []string             `yaml:"ignoredAttributes"`
}

type Metrics struct {
	Name                 string               `yaml:"name"`
	Type                 string               `yaml:"type"`
	DefaultResolution    int                  `yaml:"defaultResolution"`
	Unit                 string               `yaml:"unit"`
	MigrationInformation MigrationInformation `yaml:"migrationInformation"`
}

type Tags struct {
	Name                 string               `yaml:"name"`
	Type                 string               `yaml:"type"`
	MigrationInformation MigrationInformation `yaml:"migrationInformation"`
	Description          string               `yaml:"description,omitempty"`
}
type InternalAttributes struct {
	Name                 string               `yaml:"name"`
	Type                 string               `yaml:"type"`
	MigrationInformation MigrationInformation `yaml:"migrationInformation"`
	Description          string               `yaml:"description,omitempty"`
}

type MigrationInformation struct {
	LegacyEventType string   `yaml:"legacyEventType"`
	LegacyNames     []string `yaml:"legacyNames"`
}

func LoadData() Specs {
	s := Specs{}

	data, err := ioutil.ReadFile("./src/spec.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	err = yaml.Unmarshal(data, &s)
	if err != nil {
		log.Fatalln(err)
	}

	return s
}

func (s Specs) SetMetrics(i *integration.Integration) {
	for _, e := range i.Entities {
		for _, mSet := range e.Metrics {
			entityType := mSet.Metrics["event_type"].(string)
			for metricName, v := range mSet.Metrics {

				var value float64

				switch t := v.(type) {
				case float64:
					value = t
				case bool:
					continue
				case int:
					value = float64(t)
				default:
					log.Println(t)
				}

				if c, ok := a[entityType+metricName]; ok {
					//if c, ok := c.(*prometheus.CounterVec); ok {
					//	c.DeleteLabelValues(e.Metadata.Name)
					//	c.WithLabelValues(e.Metadata.Name).Add(value)
					//}
					if c, ok := c.(*prometheus.GaugeVec); ok {
						c.WithLabelValues(e.Metadata.Name).Set(value)
					}
				} else {
					log.Printf("metric not present in prometheus: %q \n", metricName)
				}
			}
		}
	}
}

func (s Specs) RegisterMetrics() {
	a = map[string]prometheus.Collector{}
	for _, e := range s.Entities {
		for _, m := range e.Metrics {
			for _, l := range m.MigrationInformation.LegacyNames {
				parsedName := strings.ToLower(strings.Replace(m.Name, ".", "_", -1))
				if m.Type == "gauge" {
					g := prometheus.NewGaugeVec(prometheus.GaugeOpts{
						Namespace:   "",
						Subsystem:   "",
						Name:        parsedName,
						Help:        "",
						ConstLabels: nil,
					}, []string{"entity_name"})
					prometheus.MustRegister(g)
					a[m.MigrationInformation.LegacyEventType+l] = g
				}
				//if m.Type == "count" {
				//	c := prometheus.NewCounterVec(prometheus.CounterOpts{
				//
				//		Namespace:   "",
				//		Subsystem:   "",
				//		Name:        parsedName,
				//		Help:        "",
				//		ConstLabels: nil,
				//	}, []string{"entity_name"})
				//	prometheus.MustRegister(c)
				//	a[m.MigrationInformation.LegacyEventType+l] = c
				//}
			}
		}
	}
}

var a map[string]prometheus.Collector
