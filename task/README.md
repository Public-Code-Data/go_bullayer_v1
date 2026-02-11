# Task 任务模块

提供后台统计任务和定时任务的模块，依赖 base 基础模块。

## 模块结构

```
task/
├── cmd/              # 程序入口
│   └── main.go       # 主函数
├── internal/         # 内部代码
│   ├── config/       # 配置定义
│   ├── service/      # 任务服务
│   └── task/         # 任务实现
│       ├── task.go   # 任务接口
│       └── statstask.go # 统计任务实现
├── etc/              # 配置文件
│   └── task.yaml     # 任务服务配置
└── go.mod            # 模块定义文件
```

## 功能说明

### 1. cmd/main.go
- 程序入口
- 初始化配置和日志
- 启动任务服务
- 优雅关闭处理

### 2. internal/service
- 任务服务管理
- 任务注册和调度
- 任务生命周期管理

### 3. internal/task
- 任务接口定义
- 具体任务实现
- 当前包含统计任务示例

### 4. internal/config
- 任务服务配置
- 任务执行参数配置

## 任务类型

### 统计任务 (StatsTask)
- **功能**: 执行数据统计和报表生成
- **执行时间**: 可配置（默认每小时执行）
- **用途**: 
  - 用户数据统计
  - 订单数据统计
  - 业务指标统计
  - 报表生成

## 运行方式

### 开发环境
```bash
cd task/cmd
go run main.go -f ../etc/task.yaml
```

### 生产环境
```bash
# 编译
go build -o task-server ./cmd

# 运行
./task-server -f etc/task.yaml
```

## 配置说明

配置文件位于 `etc/task.yaml`，包含以下配置项：

- `Name`: 服务名称
- `TaskEnabled`: 是否启用任务服务
- `Interval`: 任务执行间隔（秒）
- `StatsTask`: 统计任务配置
  - `Enabled`: 是否启用统计任务
  - `Hour`: 执行时间（小时）
  - `Minute`: 执行时间（分钟）
- `Database`: 数据库配置（可选）

## 添加新任务

### 1. 实现 Task 接口

```go
package task

import "context"

type MyTask struct {
    // 任务字段
}

func (t *MyTask) Name() string {
    return "我的任务"
}

func (t *MyTask) Execute(ctx context.Context) error {
    // 实现任务逻辑
    return nil
}
```

### 2. 注册任务

在 `service/taskservice.go` 的 `registerTasks` 方法中注册：

```go
func (s *TaskService) registerTasks() {
    // 注册新任务
    myTask := task.NewMyTask(s.config)
    s.tasks = append(s.tasks, myTask)
}
```

## 依赖关系

- **依赖**: `go_bullayer_v1/base`
- **使用**: 
  - base/pkg/logger: 日志管理
  - base/pkg/config: 配置加载
  - base/pkg/db: 数据库连接（如需要）
  - base/pkg/utils: 工具函数

## 注意事项

1. **优雅关闭**: 服务支持优雅关闭，收到 SIGINT 或 SIGTERM 信号时会等待任务完成
2. **任务隔离**: 每个任务在独立的 goroutine 中运行，互不干扰
3. **错误处理**: 任务执行失败会记录日志，但不会影响其他任务
4. **资源管理**: 注意数据库连接等资源的正确释放

## 扩展建议

可以添加以下类型的任务：

1. **数据清理任务**: 定期清理过期数据
2. **报表生成任务**: 生成各种业务报表
3. **数据同步任务**: 同步数据到其他系统
4. **通知任务**: 发送通知消息
5. **缓存刷新任务**: 定期刷新缓存数据
