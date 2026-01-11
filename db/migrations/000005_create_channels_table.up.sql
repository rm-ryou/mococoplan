CREATE TABLE IF NOT EXISTS `channels` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `workspace_id` INT NOT NULL,
  `name` VARCHAR(50) NOT NULL,
  `description` VARCHAR(255) NULL,
  `is_private` BOOLEAN NOT NULL DEFAULT FALSE,
  `created_by` INT NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  PRIMARY KEY (`id`),
  UNIQUE KEY (`workspace_id`, `name`),
  INDEX idx_channels_workspace (`workspace_id`),
  FOREIGN KEY (`workspace_id`) REFERENCES workspaces(`id`) ON DELETE CASCADE
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4
