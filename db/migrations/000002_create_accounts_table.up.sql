CREATE TABLE IF NOT EXISTS `accounts` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `user_id` INT NOT NULL,
  `account_id` VARCHAR(255) NOT NULL,
  `provider_id` VARCHAR(50) NOT NULL,
  `access_token` TEXT,
  `refresh_token` TEXT,
  `access_token_expires_at` DATETIME,
  `refresh_token_expires_at` DATETIME,
  `scope` TEXT,
  `id_token` TEXT,
  `password` TEXT,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_provider_id_account_id` (`provider_id`, `account_id`),
  FOREIGN KEY (`user_id`) REFERENCES users(`id`) ON DELETE CASCADE
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4
