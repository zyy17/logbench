package generator

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/zyy17/logbench/pkg/dataset"
	"github.com/zyy17/logbench/pkg/generator/distribution"
	"github.com/zyy17/logbench/pkg/utils"
)

type Log struct {
	App                  string `parquet:"name=app, type=BYTE_ARRAY, convertedtype=UTF8" json:"app"`
	Cluster              string `parquet:"name=cluster, type=BYTE_ARRAY, convertedtype=UTF8" json:"cluster"`
	Timestamp            string `json:"greptime_timestamp,omitempty"`
	TimestampMillisecond int64  `parquet:"name=greptime_timestamp, type=INT64" json:"timestamp_millisecond,omitempty"`
	Message              string `parquet:"name=message, type=BYTE_ARRAY, convertedtype=UTF8" json:"message"`
	Region               string `parquet:"name=region, type=BYTE_ARRAY, convertedtype=UTF8" json:"region"`
	CloudProvider        string `parquet:"name=cloud-provider, type=BYTE_ARRAY, convertedtype=UTF8" json:"cloud-provider"`
	Environment          string `parquet:"name=environment, type=BYTE_ARRAY, convertedtype=UTF8" json:"environment"`
	Product              string `parquet:"name=product, type=BYTE_ARRAY, convertedtype=UTF8" json:"product"`
	SubProduct           string `parquet:"name=sub-product, type=BYTE_ARRAY, convertedtype=UTF8" json:"sub-product"`
	Service              string `parquet:"name=service, type=BYTE_ARRAY, convertedtype=UTF8" json:"service"`
}

var (
	region        = []string{"us-east-1", "cn-hangzhou", "ap-southeast-1"}
	cloudProvider = []string{"aws", "azure", "gcp", "aliyun"}
	environment   = []string{"dev", "test", "prod"}
)

type Generator struct {
	dt      *distribution.Distribution
	dataset *dataset.Dataset
	logSize int64
}

func NewGenerator(distributionConfig []string, logSize int64, datasetType dataset.DatasetType) (*Generator, error) {
	dt, err := distribution.NewDistribution(distributionConfig)
	if err != nil {
		return nil, err
	}

	return &Generator{
		dt:      dt,
		dataset: dataset.Datasets[datasetType],
		logSize: logSize,
	}, nil
}

func (g *Generator) Generate(cluster, app string, batchSize int, timestamp string, timestampMillisecond int64) ([]byte, error) {
	var output []byte
	var logSize int64

	if g.logSize == 0 {
		logSize = g.dt.RandomNumber()
	} else {
		logSize = g.logSize
	}

	for i := 0; i < batchSize; i++ {
		log := &Log{
			App:                  app,
			Cluster:              cluster,
			Timestamp:            timestamp,
			TimestampMillisecond: timestampMillisecond,
			Message:              g.generateLogs(logSize),
			Region:               region[utils.RandomNumber(0, int64(len(region)))],
			CloudProvider:        cloudProvider[utils.RandomNumber(0, int64(len(cloudProvider)))],
			Environment:          environment[utils.RandomNumber(0, int64(len(environment)))],
			Product:              app,
			SubProduct:           app,
			Service:              app,
		}

		serializedData, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		serializedData = append(serializedData, '\n')

		output = append(output, serializedData...)
	}

	return output, nil
}

func (g *Generator) GenerateLogs(cluster, app string, batchSize int, timestamp int64) ([]*Log, error) {
	var output []*Log
	var logSize int64

	if g.logSize == 0 {
		logSize = g.dt.RandomNumber()
	} else {
		logSize = g.logSize
	}

	// If timestamp is not provided, use the current timestamp
	if timestamp == 0 {
		timestamp = time.Now().UnixMilli()
	}

	for i := 0; i < batchSize; i++ {
		log := &Log{
			App:                  app,
			Cluster:              cluster,
			TimestampMillisecond: timestamp,
			Message:              g.generateLogs(logSize),
			Region:               region[utils.RandomNumber(0, int64(len(region)))],
			CloudProvider:        cloudProvider[utils.RandomNumber(0, int64(len(cloudProvider)))],
			Environment:          environment[utils.RandomNumber(0, int64(len(environment)))],
			Product:              app,
			SubProduct:           app,
			Service:              app,
		}

		output = append(output, log)
	}

	return output, nil
}

func (g *Generator) generateLogs(size int64) string {
	logs := make([]string, 0, size/int64(g.dataset.AverageSize))
	for i := 0; i < int(size/int64(g.dataset.AverageSize)); i++ {
		randomIndex := utils.RandomNumber(0, int64(len(g.dataset.Logs)))
		logs = append(logs, g.dataset.Logs[randomIndex])
	}

	return strings.Join(logs, "\n")
}
