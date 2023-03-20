CREATE TABLE `soc_id_generator_worker` (
   `id` bigint NOT NULL COMMENT 'id主键',
   `code` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '编号',
   `time_unit` tinyint NOT NULL DEFAULT '2' COMMENT '时间单位，1：毫秒，2：秒（默认），3：分钟，4：小时，5：天',
   `heartbeat_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后心跳时间',
   `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
   `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
   `version` bigint NOT NULL DEFAULT '1' COMMENT '乐观锁版本号',
   `del_flag` tinyint NOT NULL DEFAULT '2' COMMENT '是否删除，1删除，2未删除',
   PRIMARY KEY (`id`),
   KEY `idx_soc_id_generator_worker_heartbeat_time` (`heartbeat_time`),
   KEY `idx_soc_id_generator_worker_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='id生成节点';