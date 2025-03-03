package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Generates NDJSON test data.
func generateLogs(lines int, totalSize int, payloadSize int, start time.Time, interval time.Duration, appNum int) []byte {
	var results []byte
	var currentSize int
	for i := 0; i < lines; i++ {
		timestamp := start.Add(interval * time.Duration(i))
		log := newTestLog(timestamp, payloadSize, appNum) + "\n"
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
	RFC5424    = "2006-01-02T15:04:05.000Z"
)

type podMetadata struct {
	PodNamespace  string
	PodName       string
	PodNodeName   string
	PodIP         string
	PodUID        string
	AppName       string
	ContainerName string
}

func newPodMetadataFromEnv() *podMetadata {
	return &podMetadata{
		PodNamespace:  os.Getenv("POD_NAMESPACE"),
		PodName:       os.Getenv("POD_NAME"),
		PodNodeName:   os.Getenv("POD_NODE_NAME"),
		PodIP:         os.Getenv("POD_IP"),
		PodUID:        os.Getenv("POD_UID"),
		AppName:       os.Getenv("APP_NAME"),
		ContainerName: os.Getenv("CONTAINER_NAME"),
	}
}

func newTestLog(t time.Time, payloadSize int, appNum int) string {
	const (
		testLogFormat = `{"timestamp": "%s", "kubernetes.container_name":"%s", "kubernetes.pod_labels.app":"%s", "kubernetes.pod_namespace":"%s", "kubernetes.pod_node_name": "%s", "kubernetes.pod_ip": "%s", "kubernetes.pod_name":"%s", "kubernetes.pod_uid":"%s", "message": "%s"}`
	)

	podMetadata := newPodMetadataFromEnv()

	return fmt.Sprintf(
		testLogFormat,
		t.Format(RFC5424),
		podMetadata.ContainerName,
		generateRandomAppName(appNum),
		podMetadata.PodNamespace,
		podMetadata.PodNodeName,
		podMetadata.PodIP,
		podMetadata.PodName,
		podMetadata.PodUID,
		randomWords(payloadSize),
	)
}

func generateRandomAppName(appNum int) string {
	return fmt.Sprintf("app-%d", rand.Intn(appNum))
}
