CREATE TABLE `fileserver_file`
(
    `fid`       INT(11)       NOT NULL AUTO_INCREMENT,
    `file_sha1` CHAR(40)      NOT NULL DEFAULT '' COMMENT '文件hash',
    `file_name` VARCHAR(256)  NOT NULL DEFAULT '' COMMENT '文件名',
    `file_size` BIGINT(20)    NOT NULL DEFAULT '0' COMMENT '文件大小',
    `file_addr` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
    `created`   DATETIME               DEFAULT NOW() COMMENT '创建日期',
    `updated`   DATETIME               DEFAULT NOW() ON UPDATE CURRENT_TIMESTAMP() COMMENT '创建日期',
    `status`    INT(11)       NOT NULL DEFAULT '0' COMMENT '状态(可用/禁用/已删除等状态)',
    PRIMARY KEY (`fid`),
    UNIQUE KEY `idx_file_hash` (`file_sha1`),
    KEY `indx_status` (`status`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

CREATE TABLE `fileserver_user`
(
    `uid`             INT(11)      NOT NULL AUTO_INCREMENT,
    `user_name`       VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '用户名',
    `user_pwd`        VARCHAR(256) NOT NULL DEFAULT '' COMMENT '用户encode密码',
    `email`           VARCHAR(64)           DEFAULT '' COMMENT '邮箱',
    `phone`           VARCHAR(128)          DEFAULT '' COMMENT '手机号',
    `email_validated` TINYINT(1)            DEFAULT 0 COMMENT '邮箱是否已验证',
    `phone_validated` TINYINT(1)            DEFAULT 0 COMMENT '手机是否已验证',
    `signup_at`       DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册日期',
    `last_active`     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后活跃时间时间',
    `profile`         TEXT COMMENT '用户属性',
    `status`          INT(11)      NOT NULL DEFAULT '0' COMMENT '账户状态(启用/禁用/锁定/标记删除等)',
    PRIMARY KEY (`uid`),
    UNIQUE KEY `idx_user_name` (`user_name`),
    KEY `idx_status` (`status`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 5
  DEFAULT CHARSET = utf8mb4;

INSERT INTO fileserver_user (`user_name`, `user_pwd`) VALUES ('nieacler13', '159258');

SELECT user_pwd FROM fileserver_user WHERE user_name='niracler';