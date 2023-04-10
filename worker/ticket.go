package worker

import (
	"context"
	"github.com/treeyh/raindrop/consts"
	"github.com/treeyh/raindrop/db"
	"github.com/treeyh/raindrop/logger"
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
	seq := calcTimestamp(ctx, time.Now().UnixMilli(), timeUnit)
	nowTimeSeq.Store(seq)
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
	log.Info(ctx, "worker heartbeat. workerId: "+strconv.FormatInt(worker.Id, 10))
	w, err := db.Db.HeartbeatWorker(ctx, worker)
	if err != nil {
		log.Error(ctx, err.Error(), err)
	}
	if logLevel <= logger.Debug {
		log.Debug(ctx, "worker heartbeat worker: "+utils.ToJsonIgnoreError(w))
	}
	if w.UpdateTime.Unix() > w.HeartbeatTime.Unix()+consts.HeartbeatTimeInterval ||
		w.UpdateTime.Unix() < w.HeartbeatTime.Unix()-consts.HeartbeatTimeInterval {
		log.Error(ctx, consts.ErrMsgDatabaseServerTimeInterval.Error()+". worker:"+utils.ToJsonIgnoreError(w))
	}
	worker = w
	return err
}
