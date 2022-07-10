package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aliyun-sls/opentelemetry-go-provider-sls/provider"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	if len(os.Args) <= 8 {
		fmt.Printf("Usage: %v listen endpoint project instance service version AccessKeyID AccessKeySecret\n", os.Args[0])
		fmt.Printf("For example:\n")
		fmt.Printf("    %v 8089 otel.cn-beijing.log.aliyuncs.com:10010 otel ossrs stable v1.0.0 UJPI3Ad90g4Gxxxxxxxxxxxx k3ododEdFsGRdAgEwxxxxxxxxxxxxx\n", os.Args[0])
		os.Exit(-1)
	}
	listen := os.Args[1]
	endpoint := os.Args[2]
	project := os.Args[3]
	instance := os.Args[4]
	service := os.Args[5]
	version := os.Args[6]
	akID := os.Args[7]
	akSecret := os.Args[8]

	// Setup SLS Trace provider.
	slsConfig, err := provider.NewConfig(provider.WithServiceName(service),
		provider.WithServiceVersion(version),
		provider.WithTraceExporterEndpoint(endpoint),
		provider.WithMetricExporterEndpoint(endpoint),
		provider.WithSLSConfig(project, instance, akID, akSecret))
	// 如果初始化失败则panic，可以替换为其他错误处理方式
	if err != nil {
		panic(err)
	}
	if err := provider.Start(slsConfig); err != nil {
		panic(err)
	}
	defer provider.Shutdown(slsConfig)

	// HTTP Handler: /stable
	handleStable()

	// Start HTTP server.
	err = http.ListenAndServe(fmt.Sprintf(":%s", listen), nil)
	if err != nil {
		panic(err)
	}
}

func handleStable() {
	listen := os.Args[1]

	stableHandler := func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(900 * time.Millisecond)

		// 如果需要记录一些事件，可以获取Context中的span并添加Event（非必要步骤）
		ctx := req.Context()
		span := trace.SpanFromContext(ctx)
		span.AddEvent("say : Stable is SRS 4.0", trace.WithAttributes(label.KeyValue{
			Key: "label-key-3", Value: label.StringValue("label-value-4")},
		))

		// 创建新的ChildSpan
		dbRequest(ctx, span.Tracer())

		b, _ := json.Marshal(map[string]interface{}{
			"stable": "Stable is SRS 4.0!",
			"headers": req.Header,
		})
		_, _ = io.WriteString(w, string(b))
	}
	// 使用 otel net/http的自动注入方式，只需要使用otelhttp.NewHandler包裹http.Handler即可
	otelHandler := otelhttp.NewHandler(http.HandlerFunc(stableHandler), "Stable")
	http.Handle("/stable", otelHandler)
	/*
	http.Handle("/stable", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	opts := append([]trace.SpanOption{
		trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...),
		trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(r)...),
		trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest("XXX", "", r)...),
	}, h.spanStartOptions...) // start with the configured options
	}))*/
	fmt.Println(fmt.Sprintf("You can visit http://127.0.0.1:%v/stable .", listen))
}

func dbRequest(ctx context.Context, tracer trace.Tracer) {
	ctx, span := tracer.Start(ctx, "MySQL")
	span.AddEvent("DB is done")
	defer span.End()

	time.Sleep(800 * time.Millisecond)
}

