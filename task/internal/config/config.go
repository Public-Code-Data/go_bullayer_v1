package config

// Config 任务服务配置结构
type Config struct {
	Name        string `json:"name"`         // 服务名称
	TaskEnabled bool   `json:"task_enabled"` // 是否启用任务
	Interval    int    `json:"interval"`     // 任务执行间隔（秒）

	// 统计任务配置
	StatsTask struct {
		Enabled bool `json:"enabled"` // 是否启用统计任务
		Hour    int  `json:"hour"`     // 执行时间（小时，0-23）
		Minute  int  `json:"minute"`   // 执行时间（分钟，0-59）
	} `json:"stats_task"`

	// 数据库配置（可选）
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"database"`
}
