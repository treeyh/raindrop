package utils

import (
	"testing"
)

func TestFileWithLineNum(t *testing.T) {
	value := FileWithLineNum()
	t.Log(value)
}

func TestGetLocalIP(t *testing.T) {
	ip, err := GetLocalIP()
	if err != nil {
		t.Fatal("get local ip fail: "+err.Error(), err)
	}
	t.Log(ip)
}

func TestGetMacAddrs(t *testing.T) {
	macs := GetMacAddrs()

	t.Log(macs)
}
