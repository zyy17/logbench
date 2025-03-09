package ingester

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

type Ingester struct {
	endpoint   string
	db         string
	pipeline   string
	enableGzip bool
}

func NewIngester(endpoint, db, pipeline string, enableGzip bool) (*Ingester, error) {
	return &Ingester{
		endpoint:   endpoint,
		db:         db,
		pipeline:   pipeline,
		enableGzip: enableGzip,
	}, nil
}

func (i *Ingester) Ingest(table string, input []byte) error {
	url := i.eventURL(table)

	var buffer bytes.Buffer
	if i.enableGzip {
		writer := gzip.NewWriter(&buffer)
		if _, err := writer.Write(input); err != nil {
			return err
		}
		writer.Close()
	} else {
		buffer.Write(input)
	}

	req, err := http.NewRequest("POST", url, &buffer)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if i.enableGzip {
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

func (i *Ingester) eventURL(table string) string {
	return fmt.Sprintf("%s/v1/events/logs?db=%s&table=%s&pipeline_name=%s", i.endpoint, i.db, table, i.pipeline)
}
