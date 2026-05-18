文档版本：V2.0
创建日期：2026-05-17
文档状态：终稿
所属项目：超频大师（OCMaster）

一、技术选型与MVP策略

| 模块 | 技术 |
|------|------|
| C端框架 | Electron 28+ / Vue 3 / Vite / Pinia |
| C端硬件采集 | C++ Node Native Addon (node-addon-api) 分平台 |
| S端后端 | Golang / chi / GORM v2 / zerolog（手动DI，待Wire） |
| S端前端 | Vue 3 / Element Plus / TypeScript / Axios |
| 本地数据库 | SQLite (gorm.io/driver/sqlite)，零 Docker 启动 |
| 生产数据库 | MySQL 8.0 (Docker Compose) |
| 测试 | testify + Vitest + Vue Test Utils |

拓扑：C端(Electron+Vue3) + S端(Vue3 SPA) → HTTPS → Nginx → Backend:8080 → DB。
C端 IPC：Main(Native Addon) ←→ contextBridge ←→ Renderer(Vue 3)。

MVP 分两层：C端先行（可独立使用），S端后继。`DB_TYPE=sqlite` 本地零依赖启动。

二、目录结构与双数据库

S端后端 DDD 四层：
```
s-end/backend/
├── cmd/server/          # main.go (手动 DI，等效 wire_gen.go)
├── domain/              # hardware/ sharecode/ merchant/ suggestion/ reference/
├── application/         # Service 层
├── infrastructure/      # persistence/mysql/ (GORM models + 5 repos)
│                        # auth/jwt.go  pdf/chromedp_gen.go
├── interfaces/http/     # router + middleware + handler
├── pkg/config/ + apperrors/
└── Dockerfile           # 仅生产用
```

C端：
```
c-end/
├── electron/
│   ├── main.ts / preload.ts
│   └── native/
│       ├── binding.gyp            # node-gyp 编译配置
│       ├── build/Release/          # 编译产物 .node (不在 git)
│       └── src/                    # C++: hw_win/linux/mac.cpp
├── src/                 # Vue: components/ stores/ config/
└── electron-builder.yml # extraResources 打包 .node 到 resources/native/
```

C端加固：vite-plugin-javascript-obfuscator + DevTools禁用 + F12拦截
Native Addon：node-gyp 编译 → .node 打包到 extraResources，非 mock

S端前端路由：/login /register /dashboard /hardware/:code /suggestion/... /history /profile /settings

双数据库 gorm 切换：
```go
switch cfg.DBType {
case "mysql": gorm.Open(mysql.Open(cfg.DBDSN), ...)
default:      gorm.Open(sqlite.Open(cfg.DBPath), ...)
}
```

| 环境 | DB_TYPE | 路径 | 启动命令 |
|------|---------|------|---------|
| 本地开发 | sqlite | ./data/ocmaster.db | go run ./cmd/server |
| 集成测试 | sqlite | :memory: | go test ./... |
| 生产 | mysql | MySQL 容器 | docker compose up |

三、数据库与API

五张表：merchants, hardware_uploads, suggestions, query_history, reference_data。
GORM AutoMigrate 自动建表。SQLite 不支持 ENUM → VARCHAR + 代码校验。
定时清理：每小时删 expires_at 过期硬件 + 30 天前查询历史。

```
POST   /api/v1/hardware/upload        GET    /api/v1/hardware/:code
DELETE /api/v1/hardware/:code
POST   /api/v1/merchant/register      POST   /api/v1/merchant/login
GET    /api/v1/merchant/profile       PUT    /api/v1/merchant/profile
PUT    /api/v1/merchant/password
POST   /api/v1/suggestions            GET    /api/v1/suggestions/:id
GET    /api/v1/suggestions/history    GET    /api/v1/suggestions/:id/pdf
GET    /api/v1/reference/cpu|ram|cooler
```

四、测试策略

| 模块 | 单元测试 | 集成测试 | 工具 |
|------|---------|---------|------|
| S端 Domain | 实体+值对象+Service(mock) | — | testify |
| S端 Infra | — | SQLite :memory: | testify + gorm |
| S端 HTTP | Handler(mock) | 端到端 + SQLite | httptest |
| S端前端 | Stores+Utils | 组件测试 | Vitest+Vue Test Utils |
| C端 | Stores+Composables | 组件+E2E | Vitest+electron-mocha |

原则：测试用代码默认配置，SQLite :memory: 免外部依赖。生产测试用 Docker MySQL。

五、部署与配置

生产部署 (Docker Compose):
```yaml
services:
  mysql:   { image: mysql:8.0, environment: { MYSQL_ROOT_PASSWORD, MYSQL_DATABASE: ocmaster } }
  backend: { build: ./s-end/backend, environment: { DB_TYPE: mysql, DB_DSN, JWT_SECRET } }
  nginx:   { image: nginx:alpine, ports: [80,443], volumes: [nginx.conf, frontend/dist] }
```

配置：
- S端后端：viper SetDefault → config.yaml(可选) → env。默认 sqlite 零配置启动
- C端：electron-store → ~/.ocmaster/config.json
- S端前端：Vite env 配置 API 地址。模板存后端

C端压缩策略（135MB portable / 129MB 7z ultra）：
- Electron 28 运行时 ~110MB 已预压缩，是体积瓶颈
- asar 打包 + locale 裁剪（仅 zh-CN/en-US）
- vite 分包：vue/pinia 独立 vendor chunk
- 如需 <80MB：改用系统 WebView2 (Win10+) 或 nsis-web（按需下载 Electron）
