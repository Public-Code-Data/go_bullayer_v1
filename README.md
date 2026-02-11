# Go Bullayer V1 - 模块化微服务项目

这是一个采用模块化架构的 Go 微服务项目，类似于 Java 的 Maven 多模块项目。项目使用 Go Modules 和 Go Workspace 实现模块化管理。

## 📁 项目结构

```
go_bullayer_v1/
├── base/          # 基础模块 - 提供通用工具和基础服务
│   ├── pkg/        # 公共包
│   │   ├── common/     # 通用功能：错误处理、响应结构
│   │   ├── logger/     # 日志管理
│   │   ├── config/     # 配置管理
│   │   ├── db/         # 数据库连接管理
│   │   └── utils/      # 工具函数：字符串、时间等
│   └── go.mod
│
├── api/           # API模块 - 提供RESTful API服务
│   ├── cmd/            # 程序入口
│   ├── internal/       # 内部代码
│   │   ├── handler/    # HTTP处理器
│   │   ├── logic/      # 业务逻辑
│   │   ├── svc/        # 服务上下文
│   │   ├── types/      # 类型定义
│   │   └── config/     # 配置定义
│   ├── etc/            # 配置文件
│   └── go.mod          # 依赖 base 模块
│
├── task/          # Task模块 - 提供后台统计任务
│   ├── cmd/            # 程序入口
│   ├── internal/       # 内部代码
│   │   ├── service/    # 任务服务
│   │   ├── task/       # 任务实现
│   │   └── config/     # 配置定义
│   ├── etc/            # 配置文件
│   └── go.mod          # 依赖 base 模块
│
├── gateway/       # Gateway模块 - 服务路由、接口降级、QPS限制
│   ├── cmd/            # 程序入口
│   ├── internal/       # 内部代码
│   │   ├── handler/    # HTTP处理器
│   │   ├── middleware/ # 中间件：限流、熔断
│   │   ├── limiter/    # 限流器实现
│   │   ├── router/     # 路由管理器
│   │   └── config/     # 配置定义
│   ├── etc/            # 配置文件
│   └── go.mod          # 依赖 base 模块
│
├── go.work        # Go工作区文件，统一管理所有模块
└── README.md      # 项目说明文档
```

## 🚀 快速开始

### 1. 环境要求

- Go 1.25.7 或更高版本
- 支持 Go Workspace（Go 1.18+）

### 2. 初始化工作区

项目已经配置了 `go.work` 文件，可以直接使用：

```bash
# 同步工作区依赖
go work sync
```

### 3. 运行各个服务

#### 运行 API 服务
```bash
cd api/cmd
go run main.go -f ../etc/api.yaml
```

#### 运行 Task 服务
```bash
cd task/cmd
go run main.go -f ../etc/task.yaml
```

#### 运行 Gateway 服务
```bash
cd gateway/cmd
go run main.go -f ../etc/gateway.yaml
```

## 📦 模块说明

### Base 基础模块

**位置**: `base/`

**功能**: 
- 通用错误处理和响应结构
- 日志管理
- 配置加载
- 数据库连接管理
- 工具函数（字符串、时间等）

**依赖**: 无（基础模块）

**使用示例**:
```go
import (
    "go_bullayer_v1/base/pkg/common"
    "go_bullayer_v1/base/pkg/logger"
    "go_bullayer_v1/base/pkg/config"
)
```

### API 模块

**位置**: `api/`

**功能**: 
- RESTful API 接口服务
- 业务逻辑处理
- 数据验证和处理

**依赖**: `base` 模块

**API接口**:
- `GET /api/health` - 健康检查
- `GET /api/user/:id` - 获取用户信息

### Task 模块

**位置**: `task/`

**功能**: 
- 后台统计任务
- 定时任务执行
- 数据报表生成

**依赖**: `base` 模块

**任务类型**:
- 统计任务：数据统计和报表生成

### Gateway 模块

**位置**: `gateway/`

**功能**: 
- 服务路由和转发
- QPS限制（限流）
- 接口降级（熔断器）
- 请求代理

**依赖**: `base` 模块

**特性**:
- 基于令牌桶算法的限流
- 基于熔断器模式的降级
- 可配置的路由规则

## 🔧 开发指南

### 模块间依赖管理

每个模块的 `go.mod` 文件中使用 `replace` 指令引用本地 base 模块：

```go
require (
    go_bullayer_v1/base v0.0.0
)

replace go_bullayer_v1/base => ../base
```

### 添加新模块

1. 创建新模块目录和 `go.mod` 文件
2. 在 `go.work` 中添加新模块路径
3. 使用 `replace` 引用 base 模块

### 添加新功能

1. **Base 模块**: 添加公共功能
2. **API 模块**: 添加新的 API 接口
3. **Task 模块**: 添加新的后台任务
4. **Gateway 模块**: 添加新的路由规则

## 📝 代码规范

### 注释规范

- 所有导出的函数、类型、变量都需要注释
- 注释使用中文，清晰说明功能和参数
- 复杂逻辑需要添加行内注释

### 命名规范

- 包名：小写，简短
- 函数名：驼峰命名，首字母大写表示导出
- 变量名：驼峰命名，见名知意

### 错误处理

- 使用 `base/pkg/common` 中的错误类型
- 统一错误码定义
- 错误信息清晰明确

## 🏗️ 构建和部署

### 构建各个模块

```bash
# 构建 base 模块
cd base
go build ./...

# 构建 api 模块
cd api
go build ./cmd

# 构建 task 模块
cd task
go build ./cmd

# 构建 gateway 模块
cd gateway
go build ./cmd
```

### 生产环境部署

每个模块可以独立编译和部署：

```bash
# 编译为二进制文件
go build -o api-server ./api/cmd
go build -o task-server ./task/cmd
go build -o gateway-server ./gateway/cmd

# 运行
./api-server -f api/etc/api.yaml
./task-server -f task/etc/task.yaml
./gateway-server -f gateway/etc/gateway.yaml
```

## 🔍 与 Java Maven 对比

| 特性 | Java Maven | Go Modules |
|------|-----------|------------|
| 模块定义 | `pom.xml` | `go.mod` |
| 工作区管理 | 父 pom.xml | `go.work` |
| 本地依赖 | `relativePath` | `replace` 指令 |
| 模块路径 | `groupId:artifactId` | URL 风格路径 |
| 版本管理 | 版本号 | 版本标签或 commit |

## ✨ 项目优势

1. **模块化**: 清晰的模块划分，职责明确
2. **依赖管理**: 通过 base 模块统一管理公共依赖
3. **独立部署**: 每个模块可以独立编译和部署
4. **代码复用**: base 模块提供通用功能，避免重复代码
5. **易于维护**: 模块化结构便于团队协作和维护

## 📚 相关文档

- [Base 模块文档](./base/README.md)
- [API 模块文档](./api/README.md)
- [Task 模块文档](./task/README.md)
- [Gateway 模块文档](./gateway/README.md)

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证。

## 👥 作者

- 项目维护者

## 🙏 致谢

感谢所有为本项目做出贡献的开发者！

---

**注意**: 这是一个示例项目结构，实际使用时请根据业务需求进行调整和完善。
