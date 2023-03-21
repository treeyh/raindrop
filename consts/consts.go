package consts

type TimeUnit int

const (
	TimeUnitMillisecond TimeUnit = iota + 1
	TimeUnitSecond
	TimeUnitMinute
	TimeUnitHour
	TimeUnitDay
)

const (
	ProjectName = "raindrop"

	SkipHeartbeat = "skip_heartbeat"
)

const (
	DbTypeMySql = "mysql"

	DbTypePostgreSQL = "postgresql"

	DbMaxOpenConns = 3
	DbMaxIdleConns = 2
)

const (
	ErrMsgDatabaseInitFail = "database initialization failed"

	ErrMsgDatabaseGetNowTimeFail = "get now time database fail"

	ErrMsgDatabaseCreateTableFail = "create table fail"

	ErrMsgDatabaseInitWorkersFail = "initialization workers fail"

	ErrMsgDatabaseServerTimeInterval = "Server and database time differences of more than %s seconds"

	ErrMsgWorkersNotAvailable = "Workers not available"

	ErrMsgStartTimeStampError = "StartTimeStamp is greater than the current time"
)

const (
	// DatabaseTimeInterval 服务器与DB时间允许间隔，秒
	DatabaseTimeInterval = 30

	// HeartbeatTimeInterval 数据库心跳时间间隔，秒
	HeartbeatTimeInterval = 30
)
