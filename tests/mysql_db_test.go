package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/treeyh/raindrop/db"
	"github.com/treeyh/raindrop/logger"
	"testing"
)

func TestMySqlDb_GetNowTime(t *testing.T) {
	ctx := GetContext()
	db.InitMySqlDb(ctx, GetMySqlConfig(), logger.NewDefault())

	now, err := db.Db.GetNowTime(ctx)

	assert.NoError(t, err)
	t.Log(now)
}
