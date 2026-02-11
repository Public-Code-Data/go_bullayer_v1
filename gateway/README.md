# Gateway 网关模块

提供服务路由、接口降级和用户 QPS 限制的网关模块，依赖 base 基础模块。

## 模块结构

```
gateway/
├── cmd/              # 程序入口
│   └── main.go       # 主函数
├── internal/         # 内部代码
│   ├── config/       # 配置定义
│   ├── handler/      # HTTP处理器
│   ├── middleware/   # 中间件
│   │   ├── ratelimit.go      # QPS限制中间件
│   │   └── circuitbreaker.go # 熔断降级中间件
│   ├── limiter/      # 限流器实现
│   └── router/       # 路由管理器
├── etc/              # 配置文件
│   └── gateway.yaml   # Gateway服务配置
└── go.mod            # 模块定义文件
```

## 功能说明

### 1. 服务路由
- 根据配置将请求转发到后端服务
- 支持路径匹配和通配符
- 支持多种HTTP方法

### 2. QPS限制
- 基于令牌桶算法的限流实现
- 按客户端IP进行限流
- 可配置每秒请求数限制

### 3. 接口降级
- 基于熔断器模式的降级机制
- 后端服务异常时自动降级
- 返回预设的降级响应

### 4. 中间件系统
- 可插拔的中间件架构
- 支持中间件链式调用
- 中间件执行顺序可配置

## 核心组件

### RateLimitMiddleware - QPS限制中间件
- **功能**: 限制每个客户端的请求频率
- **算法**: 令牌桶算法
- **配置**: 通过配置文件设置QPS限制

### CircuitBreakerMiddleware - 熔断降级中间件
- **功能**: 后端服务异常时自动熔断
- **机制**: 基于错误率和超时时间
- **降级**: 返回预设的降级响应

### Router - 路由管理器
- **功能**: 管理路由规则和请求转发
- **匹配**: 支持路径匹配和通配符
- **转发**: HTTP请求转发到后端服务

## 运行方式

### 开发环境
```bash
cd gateway/cmd
go run main.go -f ../etc/gateway.yaml
```

### 生产环境
```bash
# 编译
go build -o gateway-server ./cmd

# 运行
./gateway-server -f etc/gateway.yaml
```

## 配置说明

配置文件位于 `etc/gateway.yaml`，包含以下配置项：

### 基础配置
- `Name`: 服务名称
- `Host`: 监听地址
- `Port`: 监听端口
- `Mode`: 运行模式（dev/test/prod）

### QPS限制配置
- `RateLimit.Enabled`: 是否启用QPS限制
- `RateLimit.QPS`: 每秒允许的请求数

### 熔断降级配置
- `CircuitBreaker.Enabled`: 是否启用熔断器
- `CircuitBreaker.Timeout`: 超时时间（秒）

### 路由配置
- `Routes`: 路由规则列表
  - `Path`: 请求路径（支持通配符）
  - `Backend`: 后端服务地址
  - `Method`: HTTP方法
  - `Fallback`: 降级响应内容
  - `Timeout`: 超时时间（秒）

## API 接口

### 健康检查
- **路径**: `GET /gateway/health`
- **说明**: 检查Gateway服务运行状态
- **响应**: 
```json
{
  "status": "ok",
  "message": "Gateway服务运行正常",
  "routes": 2
}
```

## 使用示例

### 1. 配置路由
```yaml
Routes:
  - Path: /api/*
    Backend: http://localhost:8888
    Method: GET
    Fallback: "API服务暂时不可用"
    Timeout: 5
```

### 2. 启用QPS限制
```yaml
RateLimit:
  Enabled: true
  QPS: 100  # 每秒最多100个请求
```

### 3. 启用熔断降级
```yaml
CircuitBreaker:
  Enabled: true
  Timeout: 5  # 5秒超时
```

## 依赖关系

- **依赖**: `go_bullayer_v1/base`
- **使用**: 
  - base/pkg/logger: 日志管理
  - base/pkg/config: 配置加载
  - base/pkg/common: 错误处理

## 架构说明

### 请求流程
1. 客户端请求到达Gateway
2. QPS限制中间件检查请求频率
3. 熔断降级中间件检查服务状态
4. 路由管理器匹配路由规则
5. 转发请求到后端服务
6. 返回响应给客户端

### 降级流程
1. 后端服务异常或超时
2. 熔断器触发
3. 返回预设的降级响应
4. 记录错误日志

## 扩展建议

1. **负载均衡**: 添加负载均衡功能，支持多个后端服务
2. **服务发现**: 集成服务发现机制（如Consul、Etcd）
3. **认证授权**: 添加API密钥、JWT等认证机制
4. **监控告警**: 集成Prometheus等监控系统
5. **日志追踪**: 添加分布式追踪功能（如Jaeger）
6. **缓存**: 添加响应缓存功能
