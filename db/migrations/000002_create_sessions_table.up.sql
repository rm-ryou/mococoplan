CREATE TABLE IF NOT EXISTS `sessions` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `user_id` INT NOT NULL,
  `token` VARBINARY(32) NOT NULL,
  `expires_at` DATETIME NOT NULL,
  `ip_address` VARBINARY(16),
  `user_agent` VARCHAR(512),
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  PRIMARY KEY (`id`),
  UNIQUE KEY (`token`),
  CONSTRAINT fk_sessions_user FOREIGN KEY (`user_id`) REFERENCES users(`id`) ON DELETE CASCADE,
  INDEX idx_sessions_user_id (`user_id`),
  INDEX idx_sessions_expires_at (`expires_at`),
  INDEX idx_sessions_user_expires (`user_id`, `expires_at`)
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4
