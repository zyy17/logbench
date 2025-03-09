package dataset

import (
	"bufio"
	"bytes"
	_ "embed"
	"strings"
)

type Dataset struct {
	Logs        []string
	AverageSize int // in bytes
}

type DatasetType string

const (
	Apache2kDatasetType    DatasetType = "Apache_2k"
	Zookeeper2kDatasetType DatasetType = "Zookeeper_2k"
)

var Datasets = make(map[DatasetType]*Dataset)

//go:embed Apache_2k.log
var apacheLogs []byte

//go:embed ZooKeeper_2k.log
var zookeeperLogs []byte

func init() {
	apacheLogs := normalizeLog(apacheLogs)
	zookeeperLogs := normalizeLog(zookeeperLogs)

	Datasets[Apache2kDatasetType] = &Dataset{
		Logs:        apacheLogs,
		AverageSize: calculateAverageSize(apacheLogs),
	}
	Datasets[Zookeeper2kDatasetType] = &Dataset{
		Logs:        zookeeperLogs,
		AverageSize: calculateAverageSize(zookeeperLogs),
	}
}

func calculateAverageSize(logs []string) int {
	var totalSize int
	for _, log := range logs {
		totalSize += len(log)
	}
	return totalSize / len(logs)
}

func normalizeLog(rawLogs []byte) []string {
	buf := bytes.NewBuffer(rawLogs)
	scanner := bufio.NewScanner(buf)
	var logs []string
	for scanner.Scan() {
		// Remove the leading and trailing spaces and the trailing newline.
		logs = append(logs, strings.TrimSpace(scanner.Text()))
	}
	return logs
}
