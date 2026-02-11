# Base 基础模块

这是项目的核心基础模块，提供所有其他模块依赖的通用功能。

## 模块结构

```
base/
├── pkg/              # 公共包
│   ├── common/       # 通用功能：错误处理、响应结构
│   ├── logger/       # 日志管理
│   ├── config/       # 配置管理
│   ├── db/           # 数据库连接管理
│   └── utils/        # 工具函数：字符串、时间等
├── internal/         # 内部代码（可选）
│   ├── model/        # 数据模型
│   └── service/      # 基础服务
└── go.mod            # 模块定义文件
```

## 功能说明

### 1. common - 通用功能
- **错误处理**: `BaseError` 统一错误类型
- **响应结构**: `Response` 统一API响应格式
- **错误码**: 预定义的通用错误码

### 2. logger - 日志管理
- 基于 go-zero 的日志系统
- 支持文件日志和标准输出
- 提供 Info、Error、Debug 等便捷方法

### 3. config - 配置管理
- 统一的配置加载接口
- 支持 YAML 配置文件
- 基础配置结构定义

### 4. db - 数据库管理
- MySQL 数据库连接
- 连接池管理
- 连接参数配置

### 5. utils - 工具函数
- 字符串工具函数
- 时间工具函数
- 其他通用工具

## 使用示例

### 在其他模块中引用

```go
import (
    "go_bullayer_v1/base/pkg/common"
    "go_bullayer_v1/base/pkg/logger"
    "go_bullayer_v1/base/pkg/config"
    "go_bullayer_v1/base/pkg/db"
    "go_bullayer_v1/base/pkg/utils"
)

// 使用错误处理
err := common.NewError(common.ErrCodeInvalidParam, "参数错误")

// 使用日志
logger.Info("服务启动成功")

// 使用配置
config.MustLoadConfig("config.yaml", &cfg)

// 使用数据库
dbConfig := db.DBConfig{
    Host:     "localhost",
    Port:     3306,
    User:     "root",
    Password: "password",
    Database: "test",
}
database, err := db.NewDB(dbConfig)

// 使用工具函数
if utils.IsEmpty(str) {
    // 处理空字符串
}
```

## 注意事项

1. 此模块是基础模块，不应依赖其他业务模块
2. 所有公共功能都应放在此模块中
3. 保持接口简洁，避免过度设计
4. 遵循 Go 语言的最佳实践
