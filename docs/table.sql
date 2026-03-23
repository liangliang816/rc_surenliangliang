-- 创建请求记录表
CREATE TABLE IF NOT EXISTS `api_job` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `name` VARCHAR(255) NOT NULL COMMENT '名称',
  `code` VARCHAR(32) NOT NULL UNIQUE COMMENT '随机code',
  `url` VARCHAR(2048) NOT NULL COMMENT '请求URL',
  `request_header` TEXT COMMENT '请求头信息',
  `request_body` TEXT COMMENT '请求体信息',
  `is_executed` TINYINT(1) DEFAULT 0 COMMENT '是否执行(0:未执行,1:已执行)',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  INDEX `idx_code` (`code`),
  INDEX `idx_is_executed` (`is_executed`),
  INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='请求记录表';

-- 创建执行记录表
CREATE TABLE IF NOT EXISTS `api_run_record` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `api_code` VARCHAR(32) NOT NULL COMMENT '关联的请求任务code',
  `status` VARCHAR(50) DEFAULT 'pending' COMMENT '执行状态(pending/running/success/failed)',
  `execution_count` INT UNSIGNED DEFAULT 0 COMMENT '执行次数',
  `execution_time` DATETIME COMMENT '执行时间',
  `response_result` TEXT COMMENT '返回结果',
  `is_success` TINYINT(1) DEFAULT 0 COMMENT '是否成功(0:失败,1:成功)',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  INDEX `idx_api_code` (`api_code`),
  INDEX `idx_status` (`status`),
  INDEX `idx_created_at` (`created_at`),
  CONSTRAINT `fk_api_record_job` FOREIGN KEY (`api_code`) REFERENCES `api_job` (`code`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='执行记录表';

-- 插入示例数据
INSERT INTO `api_job` (`name`, `code`, `url`, `request_header`, `request_body`, `is_executed`) VALUES
('示例请求1', 'abc123', 'https://example.com/api', '{"Content-Type": "application/json"}', '{"key": "value"}', 1),
('示例请求2', 'def456', 'https://example.com/api/users', '{"Authorization": "Bearer token"}', '{"name": "test", "age": 20}', 1);

-- 插入执行记录示例数据
INSERT INTO `api_run_record` (`api_code`, `status`, `execution_count`, `execution_time`, `response_result`, `is_success`) VALUES
('abc123', 'success', 1, '2026-03-23 10:00:00', '{"code": 200, "message": "success"}', 1),
('def456', 'failed', 2, '2026-03-23 10:05:00', '{"code": 401, "message": "unauthorized"}', 0);
