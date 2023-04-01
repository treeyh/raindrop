package consts

import "errors"

var (
	// ErrMsgDatabaseInitFail 数据库初始化失败
	ErrMsgDatabaseInitFail = errors.New("Database initialization failed")

	// ErrMsgDatabaseGetNowTimeFail 获取数据库当前时间失败
	ErrMsgDatabaseGetNowTimeFail = errors.New("Failed to get the current time of the database")

	// ErrMsgDatabaseInitTableFail 创建表失败
	ErrMsgDatabaseInitTableFail = errors.New("Initialization table fail")

	// ErrMsgDatabaseInitWorkersFail 初始化workers失败
	ErrMsgDatabaseInitWorkersFail = errors.New("Initialization workers fail")

	// ErrMsgDatabaseServerTimeInterval 服务器和数据库时间差异过大
	ErrMsgDatabaseServerTimeInterval = errors.New("Server and database time gap exceeds threshold")

	// ErrMsgWorkersNotAvailable 没有有效的worker分配
	ErrMsgWorkersNotAvailable = errors.New("Workers not available")

	// ErrMsgStartTimeStampError 开始时间戳大于当前时间
	ErrMsgStartTimeStampError = errors.New("StartTimeStamp is greater than the current time")

	// ErrMsgServerClockBackwardsError 服务器时钟向后移动，无法生成id。
	ErrMsgServerClockBackwardsError = errors.New("Server clock moved backwards and id could not be generated")

	// ErrMsgIdSeqReachesMaxValueError 当前时刻id生成序列达到最大值
	ErrMsgIdSeqReachesMaxValueError = errors.New("The current moment id generation sequence reaches its maximum value")

	// ErrMsgGetCodeLockFail 获取编号锁失败
	ErrMsgGetCodeLockFail = errors.New("Failed to get code lock")
)
