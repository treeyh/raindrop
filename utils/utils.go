package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var raindropSourceDir string

func init() {
	_, file, _, _ := runtime.Caller(0)
	raindropSourceDir = regexp.MustCompile(`utils.utils\.go`).ReplaceAllString(file, "")
}

func FileWithLineNum() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		//  && strings.HasPrefix(file, raindropSourceDir)
		if ok && !strings.HasSuffix(file, "_test.go") {
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}
	return ""
}

func ToJson(obj interface{}) (string, error) {
	bs, err := json.Marshal(obj)
	return string(bs), err
}

func ToJsonIgnoreError(obj interface{}) string {
	bs, _ := json.Marshal(obj)
	return string(bs)
}

// GetLocalIP 获取内网ip
func GetLocalIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		return ipAddr.IP.String(), nil
	}
	return
}

func GetMacAddrs() (macAddrs []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v", err)
		return macAddrs
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}

		macAddrs = append(macAddrs, macAddr)
	}
	return macAddrs
}

func GetFirstMacAddr() string {
	macs := GetMacAddrs()
	if len(macs) > 0 {
		return macs[0]
	}
	return ""
}
