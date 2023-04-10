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
	logLevel   logger.LogLevel
	log        logger.ILogger
	worker     *model.RaindropWorker

	// idMode id模式
	idMode   string
	workerId int64
	// timeBackInitValue 时间回拨初始值
	timeBackInitValue int64

	// endBitsValue 最后预留位bit值
	endBitsValue int64

	// timeStampShift 时间戳位移位数
	timeStampShift int
	// workerIdShift 位移位数
	workerIdShift int
	// timeBackShift 时间回拨位移位数
	timeBackShift int
	// seqShift 流水号移位数
	seqShift int

	// maxIdSeq 最大的id序列值
	maxIdSeq int64
	// startTime 开始计算时间戳，毫秒
	startTime int64

	// nowTimeSeq 当前时间流水，当前时刻毫秒 - startTime,换算时间单位取整
	nowTimeSeq atomic.Int64

	// newIdLock 获取新id的锁
	newIdLock sync.Mutex
	// newIdLastTimeSeq 上次的获取新id时间序列
	newIdLastTimeSeq atomic.Int64
	// timeBackBitValue 时间回拨值
	timeBackBitValue atomic.Int64
	// newIdSeq 获取新id同一时间的自增序列
	newIdSeq atomic.Int64

	// newCodeLockLock 基于code获取新id的生成code锁的锁
	newCodeLockLock sync.Mutex
	// newIdByCodeLockMap 获取基于code新id的锁
	newIdByCodeLockMap = make(map[string]sync.Mutex)
	// newIdByCodeTimeSeqMap 上次的获取基于code新id时间序列
	newIdByCodeTimeSeqMap = make(map[string]atomic.Int64)
	// newIdByCodeTimeBackValueMap 基于code 时间回拨值
	newIdByCodeTimeBackValueMap = make(map[string]atomic.Int64)
	// newIdByCodeSeqMap 获取基于code新id同一时间的自增序列
	newIdByCodeSeqMap = make(map[string]atomic.Int64)
)

