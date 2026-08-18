# 校准样本流转中心

## 项目目标

校准样本流转中心为实验室运营人员提供一个本地业务程序，用于登记待校准样本、记录审核意见、执行放行并追踪完整事件链。程序提供 CLI 工作流检查和 HTTP 适配层，使用内存存储，便于在开发和标注环境中重复运行，不依赖外部数据库或网络服务。

## 用户角色

- 接样员：登记样本来源、优先级和批次信息，查看当前处理状态。
- 复核员：检查样本元数据，记录复核意见，批准或退回样本。
- 放行员：确认已批准样本的放行条件，生成放行事件。

## 核心实体

- Sample：待校准样本，包含编号、来源、优先级、状态和标签。
- Review：复核记录，包含复核员、结论、意见和时间。
- Release：放行记录，包含放行员、目的地和时间。
- Event：不可变的业务事件，用于展示样本状态变化轨迹。

## 业务流程

### 登记样本

接样员通过 `POST /v1/intakes` 提交样本编号、来源、优先级和标签。服务校验编号唯一、来源非空、优先级合法后创建待复核样本，并写入接收事件。成功结果是返回样本及其初始事件；重复编号或非法字段返回明确的客户端错误。

### 审核样本

复核员通过 `POST /v1/intakes/{id}/review` 提交结论和意见。服务读取样本、确认当前状态为待复核，再写入复核记录并把样本转为已批准或已退回，同时追加审核事件。不存在的样本、重复审核或非法结论返回错误。

### 放行样本

放行员通过 `POST /v1/intakes/{id}/release` 提交目的地。服务确认样本已批准、目的地合法且尚未放行，创建放行记录并把样本转为已放行，同时追加放行事件。状态不满足条件时拒绝请求。

### 查询与追踪

运营人员通过 `GET /v1/intakes` 按状态和优先级筛选样本，通过 `GET /v1/intakes/{id}/timeline` 查看事件时间线。查询结果按登记时间稳定排序，时间线保持事件发生顺序。

## 状态与规则

- 样本状态依次为 `pending_review`、`approved`、`rejected`、`released`。
- 只有 `pending_review` 样本可以审核；只有 `approved` 样本可以放行。
- 样本编号在整个服务生命周期内唯一且大小写敏感。
- 优先级只能是 `routine`、`priority` 或 `urgent`。
- 审核意见不能为空；放行目的地必须是受支持的实验室区域名称。
- 每次状态变化都必须产生一个事件，事件只允许追加不允许修改。

## 接口与验证

- `go run ./cmd/calibrationd describe`：读取操作员名称并返回准备状态。
- `go run ./cmd/calibrationd check-intake`：验证登记工作流。
- `go run ./cmd/calibrationd check-review`：验证审核工作流。
- `go run ./cmd/calibrationd check-release`：验证放行工作流。
- `go run ./cmd/calibrationd check-query`：验证时间线查询工作流。
- `GET /healthz`：HTTP 适配层的健康检查。
- `POST /v1/intakes`：登记样本并返回 `201`。
- `POST /v1/intakes/{id}/review`：完成审核并返回 `200`。
- `POST /v1/intakes/{id}/release`：完成放行并返回 `200`。
- `GET /v1/intakes`：按状态、优先级筛选样本。
- `GET /v1/intakes/{id}/timeline`：读取样本事件时间线。

所有接口都通过公开的 `go test ./...` 覆盖核心规则，并由 `.go-annotation-runtime.json` 声明真实 HTTP 运行检查。
