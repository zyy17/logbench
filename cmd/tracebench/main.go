package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/zyy17/logbench/pkg/ingester"
	"github.com/zyy17/logbench/pkg/trace"
)

type Options struct {
	BatchSize int
	Ops       int
	Endpoint  string
	AppNum    int
}

func main() {
	batchSize := flag.Int("batch-size", 0, "number of logs to ingest")
	ops := flag.Int("ops", 0, "number of operations to perform")
	endpoint := flag.String("endpoint", "http://localhost:4000", "endpoint to ingest logs to")
	appNum := flag.Int("app-num", 50, "number of apps to generate logs")
	flag.Parse()

	if *batchSize == 0 {
		fmt.Println("batch-size is required")
		os.Exit(1)
	}

	if *ops == 0 {
		fmt.Println("ops is required")
		os.Exit(1)
	}

	if *ops < *batchSize {
		fmt.Println("ops must be greater than batch-size")
		os.Exit(1)
	}

	options := &Options{
		BatchSize: *batchSize,
		Ops:       *ops,
		Endpoint:  *endpoint,
		AppNum:    *appNum,
	}

	traceGenerator := trace.NewGenerator()
	ingester, err := ingester.NewIngester(options.Endpoint, "", true)
	if err != nil {
		panic(err)
	}

	// Calculate the number of concurrent tasks
	concurrencyNum := options.Ops / options.BatchSize

	fmt.Printf("Starting benchmark with options: %+v\n", options)
	for {
		doBenchmark(options, traceGenerator, ingester, concurrencyNum)
	}
}

func doBenchmark(opts *Options, generator *trace.Generator, ingester *ingester.Ingester, concurrencyNum int) {
	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < concurrencyNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			appIndex := rand.Intn(opts.AppNum) + 1
			app := fmt.Sprintf("app%d", appIndex)

			data, err := generator.GenerateTraceRequest(app, time.Now().UnixNano(), opts.BatchSize/50, 50)
			if err != nil {
				fmt.Printf("generate traces failed: %v\n", err)
			}

			if err := ingester.IngestTraces(data); err != nil {
				fmt.Printf("ingest traces failed: %v\n", err)
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)

	// Sleep if the elapsed time is less than 1 second.
	if elapsed < time.Second*1 {
		time.Sleep(time.Second*1 - elapsed)
	}
}
