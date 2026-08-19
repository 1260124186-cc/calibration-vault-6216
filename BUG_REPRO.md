# 修复前故障复现（Docker）

## 项目与标准命令
项目为校准样本流转中心，标准命令为 `go test -count=1 -run TestTimelineForUnknownSampleReturnsNotFound ./internal/httpapi`。

## 环境构建与编译
使用 `golang:1.26.2` 官方镜像在 linux/arm64 平台验证：`docker build` 通过，镜像内 `go version` 为 go1.26.2；`go build ./...` 通过；`go test ./...` 如预期失败。

## 故障触发步骤
请求一个不存在样本的时间线接口。

## 实际错误输出
```
?   	example.com/calibration-vault/cmd/calibrationd	[no test files]
?   	example.com/calibration-vault/internal/domain	[no test files]
--- FAIL: TestTimelineForUnknownSampleReturnsNotFound (0.00s)
    timeline_not_found_test.go:18: status = 500, want 404; body={"code":"internal_error","message":"unexpected service error"}
FAIL
FAIL	example.com/calibration-vault/internal/httpapi	0.010s
?   	example.com/calibration-vault/internal/service	[no test files]
?   	example.com/calibration-vault/internal/store	[no test files]
FAIL
```

## 期望行为
不存在的样本应返回资源不存在，避免被监控误判为内部服务故障。
