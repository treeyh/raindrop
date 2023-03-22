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
	"sync"
	"time"
)

var (
	workerCode string
	timeUnit   consts.TimeUnit
	log        logger.ILogger
	worker     *model.IdGeneratorWorker

	// 开始计算时间戳，毫秒
	startTime int64
	// 当前时间流水，当前时刻毫秒 - startTime,换算时间单位取整
	nowTimeSeq int64

	_globalLock  sync.Mutex
	_codeSeqMap  = make(map[string]int64)
	_codeLockMap = make(map[string]sync.Mutex)
)

// Init 初始化worker
func Init(ctx context.Context, conf config.RainDropConfig) error {
	log = conf.Logger
	ip, err := utils.GetLocalIP()
	if err != nil {
		log.Error(ctx, "get local ip fail", err)
		return err
	}
	workerCode = ip + ":" + strconv.Itoa(conf.ServicePort) + "#" + utils.GetFirstMacAddr()
	timeUnit = conf.TimeUnit

	worker, err = activateWorker(ctx, conf)
	if worker == nil {
		if err != nil {
			return err
		}
		log.Error(ctx, consts.ErrMsgWorkersNotAvailable)
		return errors.New(consts.ErrMsgWorkersNotAvailable)
	}

	startTime = conf.StartTimeStamp.UnixMilli()
	err = calcNowTimeSeq(ctx)
	if err != nil {
		return err
	}

	if v := ctx.Value(consts.ProjectName); v != nil {
		// 支持单元测试，跳过启动心跳线程
		if consts.SkipHeartbeat == v.(string) {
			return nil
		}
	}

	go startCalcNowTimeSeq(ctx)
	go startHeartbeat(ctx)

	return nil
}

// GetWorkerId 获得WorkerId
func GetWorkerId(ctx context.Context) int64 {
	return worker.Id
}

// GetNowTimeSeq 获得NowTimeSeq
func GetNowTimeSeq(ctx context.Context) int64 {
	return nowTimeSeq
}

// activateWorker 激活worker
func activateWorker(ctx context.Context, conf config.RainDropConfig) (*model.IdGeneratorWorker, error) {

	if conf.PriorityEqualCodeWorkId && (timeUnit == consts.TimeUnitMillisecond || timeUnit == consts.TimeUnitSecond) {
		w, err := db.Db.GetBeforeWorker(ctx, workerCode, int(timeUnit))

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

	for _, w := range workers {
		w2, e := db.Db.ActivateWorker(ctx, w.Id, workerCode, int(timeUnit), w.Version)
		if w2 != nil {
			return w2, e
		}
	}
	return nil, nil
}
