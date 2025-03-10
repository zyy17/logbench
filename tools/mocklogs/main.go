package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go/writer"

	"github.com/zyy17/logbench/pkg/dataset"
	"github.com/zyy17/logbench/pkg/generator"
)

type Options struct {
	StartTime     string
	EndTime       string
	IntervalCount int
	Cluster       string
	App           string
	OutputDir     string
	Format        string
	Interval      string
	Precision     string
}

var distributionConfig = []string{
	"20%:200bytes-1kb",
	"76%:1kb-2kb",
	"3.5%:2kb-100kb",
	"0.4%:200kb-1mb",
	"0.1%:1mb-5mb",
}

type Format string

const (
	FormatJSON    Format = "json"
	FormatParquet Format = "parquet"
)

func main() {
	options := &Options{}
	flag.StringVar(&options.StartTime, "start-time", "", "start time")
	flag.StringVar(&options.EndTime, "end-time", "", "end time")
	flag.IntVar(&options.IntervalCount, "interval-count", 1, "interval count(seconds)")
	flag.StringVar(&options.Cluster, "cluster", "", "cluster")
	flag.StringVar(&options.App, "app", "", "app")
	flag.StringVar(&options.OutputDir, "output-dir", "", "output directory")
	flag.StringVar(&options.Format, "format", "json", "output format, can be json or parquet")
	flag.StringVar(&options.Interval, "interval", "1s", "interval")
	flag.StringVar(&options.Precision, "precision", "ms", "precision, can be ms or ns")
	flag.Parse()

	startTime, err := time.Parse(time.RFC3339, options.StartTime)
	if err != nil {
		log.Fatalf("failed to parse start time: %v", err)
	}

	endTime, err := time.Parse(time.RFC3339, options.EndTime)
	if err != nil {
		log.Fatalf("failed to parse end time: %v", err)
	}

	interval, err := time.ParseDuration(options.Interval)
	if err != nil {
		log.Fatalf("failed to parse interval: %v", err)
	}

	if endTime.Before(startTime) {
		log.Fatalf("end time must be after start time")
	}

	if options.Cluster == "" {
		log.Fatalf("cluster is required")
	}

	if options.App == "" {
		log.Fatalf("app is required")
	}

	if options.OutputDir == "" {
		log.Fatalf("log output directory is required")
	}

	if options.Format != string(FormatJSON) && options.Format != string(FormatParquet) {
		log.Fatalf("unsupported format: %s", options.Format)
	}

	if options.Precision == "" {
		options.Precision = "ns"
	}

	if err := os.MkdirAll(options.OutputDir, 0755); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	log.Printf("Generating %s logs cluster: '%s', app: '%s', precision: '%s', from [%s] to [%s], interval: '%s', interval count: '%d', output directory: '%s'",
		options.Format, options.Cluster, options.App, options.Precision, options.StartTime, options.EndTime, options.Interval, options.IntervalCount, options.OutputDir)

	var (
		fileWriter    io.Writer
		parquetWriter *writer.ParquetWriter
		parquetFile   source.ParquetFile
	)

	if options.Format == string(FormatJSON) {
		outputFile := filepath.Join(options.OutputDir, fileName(options))
		fileWriter, err = os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("failed to create output file: %v", err)
		}
	} else {
		outputFile := filepath.Join(options.OutputDir, fileName(options))
		parquetFile, err = local.NewLocalFileWriter(outputFile)
		if err != nil {
			log.Fatalf("failed to create output file: %v", err)
		}
		parquetWriter, err = writer.NewParquetWriter(parquetFile, new(generator.Log), 4)
		if err != nil {
			log.Fatalf("failed to create parquet writer: %v", err)
		}
		parquetWriter.RowGroupSize = 1024 * 1024 * 10 // 10MB
		parquetWriter.PageSize = 8 * 1024             // 8K
		parquetWriter.CompressionType = parquet.CompressionCodec_SNAPPY
	}

	generator, err := generator.NewGenerator(distributionConfig, 0, dataset.DatasetType(dataset.Zookeeper2kDatasetType))
	if err != nil {
		log.Fatalf("failed to create generator: %v", err)
	}

	count := 0
	start := time.Now()
	size := 0
	for startTime.Before(endTime) {
		if options.Format == string(FormatJSON) {
			logs, err := generator.Generate(options.Cluster, options.App, options.IntervalCount, startTime.Format(time.RFC3339), 0)
			if err != nil {
				log.Fatalf("failed to generate logs: %v", err)
			}
			_, err = fileWriter.Write(logs)
			if err != nil {
				log.Fatalf("failed to write logs: %v", err)
			}
			size += len(logs)
		} else {
			logs, err := generator.GenerateLogs(options.Cluster, options.App, options.IntervalCount, timestamp(options.Precision, startTime))
			if err != nil {
				log.Fatalf("failed to generate logs: %v", err)
			}
			for _, generatedLog := range logs {
				if err := parquetWriter.Write(generatedLog); err != nil {
					log.Fatalf("failed to write logs: %v", err)
				}
				size += len(generatedLog.Message)
			}
		}

		count += options.IntervalCount
		startTime = startTime.Add(interval)
	}

	if options.Format == string(FormatParquet) {
		if err := parquetWriter.WriteStop(); err != nil {
			log.Fatalf("failed to write stop: %v", err)
		}
		if err := parquetFile.Close(); err != nil {
			log.Fatalf("failed to close parquet file: %v", err)
		}
	}

	log.Printf("Generated '%d' file in '%v', approximate size: '%d MB, output file: '%s'",
		count, time.Since(start), size/1024/1024, fileName(options))
}

// The file name will be like:
// cluster1_app1-[2025-03-02T00:00:00Z-2025-03-02T01:00:00Z].json
// cluster1_app1-[2025-03-02T00:00:00Z-2025-03-02T01:00:00Z].parquet
func fileName(options *Options) string {
	return fmt.Sprintf("%s.%s-[%s-%s].%s", options.Cluster, options.App, options.StartTime, options.EndTime, options.Format)
}

func timestamp(precision string, timestamp time.Time) int64 {
	switch precision {
	case "s":
		return timestamp.Unix()
	case "ms":
		return timestamp.UnixMilli()
	case "us":
		return timestamp.UnixMicro()
	case "ns":
		return timestamp.UnixNano()
	default:
		return timestamp.UnixNano()
	}
}
