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
	IdModeSnowflake     = "snowflake"
	IdModeNumberSection = "numbersection"

	IdBitLength       = 63
	TimeBackBitLength = 1
)

const (
	DbTypeMySql = "mysql"

	DbTypePostgreSQL = "postgresql"

	DbMaxOpenConns = 3
	DbMaxIdleConns = 2
)

const (
	// DatabaseTimeInterval 服务器与DB时间允许间隔，秒
	DatabaseTimeInterval = 30

	// HeartbeatTimeInterval 数据库心跳时间间隔，秒
	HeartbeatTimeInterval = 30
)
