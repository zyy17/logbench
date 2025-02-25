package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/openzipkin/zipkin-go/idgenerator"
)

const (
	JSONLogFormatWithTraceExtraBytes = `{"host":"%s", "user-identifier":"%s", "datetime":"%s", "method": "%s", "request": "%s", "protocol":"%s", "status":%d, "bytes":%d, "referer": "%s", "traceId": "%s", "spanId": "%s", "dummyBytes": "%s"}`
)

// Generates NDJSON test data.
func generateLogs(lines int, totalSize int, payloadSize int, start time.Time, interval time.Duration) []byte {
	var results []byte
	var currentSize int
	for i := 0; i < lines; i++ {
		log := newJSONLogFormatWithTrace(start.Add(interval*time.Duration(i)), payloadSize) + "\n"
		results = append(results, []byte(log)...)
		currentSize += len(log)
		if currentSize >= totalSize {
			break
		}
	}
	return results
}

const (
	ClickHouse = "02/Jan/2006:15:04:05 -0700"
)

func newJSONLogFormatWithTrace(t time.Time, payloadSize int) string {
	g := idgenerator.NewRandom128()
	traceId := g.TraceID()

	return fmt.Sprintf(
		JSONLogFormatWithTraceExtraBytes,
		gofakeit.IPv4Address(),
		randAuthUserID(),
		t.Format(ClickHouse),
		gofakeit.HTTPMethod(),
		randResourceURI(),
		randHTTPVersion(),
		gofakeit.StatusCode(),
		gofakeit.Number(0, 30000),
		gofakeit.URL(),
		traceId.String(),
		g.SpanID(traceId).String(),
		randomWords(payloadSize),
	)
}

// RandResourceURI generates a random resource URI
func randResourceURI() string {
	var uri string
	num := gofakeit.Number(1, 4)
	for i := 0; i < num; i++ {
		uri += "/" + url.QueryEscape(gofakeit.BS())
	}
	uri = strings.ToLower(uri)
	return uri
}

// RandAuthUserID generates a random auth user id
func randAuthUserID() string {
	candidates := []string{"-", strings.ToLower(gofakeit.Username())}
	return candidates[rand.Intn(2)]
}

// RandHTTPVersion returns a random http version
func randHTTPVersion() string {
	versions := []string{"HTTP/1.0", "HTTP/1.1", "HTTP/2.0"}
	return versions[rand.Intn(3)]
}
