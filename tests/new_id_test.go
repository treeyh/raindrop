package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/treeyh/raindrop"
	"github.com/treeyh/raindrop/worker"
	"testing"
	"time"
)

func TestTicketInterval(t *testing.T) {
	ticket := worker.NewTicket(time.Duration(1)*time.Millisecond, func(ctx context.Context) error {
		t.Log(time.Now().UnixMilli())
		return nil
	})
	ticket.Start(getTestContext())
}

// TestSimpleNewId 获取id
func TestSimpleNewId(t *testing.T) {
	ctx := getTestContext()
	//conf := getTestSimpleMillisecondConfig()
	conf := getTestSecondConfig()

	dropTestWorkerTable(ctx, conf.DbConfig.DbType)

	raindrop.Init(ctx, conf)
	if worker.GetWorkerId(ctx) != minWorkerId {
		t.Fatalf("%s worker id get error.", t.Name())
	}

	t.Logf("%s pass.", t.Name())

	batchNewId(ctx, t, 0, true)

	time.Sleep(time.Duration(10) * time.Second)

	t.Log("End")
}

// TestBenchmarkNewId 压力测试
func TestBenchmarkNewId(t *testing.T) {
	ctx := getTestContext()
	//conf := getTestSimpleMillisecondConfig()
	conf := getTestSecondConfig()

	dropTestWorkerTable(ctx, conf.DbConfig.DbType)

	raindrop.Init(ctx, conf)
	if worker.GetWorkerId(ctx) != minWorkerId {
		t.Fatalf("%s worker id get error.", t.Name())
	}

	t.Logf("%s pass.", t.Name())

	for i := 0; i < 16; i++ {
		go batchNewId(ctx, t, i, false)
	}

	time.Sleep(time.Duration(60) * time.Second)
}

// TestSimpleLongTimeNewId 获取长时间获取id
func TestSimpleLongTimeNewId(t *testing.T) {
	ctx := getTestContext()
	//conf := getTestSimpleMillisecondConfig()
	conf := getTestSecondConfig()

	dropTestWorkerTable(ctx, conf.DbConfig.DbType)

	raindrop.Init(ctx, conf)
	if worker.GetWorkerId(ctx) != minWorkerId {
		t.Fatalf("%s worker id get error.", t.Name())
	}

	t.Logf("%s pass.", t.Name())

	batchNewId(ctx, t, 0, true)

	time.Sleep(time.Duration(10) * time.Second)

	endTime := time.Now().Add(time.Duration(6) * time.Hour)

	for true {
		id, err := worker.NewId(ctx)
		assert.NoError(t, err)
		t.Logf("new id: %d", id)

		if endTime.Unix() < time.Now().Unix() {
			break
		}
		time.Sleep(time.Duration(5) * time.Second)
	}

	t.Log("End")
}

func batchNewId(ctx context.Context, t *testing.T, index int, logFlag bool) {

	idMap := make(map[int64]bool)
	start := time.Now().UnixMilli()
	for i := 0; i < 10000000; i++ {
		id, err := raindrop.NewId()
		if err != nil {
			t.Fatalf("%s newId get fail. %s", t.Name(), err.Error())
		}
		if _, ok := idMap[id]; ok {
			t.Errorf("%s duplicate id generated: %d", t.Name(), id)
		}
		idMap[id] = true
		if logFlag && i%100000 == 0 {
			t.Logf("%s new id index: %d id: %d ", t.Name(), i, id)
		}
	}
	end := time.Now().UnixMilli()
	t.Logf("index:%d new id start:%d, end:%d, time: %d", index, start, end, end-start)
}
