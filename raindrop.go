package raindrop

import (
	"context"
	"errors"
	"fmt"
	"github.com/treeyh/raindrop/config"
	"github.com/treeyh/raindrop/consts"
	"github.com/treeyh/raindrop/db"
	"github.com/treeyh/raindrop/logger"
	"github.com/treeyh/raindrop/utils"
	"strconv"
	"time"
)

var (
	log logger.ILogger
)

// Init 初始化
func Init(ctx context.Context, conf config.RainDropConfig) {

	initLogger(ctx, conf)
	log.Info(ctx, "raindrop init. config: "+utils.ToJsonIgnoreError(conf))

	err := config.CheckConfig(ctx, conf)
	if err != nil {
		log.Fatal(ctx, "config check fail", err)
	}

	initDb(ctx, conf)

	initRaindrop(ctx, conf)
}

// initLogger 初始化日志
func initLogger(ctx context.Context, conf config.RainDropConfig) {
	if conf.Logger != nil {
		log = conf.Logger
	} else {
		log = logger.New(&logger.DefaultWriter{}, logger.Info, true)
	}
}

// initDb 初始化数据库
func initDb(ctx context.Context, conf config.RainDropConfig) {
	var err error
	if consts.DbTypeMySql == conf.DbConfig.DbType {
		err = db.InitMySqlDb(ctx, conf.DbConfig, log)
	} else {
		log.Fatal(ctx, "raindrop not support ["+conf.DbConfig.DbType+"] db type.")
		err = errors.New(consts.ErrMsgDatabaseInitFail)
	}
	if err != nil {
		log.Fatal(ctx, err.Error())
	}

	log.Debug(ctx, "raindrop database initialization completed.")
}

// initRaindrop 初始化雨滴
func initRaindrop(ctx context.Context, conf config.RainDropConfig) {
	checkDbTimeInterval(ctx)
	err := db.InitTableWorkers(ctx, conf.ServiceMinWorkId, conf.ServiceMaxWorkId)
	if err != nil {
		log.Fatal(ctx, err.Error(), err)
	}
	activateWorker(ctx, conf)
}

// checkDbTimeInterval 校验服务器时间和db时间间隔
func checkDbTimeInterval(ctx context.Context) {
	now := time.Now()
	dbNow, err := db.Db.GetNowTime(ctx)

	if err != nil {
		log.Fatal(ctx, "get database now time fail", err)
	}

	if now.Unix() > (dbNow.Unix()+consts.DatabaseTimeInterval) || now.Unix() < (dbNow.Unix()-consts.DatabaseTimeInterval) {
		log.Fatal(ctx, fmt.Sprintf(consts.ErrMsgDatabaseServerTimeInterval, strconv.Itoa(consts.DatabaseTimeInterval)))
	}
}

// activateWorker 激活worker
func activateWorker(ctx context.Context, conf config.RainDropConfig) {
	heartbeatMaxTime := time.Now().Add(time.Duration(consts.HeartbeatTimeInterval*-3) * time.Second)

	if conf.TimeUnit == consts.TimeUnitHour {
		heartbeatMaxTime = time.Now().Add(time.Duration(-1) * time.Hour)
	} else if conf.TimeUnit == consts.TimeUnitDay {
		heartbeatMaxTime = time.Now().Add(time.Duration(-24) * time.Hour)
	}

	workers, err := db.Db.QueryFreeWorkers(ctx, heartbeatMaxTime)

	if err != nil {
		log.Fatal(ctx, err.Error(), err)
	}
	if len(workers) <= 0 {
		log.Fatal(ctx, consts.ErrMsgNoWorkerAvailable)
	}

}
