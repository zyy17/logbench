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
	App       string `parquet:"name=app, type=BYTE_ARRAY, convertedtype=UTF8" json:"app"`
	Cluster   string `parquet:"name=cluster, type=BYTE_ARRAY, convertedtype=UTF8" json:"cluster"`
	Timestamp int64  `parquet:"name=greptime_timestamp, type=INT64" json:"greptime_timestamp"`
	Message   string `parquet:"name=message, type=BYTE_ARRAY, convertedtype=UTF8" json:"message"`
}

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

func (g *Generator) Generate(cluster, app string, batchSize int, timestamp int64) ([]byte, error) {
	var output []byte
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
			App:       app,
			Cluster:   cluster,
			Timestamp: timestamp,
			Message:   g.generateLogs(logSize),
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
			App:       app,
			Cluster:   cluster,
			Timestamp: timestamp,
			Message:   g.generateLogs(logSize),
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
