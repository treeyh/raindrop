package tests

import (
	"github.com/treeyh/raindrop"
	"github.com/treeyh/raindrop/worker"
	"testing"
)

func TestGetDefaultWorkerId(t *testing.T) {
	ctx := getTestSkipHeartbeatContext()
	conf := getTestMillisecondConfig()

	dropTestWorkerTable(ctx)

	raindrop.Init(ctx, conf)
	t.Log(worker.GetWorkerId(ctx))
	t.Log(worker.GetNowTimeSeq(ctx))
	if worker.GetWorkerId(ctx) == minWorkerId {
		t.Fatal("worker id get error.")
	}
}
