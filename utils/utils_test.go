package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileWithLineNum(t *testing.T) {
	value := FileWithLineNum()
	t.Log(value)
	//assert.NoError(t, err)
}

func TestGetLocalIP(t *testing.T) {
	ip, err := GetLocalIP()
	assert.NoError(t, err)
	t.Log(ip)
}
