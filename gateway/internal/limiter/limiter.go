package limiter

import (
	"sync"
	"time"
)

// Limiter QPS限制器
// 使用令牌桶算法实现限流功能
type Limiter struct {
	limit    int                      // 每秒允许的请求数
	interval time.Duration            // 时间间隔
	tokens   map[string]*tokenBucket  // 每个key的令牌桶
	mu       sync.RWMutex             // 读写锁，保护 tokens map
}

// tokenBucket 令牌桶
// 为每个客户端IP维护一个独立的令牌桶
type tokenBucket struct {
	tokens     int       // 当前令牌数
	lastUpdate time.Time // 上次更新时间
}

// NewLimiter 创建限流器
// qps: 每秒允许的请求数
// 返回限流器实例
func NewLimiter(qps int) *Limiter {
	return &Limiter{
		limit:    qps,
		interval: time.Second,
		tokens:   make(map[string]*tokenBucket),
	}
}

// Allow 检查是否允许请求
// key: 限流的key（通常是客户端IP）
// 返回 true 表示允许请求，false 表示超过限制
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	bucket, exists := l.tokens[key]
	now := time.Now()

	// 如果不存在，创建新的令牌桶
	if !exists {
		l.tokens[key] = &tokenBucket{
			tokens:     l.limit - 1, // 初始令牌数为 limit-1（因为这次请求消耗一个）
			lastUpdate: now,
		}
		return true
	}

	// 计算应该补充的token数量
	elapsed := now.Sub(bucket.lastUpdate)
	tokensToAdd := int(elapsed / l.interval * time.Duration(l.limit))

	// 补充令牌
	if tokensToAdd > 0 {
		bucket.tokens = min(bucket.tokens+tokensToAdd, l.limit)
		bucket.lastUpdate = now
	}

	// 检查是否有可用令牌
	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	// 没有可用令牌，拒绝请求
	return false
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
