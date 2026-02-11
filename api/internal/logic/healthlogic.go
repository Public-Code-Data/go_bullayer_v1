package logic

import (
	"context"

	"go_bullayer_v1/api/internal/svc"
	"go_bullayer_v1/api/internal/types"
	"go_bullayer_v1/base/pkg/logger"
	"go_bullayer_v1/base/pkg/utils"
)

// HealthLogic 健康检查业务逻辑
type HealthLogic struct {
	ctx    context.Context // 上下文
	svcCtx *svc.ServiceContext // 服务上下文
}

// NewHealthLogic 创建健康检查逻辑处理器
// ctx: 上下文
// svcCtx: 服务上下文
// 返回健康检查逻辑处理器实例
func NewHealthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthLogic {
	return &HealthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Health 执行健康检查
// req: 健康检查请求
// 返回健康检查响应和错误信息
func (l *HealthLogic) Health(req *types.HealthRequest) (resp *types.HealthResponse, err error) {
	// 记录日志
	logger.Info("执行健康检查")

	// 获取当前时间
	now := utils.Now()
	currentTime := utils.FormatDateTime(now)

	// 返回健康检查结果
	return &types.HealthResponse{
		Status:  "ok",
		Message: "API服务运行正常",
		Time:    currentTime,
	}, nil
}
