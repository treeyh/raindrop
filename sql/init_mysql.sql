CREATE TABLE `soc_id_generator_worker` (
    `id` bigint NOT NULL,
    `code` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `time_unit` tinyint NOT NULL DEFAULT '2',
    `heartbeat_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `version` bigint NOT NULL DEFAULT '1',
    `del_flag` tinyint NOT NULL DEFAULT '2',
    PRIMARY KEY (`id`),
    KEY `idx_soc_id_generator_worker_heartbeat_time` (`heartbeat_time`),
    KEY `idx_soc_id_generator_worker_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='id生成节点';