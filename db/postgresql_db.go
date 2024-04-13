package db

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/treeyh/raindrop/consts"
	"github.com/treeyh/raindrop/model"
)

type PostgreSqlDb struct {
	tableName      string
	preSelectSql   string
	createTableSql string
}

func (m *PostgreSqlDb) InitSql(tableName string) {
	m.tableName = tableName
	m.preSelectSql = "SELECT \"id\", \"code\", \"time_unit\", \"heartbeat_time\", \"create_time\", \"update_time\", \"version\", \"del_flag\" FROM \"" + m.tableName + "\" WHERE \"del_flag\" = 2 "
	m.createTableSql = "CREATE TABLE IF NOT EXISTS \"" + m.tableName + "\" (\n" +
		"\t\"id\"                   bigint               not null,\n" +
		"\t\"code\"                 varchar(128)         not null default '',\n" +
		"\t\"lang_code\"            varchar(128)         not null default '',\n" +
		"\t\"time_unit\"            smallint             not null default '2',\n" +
		"\t\"heartbeat_time\"       TIMESTAMP WITH TIME ZONE not null default CURRENT_TIMESTAMP,\n" +
		"\t\"create_time\"          TIMESTAMP WITH TIME ZONE not null default CURRENT_TIMESTAMP,\n" +
		"\t\"update_time\"          TIMESTAMP WITH TIME ZONE not null default CURRENT_TIMESTAMP,\n" +
		"\t\"version\"              bigint               not null default '1',\n" +
		"\t\"del_flag\"             bool                 not null default false,\n" +
		"\tconstraint \"PK_SOC_RAINDROP_WORKER\" primary key (\"id\")\n" +
		"\t);\n" +
		"\tCREATE INDEX \"idx_soc_raindrop_worker_hb_time\" on \"soc_raindrop_worker\" (\n" +
		"\t\"heartbeat_time\"\n" +
		"\t);\n" +
		"\tCREATE INDEX \"idx_soc_raindrop_worker_code\" on \"soc_raindrop_worker\" (\n" +
		"\t\"code\"\n" +
		"\t);\n"
}

// GetNowTime 获取数据库当前时间
func (m *PostgreSqlDb) GetNowTime(ctx context.Context) (time.Time, error) {
	var now time.Time
	err := _db.QueryRowContext(ctx, "SELECT NOW() as now;").Scan(&now)

	if err != nil {
		log.Error(ctx, consts.ErrMsgDatabaseGetNowTimeFail.Error(), err)
		return time.Now(), err
	}
	return now, err
}

func (m *PostgreSqlDb) getDatabaseName(ctx context.Context) string {
	var dbName string
	_db.QueryRowContext(ctx, "SELECT current_database();").Scan(&dbName)
	return dbName
}

// ExistTable 表是否存在
func (m *PostgreSqlDb) ExistTable(ctx context.Context) (bool, error) {
	// dbName := m.getDatabaseName(ctx)

	var count int
	err := _db.QueryRowContext(ctx, "select count(*) from \"pg_tables\" where \"tablename\" = ?;", m.tableName).Scan(&count)

	if err != nil {
		log.Error(ctx, err.Error(), err)
		return false, err
	}

	return count == 1, nil
}

// InitTableWorkers 初始化数据
func (m *PostgreSqlDb) InitTableWorkers(ctx context.Context, beginId int64, endId int64) error {
	if beginId > endId {
		err := errors.New("endId must be greater than beginId")
		log.Error(ctx, err.Error(), err)
		return err
	}

	values := make([]string, 0)

	for i := beginId; i <= endId; i++ {
		values = append(values, "("+strconv.FormatInt(i, 10)+", '2023-01-01 00:00:00')")
	}

	rowsSql := "INSERT INTO \"" + m.tableName + "\"(\"id\", \"heartbeat_time\") VALUES " + strings.Join(values, ",") + ";"

	tx, err := _db.Begin()
	if err != nil {
		log.Error(ctx, err.Error(), err)
		if tx != nil {
			tx.Rollback()
		}
		return err
	}

	_, err = tx.ExecContext(ctx, m.createTableSql)
	if err != nil {
		log.Error(ctx, err.Error(), err)
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, rowsSql)
	if err != nil {
		log.Error(ctx, err.Error(), err)
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

// GetBeforeWorker 找到该节点之前的worker
func (m *PostgreSqlDb) GetBeforeWorker(ctx context.Context, code string) (*model.RaindropWorker, error) {
	var worker model.RaindropWorker
	s := m.preSelectSql + "AND \"code\" = ? ORDER BY \"id\" asc LIMIT 1 ;"
	err := _db.QueryRowContext(ctx, s, code).Scan(&worker.Id, &worker.Code,
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
func (m *PostgreSqlDb) QueryFreeWorkers(ctx context.Context, heartbeatTime time.Time) ([]model.RaindropWorker, error) {
	workers := make([]model.RaindropWorker, 0)
	s := m.preSelectSql + " AND \"heartbeat_time\" < ? ORDER BY \"heartbeat_time\" ASC ;"
	rows, err := _db.QueryContext(ctx, s, heartbeatTime)
	if err != nil {
		log.Error(ctx, "query workers fail", err)
		return nil, err
	}
	for rows.Next() {
		var worker model.RaindropWorker
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
func (m *PostgreSqlDb) ActivateWorker(ctx context.Context, id int64, code string, timeUnit int, version int64) (*model.RaindropWorker, error) {
	sql := "UPDATE \"" + m.tableName + "\" SET \"code\" = ?, \"time_unit\" = ?, \"version\" = \"version\" + 1, \"heartbeat_time\" = ? WHERE \"id\" = ? AND \"version\" = ? "

	result, err := _db.ExecContext(ctx, sql, code, timeUnit, time.Now(), id, version)
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
		return &model.RaindropWorker{
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
func (m *PostgreSqlDb) HeartbeatWorker(ctx context.Context, worker *model.RaindropWorker) (*model.RaindropWorker, error) {
	sql := "UPDATE \"" + m.tableName + "\" SET \"version\" = \"version\" + 1, \"heartbeat_time\" = ? WHERE \"id\" = ? AND \"version\" = ? "

	result, err := _db.ExecContext(ctx, sql, time.Now(), worker.Id, worker.Version)
	if err != nil {
		log.Error(ctx, "heartbeat worker fail!!!", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Error(ctx, "heartbeat worker fail!!!", err)
	}
	if count != 1 {
		log.Error(ctx, "heartbeat worker fail!!! id:"+strconv.FormatInt(worker.Id, 10)+" result: "+strconv.FormatInt(count, 10))
	}

	w, _ := m.GetWorkerById(ctx, worker.Id)

	if w != nil {
		return w, nil
	}
	worker.Version += 1
	return worker, nil
}

// GetWorkerById 根据id获取worker
func (m *PostgreSqlDb) GetWorkerById(ctx context.Context, id int64) (*model.RaindropWorker, error) {
	s := m.preSelectSql + " AND \"id\" = ? "
	var worker model.RaindropWorker

	err := _db.QueryRowContext(ctx, s, id).Scan(&worker.Id, &worker.Code, &worker.TimeUnit, &worker.HeartbeatTime,
		&worker.CreateTime, &worker.UpdateTime, &worker.Version, &worker.DelFlag)

	if err != nil {
		log.Error(ctx, "get worker by id fail. id: "+strconv.FormatInt(id, 10), err)
		return nil, err
	}

	return &worker, nil
}
