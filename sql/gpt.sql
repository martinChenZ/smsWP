/*
 Navicat Premium Data Transfer

 Source Server         : home
 Source Server Type    : MySQL
 Source Server Version : 80019
 Source Host           : localhost:53306
 Source Schema         : gpt

 Target Server Type    : MySQL
 Target Server Version : 80019
 File Encoding         : 65001

 Date: 05/04/2023 17:46:26
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for gpt_log
-- ----------------------------
DROP TABLE IF EXISTS `gpt_log`;
CREATE TABLE `gpt_log`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT,
  `question` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `response` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `update_time` timestamp(0) NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0),
  `api_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `request_ip` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of gpt_log
-- ----------------------------
INSERT INTO `gpt_log` VALUES (8, 'Ai会替代人类工作吗。', '请求失败', '2023-04-05 17:22:52', '666', '127.0.0.1');
INSERT INTO `gpt_log` VALUES (9, 'Ai会替代人类工作吗。', '请求失败', '2023-04-05 17:33:08', '666', '127.0.0.1');
INSERT INTO `gpt_log` VALUES (11, 'Ai会替代人类工作吗。', '请求失败', '2023-04-05 17:43:29', '91573eda-c758-4188-92af-a8e85f711c02', '');
INSERT INTO `gpt_log` VALUES (12, 'Ai会替代人类工作吗。', '请求失败', '2023-04-05 17:43:54', '91573eda-c758-4188-92af-a8e85f711c02', '');
INSERT INTO `gpt_log` VALUES (13, 'Ai会替代人类工作吗。', '请求失败', '2023-04-05 17:43:59', '91573eda-c758-4188-92af-a8e85f711c02', '');
INSERT INTO `gpt_log` VALUES (14, 'Ai会替代人类工作吗。', '请求失败', '2023-04-05 17:44:34', '91573eda-c758-4188-92af-a8e85f711c02', '127.0.0.1');

-- ----------------------------
-- Table structure for gpt_user
-- ----------------------------
DROP TABLE IF EXISTS `gpt_user`;
CREATE TABLE `gpt_user`  (
  `api_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `balance` int(0) NULL DEFAULT NULL COMMENT '次数\r\n次数',
  `order_no` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '商单号码',
  `update_time` timestamp(0) NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0)
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of gpt_user
-- ----------------------------
INSERT INTO `gpt_user` VALUES ('dbf24106-ec52-45c3-b0e4-df8ab79774fa', 1, NULL, '2023-04-01 15:24:44');
INSERT INTO `gpt_user` VALUES ('91573eda-c758-4188-92af-a8e85f711c02', 0, '222', '2023-04-05 09:43:59');

SET FOREIGN_KEY_CHECKS = 1;
