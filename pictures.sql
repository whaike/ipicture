/*
 Navicat Premium Data Transfer

 Source Server         : ipicture
 Source Server Type    : SQLite
 Source Server Version : 3030001
 Source Schema         : main

 Target Server Type    : SQLite
 Target Server Version : 3030001
 File Encoding         : 65001

 Date: 05/01/2023 17:33:19
*/

PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for pictures
-- ----------------------------
DROP TABLE IF EXISTS "pictures";
CREATE TABLE "pictures" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "name" text(255) NOT NULL,
  "path" text NOT NULL DEFAULT '',
  "type" text(50) NOT NULL,
  "suffix" text(20) NOT NULL,
  "tags" text NOT NULL,
  "shoot_at" TEXT(50) NOT NULL,
  "lng" TEXT(50) NOT NULL,
  "lat" TEXT(50) NOT NULL
);

-- ----------------------------
-- Auto increment value for pictures
-- ----------------------------

PRAGMA foreign_keys = true;
