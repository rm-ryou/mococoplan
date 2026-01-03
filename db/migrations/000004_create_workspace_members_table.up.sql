CREATE TABLE IF NOT EXISTS `workspace_members` (
  `workspace_id` INT NOT NULL,
  `user_id` INT NOT NULL,
  `role` ENUM('owner', 'admin', 'member') NOT NULL DEFAULT 'member',
  `joined_at` DATETIME NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  PRIMARY KEY (`workspace_id`, `user_id`),
  CONSTRAINT fk_workspace_members_workspace FOREIGN KEY (`workspace_id`) REFERENCES workspaces(`id`) ON DELETE CASCADE,
  CONSTRAINT fk_workspace_members_user FOREIGN KEY (`user_id`) REFERENCES users(`id`) ON DELETE CASCADE,
  INDEX idx_workspace_members_user (`user_id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4
