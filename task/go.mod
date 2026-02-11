module go_bullayer_v1/task

go 1.25.7

require (
	go_bullayer_v1/base v0.0.0
	github.com/zeromicro/go-zero v1.6.0
)

// 使用 replace 指令引用本地 base 模块
replace go_bullayer_v1/base => ../base
