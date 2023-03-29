package tests

import (
	"github.com/treeyh/raindrop"
	"github.com/treeyh/raindrop/consts"
	"github.com/treeyh/raindrop/db"
	"github.com/treeyh/raindrop/worker"
	"testing"
	"time"
)

// TestGetDefaultWorkerId 正常获取workerId
func TestGetDefaultWorkerId(t *testing.T) {
	ctx := getTestSkipHeartbeatContext()
	conf := getTestMillisecondConfig()

	dropTestWorkerTable(ctx)

	raindrop.Init(ctx, conf)
	t.Log(worker.GetWorkerId(ctx))
	t.Log(worker.GetNowTimeSeq(ctx))
	if worker.GetWorkerId(ctx) != minWorkerId {
		t.Fatalf("%s worker id get error.", t.Name())
	}
	t.Logf("%s pass.", t.Name())
}

// TestWorkerHeartbeat 检查心跳是否生效
func TestWorkerHeartbeat(t *testing.T) {
	ctx := getTestContext()
	conf := getTestMillisecondConfig()

	dropTestWorkerTable(ctx)

	raindrop.Init(ctx, conf)
	t.Log(worker.GetWorkerId(ctx))
	t.Log(worker.GetNowTimeSeq(ctx))
	if worker.GetWorkerId(ctx) != minWorkerId {
		t.Fatalf("%s worker id get error.", t.Name())
	}

	workerId := worker.GetWorkerId(ctx)
	seq := worker.GetNowTimeSeq(ctx)

	w, err := db.Db.GetWorkerById(ctx, workerId)
	if err != nil {
		t.Fatalf("%s worker %d get by db error. %s", t.Name(), workerId, err.Error())
	}
	t.Logf("%s workerId: %d, heartTime: %s, seq: %d.", t.Name(), workerId, w.HeartbeatTime, seq)

	time.Sleep(time.Duration(consts.HeartbeatTimeInterval+2) * time.Second)

	w, err = db.Db.GetWorkerById(ctx, workerId)
	seq = worker.GetNowTimeSeq(ctx)
	if err != nil {
		t.Fatalf("%s worker %d get by db error. %s", t.Name(), workerId, err.Error())
	}
	t.Logf("%s workerId: %d, heartTime: %s, seq: %d.", t.Name(), workerId, w.HeartbeatTime, seq)

	time.Sleep(time.Duration(consts.HeartbeatTimeInterval+2) * time.Second)

	w2, err2 := db.Db.GetWorkerById(ctx, workerId)
	seq = worker.GetNowTimeSeq(ctx)
	if err2 != nil {
		t.Fatalf("%s worker %d get by db error. %s", t.Name(), workerId, err.Error())
	}
	t.Logf("%s workerId: %d, heartTime: %s, seq: %d.", t.Name(), workerId, w2.HeartbeatTime, seq)

	if !w2.HeartbeatTime.After(w.HeartbeatTime) {
		t.Fatalf("%s worker heartbeat error. w: %s  w2: %s;", t.Name(), w.HeartbeatTime, w2.HeartbeatTime)
	}

	raindrop.Init(ctx, conf)
	wId := worker.GetWorkerId(ctx)
	if wId != workerId+1 {
		t.Fatalf("%s get next worker error. wId: %d  w2Id: %d;", t.Name(), workerId, wId)
	}

	t.Logf("%s pass.", t.Name())
}

