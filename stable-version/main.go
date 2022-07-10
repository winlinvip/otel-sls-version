package main

import (
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
		if time.Now().Unix()%3 == 0 {
			_, _ = io.WriteString(w, "Stable, world!\n")
		} else {
			// 如果需要记录一些事件，可以获取Context中的span并添加Event（非必要步骤）
			ctx := req.Context()
			span := trace.SpanFromContext(ctx)
			span.AddEvent("say : Stable is SRS 4.0", trace.WithAttributes(label.KeyValue{
				Key: "label-key-3", Value: label.StringValue("label-value-4")},
			))

			_, _ = io.WriteString(w, "Stable is SRS 4.0!\n")
		}
	}
	// 使用 otel net/http的自动注入方式，只需要使用otelhttp.NewHandler包裹http.Handler即可
	otelHandler := otelhttp.NewHandler(http.HandlerFunc(stableHandler), "Stable")
	http.Handle("/stable", otelHandler)
	fmt.Println(fmt.Sprintf("You can visit http://127.0.0.1:%v/stable .", listen))
}

