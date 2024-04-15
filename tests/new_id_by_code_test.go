package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/treeyh/raindrop"
	"github.com/treeyh/raindrop/worker"
	"testing"
	"time"
)

const (
	accountCode = "account"

	orderCode = "order"
)

// TestSimpleNewId 获取id
func TestSimpleNewIdByCode(t *testing.T) {
	ctx := getTestContext()
	conf := getTestMinuteConfig()

	dropTestWorkerTable(ctx, conf.DbConfig.DbType)

	raindrop.Init(ctx, conf)
	if worker.GetWorkerId(ctx) != minWorkerId {
		t.Fatalf("%s worker id get error.", t.Name())
	}

	t.Logf("%s pass.", t.Name())

	idMap1 := batchNewIdByCode(ctx, t, accountCode)
	idMap2 := batchNewIdByCode(ctx, t, orderCode)

	for k, _ := range idMap1 {
		if _, ok := idMap2[k]; ok {
			t.Logf("%d id exist.", k)
			return
		}
	}
	assert.Fail(t, "The same id should exist for different code")
	t.Log("End")
}

func batchNewIdByCode(ctx context.Context, t *testing.T, code string) map[int64]bool {

	idMap := make(map[int64]bool)
	start := time.Now().UnixMilli()
	for i := 0; i < 100000; i++ {
		id, err := raindrop.NewIdByCode(code)
		if err != nil {
			t.Fatalf("%s newId get fail. %s", t.Name(), err.Error())
		}
		if _, ok := idMap[id]; ok {
			t.Errorf("%s duplicate id generated: %d", t.Name(), id)
		}
		idMap[id] = true
	}
	end := time.Now().UnixMilli()
	t.Logf("code:%s new id start:%d, end:%d, time: %d", code, start, end, end-start)
	return idMap
}
