# 统一运维管理平台 API 接口层设计

## 文档信息

| 项目 | 内容 |
|------|------|
| 文档版本 | v1.0 |
| 更新日期 | 2026-06-11 |
| 当前实现状态 | M0-M2 基本功能已实现；M3-M7 尚未开始开发 |
| 适用范围 | 一期统一运维管理平台 API 设计 |

## 1. 设计原则

- API 基础路径统一为 `/api/v1`，健康检查使用 `/healthz`。
- 业务 API 默认使用 JWT Bearer Token 认证；CMDB 外部同步接口使用 API Key + HMAC 签名。
- 返回体使用统一结构，列表接口支持分页、筛选和排序。
- M0-M2 接口以当前代码和迁移为准；M3-M7 接口为设计稿，开发时需结合实现计划落地。

## 2. 通用约定

### 2.1 认证方式

| 场景 | Header | 说明 |
|------|--------|------|
| 管理端 API | `Authorization: Bearer <token>` | 登录后获取 JWT |
| CMDB 同步 | `X-API-Key: <key_id>`、`X-Signature: <hmac>` | HMAC-SHA256，请求体作为签名内容 |

### 2.2 响应格式

```json
{
  "data": {},
  "message": "ok"
}
```

列表响应建议：

```json
{
  "data": {
    "items": [],
    "total": 0,
    "page": 1,
    "page_size": 20
  },
  "message": "ok"
}
```

### 2.3 通用状态码

| 状态码 | 含义 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或认证失败 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 409 | 唯一约束或状态冲突 |
| 500 | 服务端错误 |

## 3. M0-M2 已实现接口

### 3.1 健康检查

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/healthz` | 无 | 服务健康检查 |

### 3.2 认证与用户

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| POST | `/api/v1/auth/login` | 无 | 用户登录，返回 JWT |
| GET | `/api/v1/auth/users` | `auth.user.read` | 查询用户列表 |
| POST | `/api/v1/auth/users` | `auth.user.write` | 创建用户 |
| PUT | `/api/v1/auth/users/{id}` | `auth.user.write` | 更新用户资料和角色 |
| DELETE | `/api/v1/auth/users/{id}` | `auth.user.write` | 删除用户 |

### 3.3 角色与权限

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/v1/auth/roles` | `auth.role.read` | 查询角色列表 |
| POST | `/api/v1/auth/roles` | `auth.role.write` | 创建角色 |
| PUT | `/api/v1/auth/roles/{name}` | `auth.role.write` | 更新角色和权限 |
| DELETE | `/api/v1/auth/roles/{name}` | `auth.role.write` | 删除角色 |
| GET | `/api/v1/auth/permissions` | `auth.role.read` | 查询权限列表 |

### 3.4 审计日志与系统配置

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/v1/audit/logs` | `audit.log.read` | 查询 API 操作审计日志 |
| GET | `/api/v1/system/configs` | `system.config.read` | 查询系统配置 |
| PUT | `/api/v1/system/configs/{key}` | `system.config.write` | 更新系统配置 |

### 3.5 CMDB 模型管理

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/v1/cmdb/model-groups` | `cmdb.model.read` | 查询模型组 |
| POST | `/api/v1/cmdb/model-groups` | `cmdb.model.write` | 创建模型组 |
| GET | `/api/v1/cmdb/model-groups/{id}/models` | `cmdb.model.read` | 查询模型组下模型 |
| POST | `/api/v1/cmdb/model-groups/{id}/models` | `cmdb.model.write` | 创建模型 |
| GET | `/api/v1/cmdb/models/{id}` | `cmdb.model.read` | 查询模型详情 |
| PUT | `/api/v1/cmdb/models/{id}` | `cmdb.model.write` | 更新模型、字段和关系 |
| DELETE | `/api/v1/cmdb/models/{id}` | `cmdb.model.write` | 删除模型 |

