# 交易平台需求分析文档

## 一、用户账户系统

### 1.1 钱包登录与账户创建

**功能描述：**
- 钱包首次登录时，以钱包地址为KEY创建中心化账户
- 通过签名验证后，账户信息落库
- 支持Metamask和OKX钱包（第一版）

**数据库设计：**
```sql
CREATE TABLE accounts (
    account_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    address VARCHAR(42) UNIQUE NOT NULL COMMENT '钱包地址',
    status TINYINT DEFAULT 1 COMMENT '账户状态：1-正常，2-冻结，3-禁用',
    invitation_code VARCHAR(20) UNIQUE COMMENT '邀请码',
    wallet_type VARCHAR(20) COMMENT '钱包类型：metamask, okx',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_address (address),
    INDEX idx_invitation_code (invitation_code)
);
```

**API接口：**
- `POST /api/v1/auth/login` - 钱包登录（签名验证）
- `POST /api/v1/auth/refresh` - Token刷新

**技术要点：**
- 签名验证：使用EIP-191标准签名验证
- Token生成：JWT Token，设置过期时间
- 钱包类型识别：通过签名消息格式或前端标识

### 1.2 认证与鉴权

**功能描述：**
- 每次登录验证签名，获取时效Token
- 后续下单等操作使用Token进行鉴权

**技术实现：**
- JWT Token机制
- Token过期时间：建议24小时
- 刷新Token机制：使用Refresh Token延长会话

**中间件：**
- 认证中间件：验证Token有效性
- 权限中间件：验证用户权限

---

## 二、Bridge/充值系统

### 2.1 资产配置

**功能描述：**
- 支持USDT、ETH、BTC三种资产
- 配置最小充值金额限制
- 在以太坊Sepolia测试网部署

**数据库设计：**
```sql
CREATE TABLE asset_configs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    coin VARCHAR(10) NOT NULL COMMENT '币种：BTC, ETH, USDT',
    contract_address VARCHAR(42) COMMENT '合约地址（ERC20）',
    min_deposit DECIMAL(36,18) NOT NULL COMMENT '最小充值金额',
    status TINYINT DEFAULT 1 COMMENT '状态：1-启用，2-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_coin (coin)
);

-- 初始数据
INSERT INTO asset_configs (coin, min_deposit) VALUES
('BTC', 0.001),
('ETH', 0.01),
('USDT', 10);
```

### 2.2 充值监听

**功能描述：**
- 使用一个EOA地址作为充值地址
- 监听USDT、ETH、BTC三种资产的转账
- 验证最小充值金额
- 充值成功后自动入账到用户账户

**技术实现：**
- 区块链监听服务：
  - 使用Web3客户端连接Sepolia测试网
  - 监听ERC20 Transfer事件（USDT）
  - 监听ETH转账（原生币）
  - 监听BTC转账（需要BTC节点或第三方API）
- 充值处理流程：
  1. 监听链上交易
  2. 验证交易状态（确认数）
  3. 验证最小金额
  4. 查询用户账户（通过地址）
  5. 更新用户资产
  6. 记录充值记录

**数据库设计：**
```sql
CREATE TABLE deposits (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    account_id BIGINT NOT NULL COMMENT '账户ID',
    coin VARCHAR(10) NOT NULL COMMENT '币种',
    amount DECIMAL(36,18) NOT NULL COMMENT '充值金额',
    tx_hash VARCHAR(66) UNIQUE NOT NULL COMMENT '交易哈希',
    from_address VARCHAR(42) COMMENT '发送地址',
    to_address VARCHAR(42) COMMENT '接收地址',
    block_number BIGINT COMMENT '区块号',
    confirmations INT DEFAULT 0 COMMENT '确认数',
    status TINYINT DEFAULT 0 COMMENT '状态：0-待确认，1-成功，2-失败',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_account_id (account_id),
    INDEX idx_tx_hash (tx_hash),
    INDEX idx_status (status)
);
```

### 2.3 提现功能

