# 修复前故障复现（Docker）

## 项目与标准命令

项目使用 Go 1.26。标准验证命令为：

```bash
go test -count=20 ./...
```

## 环境构建与编译

本地 `go build ./...` 可以通过，故障出现在已取消的批量接样请求中。

## 故障触发步骤

执行标准验证命令，取消批量接样请求后检查返回错误和内存仓储数量。

## 实际错误输出

```text
batch error = <nil>, want context canceled
```

## 期望行为

请求已经取消时，服务应立即返回 `context canceled`，且不应写入任何样本。
