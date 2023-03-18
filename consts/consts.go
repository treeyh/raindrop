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
	DbTypeMySql = "mysql"

	DbTypePostgreSQL = "postgresql"

	DbMaxOpenConns = 3
	DbMaxIdleConns = 2
)

const (
	ErrMsgDatabaseInitFail = "Database initialization failed."
)
