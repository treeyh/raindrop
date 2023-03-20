package db

import (
	"github.com/stretchr/testify/assert"
	"github.com/treeyh/raindrop/logger"
	"github.com/treeyh/raindrop/tests"
	"testing"
	"time"
)

func TestMySqlDb_GetNowTime(t *testing.T) {
	ctx := tests.GetContext()
	l := logger.NewDefault()
	InitMySqlDb(ctx, tests.GetMySqlConfig(), l)

	now, err := Db.GetNowTime(ctx)

	assert.NoError(t, err)
	t.Log(now)
	t.Log(now.Unix())
}

func TestMySqlDb_QueryFreeWorkers(t *testing.T) {

	ctx := tests.GetContext()
	InitMySqlDb(ctx, tests.GetMySqlConfig(), logger.NewDefault())

	workers, err := Db.QueryFreeWorkers(ctx, time.Now())
	assert.NoError(t, err)
	t.Log(workers)
}
