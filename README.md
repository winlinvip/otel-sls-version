# otel-sls-version

OpenTelemetry demo for SLS Trace as version service.

## Call Graph

* api-server: http://localhost:8088/hello

## Usage

Build and start api-server:

```bash
(cd api-server && go build . && ./api-server 8088 otel.cn-beijing.log.aliyuncs.com:10010 otel ossrs versions v1.0.0 UJPI3Ad90g4Gxxxxxxxxxxxx k3ododEdFsGRdAgEwxxxxxxxxxxxxx)
```

