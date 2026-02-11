package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DBConfig 数据库配置
type DBConfig struct {
	Host     string // 数据库主机
	Port     int    // 数据库端口
	User     string // 数据库用户名
	Password string // 数据库密码
	Database string // 数据库名称
	MaxOpen  int    // 最大打开连接数
	MaxIdle  int    // 最大空闲连接数
}

// NewDB 创建数据库连接
// config: 数据库配置
// 返回数据库连接实例和错误信息
func NewDB(config DBConfig) (*sql.DB, error) {
	// 构建 DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.Database)

	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("打开数据库连接失败: %w", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %w", err)
	}

	// 设置连接池参数
	if config.MaxOpen > 0 {
		db.SetMaxOpenConns(config.MaxOpen)
	} else {
		db.SetMaxOpenConns(100) // 默认最大100个连接
	}

	if config.MaxIdle > 0 {
		db.SetMaxIdleConns(config.MaxIdle)
	} else {
		db.SetMaxIdleConns(10) // 默认最大10个空闲连接
	}

	// 设置连接最大生存时间
	db.SetConnMaxLifetime(time.Hour)

	return db, nil
}
