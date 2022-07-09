package main

import (
	"fmt"
	"io"
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
	if len(os.Args) <= 7 {
		fmt.Printf("Usage: %v endpoint project instance service version AccessKeyID AccessKeySecret\n", os.Args[0])
		fmt.Printf("For example:\n")
		fmt.Printf("    %v otel.cn-beijing.log.aliyuncs.com:10010 otel ossrs versions v1.0.0 UJPI3Ad90g4Gxxxxxxxxxxxx k3ododEdFsGRdAgEwxxxxxxxxxxxxx\n", os.Args[0])
		os.Exit(-1)
	}
	endpoint := os.Args[1]
	project := os.Args[2]
	instance := os.Args[3]
	service := os.Args[4]
	version := os.Args[5]
	akID := os.Args[6]
	akSecret := os.Args[7]

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

	// 注册一个Metric指标（非必要步骤）
	labels := []label.KeyValue{
		label.String("label1", "value1"),
	}
	meter := otel.Meter("aliyun.sls")
	sayDavidCount := metric.Must(meter).NewInt64Counter("say_david_count")

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		if time.Now().Unix()%10 == 0 {
			_, _ = io.WriteString(w, "Hello, world!\n")
		} else {
			// 如果需要记录一些事件，可以获取Context中的span并添加Event（非必要步骤）
			ctx := req.Context()
			span := trace.SpanFromContext(ctx)
			span.AddEvent("say : Hello, I am david", trace.WithAttributes(label.KeyValue{
				Key:   "label-key-1",
				Value: label.StringValue("label-value-1"),
			}))

			_, _ = io.WriteString(w, "Hello, I am david!\n")
			sayDavidCount.Add(req.Context(), 1, labels...)
		}
	}

	// 使用 otel net/http的自动注入方式，只需要使用otelhttp.NewHandler包裹http.Handler即可
	otelHandler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "Hello")

	http.Handle("/hello", otelHandler)
	fmt.Println("Now listen port 8088, you can visit 127.0.0.1:8088/hello .")
	err = http.ListenAndServe(":8088", nil)
	if err != nil {
		panic(err)
	}
}

