# 修复前故障复现（Docker）

## 项目与标准命令
项目为校准样本流转中心，标准命令为 `go test -count=1 -run TestOperationsReportIncludesPendingUrgentSample ./internal/service`。

## 环境构建与编译
使用 `golang:1.26.2` 官方镜像在 linux/arm64 平台验证：`docker build` 通过，镜像内 `go version` 为 go1.26.2；`go build ./...` 通过；`go test ./...` 如预期失败。

## 故障触发步骤
登记一个 priority 为 urgent 的样本但不审核，然后生成运营报告。

## 实际错误输出
```
?   	example.com/calibration-vault/cmd/calibrationd	[no test files]
?   	example.com/calibration-vault/internal/domain	[no test files]
?   	example.com/calibration-vault/internal/httpapi	[no test files]
--- FAIL: TestOperationsReportIncludesPendingUrgentSample (0.00s)
    pending_urgent_report_test.go:26: UrgentSamples = []domain.Sample(nil), want pending urgent sample
FAIL
FAIL	example.com/calibration-vault/internal/service	0.001s
?   	example.com/calibration-vault/internal/store	[no test files]
FAIL
```

## 期望行为
仍处于开放流转状态的紧急样本应出现在运营报告中，并计入开放样本数量。
