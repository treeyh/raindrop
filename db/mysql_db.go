package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/treeyh/raindrop/consts"
	"github.com/treeyh/raindrop/model"
	"strconv"
	"strings"
	"time"
)

const (
	mysqlTableName = "soc_id_generator_worker"

	mysqlPreSelectSql = "SELECT `id`, `code`, `time_unit`, `heartbeat_time`, `create_time`, `update_time`, `version`, `del_flag` FROM " + mysqlTableName + " WHERE `del_flag` = 2 "
)

type MySqlDb struct {
}

// GetNowTime 获取数据库当前时间
func (m *MySqlDb) GetNowTime(ctx context.Context) (time.Time, error) {
	var now time.Time
	err := _dbConn.QueryRowContext(ctx, "SELECT NOW() as now;").Scan(&now)

	if err != nil {
		log.Error(ctx, consts.ErrMsgDatabaseGetNowTimeFail, err)
		return time.Now(), err
	}
	return now, err
}

func (m *MySqlDb) getDatabaseName(ctx context.Context) string {
	var dbName string
	_dbConn.QueryRowContext(ctx, "SELECT DATABASE();").Scan(&dbName)
	return dbName
}

// ExistTable 表是否存在
func (m *MySqlDb) ExistTable(ctx context.Context) (bool, error) {
	dbName := m.getDatabaseName(ctx)

	var count int
	err := _dbConn.QueryRowContext(ctx, "SELECT count(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ? AND table_type = ?", dbName, mysqlTableName, "BASE TABLE").Scan(&count)

	if err != nil {
		log.Error(ctx, err.Error(), err)
		return false, err
	}

	return count == 1, nil
}

// CreateTable 创建依赖表
func (m *MySqlDb) CreateTable(ctx context.Context) error {

	sql := "CREATE TABLE `" + mysqlTableName + "` (\n" +
		"\t`id` bigint NOT NULL,\n" +
		"\t`code` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',\n" +
		"\t`time_unit` tinyint NOT NULL DEFAULT '2',\n" +
		"\t`heartbeat_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
		"\t`create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
		"\t`update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n" +
		"\t`version` bigint NOT NULL DEFAULT '1',\n" +
		"\t`del_flag` tinyint NOT NULL DEFAULT '2',\n" +
		"\tPRIMARY KEY (`id`),\n" +
		"\tKEY `idx_soc_id_generator_worker_heartbeat_time` (`heartbeat_time`),\n" +
		"\tKEY `idx_soc_id_generator_worker_code` (`code`)\n" +
		"\t) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;"

	_, err := _dbConn.ExecContext(ctx, sql)
	if err != nil {
		log.Error(ctx, consts.ErrMsgDatabaseCreateTableFail, err)
	}
	return err
}

// InitWorkers 初始化数据
func (m *MySqlDb) InitWorkers(ctx context.Context, beginId int64, endId int64) error {
	if beginId > endId {
		err := errors.New("endId must be greater than beginId")
		log.Error(ctx, err.Error(), err)
		return err
	}
	values := make([]string, 0)

	for i := beginId; i <= endId; i++ {
		values = append(values, "("+strconv.FormatInt(i, 10)+", '2023-01-01 00:00:00')")
	}

	sql := "INSERT INTO " + mysqlTableName + "(`id`, `heartbeat_time`) VALUES " + strings.Join(values, ",") + ";"

	_, err := _dbConn.ExecContext(ctx, sql)
	if err != nil {
		log.Error(ctx, consts.ErrMsgDatabaseInitWorkersFail, err)
	}
	return err
}

// GetBeforeWorker 找到该节点之前的worker
func (m *MySqlDb) GetBeforeWorker(ctx context.Context, code string, timeUnit int) (*model.IdGeneratorWorker, error) {
	var worker model.IdGeneratorWorker
	s := mysqlPreSelectSql + "AND `code` = ? AND `time_unit` = ? ORDER BY `id` asc LIMIT 0,1 "
	err := _dbConn.QueryRowContext(ctx, s, code, timeUnit).Scan(&worker.Id, &worker.Code,
		&worker.TimeUnit, &worker.HeartbeatTime, &worker.CreateTime, &worker.UpdateTime, &worker.Version, &worker.DelFlag)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Error(ctx, "find before worker fail", err)
		return nil, err
	}

	return &worker, nil
}

