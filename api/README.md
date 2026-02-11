# API 模块

提供 RESTful API 服务的模块，依赖 base 基础模块。

## 模块结构

```
api/
├── cmd/              # 程序入口
│   └── main.go       # 主函数
├── internal/         # 内部代码
│   ├── config/       # 配置定义
│   ├── handler/      # HTTP处理器
│   ├── logic/        # 业务逻辑
│   ├── svc/          # 服务上下文
│   └── types/        # 类型定义
├── etc/              # 配置文件
│   └── api.yaml      # API服务配置
└── go.mod            # 模块定义文件
```

## 功能说明

### 1. cmd/main.go
- 程序入口
- 初始化配置和日志
- 启动 HTTP 服务器

### 2. internal/config
- 服务配置结构定义
- 包含数据库等业务配置

### 3. internal/handler
- HTTP 请求处理器
- 负责请求解析和响应返回
- 调用 logic 层处理业务逻辑

### 4. internal/logic
- 业务逻辑处理
- 数据验证和处理
- 调用数据库或其他服务

### 5. internal/svc
- 服务上下文
- 管理服务依赖（数据库、缓存等）

### 6. internal/types
- 请求和响应类型定义
- API 数据结构

## API 接口

### 健康检查
- **路径**: `GET /api/health`
- **说明**: 检查服务运行状态
- **响应**: 
```json
{
  "status": "ok",
  "message": "API服务运行正常",
  "time": "2026-02-09 10:00:00"
}
```

### 获取用户信息
- **路径**: `GET /api/user/:id`
- **说明**: 根据ID获取用户信息
- **参数**: 
  - `id` (路径参数): 用户ID
- **响应**:
```json
{
  "id": 1,
  "username": "test_user",
  "email": "test@example.com"
}
```

## 运行方式

### 开发环境
```bash
cd api/cmd
go run main.go -f ../etc/api.yaml
```

### 生产环境
```bash
# 编译
go build -o api-server ./cmd

# 运行
./api-server -f etc/api.yaml
```

## 配置说明

配置文件位于 `etc/api.yaml`，包含以下配置项：

- `Name`: 服务名称
- `Host`: 监听地址
- `Port`: 监听端口
- `Mode`: 运行模式（dev/test/prod）
- `Database`: 数据库配置（可选）

## 依赖关系

- **依赖**: `go_bullayer_v1/base`
- **使用**: 
  - base/pkg/common: 错误处理和响应结构
  - base/pkg/logger: 日志管理
  - base/pkg/config: 配置加载
  - base/pkg/db: 数据库连接
  - base/pkg/utils: 工具函数

## 开发指南

1. **添加新接口**:
   - 在 `types` 中定义请求和响应类型
   - 在 `handler` 中添加处理器
   - 在 `logic` 中实现业务逻辑

2. **错误处理**:
   - 使用 `base/pkg/common` 中的错误类型
   - 返回统一的错误响应

3. **日志记录**:
   - 使用 `base/pkg/logger` 记录日志
   - 关键操作都要记录日志

4. **配置管理**:
   - 配置项定义在 `config.go` 中
   - 配置文件使用 YAML 格式
