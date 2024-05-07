package tests

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/treeyh/raindrop/config"
	"github.com/treeyh/raindrop/consts"
	"github.com/treeyh/raindrop/logger"
	"time"
)

const (
	port = 8888

	workerLength = 4

	minWorkerId = 10

	maxWorkerId = 15
)

var (
	_mysqlDbConn *sql.DB

	_pgDbPool *pgxpool.Pool

	tableName = "soc_raindrop_worker"
)

func init() {
	ctx := getTestContext()

	dbConfig := getTestConfig()
	if consts.DbTypeMySql == dbConfig.DbType {
		initTestMySqlDb(ctx)
	} else {
		initTestPostgreSqlDb(ctx)
	}
}

func getTestSkipHeartbeatContext() context.Context {
	ctx := context.WithValue(context.Background(), consts.ProjectName, consts.SkipHeartbeat)
	return ctx
}

func getTestContext() context.Context {
	ctx := context.WithValue(context.Background(), consts.ProjectName, consts.ProjectName)
	return ctx
}

func getTestConfig() config.RainDropDbConfig {

	return config.RainDropDbConfig{
		DbType:    consts.DbTypePostgreSQL,
		DbUrl:     "postgres://proot:4pVmsxTuB_5ZlnSX@127.0.0.1:5432/soc_expense_tracker_utest_db",
		TableName: tableName,
	}

	//return config.RainDropDbConfig{
	//	DbType: "mysql",
	//	//DbUrl:     "root:mysqlpw@(172.25.100.40:3306)/test?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai",
	//	TableName: tableName,
	//	DbUrl:     "dev_account:9CrgLlsDN9QlitQFRNW9@(rm-uf6cl3tt9t814wv84.mysql.rds.aliyuncs.com)/test?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai",
	//}
}

func getTestStdoutLogger() logger.ILogger {
	d := logger.DefaultWriter{}
	return logger.New(&d, logger.Info, true)
}

func getTestStdoutDebugLogger() logger.ILogger {
	d := logger.DefaultWriter{}
	return logger.New(&d, logger.Debug, true)
}

func getTestSecondConfig() config.RainDropConfig {
	return config.RainDropConfig{
		IdMode:                  consts.IdModeSnowflake,
		DbConfig:                getTestConfig(),
		Logger:                  getTestStdoutLogger(),
		ServicePort:             port,
		TimeUnit:                consts.TimeUnitSecond,
		StartTimeStamp:          time.Date(2023, 3, 1, 0, 0, 0, 0, time.Local),
		TimeStampLength:         31,
		PriorityEqualCodeWorkId: false,
		WorkIdLength:            workerLength,
		ServiceMinWorkId:        minWorkerId,
		ServiceMaxWorkId:        maxWorkerId,
		TimeBackBitValue:        0,
		EndBitsLength:           0,
		EndBitsValue:            0,
	}
}

func getTestMinuteConfig() config.RainDropConfig {
	return config.RainDropConfig{
		IdMode:                  consts.IdModeSnowflake,
		DbConfig:                getTestConfig(),
		Logger:                  getTestStdoutLogger(),
		ServicePort:             port,
		TimeUnit:                consts.TimeUnitMinute,
		StartTimeStamp:          time.Date(2023, 3, 1, 0, 0, 0, 0, time.Local),
		TimeStampLength:         26,
		PriorityEqualCodeWorkId: true,
		WorkIdLength:            workerLength,
		ServiceMinWorkId:        minWorkerId,
		ServiceMaxWorkId:        maxWorkerId,
		TimeBackBitValue:        0,
		EndBitsLength:           0,
		EndBitsValue:            0,
	}
}

func getTestSimpleMillisecondConfig() config.RainDropConfig {
	return config.RainDropConfig{
		IdMode:                  consts.IdModeSnowflake,
		DbConfig:                getTestConfig(),
		Logger:                  getTestStdoutDebugLogger(),
		ServicePort:             port,
		TimeUnit:                consts.TimeUnitMillisecond,
		StartTimeStamp:          time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local),
		TimeStampLength:         41,
		PriorityEqualCodeWorkId: false,
		WorkIdLength:            workerLength,
		ServiceMinWorkId:        minWorkerId,
		ServiceMaxWorkId:        maxWorkerId,
		TimeBackBitValue:        0,
		EndBitsLength:           0,
		EndBitsValue:            0,
	}
}

func getTestMillisecondConfig() config.RainDropConfig {
	return config.RainDropConfig{
		IdMode:                  consts.IdModeSnowflake,
		DbConfig:                getTestConfig(),
		Logger:                  getTestStdoutLogger(),
		ServicePort:             port,
		TimeUnit:                consts.TimeUnitMillisecond,
		StartTimeStamp:          time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local),
		TimeStampLength:         44,
		PriorityEqualCodeWorkId: false,
		WorkIdLength:            workerLength,
		ServiceMinWorkId:        minWorkerId,
		ServiceMaxWorkId:        maxWorkerId,
		TimeBackBitValue:        0,
		EndBitsLength:           1,
		EndBitsValue:            0,
	}
}

// initTestMySqlDb 初始化MySql
func initTestMySqlDb(ctx context.Context) error {
	dbConfig := getTestConfig()
	var err error
	_mysqlDbConn, err = sql.Open(dbConfig.DbType, dbConfig.DbUrl)
	if err != nil {
		return err
	}

	_mysqlDbConn.SetMaxOpenConns(consts.DbMaxOpenConns)
	_mysqlDbConn.SetMaxIdleConns(consts.DbMaxIdleConns)

	err = _mysqlDbConn.Ping()
	if err != nil {
		return err
	}
	return nil
}

// initTestPostgreSqlDb 初始化PostgreSql
func initTestPostgreSqlDb(ctx context.Context) error {
	dbConfig := getTestConfig()
	var err error
	dbUrl := dbConfig.DbUrl
	_pgDbPool, err = pgxpool.New(ctx, dbUrl)
	if err != nil {
		return err
	}

	err = _pgDbPool.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}

// dropTestWorkerTable 删除表
func dropTestWorkerTable(ctx context.Context, dbType string) error {

	s := "DROP TABLE " + tableName + ";"

	if dbType == consts.DbTypeMySql {
		_, err := _mysqlDbConn.ExecContext(ctx, s)
		return err
	} else {
		_, err := _pgDbPool.Exec(ctx, s)
		return err
	}

}

// updateWorker 更新 Worker
func updateWorker(ctx context.Context, id int64, code string, timeUnit int, heartbeatTime time.Time) error {
	s := "UPDATE " + tableName + " SET code = ?, time_unit = ?, heartbeat_time = ? WHERE id = ? "
	_, err := _mysqlDbConn.ExecContext(ctx, s)
	return err
}
