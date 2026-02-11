package config

import "github.com/zeromicro/go-zero/rest"

// Config API服务配置结构
// 包含 RESTful 服务器配置和业务配置
type Config struct {
	rest.RestConf // go-zero RESTful 服务器配置

	// 业务配置可以在这里添加
	// 例如：数据库配置、Redis配置、第三方服务配置等
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"database"`
}
