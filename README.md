# 校准样本流转中心

校准样本流转中心是一个使用 Go 1.26 编写的本地业务程序，帮助接样员、复核员和放行员管理样本从登记到放行的状态流转，并提供可追踪的事件时间线。程序同时提供 HTTP 适配层和用于本地流程检查的 CLI 入口。

## 目录结构

- `cmd/calibrationd`：CLI 与 HTTP 服务启动入口。
- `internal/domain`：样本、审核、放行和事件实体，以及状态规则。
- `internal/store`：线程安全的内存存储。
- `internal/service`：登记、审核、放行和查询业务流程。
- `internal/httpapi`：HTTP 路由、请求解析和响应编码。

## 运行

```bash
go run ./cmd/calibrationd describe
go run ./cmd/calibrationd check-intake
```

`describe` 从标准输入读取操作员名称并输出准备状态。无子命令启动时，程序监听 `127.0.0.1:18080`，可使用 `CALIBRATION_ADDR` 覆盖监听地址。

## 验证

```bash
go build ./...
go test ./...
```

健康检查地址为 `GET /healthz`。业务入口包括 `POST /v1/intakes`、`POST /v1/intakes/{id}/review`、`POST /v1/intakes/{id}/release`、`GET /v1/intakes` 和 `GET /v1/intakes/{id}/timeline`。

## 数据约束

样本优先级只能使用 `routine`、`priority` 或 `urgent`。样本需要先审核批准，才能被放行到 `lab-east`、`lab-west` 或 `lab-central`。
