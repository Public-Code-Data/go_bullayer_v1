# Processor 数据处理模块

负责追踪链上区块高度，并解析区块交易和事件数据。

## 目录结构

```text
processor/
├── cmd/                 # 程序入口
├── internal/
│   ├── core/            # 核心解析能力（回执拉取、并发交易解析、过滤）
│   ├── config/          # 配置定义
│   ├── service/         # 处理服务
│   └── processor/       # 处理任务实现
├── etc/                 # 配置文件
└── go.mod               # 模块依赖
```

## 功能说明

- 周期拉取链上最新区块高度
- 根据确认数计算可安全处理高度
- 按批次推进处理高度
- 解析区块交易与事件（当前为骨架，预留真实链 RPC 和落库逻辑）

## 运行方式

```bash
cd processor/cmd
go run main.go -f ../etc/processor.yaml
```

## 配置项

- `ProcessorEnabled`: 是否启用处理服务
- `Interval`: 轮询间隔（秒）
- `Chain`: 链相关配置（RPC、起始高度、确认数、单轮处理上限）
- `BlockProcessor`: 区块解析配置（是否启用、解析项开关、并发数、目标地址、代币白名单）
- `Database`: 数据库配置（可选，用于解析结果落库）
