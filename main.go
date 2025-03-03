package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Options struct {
	BatchSize   int
	Ops         int
	Endpoint    string
	Db          string
	Table       string
	Pipeline    string
	Interval    time.Duration
	PayloadSize int
	TotalSize   int
	AppNum      int
	LogSize     int
}

func main() {
	batchSize := flag.Int("batch-size", 1000, "number of logs to ingest")
	ops := flag.Int("ops", 0, "number of operations to perform")
	endpoint := flag.String("endpoint", "http://localhost:4000", "endpoint to ingest logs to")
	db := flag.String("db", "public", "database to ingest logs to")
	table := flag.String("table", "testlogs", "table to ingest logs to")
	pipeline := flag.String("pipeline", "greptime_identity", "pipeline name to ingest logs to")
	interval := flag.Duration("interval", time.Second*1, "interval to generate logs")
	payloadSize := flag.Int("payload-size", 100, "payload size to generate logs(bytes)")
	totalSize := flag.Int("total-size", 5000100, "total size to generate logs(bytes)")
	appNum := flag.Int("app-num", 1000, "number of apps to generate logs")
	logSize := flag.Int("log-size", 0, "log size to generate logs(bytes)")
	flag.Parse()

	options := &Options{
		BatchSize:   *batchSize,
		Ops:         *ops,
		Endpoint:    *endpoint,
		Db:          *db,
		Table:       *table,
		Pipeline:    *pipeline,
		Interval:    *interval,
		PayloadSize: *payloadSize,
		TotalSize:   *totalSize,
		AppNum:      *appNum,
		LogSize:     *logSize,
	}

	fmt.Printf("Starting benchmark with options: %+v\n", options)
	for {
		doBenchmark(options)
	}
}

func doBenchmark(opts *Options) {
	var wg sync.WaitGroup

	taskNum := opts.Ops / opts.BatchSize
	start := time.Now()
	for i := 0; i < taskNum; i++ {
		startTime := start.Add(time.Duration(i) * opts.Interval * time.Duration(opts.BatchSize))
		wg.Add(1)
		go func() {
			defer wg.Done()
			var data []byte
			if opts.LogSize > 0 {
				data = generateLogs(opts.BatchSize, opts.TotalSize, opts.LogSize, startTime, opts.Interval, opts.AppNum)
			} else {
				data = generateLogs(opts.BatchSize, opts.TotalSize, payloadSizeByProbabilityDistribution(), startTime, opts.Interval, opts.AppNum)
			}
			if err := ingestLogs(opts.Endpoint, opts.Db, opts.Table, opts.Pipeline, data, true); err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)

	if elapsed < time.Second*1 {
		time.Sleep(time.Second*1 - elapsed)
	}
}

func ingestLogs(endpoint, db, table, pipeline string, logs []byte, enableGzipCompression bool) error {
	url := fmt.Sprintf("%s/v1/events/logs?db=%s&table=%s&pipeline_name=%s", endpoint, db, table, pipeline)

	var buffer bytes.Buffer
	if enableGzipCompression {
		writer := gzip.NewWriter(&buffer)
		if _, err := writer.Write(logs); err != nil {
			return err
		}
		writer.Close()
	}

	req, err := http.NewRequest("POST", url, &buffer)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if enableGzipCompression {
		req.Header.Set("Content-Encoding", "gzip")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// TODO: Make it configurable.
func payloadSizeByProbabilityDistribution() int {
	r := rand.Float64() * 100

	switch {
	case r < 20: // 20%
		return rand.Intn(800) + 200 // 200-1000
	case r < 95: // 75%
		return rand.Intn(1000) + 1000 // 1000-2000
	case r < 99: // 4%
		return rand.Intn(98000) + 2000 // 2000-100000
	case r < 99.5: // 0.5%
		return rand.Intn(900000) + 100000 // 100000-1000000
	default: // 0.5%
		return rand.Intn(4000000) + 1000000 // 1000000-5000000
	}
}
