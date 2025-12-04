/*
Navicat MySQL Data Transfer

Source Server         : localhost
Source Server Version : 80027
Source Host           : 127.0.0.1:3306
Source Database       : go_gateway

Target Server Type    : MYSQL
Target Server Version : 80027
File Encoding         : 65001

Date: 2022-11-10 19:43:42
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for gateway_admin
-- ----------------------------
DROP TABLE IF EXISTS `gateway_admin`;
CREATE TABLE `gateway_admin` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `user_name` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名',
  `salt` varchar(50) NOT NULL DEFAULT '' COMMENT '盐',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '密码',
  `create_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '新增时间',
  `update_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '更新时间',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb3 COMMENT='管理员表';

-- ----------------------------
-- Records of gateway_admin
-- ----------------------------
INSERT INTO `gateway_admin` VALUES ('1', 'admin', 'admin', '2823d896e9822c0833d41d4904f0c00756d718570fce49b9a379a62c804689d3', '2020-04-10 16:42:05', '2022-11-09 20:44:20', '0');

-- ----------------------------
-- Table structure for gateway_app
-- ----------------------------
DROP TABLE IF EXISTS `gateway_app`;
CREATE TABLE `gateway_app` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `app_id` varchar(255) NOT NULL DEFAULT '' COMMENT '租户id',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '租户名称',
  `secret` varchar(255) NOT NULL DEFAULT '' COMMENT '密钥',
  `white_ips` varchar(1000) NOT NULL DEFAULT '' COMMENT 'ip白名单，支持前缀匹配',
  `qpd` bigint NOT NULL DEFAULT '0' COMMENT '日请求量限制',
  `qps` bigint NOT NULL DEFAULT '0' COMMENT '每秒请求量限制',
  `create_at` datetime NOT NULL COMMENT '添加时间',
  `update_at` datetime NOT NULL COMMENT '更新时间',
  `is_delete` tinyint NOT NULL DEFAULT '0' COMMENT '是否删除 1=删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8mb3 COMMENT='网关租户表';

-- ----------------------------
-- Records of gateway_app
-- ----------------------------
INSERT INTO `gateway_app` VALUES ('31', 'app_id_a', '租户A', '449441eb5e72dca9c42a12f3924ea3a2', 'white_ips', '100000', '100', '2022-11-09 20:44:20', '2022-11-09 20:44:20', '0');
INSERT INTO `gateway_app` VALUES ('32', 'app_id_b', '租户B', '8d7b11ec9be0e59a36b52f32366c09cb', '', '20', '0', '2022-11-09 20:44:20', '2022-11-09 20:44:20', '0');

-- ----------------------------
-- Table structure for gateway_service_access_control
-- ----------------------------
DROP TABLE IF EXISTS `gateway_service_access_control`;
CREATE TABLE `gateway_service_access_control` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `service_id` bigint NOT NULL DEFAULT '0' COMMENT '服务id',
  `open_auth` tinyint NOT NULL DEFAULT '0' COMMENT '是否开启权限 1=开启',
  `black_list` varchar(1000) NOT NULL DEFAULT '' COMMENT '黑名单ip',
  `white_list` varchar(1000) NOT NULL DEFAULT '' COMMENT '白名单ip',
  `white_host_name` varchar(1000) NOT NULL DEFAULT '' COMMENT '白名单主机',
  `clientip_flow_limit` int NOT NULL DEFAULT '0' COMMENT '客户端ip限流',
  `service_flow_limit` int NOT NULL DEFAULT '0' COMMENT '服务端限流',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=191 DEFAULT CHARSET=utf8mb3 COMMENT='网关权限控制表';

-- ----------------------------
-- Records of gateway_service_access_control
-- ----------------------------
INSERT INTO `gateway_service_access_control` VALUES ('184', '56', '0', '192.168.1.0', '', '', '0', '0');
INSERT INTO `gateway_service_access_control` VALUES ('185', '57', '0', '', '127.0.0.1,127.0.0.2', '', '0', '0');
INSERT INTO `gateway_service_access_control` VALUES ('186', '58', '1', '', '', '', '0', '0');
INSERT INTO `gateway_service_access_control` VALUES ('187', '59', '1', '127.0.0.1', '', '', '2', '3');
INSERT INTO `gateway_service_access_control` VALUES ('188', '60', '1', '', '', '', '11', '12');
INSERT INTO `gateway_service_access_control` VALUES ('189', '61', '0', '', '', '', '45', '34');
INSERT INTO `gateway_service_access_control` VALUES ('190', '62', '0', '', '', '', '0', '0');

-- ----------------------------
-- Table structure for gateway_service_grpc_rule
-- ----------------------------
DROP TABLE IF EXISTS `gateway_service_grpc_rule`;
CREATE TABLE `gateway_service_grpc_rule` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `service_id` bigint NOT NULL DEFAULT '0' COMMENT '服务id',
  `port` int NOT NULL DEFAULT '0' COMMENT '端口',
  `header_transfor` varchar(5000) NOT NULL DEFAULT '' COMMENT 'header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多个逗号间隔',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=174 DEFAULT CHARSET=utf8mb3 COMMENT='网关路由匹配表';

-- ----------------------------
-- Records of gateway_service_grpc_rule
-- ----------------------------
INSERT INTO `gateway_service_grpc_rule` VALUES ('173', '58', '8012', 'add meta_name meta_value');

-- ----------------------------
-- Table structure for gateway_service_http_rule
-- ----------------------------
DROP TABLE IF EXISTS `gateway_service_http_rule`;
CREATE TABLE `gateway_service_http_rule` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `service_id` bigint NOT NULL COMMENT '服务id',
  `rule_type` tinyint NOT NULL DEFAULT '0' COMMENT '匹配类型 0=url前缀url_prefix 1=域名domain ',
  `rule` varchar(255) NOT NULL DEFAULT '' COMMENT 'type=domain表示域名，type=url_prefix时表示url前缀',
  `need_https` tinyint NOT NULL DEFAULT '0' COMMENT '支持https 1=支持',
  `need_strip_uri` tinyint NOT NULL DEFAULT '0' COMMENT '启用strip_uri 1=启用',
  `need_websocket` tinyint NOT NULL DEFAULT '0' COMMENT '是否支持websocket 1=支持',
  `url_rewrite` varchar(5000) NOT NULL DEFAULT '' COMMENT 'url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔',
  `header_transfor` varchar(5000) NOT NULL DEFAULT '' COMMENT 'header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多个逗号间隔',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=182 DEFAULT CHARSET=utf8mb3 COMMENT='网关路由匹配表';

-- ----------------------------
-- Records of gateway_service_http_rule
-- ----------------------------
INSERT INTO `gateway_service_http_rule` VALUES ('177', '56', '0', '/test_http_service', '1', '1', '1', '^/test_http_service/abb/(.*) /test_http_service/bba/$1', 'add header_name header_value');
INSERT INTO `gateway_service_http_rule` VALUES ('178', '59', '1', 'test.com', '0', '1', '1', '', 'add headername headervalue');
INSERT INTO `gateway_service_http_rule` VALUES ('179', '60', '0', '/test_strip_uri', '0', '1', '0', '^/aaa/(.*) /bbb/$1', '');
INSERT INTO `gateway_service_http_rule` VALUES ('180', '61', '0', '/test_https_server', '1', '1', '0', '', '');
INSERT INTO `gateway_service_http_rule` VALUES ('181', '62', '0', '/test_httpservice_lwzy', '1', '0', '0', '', '');

-- ----------------------------
-- Table structure for gateway_service_info
-- ----------------------------
DROP TABLE IF EXISTS `gateway_service_info`;
CREATE TABLE `gateway_service_info` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `load_type` tinyint NOT NULL DEFAULT '0' COMMENT '负载类型 0=http 1=tcp 2=grpc',
  `service_name` varchar(255) NOT NULL DEFAULT '' COMMENT '服务名称 6-128 数字字母下划线',
  `service_desc` varchar(255) NOT NULL DEFAULT '' COMMENT '服务描述',
  `create_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '添加时间',
  `update_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '更新时间',
  `is_delete` tinyint DEFAULT '0' COMMENT '是否删除 1=删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=63 DEFAULT CHARSET=utf8mb3 COMMENT='网关基本信息表';

-- ----------------------------
-- Records of gateway_service_info
-- ----------------------------
INSERT INTO `gateway_service_info` VALUES ('56', '0', 'test_http_service', '测试HTTP代理', '2022-11-09 00:54:45', '2022-11-09 20:44:20', '0');
INSERT INTO `gateway_service_info` VALUES ('57', '1', 'test_tcp_service', '测试TCP代理', '2022-11-09 14:03:09', '2022-11-09 20:44:20', '0');
INSERT INTO `gateway_service_info` VALUES ('58', '2', 'test_grpc_service', '测试GRPC服务', '2022-11-09 07:20:16', '2022-11-09 20:44:20', '0');
INSERT INTO `gateway_service_info` VALUES ('59', '0', 'test.com:8080', '测试域名接入', '2022-11-09 22:54:14', '2022-11-09 20:44:20', '0');
INSERT INTO `gateway_service_info` VALUES ('60', '0', 'test_strip_uri', '测试路径接入', '2022-11-09 06:55:26', '2022-11-09 20:44:20', '0');
INSERT INTO `gateway_service_info` VALUES ('61', '0', 'test_https_server', '测试https服务', '2022-11-09 12:22:33', '2022-11-09 20:44:20', '0');
INSERT INTO `gateway_service_info` VALUES ('62', '0', 'httpservice_lwzy', '测试http服务_lwzy', '2022-11-09 20:49:29', '2022-11-09 20:44:20', '0');

-- ----------------------------
-- Table structure for gateway_service_load_balance
-- ----------------------------
DROP TABLE IF EXISTS `gateway_service_load_balance`;
CREATE TABLE `gateway_service_load_balance` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `service_id` bigint NOT NULL DEFAULT '0' COMMENT '服务id',
  `check_method` tinyint NOT NULL DEFAULT '0' COMMENT '检查方法 0=tcpchk,检测端口是否握手成功',
  `check_timeout` int NOT NULL DEFAULT '0' COMMENT 'check超时时间,单位s',
  `check_interval` int NOT NULL DEFAULT '0' COMMENT '检查间隔, 单位s',
  `round_type` tinyint NOT NULL DEFAULT '2' COMMENT '轮询方式 0=random 1=round-robin 2=weight_round-robin 3=ip_hash',
  `ip_list` varchar(2000) NOT NULL DEFAULT '' COMMENT 'ip列表',
  `weight_list` varchar(2000) NOT NULL DEFAULT '' COMMENT '权重列表',
  `forbid_list` varchar(2000) NOT NULL DEFAULT '' COMMENT '禁用ip列表',
  `upstream_connect_timeout` int NOT NULL DEFAULT '0' COMMENT '建立连接超时, 单位s',
  `upstream_header_timeout` int NOT NULL DEFAULT '0' COMMENT '获取header超时, 单位s',
  `upstream_idle_timeout` int NOT NULL DEFAULT '0' COMMENT '链接最大空闲时间, 单位s',
  `upstream_max_idle` int NOT NULL DEFAULT '0' COMMENT '最大空闲链接数',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=191 DEFAULT CHARSET=utf8mb3 COMMENT='网关负载表';

-- ----------------------------
-- Records of gateway_service_load_balance
-- ----------------------------
INSERT INTO `gateway_service_load_balance` VALUES ('185', '57', '0', '2', '5', '2', '127.0.0.1:6379', '50', '', '0', '0', '0', '0');
INSERT INTO `gateway_service_load_balance` VALUES ('186', '58', '0', '2', '5', '2', '127.0.0.1:8005', '50', '', '0', '0', '0', '0');
INSERT INTO `gateway_service_load_balance` VALUES ('190', '62', '0', '2', '5', '2', '127.0.0.1:8080', '100', '', '0', '0', '0', '0');

-- ----------------------------
-- Table structure for gateway_service_tcp_rule
-- ----------------------------
DROP TABLE IF EXISTS `gateway_service_tcp_rule`;
CREATE TABLE `gateway_service_tcp_rule` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `service_id` bigint NOT NULL COMMENT '服务id',
  `port` int NOT NULL DEFAULT '0' COMMENT '端口号',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=182 DEFAULT CHARSET=utf8mb3 COMMENT='网关路由匹配表';

-- ----------------------------
-- Records of gateway_service_tcp_rule
-- ----------------------------
INSERT INTO `gateway_service_tcp_rule` VALUES ('181', '57', '8011');
