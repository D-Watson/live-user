use live_chat;
CREATE TABLE `users` (
     `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID，主键',
     `username` varchar(50) NOT NULL COMMENT '用户名，用于登录',
     `password_hash` varchar(255) NOT NULL COMMENT '加盐值加密后的密码',
     `email` varchar(100) NOT NULL COMMENT '电子邮箱',
     `phone` varchar(20) DEFAULT NULL COMMENT '手机号码',
     `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-正常',
     `avatar` varchar(255) DEFAULT NULL COMMENT '头像URL',
     `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
     `last_login_ip` varchar(50) DEFAULT NULL COMMENT '最后登录IP',
     `extra` json not null comment '扩展信息,包括设备信息',
     `created_time` bigint(20) NOT NULL DEFAULT 0 COMMENT '创建时间',
     `updated_time` bigint(20) NOT NULL DEFAULT 0 COMMENT '更新时间',
     PRIMARY KEY (`id`),
     UNIQUE KEY `idx_username` (`username`),
     UNIQUE KEY `idx_email` (`email`),
     KEY `idx_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户基本信息表';