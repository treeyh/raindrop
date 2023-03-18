package utils

import (
	"testing"
)

func TestFileWithLineNum(t *testing.T) {
	value := FileWithLineNum()
	t.Log(value)
	//assert.NoError(t, err)
}
