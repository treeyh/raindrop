package tests

import (
	"github.com/treeyh/raindrop"
	"testing"
)

func TestInit(t *testing.T) {
	ctx := getTestContext()
	//conf := tests.GetMillisecondConfig()
	conf := getTestSecondConfig()
	raindrop.Init(ctx, conf)
	//time.Sleep(time.Duration(20) * time.Minute)
}
