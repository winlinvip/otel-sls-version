# otel-sls-version

OpenTelemetry demo for SLS Trace as version service.

## Call Graph

* api-server: http://localhost:8088/hello
  * stable-version: http://localhost:8089/stable

## Usage

Build and start api-server service:

```bash
(cd api-server && go build . && ./api-server 8088 otel.cn-beijing.log.aliyuncs.com:10010 otel ossrs versions v1.0.0 UJPI3Ad90g4Gxxxxxxxxxxxx k3ododEdFsGRdAgEwxxxxxxxxxxxxx http://127.0.0.1:8089/stable)
```

Build and start stable-version service:

```bash
(cd stable-version && go build . && ./stable-version 8089 otel.cn-beijing.log.aliyuncs.com:10010 otel ossrs stable v1.0.0 UJPI3Ad90g4Gxxxxxxxxxxxx k3ododEdFsGRdAgEwxxxxxxxxxxxxx)
```

