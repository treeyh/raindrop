package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/treeyh/raindrop"
	"github.com/treeyh/raindrop/db"
	"github.com/treeyh/raindrop/logger"
	"testing"
	"time"
)

func TestMySqlDb_GetNowTime(t *testing.T) {
	ctx := getTestContext()
	l := logger.NewDefault()
	db.InitMySqlDb(ctx, getTestMySqlConfig(), l)

	now, err := db.Db.GetNowTime(ctx)

	assert.NoError(t, err)
	t.Log(now)
	t.Log(now.Unix())
}

func TestMySqlDb_QueryFreeWorkers(t *testing.T) {

	ctx := getTestContext()

	conf := getTestSecondConfig()

	raindrop.Init(ctx, conf)

	workers, err := db.Db.QueryFreeWorkers(ctx, time.Now())
	assert.NoError(t, err)
	t.Log(workers)
}
