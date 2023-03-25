package worker

import (
	"context"
	"fmt"
	"github.com/treeyh/raindrop/config"
	"github.com/treeyh/raindrop/consts"
	"github.com/treeyh/raindrop/db"
	"github.com/treeyh/raindrop/logger"
	"github.com/treeyh/raindrop/model"
	"github.com/treeyh/raindrop/utils"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	workerCode string
	timeUnit   consts.TimeUnit
	log        logger.ILogger
	worker     *model.IdGeneratorWorker

	// idMode id模式
	idMode   string
	workerId int64
	// timeBackBitValue 时间回拨值
	timeBackBitValue atomic.Int64
	// endBitsValue 最后预留位bit值
	endBitsValue int64

	// timeStampShift 时间戳位移位数
	timeStampShift int
	// workerId 位移位数
	workerIdShift int
	// timeBackShift 时间回拨位移位数
	timeBackShift int
	// seqShift 流水号移位数
	seqShift int

	// maxIdSeq 最大的id序列值
	maxIdSeq int64
	// startTime 开始计算时间戳，毫秒
	startTime int64

	// 当前时间流水，当前时刻毫秒 - startTime,换算时间单位取整
	nowTimeSeq atomic.Int64

	// 获取新id的锁
	newIdLock sync.Mutex
	// 上次的获取新id时间序列
	newIdLastTimeSeq atomic.Int64
	// 获取新id同一时间的自增序列
	newIdSeq atomic.Int64

	// 基于code获取新id的生成code锁的锁
	newCodeLockLock sync.Mutex
	// 获取基于code新id的锁
	newIdByCodeLockMap = make(map[string]sync.Mutex)
	// 上次的获取基于code新id时间序列
	newIdByCodeTimeSeqMap = make(map[string]atomic.Int64)
	// 获取基于code新id同一时间的自增序列
	newIdByCodeSeqMap = make(map[string]atomic.Int64)
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
		log.Error(ctx, consts.ErrMsgWorkersNotAvailable.Error())
		return consts.ErrMsgWorkersNotAvailable
	}

	startTime = conf.StartTimeStamp.UnixMilli()
	err = calcNowTimeSeq(ctx)
	if err != nil {
		return err
	}

	initParams(ctx, conf)

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
	return nowTimeSeq.Load()
}

// initParams 初始化参数
func initParams(ctx context.Context, conf config.RainDropConfig) {
	idMode = strings.ToLower(conf.IdMode)
	timeBackBitValue.Store(int64(conf.TimeBackBitValue))
	endBitsValue = int64(conf.EndBitsValue)

	workerId = worker.Id
	seqLength := consts.IdBitLength - conf.TimeStampLength - conf.WorkIdLength - consts.TimeBackBitLength - conf.EndBitsLength

	// 计算同一时刻最大流水号
	maxIdSeq = (1 << seqLength) - 1

	seqShift = conf.EndBitsLength
	timeBackShift = seqLength + seqShift
	workerIdShift = timeBackShift + consts.TimeBackBitLength
	timeStampShift = workerIdShift + conf.WorkIdLength

	log.Info(ctx, fmt.Sprintf("idMode:%s, timeBackBitValue:%d, endBitsValue:%d, workerId:%d, seqLength:%d, "+
		"workerLength:%d, timeLength:%d, maxIdSeq:%d, seqShift: %d, timeBackShift: %d, workerIdShift: %d, timeStampShift:%d",
		idMode, timeBackBitValue.Load(), endBitsValue, workerId, seqLength,
		conf.WorkIdLength, conf.TimeStampLength, maxIdSeq, seqShift, timeBackShift, workerIdShift, timeStampShift))
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
		log.Error(ctx, consts.ErrMsgWorkersNotAvailable.Error())
		return nil, consts.ErrMsgWorkersNotAvailable
	}

	for _, w := range workers {
		w2, e := db.Db.ActivateWorker(ctx, w.Id, workerCode, int(timeUnit), w.Version)
		if w2 != nil {
			return w2, e
		}
	}
	return nil, nil
}

func NewId(ctx context.Context) (int64, error) {
	newIdLock.Lock()
	defer newIdLock.Unlock()

	timeBackValue := timeBackBitValue.Load()
	timestamp := nowTimeSeq.Load()
	lastTimeSeq := newIdLastTimeSeq.Load()

	var seq int64
	if lastTimeSeq == timestamp {
		// 时间戳未发生变化，需要增加newIdSeq
		seq = newIdSeq.Add(1)
		if seq > maxIdSeq {
			// 超过了序列最大值
			if timeUnit != consts.TimeUnitMillisecond && timeUnit != consts.TimeUnitSecond {
				// 不是毫秒或秒时间单位，不等待直接返回错误
				log.Error(ctx, fmt.Sprintf("timeUnit: %d, timeSeq: %d, seq: %d, maxIdSeq: %d",
					int(timeUnit), timestamp, seq, maxIdSeq))
				return 0, consts.ErrMsgIdSeqReachesMaxValueError
			}

			// 毫秒，秒还能抢救一下
			if timeUnit == consts.TimeUnitMillisecond {
				log.Debug(ctx, fmt.Sprintf("millisecond unit sleep %d, seq: %d, maxIdSeq: %d", timestamp, seq, maxIdSeq))
				for {
					timestamp = nowTimeSeq.Load()
					if timestamp > lastTimeSeq {
						break
					}
				}
			} else {
				log.Debug(ctx, fmt.Sprintf("second unit sleep %d, seq: %d, maxIdSeq: %d", timestamp, seq, maxIdSeq))
				for {
					time.Sleep(time.Duration(10) * time.Millisecond)
					timestamp = nowTimeSeq.Load()
					if timestamp > lastTimeSeq {
						break
					}
				}
			}
			seq = 0
			newIdSeq.Store(0)
		}
	} else {
		seq = 0
		newIdSeq.Store(0)
	}

	if lastTimeSeq > timestamp {
		log.Error(ctx, fmt.Sprintf("timeUnit:%d, lastTimeSeq: %d, timestamp: %d ", int(timeUnit), lastTimeSeq, timestamp),
			consts.ErrMsgServerClockBackwardsError)
		// timeBackValue 取反，避免重复
		timeBackValue = timeBackValue ^ 1
		timeBackBitValue.Store(timeBackValue)
	}
	newIdLastTimeSeq.Store(timestamp)

	return ((timestamp - startTime) << timeStampShift) |
		(workerId << workerIdShift) |
		(timeBackValue << timeBackShift) |
		(seq << seqShift) |
		endBitsValue, nil
}

func NewIdByCode(ctx context.Context, code string) (int64, error) {
	return 0, nil
}