### 3.6 CMDB 资产管理

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/v1/cmdb/asset-groups` | `cmdb.asset.read` | 查询资产分组 |
| POST | `/api/v1/cmdb/asset-groups` | `cmdb.asset.write` | 创建资产分组 |
| GET | `/api/v1/cmdb/assets` | `cmdb.asset.read` | 查询资产列表，支持模型、分组、状态筛选 |
| POST | `/api/v1/cmdb/assets` | `cmdb.asset.write` | 创建资产 |
| GET | `/api/v1/cmdb/assets/{id}` | `cmdb.asset.read` | 查询资产详情 |
| PUT | `/api/v1/cmdb/assets/{id}` | `cmdb.asset.write` | 更新资产 |
| DELETE | `/api/v1/cmdb/assets/{id}` | `cmdb.asset.write` | 删除资产 |
| GET | `/api/v1/cmdb/assets/{id}/change-logs` | `cmdb.asset.read` | 查询资产变更历史 |
| POST | `/api/v1/cmdb/assets/import/preview` | `cmdb.asset.write` | CSV 导入预览和校验 |

### 3.7 CMDB API Key 与外部同步

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/v1/cmdb/api-keys` | `cmdb.apikey.read` | 查询 API Key 列表 |
| POST | `/api/v1/cmdb/api-keys` | `cmdb.apikey.write` | 创建 API Key |
| DELETE | `/api/v1/cmdb/api-keys/{id}` | `cmdb.apikey.write` | 吊销 API Key |
| POST | `/api/v1/cmdb/assets/sync` | API Key 签名 | 外部系统增量/全量同步资产 |

## 4. M3 监控系统接口设计

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/v1/monitoring/targets` | `monitoring.target.read` | 查询监控目标 |
| POST | `/api/v1/monitoring/targets` | `monitoring.target.write` | 创建监控目标 |
| GET | `/api/v1/monitoring/targets/{id}/config` | `monitoring.target.read` | 查询采集配置 |
| PUT | `/api/v1/monitoring/targets/{id}/config` | `monitoring.target.write` | 更新采集配置 |
| POST | `/api/v1/monitoring/ingest` | Agent Token | Agent/Collector 上报指标样本 |
| GET | `/api/v1/monitoring/metrics` | `monitoring.metric.read` | 按资产查询可用指标 |
| GET | `/api/v1/monitoring/query` | `monitoring.metric.read` | 查询时序指标数据 |
| GET | `/api/v1/monitoring/dashboards` | `monitoring.dashboard.read` | 查询仪表盘列表 |
| POST | `/api/v1/monitoring/dashboards` | `monitoring.dashboard.write` | 创建仪表盘 |
| GET | `/api/v1/monitoring/dashboards/{id}` | `monitoring.dashboard.read` | 查询仪表盘详情 |
| PUT | `/api/v1/monitoring/dashboards/{id}` | `monitoring.dashboard.write` | 更新仪表盘 |
| DELETE | `/api/v1/monitoring/dashboards/{id}` | `monitoring.dashboard.write` | 删除仪表盘 |

## 5. M4 告警平台接口设计

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/v1/alerting/rules` | `alerting.rule.read` | 查询告警规则 |
| POST | `/api/v1/alerting/rules` | `alerting.rule.write` | 创建告警规则 |
| GET | `/api/v1/alerting/rules/{id}` | `alerting.rule.read` | 查询规则详情 |
| PUT | `/api/v1/alerting/rules/{id}` | `alerting.rule.write` | 更新规则 |
| DELETE | `/api/v1/alerting/rules/{id}` | `alerting.rule.write` | 删除规则 |
| POST | `/api/v1/alerting/rules/{id}/enable` | `alerting.rule.write` | 启用规则 |
| POST | `/api/v1/alerting/rules/{id}/disable` | `alerting.rule.write` | 禁用规则 |
| GET | `/api/v1/alerting/events` | `alerting.event.read` | 查询告警事件列表 |
| GET | `/api/v1/alerting/events/{id}` | `alerting.event.read` | 查询告警详情 |
| POST | `/api/v1/alerting/events/{id}/ack` | `alerting.event.write` | 认领告警 |
| POST | `/api/v1/alerting/events/{id}/close` | `alerting.event.write` | 手工关闭告警 |
| POST | `/api/v1/alerting/events/{id}/notes` | `alerting.event.write` | 添加告警备注 |
| GET | `/api/v1/alerting/silences` | `alerting.silence.read` | 查询静默规则 |
| POST | `/api/v1/alerting/silences` | `alerting.silence.write` | 创建静默规则 |
| DELETE | `/api/v1/alerting/silences/{id}` | `alerting.silence.write` | 删除静默规则 |
| GET | `/api/v1/alerting/overview` | `alerting.event.read` | 告警大盘统计 |

