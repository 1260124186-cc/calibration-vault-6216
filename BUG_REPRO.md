# 修复前故障复现（Docker）

## 项目与标准命令
项目为校准样本流转中心，标准命令为 `go test -count=1 -run TestCanceledOperationsReportReturnsRequestCancelled ./internal/httpapi`。

## 环境构建与编译
使用 `golang:1.26.2` 官方镜像在 linux/arm64 平台验证：`docker build` 通过，镜像内 `go version` 为 go1.26.2；`go build ./...` 通过；`go test ./...` 如预期失败。

## 故障触发步骤
创建已取消的请求上下文，访问运营报告接口。

## 实际错误输出
```
?   	example.com/calibration-vault/cmd/calibrationd	[no test files]
?   	example.com/calibration-vault/internal/domain	[no test files]
--- FAIL: TestCanceledOperationsReportReturnsRequestCancelled (0.00s)
    canceled_report_test.go:21: status = 200, want 408; body={"metrics":{"summary":{"total":0,"by_status":{},"by_priority":{},"open_count":0,"released_count":0},"tag_counts":{},"oldest_open_age":0,"generated_at":"2026-08-19T01:51:59.926044801Z"},"urgent_samples":null,"pending_review":null}
FAIL
FAIL	example.com/calibration-vault/internal/httpapi	0.010s
?   	example.com/calibration-vault/internal/service	[no test files]
?   	example.com/calibration-vault/internal/store	[no test files]
FAIL
```

## 期望行为
请求取消后应尽快停止业务处理并返回 request cancelled。
