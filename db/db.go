package db

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/treeyh/raindrop/config"
	"github.com/treeyh/raindrop/consts"
	"github.com/treeyh/raindrop/logger"
	"github.com/treeyh/raindrop/model"
	"time"
)

var (
	//Db      *IDb
	Db  *MySqlDb
	_db *sql.DB

	log logger.ILogger

	tableName = "soc_raindrop_worker"
)

type IDb interface {
	// GetNowTime 获取数据库当前时间
	GetNowTime(ctx context.Context) (time.Time, error)

	// ExistTable 表是否存在
	ExistTable(ctx context.Context) (bool, error)

	// InitTableWorkers 初始化workers
	InitTableWorkers(ctx context.Context, beginId int64, endId int64) error

	// GetBeforeWorker 找到该节点之前的worker
	GetBeforeWorker(ctx context.Context, code string) (*model.RaindropWorker, error)

	// QueryFreeWorkers 查询空闲的workers
	QueryFreeWorkers(ctx context.Context, heartbeatTime time.Time) ([]model.RaindropWorker, error)

	// ActivateWorker 激活启用worker
	ActivateWorker(ctx context.Context, id int64, code string, timeUnit int, version int64) (*model.RaindropWorker, error)

	// HeartbeatWorker 心跳
	HeartbeatWorker(ctx context.Context, worker *model.RaindropWorker) (*model.RaindropWorker, error)

	// GetWorkerById 根据id获取worker
	GetWorkerById(ctx context.Context, id int64) (*model.RaindropWorker, error)
}

// InitMySqlDb 初始化MySql
func InitMySqlDb(ctx context.Context, dbConfig config.RainDropDbConfig, l logger.ILogger) error {
	log = l

	if dbConfig.TableName != "" {
		tableName = dbConfig.TableName
	}

	var err error
	_db, err = sql.Open(dbConfig.DbType, dbConfig.DbUrl)
	if err != nil {
		log.Error(ctx, consts.ErrMsgDatabaseInitFail.Error(), err)
		return err
	}

	_db.SetMaxOpenConns(consts.DbMaxOpenConns)
	_db.SetMaxIdleConns(consts.DbMaxIdleConns)

	err = _db.Ping()
	if err != nil {
		log.Error(ctx, consts.ErrMsgDatabaseInitFail.Error(), err)
		return err
	}

	Db = &MySqlDb{}
	Db.InitSql(tableName)
	return nil
}

func InitTableWorkers(ctx context.Context, beginId int64, endId int64) error {
	exist, err := Db.ExistTable(ctx)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	err = Db.InitTableWorkers(ctx, beginId, endId)
	if err != nil {
		return err
	}
	return err
}