// Init 初始化worker
func Init(ctx context.Context, conf config.RainDropConfig) error {
	log = conf.Logger
	logLevel = log.GetLogLevel()

	w, err := activateWorker(ctx, conf)
	if w == nil {
		if err != nil {
			return err
		}
		log.Error(ctx, consts.ErrMsgWorkersNotAvailable.Error())
		return consts.ErrMsgWorkersNotAvailable
	}
	worker = w

	initParams(ctx, conf)

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

	go startHeartbeat(ctx)
	go startCalcNowTimeSeq(ctx)

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

// calcTimestamp 计算时间戳
func calcTimestamp(ctx context.Context, timestampMilli int64, timeUnit consts.TimeUnit) int64 {
	st := timestampMilli

	switch timeUnit {
	case consts.TimeUnitSecond:
		st = st / 1000
	case consts.TimeUnitMinute:
		st = st / (1000 * 60)
	case consts.TimeUnitHour:
		st = st / (1000 * 60 * 60)
	case consts.TimeUnitDay:
		st = st / (1000 * 60 * 60 * 24)
	}
	return st
}

// initParams 初始化参数
func initParams(ctx context.Context, conf config.RainDropConfig) {
	idMode = strings.ToLower(conf.IdMode)
	workerId = worker.Id
	timeUnit = conf.TimeUnit

	timeBackInitValue = int64(conf.TimeBackBitValue)
	timeBackBitValue.Store(timeBackInitValue)
	endBitsValue = int64(conf.EndBitsValue)

	startTime = calcTimestamp(ctx, conf.StartTimeStamp.UnixMilli(), conf.TimeUnit)

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
func activateWorker(ctx context.Context, conf config.RainDropConfig) (*model.RaindropWorker, error) {

	ip, err := utils.GetLocalIP()
	if err != nil {
		log.Error(ctx, "get local ip fail", err)
		return nil, err
	}
	timeUnit = conf.TimeUnit
	workerCode = ip + "#" + strconv.Itoa(conf.ServicePort) + "#" + strconv.Itoa(int(timeUnit)) + "#" + utils.GetFirstMacAddr()

	if conf.PriorityEqualCodeWorkId && (timeUnit == consts.TimeUnitMillisecond || timeUnit == consts.TimeUnitSecond) {
		w, e := db.Db.GetBeforeWorker(ctx, workerCode)

		if e != nil {
			return nil, err
		}
		if w != nil {
			w, e = db.Db.ActivateWorker(ctx, w.Id, workerCode, int(timeUnit), w.Version)
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
			// 毫秒，秒还能抢救一下
			if timeUnit == consts.TimeUnitMillisecond {
				if logLevel <= logger.Debug {
					log.Debug(ctx, fmt.Sprintf("millisecond unit sleep %d, seq: %d, maxIdSeq: %d", timestamp, seq, maxIdSeq))
				}
				for {
					timestamp = nowTimeSeq.Load()
					if timestamp > lastTimeSeq {
						break
					}
				}
			} else if timeUnit == consts.TimeUnitSecond {
				if logLevel <= logger.Debug {
					log.Debug(ctx, fmt.Sprintf("second unit sleep %d, seq: %d, maxIdSeq: %d", timestamp, seq, maxIdSeq))
				}
				for {
					time.Sleep(time.Duration(10) * time.Millisecond)
					timestamp = nowTimeSeq.Load()
					if timestamp > lastTimeSeq {
						break
					}
				}
			} else {
				// 不是毫秒或秒时间单位，不等待直接返回错误
				log.Error(ctx, fmt.Sprintf("timeUnit: %d, timeSeq: %d, seq: %d, maxIdSeq: %d",
					int(timeUnit), timestamp, seq, maxIdSeq))
				return 0, consts.ErrMsgIdSeqReachesMaxValueError
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
	if lastTimeSeq != timestamp {
		newIdLastTimeSeq.Store(timestamp)
	}

	//log.Debug(ctx, fmt.Sprintf("timestamp:%d, startTime:%d, timeStampShift:%d\n", timestamp, startTime, timeStampShift))
	//log.Debug(ctx, fmt.Sprintf("workerId:%d, workerIdShift:%d\n", workerId, workerIdShift))
	//log.Debug(ctx, fmt.Sprintf("timeBackValue:%d, timeBackShift:%d\n", timeBackValue, timeBackShift))
	//log.Debug(ctx, fmt.Sprintf("seq:%d, seqShift:%d\n", seq, seqShift))
	//log.Debug(ctx, fmt.Sprintf("endBitsValue：%d\n", endBitsValue))

	return ((timestamp - startTime) << timeStampShift) |
		(workerId << workerIdShift) |
		(timeBackValue << timeBackShift) |
		(seq << seqShift) |
		endBitsValue, nil
}

func NewIdByCode(ctx context.Context, code string) (int64, error) {

	if lock, ok := newIdByCodeLockMap[code]; ok {
		lock.Lock()
		lock.Unlock()
	} else {
		generateCodeLock(ctx, code)
		if lock, ok = newIdByCodeLockMap[code]; ok {
			lock.Lock()
			lock.Unlock()
		} else {
			return 0, consts.ErrMsgGetCodeLockFail
		}
	}

	timeBack, _ := newIdByCodeTimeBackValueMap[code]
	timeBackValue := timeBack.Load()
	timestamp := nowTimeSeq.Load()

	codeIdSeq, _ := newIdByCodeSeqMap[code]

	lastTime, _ := newIdByCodeTimeSeqMap[code]
	lastTimeSeq := lastTime.Load()

	var seq int64
	if lastTimeSeq == timestamp {
		// 时间戳未发生变化，需要增加newIdSeq
		seq = codeIdSeq.Add(1)
		if seq > maxIdSeq {
			// 超过了序列最大值

			// 毫秒，秒还能抢救一下
			if timeUnit == consts.TimeUnitMillisecond {
				if logLevel <= logger.Debug {
					log.Debug(ctx, fmt.Sprintf("code:%s, millisecond unit sleep %d, seq: %d, maxIdSeq: %d", code, timestamp, seq, maxIdSeq))
				}
				for {
					timestamp = nowTimeSeq.Load()
					if timestamp != lastTimeSeq {
						break
					}
				}
			} else if timeUnit == consts.TimeUnitSecond {
				if logLevel <= logger.Debug {
					log.Debug(ctx, fmt.Sprintf("code:%s, second unit sleep %d, seq: %d, maxIdSeq: %d", code, timestamp, seq, maxIdSeq))
				}
				for {
					time.Sleep(time.Duration(10) * time.Millisecond)
					timestamp = nowTimeSeq.Load()
					if timestamp != lastTimeSeq {
						break
					}
				}
			} else {
				// 不是毫秒或秒时间单位，不等待直接返回错误
				log.Error(ctx, fmt.Sprintf("code:%s, timeUnit: %d, timeSeq: %d, seq: %d, maxIdSeq: %d",
					code, int(timeUnit), timestamp, seq, maxIdSeq))
				return 0, consts.ErrMsgIdSeqReachesMaxValueError
			}
			seq = 0
			codeIdSeq.Store(0)
		}
	} else {
		seq = 0
		codeIdSeq.Store(0)
	}
	newIdByCodeSeqMap[code] = codeIdSeq

	if lastTimeSeq > timestamp {
		log.Error(ctx, fmt.Sprintf("timeUnit:%d, lastTimeSeq: %d, timestamp: %d ", int(timeUnit), lastTimeSeq, timestamp),
			consts.ErrMsgServerClockBackwardsError)
		// timeBackValue 取反，避免重复
		timeBackValue = timeBackValue ^ 1
		timeBackBitValue.Store(timeBackValue)
		newIdByCodeTimeBackValueMap[code] = timeBackBitValue
	}
	if lastTimeSeq != timestamp {
		lastTime.Store(timestamp)
		newIdByCodeTimeSeqMap[code] = lastTime
	}

	return ((timestamp - startTime) << timeStampShift) |
		(workerId << workerIdShift) |
		(timeBackValue << timeBackShift) |
		(seq << seqShift) |
		endBitsValue, nil

	return 0, nil
}

func generateCodeLock(ctx context.Context, code string) {
	newCodeLockLock.Lock()
	defer newCodeLockLock.Unlock()

	if _, ok := newIdByCodeTimeSeqMap[code]; ok {
		return
	}

	var lastTime atomic.Int64
	var seq atomic.Int64
	var timeBackValue atomic.Int64
	var codeLock sync.Mutex
	lastTime.Store(0)
	seq.Store(0)
	timeBackValue.Store(timeBackInitValue)

	newIdByCodeTimeBackValueMap[code] = timeBackBitValue
	newIdByCodeLockMap[code] = codeLock
	newIdByCodeTimeSeqMap[code] = lastTime
	newIdByCodeSeqMap[code] = seq
}