**功能描述：**
- 前端发起提现请求
- 验证资产余额
- 从EOA地址转出资产
- 保存提现记录

**数据库设计：**
```sql
CREATE TABLE withdrawals (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    account_id BIGINT NOT NULL COMMENT '账户ID',
    coin VARCHAR(10) NOT NULL COMMENT '币种',
    amount DECIMAL(36,18) NOT NULL COMMENT '提现金额',
    fee DECIMAL(36,18) DEFAULT 0 COMMENT '手续费',
    tx_hash VARCHAR(66) COMMENT '交易哈希（提现成功后）',
    to_address VARCHAR(42) NOT NULL COMMENT '提现地址',
    status TINYINT DEFAULT 0 COMMENT '状态：0-待处理，1-处理中，2-成功，3-失败',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_account_id (account_id),
    INDEX idx_status (status)
);
```

**API接口：**
- `POST /api/v1/bridge/withdraw` - 发起提现
- `GET /api/v1/bridge/withdraw/history` - 提现记录查询

**技术要点：**
- 资产验证：检查可用余额是否足够
- 风控检查：单笔限额、日限额、频率限制
- 异步处理：提现请求异步处理，避免阻塞

### 2.4 用户资产表

**数据库设计：**
```sql
CREATE TABLE user_assets (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    account_id BIGINT NOT NULL COMMENT '账户ID',
    coin VARCHAR(10) NOT NULL COMMENT '币种',
    total DECIMAL(36,18) DEFAULT 0 COMMENT '总资产',
    freeze DECIMAL(36,18) DEFAULT 0 COMMENT '冻结资产（订单占用）',
    available DECIMAL(36,18) DEFAULT 0 COMMENT '可用资产 = total - freeze',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_account_coin (account_id, coin),
    INDEX idx_account_id (account_id)
);
```

**技术要点：**
- 资产计算：available = total - freeze
- 并发控制：使用数据库事务保证资产操作的原子性
- 资产快照：定期记录资产快照用于对账

---

## 三、撮合交易系统

### 3.1 Spot交易板块

**功能描述：**
- 支持BTC-USDT、ETH-USDT交易对
- 下单、撤单功能
- 持仓查询
- 交易记录查询
- 资金费率（如有）

**数据库设计：**
```sql
-- 交易对配置
CREATE TABLE trading_pairs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    symbol VARCHAR(20) UNIQUE NOT NULL COMMENT '交易对：BTC-USDT, ETH-USDT',
    base_coin VARCHAR(10) NOT NULL COMMENT '基础币种',
    quote_coin VARCHAR(10) NOT NULL COMMENT '计价币种',
    min_order_amount DECIMAL(36,18) COMMENT '最小下单量',
    price_precision INT DEFAULT 8 COMMENT '价格精度',
    amount_precision INT DEFAULT 8 COMMENT '数量精度',
    status TINYINT DEFAULT 1 COMMENT '状态：1-交易中，2-暂停',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 订单表
CREATE TABLE spot_orders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id VARCHAR(32) UNIQUE NOT NULL COMMENT '订单ID',
    account_id BIGINT NOT NULL COMMENT '账户ID',
    symbol VARCHAR(20) NOT NULL COMMENT '交易对',
    side TINYINT NOT NULL COMMENT '方向：1-买入，2-卖出',
    order_type TINYINT NOT NULL COMMENT '订单类型：1-限价，2-市价',
    price DECIMAL(36,18) COMMENT '价格（限价单）',
    amount DECIMAL(36,18) NOT NULL COMMENT '数量',
    filled_amount DECIMAL(36,18) DEFAULT 0 COMMENT '已成交数量',
    filled_value DECIMAL(36,18) DEFAULT 0 COMMENT '已成交金额',
    status TINYINT DEFAULT 0 COMMENT '状态：0-待成交，1-部分成交，2-完全成交，3-已撤销',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_account_id (account_id),
    INDEX idx_symbol_status (symbol, status),
    INDEX idx_order_id (order_id)
);

-- 成交记录
CREATE TABLE spot_trades (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    trade_id VARCHAR(32) UNIQUE NOT NULL COMMENT '成交ID',
    order_id VARCHAR(32) NOT NULL COMMENT '订单ID',
    taker_order_id VARCHAR(32) COMMENT '吃单订单ID',
    maker_order_id VARCHAR(32) COMMENT '挂单订单ID',
    account_id BIGINT NOT NULL COMMENT '账户ID',
    symbol VARCHAR(20) NOT NULL COMMENT '交易对',
    side TINYINT NOT NULL COMMENT '方向：1-买入，2-卖出',
    price DECIMAL(36,18) NOT NULL COMMENT '成交价格',
    amount DECIMAL(36,18) NOT NULL COMMENT '成交数量',
    fee DECIMAL(36,18) DEFAULT 0 COMMENT '手续费',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_account_id (account_id),
    INDEX idx_symbol (symbol),
    INDEX idx_order_id (order_id)
);

-- 持仓表
CREATE TABLE spot_positions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    account_id BIGINT NOT NULL COMMENT '账户ID',
    symbol VARCHAR(20) NOT NULL COMMENT '交易对',
    base_coin VARCHAR(10) NOT NULL COMMENT '基础币种',
    quote_coin VARCHAR(10) NOT NULL COMMENT '计价币种',
    base_amount DECIMAL(36,18) DEFAULT 0 COMMENT '基础币种数量',
    quote_amount DECIMAL(36,18) DEFAULT 0 COMMENT '计价币种数量',
    avg_price DECIMAL(36,18) DEFAULT 0 COMMENT '平均成本价',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_account_symbol (account_id, symbol),
    INDEX idx_account_id (account_id)
);
```

