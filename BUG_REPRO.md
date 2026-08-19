# 修复前故障复现（Docker）

## 项目与标准命令
项目为校准样本流转中心，标准命令为 `go test -count=1 -run TestBatchIntakeAssignsDistinctEventIDs ./internal/service`。

## 环境构建与编译
使用 `golang:1.26.2` 官方镜像在 linux/arm64 平台验证：`docker build` 通过，镜像内 `go version` 为 go1.26.2；`go build ./...` 通过；`go test ./...` 如预期失败。

## 故障触发步骤
提交包含两个样本的批量接样请求，再分别读取两个样本的事件时间线。

## 实际错误输出
```
?   	example.com/calibration-vault/cmd/calibrationd	[no test files]
?   	example.com/calibration-vault/internal/domain	[no test files]
?   	example.com/calibration-vault/internal/httpapi	[no test files]
--- FAIL: TestBatchIntakeAssignsDistinctEventIDs (0.00s)
    batch_event_ids_test.go:34: batch event IDs are both "evt-000001", want distinct event IDs
FAIL
FAIL	example.com/calibration-vault/internal/service	0.010s
?   	example.com/calibration-vault/internal/store	[no test files]
FAIL
```

## 期望行为
同一批次中每个样本的事件编号都应唯一且可追踪。