// QueryFreeWorkers 获取空闲的worker列表
func (m *MySqlDb) QueryFreeWorkers(ctx context.Context, heartbeatTime time.Time) ([]model.IdGeneratorWorker, error) {
	workers := make([]model.IdGeneratorWorker, 0)
	s := mysqlPreSelectSql + "AND `heartbeat_time` < ? "
	rows, err := _dbConn.QueryContext(ctx, s, heartbeatTime)
	if err != nil {
		log.Error(ctx, "query workers fail", err)
		return nil, err
	}
	for rows.Next() {
		var worker model.IdGeneratorWorker
		e := rows.Scan(&worker.Id, &worker.Code, &worker.TimeUnit, &worker.HeartbeatTime, &worker.CreateTime, &worker.UpdateTime, &worker.Version, &worker.DelFlag)
		if e != nil {
			log.Error(ctx, "query workers fail", e)
			return nil, e
		}
		workers = append(workers, worker)
	}
	rows.Close()

	return workers, nil
}

// ActivateWorker 激活启用worker
func (m *MySqlDb) ActivateWorker(ctx context.Context, id int64, code string, timeUnit int, version int64) (*model.IdGeneratorWorker, error) {
	sql := "UPDATE `" + mysqlTableName + "` SET `code` = ?, `time_unit` = ?, `version` = `version` + 1, `heartbeat_time` = ? WHERE `id` = ? AND `version` = ? "

	result, err := _dbConn.ExecContext(ctx, sql, code, timeUnit, time.Now(), id, version)
	if err != nil {
		log.Error(ctx, "heartbeat worker fail!!!", err)
		return nil, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Error(ctx, "heartbeat worker fail!!!", err)
		return nil, err
	}
	if count != 1 {
		log.Error(ctx, "heartbeat worker fail!!! count: "+strconv.FormatInt(count, 10))
		return nil, err
	}

	worker, err := m.GetWorkerById(ctx, id)
	if err != nil {
		log.Error(ctx, err.Error(), err)
		return &model.IdGeneratorWorker{
			Id:            id,
			Code:          code,
			TimeUnit:      consts.TimeUnit(timeUnit),
			HeartbeatTime: time.Now(),
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
			Version:       version + 1,
			DelFlag:       2,
		}, err
	}

	return worker, nil
}

// HeartbeatWorker 心跳
func (m *MySqlDb) HeartbeatWorker(ctx context.Context, worker *model.IdGeneratorWorker) (*model.IdGeneratorWorker, error) {
	sql := "UPDATE `" + mysqlTableName + "` SET `version` = `version` + 1, `heartbeat_time` = ? WHERE `id` = ? AND `version` = ? "

	result, err := _dbConn.ExecContext(ctx, sql, time.Now(), worker.Id, worker.Version)
	if err != nil {
		log.Error(ctx, "heartbeat worker fail!!!", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Error(ctx, "heartbeat worker fail!!!", err)
	}
	if count != 1 {
		log.Error(ctx, "heartbeat worker fail!!! result: "+strconv.FormatInt(count, 10))
	}

	w, _ := m.GetWorkerById(ctx, worker.Id)

	if w != nil {
		return w, nil
	}
	worker.Version += 1
	return worker, nil
}

// GetWorkerById 根据id获取worker
func (m *MySqlDb) GetWorkerById(ctx context.Context, id int64) (*model.IdGeneratorWorker, error) {
	s := mysqlPreSelectSql + " AND `id` = ? "
	var worker model.IdGeneratorWorker

	err := _dbConn.QueryRowContext(ctx, s, id).Scan(&worker.Id, &worker.Code, &worker.TimeUnit, &worker.HeartbeatTime,
		&worker.CreateTime, &worker.UpdateTime, &worker.Version, &worker.DelFlag)

	if err != nil {
		log.Error(ctx, "get worker by id fail. id: "+strconv.FormatInt(id, 10), err)
		return nil, err
	}

	return &worker, nil
}
