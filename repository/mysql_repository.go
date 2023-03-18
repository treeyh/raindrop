package repository

import (
	"context"
	"database/sql"
	"github.com/treeyh/raindrop/consts"
	"github.com/treeyh/raindrop/logger"
	"github.com/treeyh/raindrop/model"
)

var _dbConn *sql.DB
var _logger *logger.ILogger

// InitMySqlRepository 初始化mysql
func InitMySqlRepository(ctx context.Context, dbConfig *model.RainDropDbConfig, log *logger.ILogger) {

	if dbConfig.DbType != consts.DbTypeMySql {

	}

	//var err error
	//_dbConn, err = sql.Open(dbConfig.DbType, dbConfig.DbUrl)
	//if err != nil {
	//	log.Println("open db fail:", err)
	//}
	//
	//_dbConn.SetMaxOpenConns(consts.DbMaxOpenConns)
	//_dbConn.SetMaxIdleConns(consts.DbMaxIdleConns)
	//
	//err = _dbConn.Ping()
	//if err != nil {
	//	log.Fatalln("ping db fail:", err)
	//}
}
