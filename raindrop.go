package raindrop

import (
	"context"
	"strconv"
	"time"

	"github.com/treeyh/raindrop/config"
	"github.com/treeyh/raindrop/consts"
	"github.com/treeyh/raindrop/db"
	"github.com/treeyh/raindrop/logger"
	"github.com/treeyh/raindrop/utils"
	"github.com/treeyh/raindrop/worker"
)

var (
	log logger.ILogger
)

// Init 初始化
func Init(ctx context.Context, conf config.RainDropConfig) {

	initLogger(ctx, conf)
	log.Info(ctx, "raindrop init. config: "+utils.ToJsonIgnoreError(conf))

	err := config.CheckConfig(ctx, &conf)
	if err != nil {
		log.Fatal(ctx, "config check fail: "+err.Error(), err)
	}
	log.Debug(ctx, "check config over.")
	err = initDb(ctx, conf)
	if err != nil {
		log.Fatal(ctx, "init db fail: "+err.Error(), err)
	}

	log.Debug(ctx, "init db over.")
	err = initRaindrop(ctx, conf)
	if err != nil {
		log.Fatal(ctx, "init raindrop fail: "+err.Error(), err)
	}
}

// NewId 获取新id
func NewId() (int64, error) {
	ctx := context.Background()
	return NewIdContext(ctx)
}

// NewIdContext 获取新id
func NewIdContext(ctx context.Context) (int64, error) {
	return worker.NewId(ctx)
}

// NewIdByCode 基于code获取新id
func NewIdByCode(code string) (int64, error) {
	ctx := context.Background()
	return NewIdContextByCode(ctx, code)
}

// NewIdContextByCode 基于code获取新id
func NewIdContextByCode(ctx context.Context, code string) (int64, error) {
	return worker.NewIdByCode(ctx, code)
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
	} else if consts.DbTypePostgreSQL == conf.DbConfig.DbType {
		err = db.InitPostgreSqlDb(ctx, conf.DbConfig, log)
	} else {
		err = db.InitPostgreSqlDb(ctx, conf.DbConfig, log)
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
		log.Error(ctx, "get database now time fail: "+err.Error(), err)
		return err
	}

	if now.Unix() > (dbNow.Unix()+consts.DatabaseTimeInterval) || now.Unix() < (dbNow.Unix()-consts.DatabaseTimeInterval) {
		log.Error(ctx, consts.ErrMsgDatabaseServerTimeInterval.Error()+". system now time:"+now.String()+"; system now unix: "+strconv.FormatInt(now.Unix(), 10)+" ; db now time:"+dbNow.String()+"; db now unix: "+strconv.FormatInt(dbNow.Unix(), 10))
		panic(consts.ErrMsgDatabaseServerTimeInterval.Error())
		// return consts.ErrMsgDatabaseServerTimeInterval
	}
	return nil
}
