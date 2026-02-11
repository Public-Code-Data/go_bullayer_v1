/*
 Navicat Premium Dump SQL

 Source Server         : forest-new
 Source Server Type    : MySQL
 Source Server Version : 80028 (8.0.28-231003)
 Source Host           : 192.168.2.12:3306
 Source Schema         : bullayer_test_data

 Target Server Type    : MySQL
 Target Server Version : 80028 (8.0.28-231003)
 File Encoding         : 65001

 Date: 09/02/2026 20:33:36
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for accounts
-- ----------------------------
DROP TABLE IF EXISTS `accounts`;
CREATE TABLE `accounts` (
  `account_id` bigint NOT NULL AUTO_INCREMENT,
  `address` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '钱包地址',
  `status` tinyint DEFAULT '1' COMMENT '账户状态：1-正常，2-冻结，3-禁用',
  `invitation_code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '邀请码',
  `wallet_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '钱包类型：metamask, okx',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`account_id`),
  UNIQUE KEY `address` (`address`),
  UNIQUE KEY `invitation_code` (`invitation_code`),
  KEY `idx_address` (`address`),
  KEY `idx_invitation_code` (`invitation_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='账户表';

-- ----------------------------
-- Table structure for coin_configs
-- ----------------------------
DROP TABLE IF EXISTS `coin_configs`;
CREATE TABLE `coin_configs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `coin` varchar(16) NOT NULL COMMENT '币种：BTC, ETH, USDT',
  `coin_address` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '合约地址（ERC20）',
  `min_deposit` decimal(36,18) NOT NULL COMMENT '最小充值金额',
  `status` tinyint DEFAULT '1' COMMENT '状态：1-启用，2-禁用',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_coin` (`coin`),
  UNIQUE KEY `uk_coin_address` (`coin_address`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='币配置表';

-- ----------------------------
-- Table structure for klines
-- ----------------------------
DROP TABLE IF EXISTS `klines`;
CREATE TABLE `klines` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `symbol` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '交易对',
  `intervals` varchar(10) NOT NULL COMMENT '时间间隔：1m, 5m, 15m, 1h, 4h, 1d',
  `open_time` bigint NOT NULL COMMENT '开盘时间',
  `open_price` decimal(36,18) NOT NULL COMMENT '开盘价',
  `high_price` decimal(36,18) NOT NULL COMMENT '最高价',
  `low_price` decimal(36,18) NOT NULL COMMENT '最低价',
  `close_price` decimal(36,18) NOT NULL COMMENT '收盘价',
  `volume` decimal(36,18) DEFAULT '0.000000000000000000' COMMENT '成交量',
  `close_time` bigint NOT NULL COMMENT '收盘时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_symbol_interval_time` (`symbol`,`intervals`,`open_time`),
  KEY `idx_symbol_interval` (`symbol`,`intervals`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='K线数据';

-- ----------------------------
-- Table structure for orderbook_snapshots
-- ----------------------------
DROP TABLE IF EXISTS `orderbook_snapshots`;
CREATE TABLE `orderbook_snapshots` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `symbol` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '交易对',
  `bids` text COMMENT '买盘（JSON格式）',
  `asks` text COMMENT '卖盘（JSON格式）',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_symbol_created` (`symbol`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订单簿快照';

-- ----------------------------
-- Table structure for perp_configs
-- ----------------------------
DROP TABLE IF EXISTS `perp_configs`;
CREATE TABLE `perp_configs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `symbol` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '交易对：BTC-USDT, ETH-USDT',
  `base_coin` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '基础币种',
  `quote_coin` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '计价币种',
  `max_leverage` int DEFAULT '10' COMMENT '最大杠杆倍数',
  `min_order_amount` decimal(36,18) DEFAULT NULL COMMENT '最小下单量',
  `maintenance_margin_rate` decimal(10,4) DEFAULT '0.0050' COMMENT '维持保证金率',
  `funding_rate` decimal(10,8) DEFAULT '0.00000000' COMMENT '资金费率',
  `status` tinyint DEFAULT '1' COMMENT '状态',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `symbol` (`symbol`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='永续合约配置';

-- ----------------------------
-- Table structure for perp_liquidations
-- ----------------------------
DROP TABLE IF EXISTS `perp_liquidations`;
CREATE TABLE `perp_liquidations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `account_id` bigint NOT NULL COMMENT '账户ID',
  `symbol` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '交易对',
  `position_id` bigint NOT NULL COMMENT '持仓ID',
  `liquidation_price` decimal(36,18) NOT NULL COMMENT '清算价格',
  `size` decimal(36,18) NOT NULL COMMENT '清算数量',
  `pnl` decimal(36,18) DEFAULT NULL COMMENT '盈亏',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_account_id` (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='清算记录';

-- ----------------------------
-- Table structure for perp_orders
-- ----------------------------
DROP TABLE IF EXISTS `perp_orders`;
CREATE TABLE `perp_orders` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `order_id` varchar(32) NOT NULL COMMENT '订单ID',
  `account_id` bigint NOT NULL COMMENT '账户ID',
  `symbol` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '交易对',
  `side` tinyint NOT NULL COMMENT '方向：1-开多，2-开空，3-平多，4-平空',
  `order_type` tinyint NOT NULL COMMENT '订单类型：1-限价，2-市价',
  `position_side` tinyint NOT NULL COMMENT '持仓方向：1-多，2-空',
  `price` decimal(36,18) DEFAULT NULL COMMENT '价格（限价单）',
  `amount` decimal(36,18) NOT NULL COMMENT '数量',
  `leverage` int DEFAULT '1' COMMENT '杠杆倍数',
  `filled_amount` decimal(36,18) DEFAULT '0.000000000000000000' COMMENT '已成交数量',
  `status` tinyint DEFAULT '0' COMMENT '状态',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `order_id` (`order_id`),
  KEY `idx_account_id` (`account_id`),
  KEY `idx_symbol_status` (`symbol`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='永续订单表';

-- ----------------------------
-- Table structure for perp_positions
-- ----------------------------
DROP TABLE IF EXISTS `perp_positions`;
CREATE TABLE `perp_positions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `account_id` bigint NOT NULL COMMENT '账户ID',
  `symbol` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '交易对',
  `side` tinyint NOT NULL COMMENT '方向：1-多，2-空',
  `size` decimal(36,18) DEFAULT '0.000000000000000000' COMMENT '持仓数量',
  `entry_price` decimal(36,18) DEFAULT '0.000000000000000000' COMMENT '开仓均价',
  `mark_price` decimal(36,18) DEFAULT '0.000000000000000000' COMMENT '标记价格',
  `leverage` int DEFAULT '1' COMMENT '杠杆倍数',
  `margin` decimal(36,18) DEFAULT '0.000000000000000000' COMMENT '保证金',
  `unrealized_pnl` decimal(36,18) DEFAULT '0.000000000000000000' COMMENT '未实现盈亏',
  `liquidation_price` decimal(36,18) DEFAULT NULL COMMENT '强平价格',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_account_symbol_side` (`account_id`,`symbol`,`side`),
  KEY `idx_account_id` (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='永续持仓表';

-- ----------------------------
-- Table structure for prices
-- ----------------------------
DROP TABLE IF EXISTS `prices`;
CREATE TABLE `prices` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `symbol` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '交易对',
  `price` decimal(36,18) NOT NULL COMMENT '价格',
  `source` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT 'binance' COMMENT '价格来源',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_symbol_created` (`symbol`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='价格表';

-- ----------------------------
-- Table structure for spot_orders
-- ----------------------------
DROP TABLE IF EXISTS `spot_orders`;
CREATE TABLE `spot_orders` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `order_id` varchar(32) NOT NULL COMMENT '订单ID',
  `account_id` bigint NOT NULL COMMENT '账户ID',
  `symbol` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '交易对',
  `side` tinyint NOT NULL COMMENT '方向：1-买入，2-卖出',
  `order_type` tinyint NOT NULL COMMENT '订单类型：1-限价，2-市价',
  `price` decimal(36,18) DEFAULT NULL COMMENT '价格（限价单）',
  `amount` decimal(36,18) NOT NULL COMMENT '数量',
  `filled_amount` decimal(36,18) DEFAULT '0.000000000000000000' COMMENT '已成交数量',
  `filled_value` decimal(36,18) DEFAULT '0.000000000000000000' COMMENT '已成交金额',
  `status` tinyint DEFAULT '0' COMMENT '状态：0-待成交，1-部分成交，2-完全成交，3-已撤销',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `order_id` (`order_id`),
  KEY `idx_account_id` (`account_id`),
  KEY `idx_symbol_status` (`symbol`,`status`),
  KEY `idx_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订单表';

-- ----------------------------
-- Table structure for spot_positions
-- ----------------------------
DROP TABLE IF EXISTS `spot_positions`;
CREATE TABLE `spot_positions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `account_id` bigint NOT NULL COMMENT '账户ID',
  `symbol` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '交易对',
  `base_coin` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '基础币种',
  `quote_coin` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '计价币种',
  `base_amount` decimal(36,18) DEFAULT '0.000000000000000000' COMMENT '基础币种数量',
  `quote_amount` decimal(36,18) DEFAULT '0.000000000000000000' COMMENT '计价币种数量',
  `avg_price` decimal(36,18) DEFAULT '0.000000000000000000' COMMENT '平均成本价',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_account_symbol` (`account_id`,`symbol`),
  KEY `idx_account_id` (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='持仓表';

-- ----------------------------
-- Table structure for spot_trades
-- ----------------------------
DROP TABLE IF EXISTS `spot_trades`;
CREATE TABLE `spot_trades` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `trade_id` varchar(32) NOT NULL COMMENT '成交ID',
  `order_id` varchar(32) NOT NULL COMMENT '订单ID',
  `taker_order_id` varchar(32) DEFAULT NULL COMMENT '吃单订单ID',
  `maker_order_id` varchar(32) DEFAULT NULL COMMENT '挂单订单ID',
  `account_id` bigint NOT NULL COMMENT '账户ID',
  `symbol` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '交易对',
  `side` tinyint NOT NULL COMMENT '方向：1-买入，2-卖出',
  `price` decimal(36,18) NOT NULL COMMENT '成交价格',
  `amount` decimal(36,18) NOT NULL COMMENT '成交数量',
  `fee` decimal(36,18) DEFAULT '0.000000000000000000' COMMENT '手续费',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `trade_id` (`trade_id`),
  KEY `idx_account_id` (`account_id`),
  KEY `idx_symbol` (`symbol`),
  KEY `idx_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='成交记录';

-- ----------------------------
-- Table structure for trading_pairs
-- ----------------------------
DROP TABLE IF EXISTS `trading_pairs`;
CREATE TABLE `trading_pairs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `symbol` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '交易对：BTC-USDT, ETH-USDT',
  `base_coin` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '基础币种',
  `quote_coin` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '计价币种',
  `min_order_amount` decimal(36,18) DEFAULT NULL COMMENT '最小下单量',
  `price_precision` int DEFAULT '8' COMMENT '价格精度',
  `amount_precision` int DEFAULT '8' COMMENT '数量精度',
  `status` tinyint DEFAULT '1' COMMENT '状态：1-交易中，2-暂停',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `symbol` (`symbol`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='交易对配置';

-- ----------------------------
-- Table structure for transactions
-- ----------------------------
DROP TABLE IF EXISTS `transactions`;
CREATE TABLE `transactions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `block_number` bigint DEFAULT NULL COMMENT '区块号',
  `tx_hash` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '交易哈希',
  `tx_type` tinyint DEFAULT '0' COMMENT '状态：0-未知，1-充值:deposits，2-提现:withdrawals',
  `account_id` bigint NOT NULL COMMENT '账户ID',
  `coin` varchar(16) NOT NULL COMMENT '币种',
  `coin_address` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '合约地址（ERC20）',
  `amount` decimal(36,18) NOT NULL COMMENT '充值金额',
  `gas` decimal(36,18) DEFAULT '0.000000000000000000' COMMENT 'gas',
  `fee` decimal(36,18) DEFAULT '0.000000000000000000' COMMENT '手续费',
  `from_address` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '发送地址',
  `to_address` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '接收地址',
  `confirmations` int DEFAULT '0' COMMENT '确认数',
  `status` tinyint DEFAULT '0' COMMENT '状态：0-待确认，1-成功，2-失败',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `tx_hash` (`tx_hash`),
  KEY `idx_account_id` (`account_id`),
  KEY `idx_tx_hash` (`tx_hash`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='充值提现交易表';

-- ----------------------------
-- Table structure for user_assets
-- ----------------------------
DROP TABLE IF EXISTS `user_assets`;
CREATE TABLE `user_assets` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `account_id` bigint NOT NULL COMMENT '账户ID',
  `coin` varchar(16) NOT NULL COMMENT '币种',
  `total` decimal(36,18) DEFAULT '0.000000000000000000' COMMENT '总资产',
  `freeze` decimal(36,18) DEFAULT '0.000000000000000000' COMMENT '冻结资产（订单占用）',
  `available` decimal(36,18) DEFAULT '0.000000000000000000' COMMENT '可用资产 = total - freeze',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_account_coin` (`account_id`,`coin`),
  KEY `idx_account_id` (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户资产表';

SET FOREIGN_KEY_CHECKS = 1;
