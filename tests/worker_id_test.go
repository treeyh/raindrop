package tests

import (
	"github.com/treeyh/raindrop"
	"github.com/treeyh/raindrop/worker"
	"testing"
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
