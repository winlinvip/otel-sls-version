# otel-sls-version

OpenTelemetry demo for SLS Trace as version service.

## Usage

Build and start api-server service:

```bash
(cd api-server && go build . && ./api-server 8088 otel.cn-beijing.log.aliyuncs.com:10010 otel ossrs versions v1.0.0 UJPI3Ad90g4Gxxxxxxxxxxxx k3ododEdFsGRdAgEwxxxxxxxxxxxxx http://127.0.0.1:8089/stable)
```

Build and start stable-version service:

```bash
(cd stable-version && go build . && ./stable-version 8089 otel.cn-beijing.log.aliyuncs.com:10010 otel ossrs stable v1.0.0 UJPI3Ad90g4Gxxxxxxxxxxxx k3ododEdFsGRdAgEwxxxxxxxxxxxxx)
```

Then, open the first link in browser:

* api-server: http://localhost:8088/hello
  * stable-version: http://localhost:8089/stable

Now, we're able to use the [SLS Trace](https://sls.console.aliyun.com/lognext/trace/otel/ossrs?resource=/trace/ossrs/explorer).

