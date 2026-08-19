# 修复前故障复现（Docker）

## 项目与标准命令
项目为校准样本流转中心，标准命令为 `go test -count=1 -run TestBatchValidationReturnsBadRequest ./internal/httpapi`。

## 环境构建与编译
使用 `golang:1.26.2` 官方镜像在 linux/arm64 平台验证：`docker build` 通过，镜像内 `go version` 为 go1.26.2；`go build ./...` 通过；`go test ./...` 如预期失败。

## 故障触发步骤
向批量登记接口提交空的 items 列表。

## 实际错误输出
```
?   	example.com/calibration-vault/cmd/calibrationd	[no test files]
?   	example.com/calibration-vault/internal/domain	[no test files]
--- FAIL: TestBatchValidationReturnsBadRequest (0.00s)
    batch_validation_test.go:19: status = 500, want 400; body={"code":"internal_error","message":"unexpected service error"}
FAIL
FAIL	example.com/calibration-vault/internal/httpapi	0.002s
?   	example.com/calibration-vault/internal/service	[no test files]
?   	example.com/calibration-vault/internal/store	[no test files]
FAIL
```

## 期望行为
非法批量登记输入应返回客户端请求错误，而不是内部服务错误。