// TestMultipleGetWorkerIdDiffServicePort1 不同的ServicePort获取workerId场景
func TestMultipleGetWorkerIdDiffServicePort1(t *testing.T) {
	ctx := getTestSkipHeartbeatContext()
	conf := getTestMillisecondConfig()

	dropTestWorkerTable(ctx)

	raindrop.Init(ctx, conf)

	w1 := worker.GetWorkerId(ctx)
	seq1 := worker.GetNowTimeSeq(ctx)

	conf.ServicePort = 1
	raindrop.Init(ctx, conf)

	w2 := worker.GetWorkerId(ctx)
	seq2 := worker.GetNowTimeSeq(ctx)

	t.Logf("work1: %d, seq1: %d, work2: %d, seq2: %d", w1, seq1, w2, seq2)

	if w1 != w2-1 {
		t.Fatalf("%s work w1:%d; w2:%d error.", t.Name(), w1, w2)
	}
	if seq1 >= seq2 {
		t.Fatalf("%s work seq1:%d; seq2:%d error.", t.Name(), seq1, seq2)
	}

	conf.PriorityEqualCodeWorkId = true
	conf.ServicePort = 2
	raindrop.Init(ctx, conf)

	w3 := worker.GetWorkerId(ctx)
	seq3 := worker.GetNowTimeSeq(ctx)

	t.Logf("work1: %d, seq1: %d, work2: %d, seq2: %d, work3: %d, seq3: %d", w1, seq1, w2, seq2, w3, seq3)

	if w2 != w3-1 {
		t.Fatalf("%s work w1:%d; w2:%d; w3:%d error.", t.Name(), w1, w2, w3)
	}
	if seq2 >= seq3 {
		t.Fatalf("%s work seq1:%d; seq2:%d; seq3:%d error.", t.Name(), seq1, seq2, seq3)
	}
	t.Logf("%s pass.", t.Name())
}

// TestMultipleGetWorkerIdSameServicePort1 相同的ServicePort获取workerId场景
func TestMultipleGetWorkerIdSameServicePort1(t *testing.T) {
	ctx := getTestSkipHeartbeatContext()
	conf := getTestMillisecondConfig()

	dropTestWorkerTable(ctx)

	raindrop.Init(ctx, conf)

	w1 := worker.GetWorkerId(ctx)
	seq1 := worker.GetNowTimeSeq(ctx)

	raindrop.Init(ctx, conf)

	w2 := worker.GetWorkerId(ctx)
	seq2 := worker.GetNowTimeSeq(ctx)

	t.Logf("work1: %d, seq1: %d, work2: %d, seq2: %d", w1, seq1, w2, seq2)

	if w1 != w2-1 {
		t.Fatalf("%s work w1:%d; w2:%d error.", t.Name(), w1, w2)
	}
	if seq1 >= seq2 {
		t.Fatalf("%s work seq1:%d; seq2:%d error.", t.Name(), seq1, seq2)
	}

	conf.PriorityEqualCodeWorkId = true
	raindrop.Init(ctx, conf)

	w3 := worker.GetWorkerId(ctx)
	seq3 := worker.GetNowTimeSeq(ctx)

	t.Logf("work1: %d, seq1: %d, work2: %d, seq2: %d, work3: %d, seq3: %d", w1, seq1, w2, seq2, w3, seq3)

	if w1 != w3 {
		t.Fatalf("%s work w1:%d; w2:%d; w3:%d error.", t.Name(), w1, w2, w3)
	}
	if seq2 >= seq3 {
		t.Fatalf("%s work seq1:%d; seq2:%d; seq3:%d error.", t.Name(), seq1, seq2, seq3)
	}
	t.Logf("%s pass.", t.Name())
}

// TestMultipleGetWorkerIdCountOverflow 超过数量获取worker
func TestMultipleGetWorkerIdCountOverflow(t *testing.T) {
	ctx := getTestSkipHeartbeatContext()
	conf := getTestMillisecondConfig()

	dropTestWorkerTable(ctx)

	wIdMap := make(map[int64]bool)
	for i := conf.ServiceMinWorkId; i <= conf.ServiceMaxWorkId; i++ {
		raindrop.Init(ctx, conf)

		w := worker.GetWorkerId(ctx)
		seq := worker.GetNowTimeSeq(ctx)
		t.Logf("get worker:%d seq:%d", w, seq)

		if _, ok := wIdMap[w]; ok {
			t.Fatalf(" wworkerId:%d exist.", w)
			return
		}
		wIdMap[w] = true
	}

	// Open the following comment, there should be no worker available, fatal exit
	//raindrop.Init(ctx, conf)

	t.Logf("%s not pass.", t.Name())
}
