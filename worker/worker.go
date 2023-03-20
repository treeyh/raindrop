package worker

import (
	"context"
	"errors"
	"github.com/treeyh/raindrop/config"
	"github.com/treeyh/raindrop/consts"
	"github.com/treeyh/raindrop/db"
	"github.com/treeyh/raindrop/logger"
	"github.com/treeyh/raindrop/model"
	"github.com/treeyh/raindrop/utils"
	"strconv"
	"time"
)

var (
	workerCode string
	timeUnit   consts.TimeUnit
	log        logger.ILogger
	worker     *model.IdGeneratorWorker
)

// Init 初始化worker
func Init(ctx context.Context, conf config.RainDropConfig) error {
	log = conf.Logger
	ip, err := utils.GetLocalIP()
	if err != nil {
		log.Error(ctx, "get local ip fail", err)
		return err
	}
	workerCode = ip + ":" + strconv.Itoa(conf.ServicePort)
	timeUnit = conf.TimeUnit

	worker, err = activateWorker(ctx)
	if err != nil {
		return err
	}

	return nil
}

// activateWorker 激活worker
func activateWorker(ctx context.Context) (*model.IdGeneratorWorker, error) {
	if timeUnit == consts.TimeUnitMillisecond || timeUnit == consts.TimeUnitSecond {
		w, err := db.Db.GetBeforeWorker(ctx, workerCode, int(timeUnit),
			time.Now().Add(time.Duration(consts.HeartbeatTimeInterval*-3)*time.Second))

		if err != nil {
			return nil, err
		}
		if w != nil {
			w, err = db.Db.ActivateWorker(ctx, w.Id, workerCode, int(timeUnit), w.Version)
			if w != nil {
				return w, nil
			}
		}
	}

	heartbeatMaxTime := time.Now().Add(time.Duration(consts.HeartbeatTimeInterval*-4) * time.Second)

	if timeUnit == consts.TimeUnitMinute {
		heartbeatMaxTime = heartbeatMaxTime.Add(time.Duration(-1) * time.Minute)
	} else if timeUnit == consts.TimeUnitHour {
		heartbeatMaxTime = heartbeatMaxTime.Add(time.Duration(-1) * time.Hour)
	} else if timeUnit == consts.TimeUnitDay {
		heartbeatMaxTime = heartbeatMaxTime.Add(time.Duration(-24) * time.Hour)
	}

	workers, err := db.Db.QueryFreeWorkers(ctx, heartbeatMaxTime)

	if err != nil {
		log.Error(ctx, err.Error(), err)
		return nil, err
	}
	if len(workers) <= 0 {
		log.Error(ctx, consts.ErrMsgWorkersNotAvailable)
		return nil, errors.New(consts.ErrMsgWorkersNotAvailable)
	}

}
