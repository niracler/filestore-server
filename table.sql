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