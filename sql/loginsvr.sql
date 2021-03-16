/*
Navicat MySQL Data Transfer

Source Server         : localmysql
Source Server Version : 80022
Source Host           : localhost:3306
Source Database       : loginsvr

Target Server Type    : MYSQL
Target Server Version : 80022
File Encoding         : 65001

Date: 2021-03-16 11:34:54
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for accounts
-- ----------------------------
DROP TABLE IF EXISTS `accounts`;
CREATE TABLE `accounts` (
  `userid` bigint NOT NULL AUTO_INCREMENT COMMENT '用户id，唯一',
  `username` varchar(255) NOT NULL COMMENT '用户名',
  `passwd` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
  `thirdid` varchar(255) DEFAULT NULL COMMENT '第三方id',
  `nickname` varchar(255) NOT NULL COMMENT '昵称',
  `createtime` datetime NOT NULL COMMENT '账号创建时间',
  `logintime` datetime DEFAULT NULL COMMENT '最近登入时间',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `accounttype` int NOT NULL COMMENT '账号类型，0注册用户 1微信登入 2vivo。。。',
  `status` int NOT NULL COMMENT '账号状态：0正常 1冻结 。。',
  PRIMARY KEY (`userid`) USING BTREE,
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of accounts
-- ----------------------------
INSERT INTO `accounts` VALUES ('4', 'test001', '123', '', 'test001', '2020-12-25 14:37:28', null, '', '0', '0');
INSERT INTO `accounts` VALUES ('6', 'test002', '123', '', 'test002', '2020-12-25 14:39:43', null, '', '0', '0');
INSERT INTO `accounts` VALUES ('8', 'wechat_', 'cd1988b4d73d94a6a98bffd6af0a8148', '', 'wechat_defalut', '2021-02-23 14:36:10', null, '', '1', '0');

-- ----------------------------
-- Table structure for halls
-- ----------------------------
DROP TABLE IF EXISTS `halls`;
CREATE TABLE `halls` (
  `serverid` int NOT NULL,
  `servername` varchar(255) DEFAULT NULL,
  `address` varchar(255) NOT NULL COMMENT 'ip:port',
  `channel` int NOT NULL COMMENT '渠道类型 位运算   001安卓服  010 IOS服  100pc服',
  `status` int NOT NULL COMMENT '服务器状态 0开启 1关闭 2维护',
  PRIMARY KEY (`serverid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of halls
-- ----------------------------
INSERT INTO `halls` VALUES ('1000', '大厅1', '192.168.43.22:8080', '7', '0');
INSERT INTO `halls` VALUES ('1001', '大厅-安卓', '192.168.43.22.8081', '1', '0');
INSERT INTO `halls` VALUES ('1002', '大厅-IOS', '192.168.43.22:8082', '2', '0');
INSERT INTO `halls` VALUES ('1003', '大厅-PC', '192.168.43.22:8083', '4', '0');

-- ----------------------------
-- Table structure for tokens
-- ----------------------------
DROP TABLE IF EXISTS `tokens`;
CREATE TABLE `tokens` (
  `userid` bigint NOT NULL,
  `token` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `updatetime` datetime NOT NULL,
  PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of tokens
-- ----------------------------
INSERT INTO `tokens` VALUES ('4', '80e091746cfee17834721a964819c26a', '2021-02-03 15:37:38');

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `Id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `Salt` varchar(255) DEFAULT NULL,
  `AgeId` int DEFAULT NULL,
  `Passwd` varchar(200) DEFAULT NULL,
  `create` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `created1` datetime DEFAULT NULL,
  `age_id` int DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of user
-- ----------------------------

-- ----------------------------
-- Table structure for userrealinfo
-- ----------------------------
DROP TABLE IF EXISTS `userrealinfo`;
CREATE TABLE `userrealinfo` (
  `userid` bigint NOT NULL,
  `name` varchar(32) DEFAULT NULL,
  `identity` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `gender` int DEFAULT NULL,
  `addr` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of userrealinfo
-- ----------------------------
INSERT INTO `userrealinfo` VALUES ('4', '赵4', '110101201008070221', '1', '上海市');
INSERT INTO `userrealinfo` VALUES ('5', '赵5', '110101200003078932', '0', null);
INSERT INTO `userrealinfo` VALUES ('7', '赵4', '110101199003075031', '0', '');

-- ----------------------------
-- Table structure for wechatuserinfo
-- ----------------------------
DROP TABLE IF EXISTS `wechatuserinfo`;
CREATE TABLE `wechatuserinfo` (
  `openid` varchar(255) NOT NULL,
  `accesstoken` varchar(255) NOT NULL,
  `refreshtoken` varchar(255) NOT NULL,
  `nickname` varchar(255) DEFAULT NULL COMMENT '普通用户昵称',
  `sex` int DEFAULT NULL COMMENT '普通用户性别，1 为男性，2 为女性',
  `province` varchar(50) DEFAULT NULL COMMENT '普通用户个人资料填写的省份',
  `city` varchar(50) DEFAULT NULL COMMENT '普通用户个人资料填写的城市',
  `country` varchar(50) DEFAULT NULL COMMENT '国家，如中国为 CN',
  `headimgurl` varchar(500) DEFAULT NULL COMMENT '用户头像，最后一个数值代表正方形头像大小（有 0、46、64、96、132 数值可选，0 代表 640*640 正方形头像），用户没有头像时该项为空',
  `privilege` varchar(255) DEFAULT NULL,
  `unionid` varchar(255) DEFAULT NULL COMMENT '用户统一标识。针对一个微信开放平台帐号下的应用，同一用户的 unionid 是唯一的。',
  PRIMARY KEY (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of wechatuserinfo
-- ----------------------------