**API接口：**
- `POST /api/v1/spot/order` - 下单
- `DELETE /api/v1/spot/order/:orderId` - 撤单
- `GET /api/v1/spot/position` - 查询持仓
- `GET /api/v1/spot/trades` - 查询交易记录
- `GET /api/v1/spot/orders` - 查询订单列表

**撮合引擎：**
- 订单簿管理：买卖盘订单簿
- 价格匹配：价格优先、时间优先
- 成交处理：更新订单状态、更新持仓、记录成交

### 3.2 Perp（永续合约）板块

**功能描述：**
- 下单检查（余额、杠杆、风险）
- 市价单、限价单支持
- 清算引擎
- 交易记录
- 订单管理、仓位管理

**数据库设计：**
```sql
-- 永续合约配置
CREATE TABLE perp_configs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    symbol VARCHAR(20) UNIQUE NOT NULL COMMENT '交易对：BTC-USDT, ETH-USDT',
    base_coin VARCHAR(10) NOT NULL COMMENT '基础币种',
    quote_coin VARCHAR(10) NOT NULL COMMENT '计价币种',
    max_leverage INT DEFAULT 10 COMMENT '最大杠杆倍数',
    min_order_amount DECIMAL(36,18) COMMENT '最小下单量',
    maintenance_margin_rate DECIMAL(10,4) DEFAULT 0.005 COMMENT '维持保证金率',
    funding_rate DECIMAL(10,8) DEFAULT 0 COMMENT '资金费率',
    status TINYINT DEFAULT 1 COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 永续订单表
CREATE TABLE perp_orders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id VARCHAR(32) UNIQUE NOT NULL COMMENT '订单ID',
    account_id BIGINT NOT NULL COMMENT '账户ID',
    symbol VARCHAR(20) NOT NULL COMMENT '交易对',
    side TINYINT NOT NULL COMMENT '方向：1-开多，2-开空，3-平多，4-平空',
    order_type TINYINT NOT NULL COMMENT '订单类型：1-限价，2-市价',
    position_side TINYINT NOT NULL COMMENT '持仓方向：1-多，2-空',
    price DECIMAL(36,18) COMMENT '价格（限价单）',
    amount DECIMAL(36,18) NOT NULL COMMENT '数量',
    leverage INT DEFAULT 1 COMMENT '杠杆倍数',
    filled_amount DECIMAL(36,18) DEFAULT 0 COMMENT '已成交数量',
    status TINYINT DEFAULT 0 COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_account_id (account_id),
    INDEX idx_symbol_status (symbol, status)
);

-- 永续持仓表
CREATE TABLE perp_positions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    account_id BIGINT NOT NULL COMMENT '账户ID',
    symbol VARCHAR(20) NOT NULL COMMENT '交易对',
    side TINYINT NOT NULL COMMENT '方向：1-多，2-空',
    size DECIMAL(36,18) DEFAULT 0 COMMENT '持仓数量',
    entry_price DECIMAL(36,18) DEFAULT 0 COMMENT '开仓均价',
    mark_price DECIMAL(36,18) DEFAULT 0 COMMENT '标记价格',
    leverage INT DEFAULT 1 COMMENT '杠杆倍数',
    margin DECIMAL(36,18) DEFAULT 0 COMMENT '保证金',
    unrealized_pnl DECIMAL(36,18) DEFAULT 0 COMMENT '未实现盈亏',
    liquidation_price DECIMAL(36,18) COMMENT '强平价格',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_account_symbol_side (account_id, symbol, side),
    INDEX idx_account_id (account_id)
);

-- 清算记录
CREATE TABLE liquidations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    account_id BIGINT NOT NULL COMMENT '账户ID',
    symbol VARCHAR(20) NOT NULL COMMENT '交易对',
    position_id BIGINT NOT NULL COMMENT '持仓ID',
    liquidation_price DECIMAL(36,18) NOT NULL COMMENT '清算价格',
    size DECIMAL(36,18) NOT NULL COMMENT '清算数量',
    pnl DECIMAL(36,18) COMMENT '盈亏',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_account_id (account_id)
);
```

