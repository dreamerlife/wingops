# WingOps

WingOps 是统一运维管理平台 MVP，当前包含 Go 后端、Vue 前端、PostgreSQL 数据库迁移、CMDB、RBAC、审计、系统配置和本地验收数据。

## 技术栈

- 后端：Go + Gin + GORM + Viper + Zap
- 前端：Vue 3 + Vite + Pinia + Element Plus
- 数据库：PostgreSQL
- 缓存和消息：Redis、NATS
- 本地依赖：Docker Compose

## 目录结构及作用

```text
.
├── cmd/server/              # 后端启动入口，main.go 负责加载配置、连接数据库、迁移并启动 HTTP 服务
├── configs/                 # 后端配置样板，app.example.yaml 是本地配置参考
├── docs/                    # 项目设计、接口、数据库和约束规范文档
├── docs/api/                # 细分 API 文档
├── internal/                # 后端内部代码
│   ├── audit/               # 审计日志模块
│   ├── auth/                # 登录、用户、角色、权限和 RBAC 模块
│   ├── cache/               # Redis 连接入口
│   ├── cmdb/                # CMDB 模型、资产、导入、同步和 API Key 模块
│   ├── config/              # 配置读取逻辑
│   ├── database/            # PostgreSQL 连接、迁移和种子数据
│   ├── http/                # Gin Router 组装和 HTTP 中间件
│   ├── messaging/           # NATS 连接入口
│   ├── response/            # 通用响应结构
│   └── system/              # 系统配置模块
├── migrations/              # PostgreSQL 数据库迁移 SQL
├── web/                     # 前端应用
│   └── src/                 # 前端源码，包含 API、路由、状态、组件和页面
├── docker-compose.yml       # 本地依赖服务：PostgreSQL、Redis、NATS、VictoriaMetrics
├── Makefile                 # 常用开发命令
├── go.mod                   # Go 模块和后端依赖
└── README.md                # 项目入口说明
```

## 后端启动入口

后端入口文件：

```text
cmd/server/main.go
```

启动流程：

1. 读取配置：`internal/config/config.go`。
2. 初始化 Zap 日志。
3. 连接 PostgreSQL：`internal/database/postgres.go`。
4. 执行 `migrations/*.sql` 数据库迁移。
5. 写入本地验收种子数据。
6. 初始化各模块 Repository、Service、Handler。
7. 通过 `internal/http/router.go` 组装 Gin 路由。
8. 按 `server.addr` 监听 HTTP 服务。

健康检查接口：

```text
GET /healthz
```

## 配置说明

配置样板：

```text
configs/app.example.yaml
```

本地可复制为：

```bash
cp configs/app.example.yaml configs/app.yaml
```

后端默认读取 `./configs/app.yaml` 或当前目录下的 `app.yaml`。如果配置文件不存在，会使用代码中的本地开发默认值。

当前配置样板：

```yaml
server:
  addr: ":8080"
postgres:
  dsn: "host=localhost user=wingops password=wingops dbname=wingops port=5432 sslmode=disable"
redis:
  addr: "localhost:6379"
nats:
  url: "nats://localhost:4222"
jwt:
  secret: "dev-secret-change-before-production"
  access_token_ttl_minutes: 60
```

运行环境相关配置不要写死在业务代码里，应通过 `configs/app.yaml` 或环境变量覆盖。

## 环境变量配置

后端环境变量前缀为 `WINGOPS`，配置键中的 `.` 会映射为 `_`。

| 配置项 | YAML 键 | 环境变量 | 默认值 |
|--------|---------|----------|--------|
| 后端监听地址 | `server.addr` | `WINGOPS_SERVER_ADDR` | `:8080` |
| PostgreSQL DSN | `postgres.dsn` | `WINGOPS_POSTGRES_DSN` | `host=localhost user=wingops password=wingops dbname=wingops port=5432 sslmode=disable` |
| Redis 地址 | `redis.addr` | `WINGOPS_REDIS_ADDR` | `localhost:6379` |
| NATS 地址 | `nats.url` | `WINGOPS_NATS_URL` | `nats://localhost:4222` |
| JWT 密钥 | `jwt.secret` | `WINGOPS_JWT_SECRET` | `dev-secret-change-before-production` |
| Token 有效期（分钟） | `jwt.access_token_ttl_minutes` | `WINGOPS_JWT_ACCESS_TOKEN_TTL_MINUTES` | `60` |

示例：

```bash
WINGOPS_SERVER_ADDR=":8081" make run
```

前端开发服务默认把 `/api` 代理到后端：

```text
http://127.0.0.1:8080
```

可通过 `VITE_API_PROXY_TARGET` 覆盖：

```bash
cd web
VITE_API_PROXY_TARGET="http://127.0.0.1:8081" npm run dev
```

## 本地启动

默认账号：

- 用户名：`admin`
- 密码：`admin123`

启动依赖服务：

```bash
make infra-up
```

启动后端：

```bash
make run
```

后端启动时会自动执行 `migrations/*.sql`，并写入本地验收种子数据：

- admin 用户、RBAC 角色和权限
- `platform.name = WingOps`
- CMDB 预置模型组与模型：物理服务器、虚拟机、交换机、路由器、Nginx、MySQL、Redis、业务应用
- CMDB 同步 Key：`dev-sync-key` / `dev-sync-secret`

启动前端：

```bash
cd web
npm run dev
```

浏览器访问：`http://localhost:5173`

## 常用命令

```bash
make infra-up      # 启动本地依赖服务
make infra-down    # 停止本地依赖服务
make run           # 启动后端
make test          # 运行后端测试
```

前端测试：

```bash
cd web
npm test
```

## 手工验收路径

1. 登录系统，进入平台首页，确认模型组/API Key 数量来自数据库。
2. 进入“CMDB 模型”，新增模型组，新增模型，进入“编辑字段”维护动态字段。
3. 进入“资产管理”，选择模型，录入资产并编辑生命周期状态。
4. 打开资产详情，确认属性与变更历史已记录。
5. 进入“资产导入”，上传 CSV 预览并导入。CSV 至少包含 `unique_key` 列。
6. 进入“API Key”，新增/吊销同步 Key。
7. 进入“系统配置”，修改配置后刷新页面确认仍存在。
8. 进入“审计日志”，确认受保护 API 操作被记录。

CSV 示例：

```csv
unique_key,name,management_ip,environment,owner
srv-001,web-01,10.0.1.10,prod,ops
```

API 同步签名为请求体的 HMAC-SHA256 十六进制摘要，Secret 使用 `dev-sync-secret`。

## 开发验证

运行后端测试：

```bash
make test
```

启动服务：

```bash
make run
```

健康检查：

```bash
curl http://localhost:8080/healthz
```

## 相关文档

- [项目约束规范](docs/项目约束规范.md)
- [API 接口层设计](docs/API接口层设计.md)
- [数据库表结构设计](docs/数据库表结构设计.md)
- [CMDB Asset Sync API](docs/api/cmdb-sync.md)
