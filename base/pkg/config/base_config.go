package config

import "github.com/zeromicro/go-zero/core/conf"

// BaseConfig 基础配置结构
// 所有服务的通用配置项
type BaseConfig struct {
	Name    string `json:"name"`    // 服务名称
	Version string `json:"version"` // 服务版本
	Env     string `json:"env"`     // 环境：dev, test, prod
}

// LoadConfig 加载配置文件
// configFile: 配置文件路径
// v: 配置结构体指针
// 返回错误信息
func LoadConfig(configFile string, v interface{}) error {
	return conf.Load(configFile, v)
}

// MustLoadConfig 必须加载配置，失败则panic
// configFile: 配置文件路径
// v: 配置结构体指针
func MustLoadConfig(configFile string, v interface{}) {
	conf.MustLoad(configFile, v)
}
