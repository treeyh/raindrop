package worker

import (
	"context"
	"github.com/treeyh/raindrop/consts"
	"github.com/treeyh/raindrop/db"
	"github.com/treeyh/raindrop/utils"
	"strconv"
	"time"
)

type Fun func(ctx context.Context) error

type Ticket struct {
	ticket *time.Ticker

	runner Fun
}

func NewTicket(dur time.Duration, f Fun) *Ticket {
	return &Ticket{
		ticket: time.NewTicker(dur),
		runner: f,
	}
}

// Start 启动定时器需要执行的任务
func (t *Ticket) Start(ctx context.Context) {
	for {
		select {
		case <-t.ticket.C:
			t.runner(ctx)
		}
	}
}

// calcNowTimeSeq 计算当前时间戳流水
func calcNowTimeSeq(ctx context.Context) error {
	seq := time.Now().UnixMilli() - startTime
	if seq < 0 {
		log.Error(ctx, consts.ErrMsgStartTimeStampError.Error())
	}
	switch timeUnit {
	case consts.TimeUnitMillisecond:
		nowTimeSeq = seq
	case consts.TimeUnitSecond:
		nowTimeSeq = seq / 1000
	case consts.TimeUnitMinute:
		nowTimeSeq = seq / (1000 * 60)
	case consts.TimeUnitHour:
		nowTimeSeq = seq / (1000 * 60 * 60)
	case consts.TimeUnitDay:
		nowTimeSeq = seq / (1000 * 60 * 60 * 24)
	}
	log.Debug(ctx, "nowTimeSeq: "+strconv.FormatInt(nowTimeSeq, 10))
	return nil
}

func startCalcNowTimeSeq(ctx context.Context) {
	if timeUnit == consts.TimeUnitMillisecond {
		ticket := NewTicket(time.Duration(1)*time.Millisecond, calcNowTimeSeq)
		ticket.Start(ctx)
	} else {
		ticket := NewTicket(time.Duration(1)*time.Second, calcNowTimeSeq)
		ticket.Start(ctx)
	}
}

// startHeartbeat 启动心跳
func startHeartbeat(ctx context.Context) {
	ticket := NewTicket(time.Duration(consts.HeartbeatTimeInterval)*time.Second, heartbeat)
	ticket.Start(ctx)
}

func heartbeat(ctx context.Context) error {
	log.Info(ctx, "worker heartbeat start. workerId: "+strconv.FormatInt(worker.Id, 10))
	w, err := db.Db.HeartbeatWorker(ctx, worker)
	log.Info(ctx, "worker heartbeat end. workerId: "+strconv.FormatInt(worker.Id, 10))
	if err != nil {
		log.Error(ctx, err.Error(), err)
	}
	if w != nil {
		log.Debug(ctx, "worker heartbeat worker: "+utils.ToJsonIgnoreError(w))
		if w.UpdateTime.Unix() > w.HeartbeatTime.Unix()+consts.HeartbeatTimeInterval ||
			w.UpdateTime.Unix() < w.HeartbeatTime.Unix()-consts.HeartbeatTimeInterval {
			log.Error(ctx, consts.ErrMsgDatabaseServerTimeInterval.Error()+". worker:"+utils.ToJsonIgnoreError(w))
		}
		worker = w
	}
	return err
}
