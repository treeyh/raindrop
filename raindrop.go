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
	"github.com/treeyh/raindrop/worker"
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
	log.Debug(ctx, "check config over.")
	err = initDb(ctx, conf)
	if err != nil {
		log.Fatal(ctx, "init db fail", err)
	}

	log.Debug(ctx, "init db over.")
	err = initRaindrop(ctx, conf)
	if err != nil {
		log.Fatal(ctx, "init raindrop fail", err)
	}
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
func initDb(ctx context.Context, conf config.RainDropConfig) error {
	var err error
	if consts.DbTypeMySql == conf.DbConfig.DbType {
		err = db.InitMySqlDb(ctx, conf.DbConfig, log)
	} else {
		log.Fatal(ctx, "raindrop not support ["+conf.DbConfig.DbType+"] db type.")
		err = errors.New(consts.ErrMsgDatabaseInitFail)
	}
	if err != nil {
		log.Error(ctx, err.Error(), err)
	}

	log.Debug(ctx, "raindrop database initialization completed.")
	return err
}

// initRaindrop 初始化雨滴
func initRaindrop(ctx context.Context, conf config.RainDropConfig) error {
	err := checkDbTimeInterval(ctx)
	if err != nil {
		return err
	}
	err = db.InitTableWorkers(ctx, conf.ServiceMinWorkId, conf.ServiceMaxWorkId)
	if err != nil {
		log.Error(ctx, err.Error(), err)
		return err
	}
	err = worker.Init(ctx, conf)
	if err != nil {
		log.Error(ctx, err.Error(), err)
		return err
	}
	return nil
}

// checkDbTimeInterval 校验服务器时间和db时间间隔
func checkDbTimeInterval(ctx context.Context) error {
	now := time.Now()
	dbNow, err := db.Db.GetNowTime(ctx)

	if err != nil {
		log.Error(ctx, "get database now time fail", err)
		return err
	}

	if now.Unix() > (dbNow.Unix()+consts.DatabaseTimeInterval) || now.Unix() < (dbNow.Unix()-consts.DatabaseTimeInterval) {
		log.Error(ctx, fmt.Sprintf(consts.ErrMsgDatabaseServerTimeInterval, strconv.Itoa(consts.DatabaseTimeInterval)))
		return errors.New(consts.ErrMsgDatabaseServerTimeInterval)
	}
	return nil
}
