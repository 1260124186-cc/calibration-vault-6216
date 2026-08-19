# calibration-vault__009 Docker 交付说明

## 项目概览
- 校准样本流转中心是一个使用 Go 1.26 编写的本地业务程序，帮助接样员、复核员和放行员管理样本从登记到放行的状态流转，并提供可追踪的事件时间线。程序同时提供 HTTP 适配层和用于本地流程检查的 CLI 入口。
- Go module: `example.com/calibration-vault`

## 标准命令

```bash
go build ./...
go test ./...
```

## 实际启动入口

```bash
go run ./cmd/calibrationd
```

## Docker 构建

```bash
./build_benzhi_docker.sh calibration-vault__009-benzhi linux/amd64
docker run --rm -it calibration-vault__009-benzhi bash
```

## 环境

- 基础镜像: `golang:1.26.2`
- 依赖在镜像构建阶段预下载，容器内可直接执行 Go 构建和测试命令。
- 代码目录: `/app`
- 源码中检测到的服务端口: `18080`
