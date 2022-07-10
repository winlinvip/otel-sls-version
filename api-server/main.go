package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/aliyun-sls/opentelemetry-go-provider-sls/provider"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	if len(os.Args) <= 8 {
		fmt.Printf("Usage: %v listen endpoint project instance service version AccessKeyID AccessKeySecret stable-api\n", os.Args[0])
		fmt.Printf("For example:\n")
		fmt.Printf("    %v 8088 otel.cn-beijing.log.aliyuncs.com:10010 otel ossrs versions v1.0.0 UJPI3Ad90g4Gxxxxxxxxxxxx k3ododEdFsGRdAgEwxxxxxxxxxxxxx\n", os.Args[0])
		fmt.Printf("    %v 8088 otel.cn-beijing.log.aliyuncs.com:10010 otel ossrs versions v1.0.0 UJPI3Ad90g4Gxxxxxxxxxxxx k3ododEdFsGRdAgEwxxxxxxxxxxxxx http://127.0.0.1:8089/stable\n", os.Args[0])
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

	// Upstream API is optional.
	var stableAPI string
	if len(os.Args) > 9 {
		stableAPI = os.Args[9]
	}

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

	// HTTP Handler: /hello
	handleHello(stableAPI)

	// Start HTTP server.
	err = http.ListenAndServe(fmt.Sprintf(":%s", listen), nil)
	if err != nil {
		panic(err)
	}
}

func handleHello(stableAPI string) {
	listen := os.Args[1]

	// 注册一个Metric指标（非必要步骤）
	labels := []label.KeyValue{
		label.String("label1", "value1"),
	}
	meter := otel.Meter("aliyun.sls")
	sayDavidCount := metric.Must(meter).NewInt64Counter("say_david_count")

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(700 * time.Millisecond)

		// 如果需要记录一些事件，可以获取Context中的span并添加Event（非必要步骤）
		ctx := req.Context()
		span := trace.SpanFromContext(ctx)
		span.AddEvent("say : Hello, I am david", trace.WithAttributes(label.KeyValue{
			Key: "label-key-1", Value: label.StringValue("label-value-1")},
		))

		var stableAPIBody []byte
		if stableAPI != "" {
			client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
			req, _ := http.NewRequestWithContext(ctx, "GET", stableAPI, nil)
			res, _ := client.Do(req)
			defer res.Body.Close()
			stableAPIBody, _ = ioutil.ReadAll(res.Body)
		}

		_, _ = io.WriteString(w, fmt.Sprintf("Hello, I am david! %s\n", stableAPIBody))
		sayDavidCount.Add(req.Context(), 1, labels...)
	}
	// 使用 otel net/http的自动注入方式，只需要使用otelhttp.NewHandler包裹http.Handler即可
	otelHandler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "Hello")
	http.Handle("/hello", otelHandler)
	fmt.Println(fmt.Sprintf("You can visit http://127.0.0.1:%v/hello .", listen))
}

