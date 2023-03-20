package raindrop

import (
	"github.com/treeyh/raindrop/tests"
	"testing"
)

func TestInit(t *testing.T) {
	ctx := tests.GetContext()
	conf := tests.GetConfig()
	Init(ctx, conf)
}
