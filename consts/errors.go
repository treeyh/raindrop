package consts

import "errors"

var (
	ErrMsgDatabaseInitFail = errors.New("database initialization failed")

	ErrMsgDatabaseGetNowTimeFail = errors.New("get now time database fail")

	ErrMsgDatabaseCreateTableFail = errors.New("create table fail")

	ErrMsgDatabaseInitWorkersFail = errors.New("initialization workers fail")

	ErrMsgDatabaseServerTimeInterval = errors.New("Server and database time gap exceeds threshold")

	ErrMsgWorkersNotAvailable = errors.New("Workers not available")

	ErrMsgStartTimeStampError = errors.New("StartTimeStamp is greater than the current time")
)