**技术要点：**
- 风险控制：
  - 保证金率检查：margin_rate = (margin + unrealized_pnl) / position_value
  - 维持保证金率：当margin_rate < maintenance_margin_rate 时触发清算
- 清算引擎：
  - 定时检查持仓风险
  - 计算强平价格
  - 执行清算逻辑
- ADL（自动减仓）和保险基金：第二版实现

**API接口：**
- `POST /api/v1/perp/order` - 下单
- `DELETE /api/v1/perp/order/:orderId` - 撤单
- `GET /api/v1/perp/position` - 查询持仓
- `GET /api/v1/perp/trades` - 查询交易记录

### 3.3 预言机服务

**功能描述：**
- 第一版：直接使用币安价格
- 第二版：取多家交易所价格计算均值

**技术实现：**
- 价格获取服务：
  - 定时从币安API获取BTC/USDT、ETH/USDT价格
  - 缓存价格数据
  - 提供价格查询接口
- 价格更新频率：建议1秒更新一次

**数据库设计：**
```sql
CREATE TABLE prices (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    symbol VARCHAR(20) NOT NULL COMMENT '交易对',
    price DECIMAL(36,18) NOT NULL COMMENT '价格',
    source VARCHAR(20) DEFAULT 'binance' COMMENT '价格来源',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_symbol_created (symbol, created_at)
);
```

**API接口：**
- `GET /api/v1/oracle/price/:symbol` - 获取最新价格
- `GET /api/v1/oracle/prices` - 获取所有交易对价格

### 3.4 行情服务

**功能描述：**
- 数据存储：K线数据、深度数据、成交数据
- 前端展示：TokenView组件

