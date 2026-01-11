CREATE TABLE IF NOT EXISTS `channel_members` (
  `channel_id` INT NOT NULL,
  `user_id` INT NOT NULL,
  `joined_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  PRIMARY KEY (`channel_id`, `user_id`),
  CONSTRAINT fk_channel_members_channel FOREIGN KEY (`channel_id`) REFERENCES channels(`id`) ON DELETE CASCADE,
  CONSTRAINT fk_channel_members_user FOREIGN KEY (`user_id`) REFERENCES users(`id`) ON DELETE CASCADE,
  INDEX idx_channel_members_user (`user_id`)
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4
