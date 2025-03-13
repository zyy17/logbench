package trace

import (
	"math/rand"
	"net/url"
	"strings"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	collectorpb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/protobuf/proto"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) GenerateTraceRequest(serviceName string, timestamp int64, traceNum int, spanNum int) ([]byte, error) {
	var spans []*tracepb.Span
	for i := 0; i < traceNum; i++ {
		spans = append(spans, g.generateSpans(g.generateTraceID(), spanNum, timestamp)...)
	}

	request := &collectorpb.ExportTraceServiceRequest{
		ResourceSpans: []*tracepb.ResourceSpans{
			{
				Resource: &resourcepb.Resource{
					Attributes: g.generateResourceAttributes(serviceName),
				},
				ScopeSpans: []*tracepb.ScopeSpans{
					{
						Spans: spans,
					},
				},
			},
		},
	}

	data, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (g *Generator) generateSpans(traceID string, spanNum int, startTimestamp int64) []*tracepb.Span {
	spans := make([]*tracepb.Span, 0, spanNum)
	for i := 0; i < spanNum; i++ {
		spans = append(spans, g.generateSpan(traceID, startTimestamp))
	}
	return spans
}

func (g *Generator) generateSpan(traceID string, startTimestamp int64) *tracepb.Span {
	span := &tracepb.Span{
		TraceId:           []byte(traceID),
		SpanId:            []byte(gofakeit.UUID()),
		ParentSpanId:      []byte(gofakeit.UUID()),
		Kind:              tracepb.Span_SPAN_KIND_INTERNAL,
		StartTimeUnixNano: uint64(startTimestamp),
		EndTimeUnixNano:   uint64(startTimestamp) + uint64(gofakeit.Number(1000000, 1000000000)),
		Name:              g.randomResourceURI(),
		Attributes:        g.generateSpanAttributes(),
		Status:            g.generateSpanStatus(),
	}

	return span
}

func (g *Generator) generateSpanAttributes() []*commonpb.KeyValue {
	return []*commonpb.KeyValue{
		{
			Key:   "http.method",
			Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: gofakeit.HTTPMethod()}},
		},
		{
			Key:   "http.url",
			Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: gofakeit.URL()}},
		},
		{
			Key:   "http.status_code",
			Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: int64(gofakeit.Number(200, 599))}},
		},
		{
			Key:   "http.user_agent",
			Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: gofakeit.UserAgent()}},
		},
		{
			Key:   "http.request.body.size",
			Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: int64(gofakeit.Number(100, 1000000))}},
		},
		{
			Key:   "http.response.body.size",
			Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: int64(gofakeit.Number(100, 1000000))}},
		},
		{
			Key:   "http.client_ip",
			Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: gofakeit.IPv4Address()}},
		},
		{
			Key:   "http.server_ip",
			Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: gofakeit.IPv4Address()}},
		},
		{
			Key:   "http.route",
			Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: gofakeit.URL()}},
		},
	}
}

func (g *Generator) generateResourceAttributes(serviceName string) []*commonpb.KeyValue {
	return []*commonpb.KeyValue{
		{
			Key:   "service.name",
			Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: serviceName}},
		},
		{
			Key:   "service.version",
			Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: g.chooseFromStrings([]string{"1.0.0", "1.0.1", "1.0.2", "1.0.3"})}},
		},
		{
			Key:   "deployment.environment",
			Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: g.chooseFromStrings([]string{"dev", "test", "sit", "prod"})}},
		},
		{
			Key:   "deployment.region",
			Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: g.chooseFromStrings([]string{"us-east-1", "us-east-2", "us-west-1", "us-west-2"})}},
		},
		{
			Key:   "deployment.zone",
			Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: g.chooseFromStrings([]string{"us-east-1a", "us-east-1b", "us-east-1c", "us-east-1d"})}},
		},
	}
}

func (g *Generator) generateTraceID() string {
	return uuid.New().String()
}

func (g *Generator) chooseFromStrings(strings []string) string {
	return strings[rand.Intn(len(strings))]
}

func (g *Generator) randomResourceURI() string {
	var uri string
	num := gofakeit.Number(1, 4)
	for i := 0; i < num; i++ {
		uri += "/" + url.QueryEscape(gofakeit.BS())
	}
	uri = strings.ToLower(uri)
	return uri
}

func (g *Generator) generateSpanStatus() *tracepb.Status {
	return &tracepb.Status{
		Code:    tracepb.Status_STATUS_CODE_OK,
		Message: "OK",
	}
}