**数据库设计：**
```sql
-- K线数据
CREATE TABLE klines (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    symbol VARCHAR(20) NOT NULL COMMENT '交易对',
    interval VARCHAR(10) NOT NULL COMMENT '时间间隔：1m, 5m, 15m, 1h, 4h, 1d',
    open_time BIGINT NOT NULL COMMENT '开盘时间',
    open_price DECIMAL(36,18) NOT NULL COMMENT '开盘价',
    high_price DECIMAL(36,18) NOT NULL COMMENT '最高价',
    low_price DECIMAL(36,18) NOT NULL COMMENT '最低价',
    close_price DECIMAL(36,18) NOT NULL COMMENT '收盘价',
    volume DECIMAL(36,18) DEFAULT 0 COMMENT '成交量',
    close_time BIGINT NOT NULL COMMENT '收盘时间',
    UNIQUE KEY uk_symbol_interval_time (symbol, interval, open_time),
    INDEX idx_symbol_interval (symbol, interval)
);

-- 订单簿快照
CREATE TABLE orderbook_snapshots (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    symbol VARCHAR(20) NOT NULL COMMENT '交易对',
    bids TEXT COMMENT '买盘（JSON格式）',
    asks TEXT COMMENT '卖盘（JSON格式）',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_symbol_created (symbol, created_at)
);
```

**API接口：**
- `GET /api/v1/market/klines` - 获取K线数据
- `GET /api/v1/market/depth` - 获取深度数据
- `GET /api/v1/market/trades` - 获取最新成交

### 3.5 资产配置管理

**功能描述：**
- 交易对配置管理
- 资产参数配置

**管理接口：**
- `GET /api/v1/admin/pairs` - 查询交易对配置
- `POST /api/v1/admin/pairs` - 添加交易对
- `PUT /api/v1/admin/pairs/:id` - 更新交易对配置

---

## 四、初始流动性

### 4.1 做市商合作

**功能描述：**
- 与外部做市商合作提供流动性
- 提供下单API供做市商使用

**技术要点：**
- API密钥管理：为做市商分配API Key
- 做市商标识：订单标记为做市商订单
- 费率优惠：做市商可能享受手续费优惠

---

## 五、UI设计

**功能模块：**
1. Logo设计
2. 主页
3. Trading页面（交易界面）
4. 用户资产页面
5. 开发文档

**设计资源：**
- Figma链接：https://www.figma.com/design/FU88joGHdZHAt7yu0xl6cX/helieum?node-id=105-10&t=U1tOepPud8MIwKdb-1

---

## 六、测试活动

**待补充：**
- 测试活动具体需求待明确

---

## 技术架构建议

### 后端架构
- **框架**：go-zero（已使用）
- **数据库**：MySQL（主库）+ Redis（缓存）
- **消息队列**：RabbitMQ/Kafka（订单处理、清算等异步任务）
- **区块链交互**：Web3 Go SDK

### 服务拆分建议
1. **API服务**：用户接口、订单接口
2. **撮合引擎服务**：独立的撮合服务，高性能要求
3. **清算服务**：定时清算检查
4. **区块链监听服务**：充值监听
5. **价格服务**：预言机价格获取

### 关键技术点
1. **并发控制**：订单处理、资产操作需要保证原子性
2. **性能优化**：撮合引擎需要高性能，考虑内存撮合
3. **数据一致性**：资产、订单状态需要强一致性
4. **风控系统**：下单前风险检查、清算风险监控

---

## 开发优先级建议

### Phase 1（MVP）
1. ✅ 用户账户系统（登录、鉴权）
2. ✅ 资产配置表
3. ✅ 充值监听（USDT、ETH）
4. ✅ Spot交易（下单、撤单、撮合）
5. ✅ 价格服务（币安价格）
6. ✅ 基础UI（交易页面、资产页面）

### Phase 2
1. BTC充值监听
2. 提现功能
3. Perp交易
4. 清算引擎
5. 行情服务完善

### Phase 3
1. ADL（自动减仓）
2. 保险基金
3. 多家交易所价格聚合
4. 做市商API

---

## 注意事项

1. **安全性**：
   - 签名验证必须严格
   - Token过期时间合理设置
   - 资产操作必须加锁
   - 防止重放攻击

2. **性能**：
   - 撮合引擎需要高性能设计
   - 数据库索引优化
   - 缓存热点数据

3. **监控**：
   - 充值监听服务监控
   - 撮合引擎性能监控
   - 清算风险监控
   - 资产对账

4. **测试**：
   - 单元测试
   - 集成测试
   - 压力测试（撮合引擎）
   - 安全测试
