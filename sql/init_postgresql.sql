CREATE TABLE IF NOT EXISTS "soc_raindrop_worker" (
  "id"                   bigint               not null,
  "code"                 varchar(128)         not null default '',
  "lang_code"            varchar(128)         not null default '',
  "time_unit"            smallint             not null default '2',
  "heartbeat_time"       TIMESTAMP WITH TIME ZONE not null default CURRENT_TIMESTAMP,
  "create_time"          TIMESTAMP WITH TIME ZONE not null default CURRENT_TIMESTAMP,
  "update_time"          TIMESTAMP WITH TIME ZONE not null default CURRENT_TIMESTAMP,
  "version"              bigint               not null default '1',
  "del_flag"             bool                 not null default false,
constraint "PK_SOC_RAINDROP_WORKER" primary key ("id")
);
CREATE INDEX "idx_soc_raindrop_worker_hb_time" on "soc_raindrop_worker" (
  "heartbeat_time"
);
CREATE INDEX "idx_soc_raindrop_worker_code" on "soc_raindrop_worker" (
  "code"
);

INSERT INTO "soc_raindrop_worker"("id", "heartbeat_time")
VALUES (1, '2023-01-01 00:00:00'),
      (2, '2023-01-01 00:00:00'),
      (3, '2023-01-01 00:00:00'),
      (4, '2023-01-01 00:00:00'),
      (5, '2023-01-01 00:00:00'),
      (6, '2023-01-01 00:00:00');