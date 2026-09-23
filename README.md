# 家居装修全流程管理平台

面向业主、设计师和施工方的全栈 Web 应用，覆盖装修设计阶段、施工进度、材料清单、预算管控与 3D 户型预览，用于验证模型跨文件协同改动能力。

## Docker 一键启动（推荐）

```bash
cp .env.example .env
docker compose --env-file .env up -d --build
```

访问地址：

- 前端：http://127.0.0.1:18805
- 后端 API：http://127.0.0.1:19305
- 健康检查：http://127.0.0.1:19305/healthz

停止并清理：

```bash
docker compose --env-file .env down -v --remove-orphans
```

种子账号（首次启动自动创建）：

| 账号 | 密码 | 角色 |
| --- | --- | --- |
| admin | Admin123456 | Admin |
| designer | Designer123 | Designer |
| contractor | Contractor123 | Contractor |
| owner | Owner123456 | Owner |
| pm | Manager123456 | ProjectManager |

## 技术栈

| 层级 | 技术 |
| --- | --- |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | MySQL 8.0 |
| 认证 | JWT（golang-jwt/jwt/v5） |
| 前端 | React 18 + TypeScript + Vite |
| UI | Ant Design 5 |
| 图表 | ECharts |
| 状态管理 | Zustand |
| 部署 | Docker Compose |

## 目录结构

```text
.
├── backend/          # Go 后端（cmd/internal 分层）
│   ├── cmd/server/   # 启动入口
│   ├── internal/     # config/model/repository/service/handler/router/middleware/dto/constants/errors/utils/logger
│   ├── api/          # OpenAPI 文档
│   ├── deploy/       # 部署说明
│   └── database/     # migrations/seeds
├── frontend/         # React 前端
│   └── src/          # api/stores/types/components/hooks/pages/router/utils/constants
├── database/init.sql # MySQL 初始化脚本
├── docker-compose.yml
└── .env.example
```

## 枚举位置

- 后端：`backend/internal/constants/enums.go`
- 前端：`frontend/src/types/enums.ts`

## API 清单

统一前缀 `/api/v1`，响应结构 `{ "code": 0, "message": "ok", "data": ... }`。

| 方法 | 路径 | 说明 | 角色 |
| --- | --- | --- | --- |
| POST | /api/v1/auth/login | 登录 | 公开 |
| GET | /api/v1/auth/me | 当前用户 | 登录用户 |
| GET/POST | /api/v1/projects | 项目列表/创建 | 列表所有角色，创建 Admin/PM |
| GET/PUT/DELETE | /api/v1/projects/:id | 项目详情/更新/删除 | 更新 Admin/PM，删除 Admin |
| PUT | /api/v1/projects/:id/status | 项目状态流转 | Admin/PM |
| GET/POST | /api/v1/designs | 设计列表/创建 | 列表所有角色，创建 Admin/Designer/PM |
| PUT | /api/v1/designs/:id | 更新设计 | Admin/Designer/PM |
| PUT | /api/v1/designs/:id/submit | 提交审核 | Admin/Designer/PM |
| PUT | /api/v1/designs/:id/review | 审核（通过/驳回） | Admin/Owner |
| GET/POST | /api/v1/materials | 材料列表/创建 | 列表所有角色，创建 Admin/Designer/PM |
| PUT | /api/v1/materials/:id/status | 采购状态流转 | Admin/Designer/Contractor/PM |
| GET/POST | /api/v1/budgets | 预算列表/创建 | 列表所有角色，创建 Admin/PM |
| GET/POST | /api/v1/constructions | 施工列表/创建 | 列表所有角色，创建 Admin/PM |
| PUT | /api/v1/constructions/:id/status | 施工状态流转 | Admin/Contractor/PM |
| PUT | /api/v1/constructions/:id/accept | 施工验收 | Admin/Contractor/PM |
| POST | /api/v1/upload | 文件上传 | 登录用户 |
| GET | /api/v1/audit-logs | 操作日志 | Admin |

## 本地开发

后端：

```bash
cd backend
go mod tidy
go run ./cmd/server
```

前端：

```bash
cd frontend
npm install
npm run dev
```

## License

MIT
