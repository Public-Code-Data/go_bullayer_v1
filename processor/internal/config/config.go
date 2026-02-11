package config

// Config 数据处理服务配置结构
type Config struct {
	Name             string `json:"Name"`
	ProcessorEnabled bool   `json:"ProcessorEnabled"`
	Interval         int    `json:"Interval"`

	// 链配置
	Chain struct {
		RPCURL            string `json:"RPCURL"`
		ChainID           int64  `json:"ChainID"`
		StartHeight       int64  `json:"StartHeight"`
		Confirmations     int64  `json:"Confirmations"`
		MaxBlocksPerRound int64  `json:"MaxBlocksPerRound"`
	} `json:"Chain"`

	// 区块解析任务配置
	BlockProcessor struct {
		Enabled    bool  `json:"Enabled"`
		BatchSize  int64 `json:"BatchSize"`
		ParseTx    bool  `json:"ParseTx"`
		ParseEvent bool  `json:"ParseEvent"`
	} `json:"BlockProcessor"`

	// 数据库配置（可选）
	Database struct {
		Host     string `json:"Host"`
		Port     int    `json:"Port"`
		User     string `json:"User"`
		Password string `json:"Password"`
		Database string `json:"Database"`
	} `json:"Database"`
}
