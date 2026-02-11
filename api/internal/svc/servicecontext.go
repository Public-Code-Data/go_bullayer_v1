package svc

import (
	"database/sql"

	"go_bullayer_v1/api/internal/config"
	"go_bullayer_v1/base/pkg/db"
)

// ServiceContext 服务上下文
// 包含服务运行所需的所有依赖和配置
type ServiceContext struct {
	Config config.Config // 服务配置
	DB     *sql.DB       // 数据库连接（示例，根据实际需要添加）
	// 可以在这里添加其他依赖，如：
	// Redis客户端、消息队列客户端、第三方服务客户端等
}

// NewServiceContext 创建服务上下文
// c: 服务配置
// 返回服务上下文实例
func NewServiceContext(c config.Config) *ServiceContext {
	ctx := &ServiceContext{
		Config: c,
	}

	// 初始化数据库连接（示例）
	// 如果配置了数据库，则创建连接
	if c.Database.Host != "" {
		dbConfig := db.DBConfig{
			Host:     c.Database.Host,
			Port:     c.Database.Port,
			User:     c.Database.User,
			Password: c.Database.Password,
			Database: c.Database.Database,
		}
		// 创建数据库连接
		database, err := db.NewDB(dbConfig)
		if err == nil {
			ctx.DB = database
		}
		// 如果连接失败，记录日志但不影响服务启动
	}

	return ctx
}