## 6. M5 联系人与工单流程接口设计

### 6.1 联系人与通知

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/v1/notification/contacts` | `notification.contact.read` | 查询联系人 |
| POST | `/api/v1/notification/contacts` | `notification.contact.write` | 创建联系人 |
| GET | `/api/v1/notification/contacts/{id}` | `notification.contact.read` | 查询联系人详情 |
| PUT | `/api/v1/notification/contacts/{id}` | `notification.contact.write` | 更新联系人和通知渠道 |
| DELETE | `/api/v1/notification/contacts/{id}` | `notification.contact.write` | 删除联系人 |
| GET | `/api/v1/notification/logs` | `notification.log.read` | 查询通知记录 |
| POST | `/api/v1/notification/test` | `notification.contact.write` | 测试通知渠道 |

### 6.2 工单流程管理

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/v1/workorders` | `workorder.read` | 查询工单列表 |
| GET | `/api/v1/workorders/{id}` | `workorder.read` | 查询工单详情 |
| POST | `/api/v1/workorders` | `workorder.write` | 手工创建工单 |
| POST | `/api/v1/workorders/{id}/start` | `workorder.write` | 待处理 → 处理中 |
| POST | `/api/v1/workorders/{id}/complete` | `workorder.write` | 处理中 → 已完成 |
| POST | `/api/v1/workorders/{id}/close` | `workorder.write` | 已完成 → 已闭环，并关闭关联告警 |
| POST | `/api/v1/workorders/{id}/logs` | `workorder.write` | 添加处理记录 |
| GET | `/api/v1/workorders/overview` | `workorder.read` | 工单基础统计 |

## 7. M6-M7 验证与部署接口

M6-M7 以测试、UAT、部署交接为主，不新增业务域核心接口。建议保留以下运维接口：

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/healthz` | 无 | 存活检查 |
| GET | `/readyz` | 无 | 依赖就绪检查，M7 可选新增 |
| GET | `/metrics` | 内网访问 | 平台自身 Prometheus 指标，M7 可选新增 |

## 8. 权限码规划

| 模块 | 权限码 |
|------|--------|
| 已实现认证/权限 | `auth.user.read`、`auth.user.write`、`auth.role.read`、`auth.role.write` |
| 已实现 CMDB | `cmdb.model.read`、`cmdb.model.write`、`cmdb.asset.read`、`cmdb.asset.write`、`cmdb.apikey.read`、`cmdb.apikey.write` |
| 已实现公共模块 | `audit.log.read`、`system.config.read`、`system.config.write` |
| 监控 | `monitoring.target.read`、`monitoring.target.write`、`monitoring.metric.read`、`monitoring.dashboard.read`、`monitoring.dashboard.write` |
| 告警 | `alerting.rule.read`、`alerting.rule.write`、`alerting.event.read`、`alerting.event.write`、`alerting.silence.read`、`alerting.silence.write` |
| 通知 | `notification.contact.read`、`notification.contact.write`、`notification.log.read` |
| 工单 | `workorder.read`、`workorder.write` |
