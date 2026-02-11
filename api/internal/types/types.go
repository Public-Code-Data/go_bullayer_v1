package types

// HealthRequest 健康检查请求
type HealthRequest struct {
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status  string `json:"status"`  // 状态：ok/error
	Message string `json:"message"` // 消息
	Time    string `json:"time"`    // 当前时间
}

// UserRequest 用户请求示例
type UserRequest struct {
	ID int64 `json:"id" path:"id"` // 用户ID
}

// UserResponse 用户响应示例
type UserResponse struct {
	ID       int64  `json:"id"`       // 用户ID
	Username string `json:"username"` // 用户名
	Email    string `json:"email"`    // 邮箱
}
