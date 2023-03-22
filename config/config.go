package config

import (
	"context"
	"errors"
	"github.com/treeyh/raindrop/consts"
	"github.com/treeyh/raindrop/logger"
	"time"
)

type RainDropDbConfig struct {
	// DbType 数据库类型，mysql、postgresql
	DbType string `json:"dbType"`

	// 数据库连接，{user}:{password}@({host}:{port})/{dbName}?charset=utf8mb4&parseTime=True&loc={Asia%2FShanghai}
	DbUrl string `json:"dbUrl"`
}

type RainDropConfig struct {

	// DbConfig 数据库配置
	DbConfig RainDropDbConfig `json:"dbConfig"`

	// Logger 日志
	Logger logger.ILogger `json:"logger"`

	// ServicePort 服务端口
	ServicePort int `json:"servicePort"`

	// TimeUnit 时间戳单位, 1：毫秒（可能会有闰秒问题）；2：秒，默认；3：分钟；4：小时，间隔过大不建议选择；5：天，间隔过大不建议选择；
	TimeUnit consts.TimeUnit `json:"timeUnit"`

	// StartTimeStamp 起始时间，时间戳从该时间开始计时，格式：2020-01-01T00:00:00.000+0000
	StartTimeStamp time.Time `json:"startTimeStamp"`

	// TimeLength 时间戳位数
	/*
	  - timeUnit 为 1 时，取值范围 41-50 位，值越大每毫秒支持生成的 id 数就越少；
	    - 41: 约 69.7 年，默认；
	    - 42：约 139.4 年；
	    - 43：约 278.9 年；
	    - 44：约 557.8 年；
	    - 45：约 1115.6 年；
	    - 46：约 2231.3 年；
	    - 47：约 4462.7 年；
	    - 48：约 8925.5 年；
	    - 49：约 17851 年；
	    - 50：约 35702 年；
	  - timeUnit 为 2 时，取值范围 31-40 位，值越大每毫秒支持生成的 id 数就越少；
	    - 31：约 68 年；
	    - 32：约 136.1 年；
	    - 33：约 272.3 年，默认；
	    - 34：约 544.7 年；
	    - 35：约 1089.5 年；
	    - 36：约 2179 年；
	    - 37：约 4358.1 年；
	    - 38：约 8716.3 年；
	    - 39：约 17432.6 年；
	    - 40：约 34865.2 年；
	  - timeUnit 为 3 时，取值范围 25-34 位，值越大每毫秒支持生成的 id 数就越少；
	    - 25：约 63.8 年；
	    - 26：约 127.6 年；
	    - 27：约 255.3 年，默认；
	    - 28：约 510.7 年；
	    - 29：约 1021.4 年；
	    - 30：约 2042.8 年；
	    - 31：约 4085.7 年；
	    - 32：约 8171.5 年；
	    - 33：约 16343.1 年；
	    - 34：约 32686.2 年；
	  - timeUnit 为 4 时，取值范围 19-28 位，值越大每毫秒支持生成的 id 数就越少；
	    - 19：约 59.8 年；
	    - 20：约 119.7 年；
	    - 21：约 239.4 年，默认；
	    - 22：约 478.8 年；
	    - 23：约 957.6 年；
	    - 24：约 1915.2 年；
	    - 25：约 3830.4 年；
	    - 26：约 7660.8 年；
	    - 27：约 15321.6 年；
	    - 28：约 30643.3 年；
	  - timeUnit 为 5 时，取值范围 15-24 位，值越大每毫秒支持生成的 id 数就越少；
	    - 15：约 89.7 年；
	    - 16：约 179.5 年；
	    - 17：约 359 年，默认；
	    - 18：约 718.2 年；
	    - 19：约 1436.4 年；
	    - 20：约 2872.8 年；
	    - 21：约 5745.6 年；
	    - 22：约 11491.2 年；
	    - 23：约 22982.4 年；
	    - 24：约 45964.9 年；
	*/
	TimeLength int `json:"timeLength"`

	// PriorityEqualCodeWorkId 优先相同code的workerId(毫秒，秒单位场景下生效)，默认：false。code格式为：{内网ip}:{ServicePort}#{Mac地址}
	PriorityEqualCodeWorkId bool `json:"priorityEqualCodeWorkId"`

	// WorkIdLength 工作节点 id 长度，取值范围 4 - 10 位.
	/*
	  - 4：支持 15 个工作节点，默认，取值范围：1-15；
	  - 5：支持 31 个工作节点，取值范围：1-31；
	  - 6：支持 63 个工作节点，取值范围：1-63；
	  - 7：支持 127 个工作节点，取值范围：1-127；
	  - 8：支持 255 个工作节点，取值范围：1-255；
	  - 9：支持 511 个工作节点，取值范围：1-511；
	  - 10：支持 1023 个工作节点，1-1023；
	*/
	WorkIdLength int `json:"workIdLength"`

	// ServiceMinWorkId 服务的最小工作节点 id，默认 1，需在 workIdLength 的定义范围内，最大值最小值用于不同数据中心的隔离。
	ServiceMinWorkId int64 `json:"serviceMinWorkId"`

	// ServiceMaxWorkId 服务的最大工作节点 id，默认 workIdLength 的最大值，需在 workIdLength 的定义范围内。
	ServiceMaxWorkId int64 `json:"serviceMaxWorkId"`
}

func CheckConfig(ctx context.Context, conf RainDropConfig) error {
	if conf.ServicePort < 1 || conf.ServicePort > 65535 {
		return errors.New("ServicePort range between 1 and 65535")
	}

	switch conf.TimeUnit {
	case consts.TimeUnitMillisecond:
		if conf.TimeLength < 41 || conf.TimeLength > 50 {
			return errors.New("When TimeUnit is millisecond, TimeLength must be between 41 and 50")
		}
	case consts.TimeUnitSecond:
		if conf.TimeLength < 31 || conf.TimeLength > 40 {
			return errors.New("When TimeUnit is second, TimeLength must be between 31 and 40")
		}
	case consts.TimeUnitMinute:
		if conf.TimeLength < 25 || conf.TimeLength > 34 {
			return errors.New("When TimeUnit is minute, TimeLength must be between 25 and 34")
		}
	case consts.TimeUnitHour:
		if conf.TimeLength < 19 || conf.TimeLength > 28 {
			return errors.New("When TimeUnit is hour, TimeLength must be between 19 and 28")
		}
	case consts.TimeUnitDay:
		if conf.TimeLength < 15 || conf.TimeLength > 24 {
			return errors.New("When TimeUnit is day, TimeLength must be between 15 and 24")
		}
	}

	if conf.WorkIdLength < 4 || conf.WorkIdLength > 10 {
		return errors.New("WorkIdLength takes values between 4 and 10")
	}

	if conf.ServiceMinWorkId > conf.ServiceMaxWorkId {
		return errors.New("ServiceMaxWorkId must be greater than ServiceMinWorkId")
	}

	switch conf.WorkIdLength {
	case 4:
		if conf.ServiceMinWorkId < 1 || conf.ServiceMaxWorkId > 15 {
			return errors.New("When WorkIdLength is 4, ServiceMinWorkId and ServiceMaxWorkId take values in the range of 1 to 15.")
		}
	case 5:
		if conf.ServiceMinWorkId < 1 || conf.ServiceMaxWorkId > 31 {
			return errors.New("When WorkIdLength is 5, ServiceMinWorkId and ServiceMaxWorkId take values in the range of 1 to 31.")
		}
	case 6:
		if conf.ServiceMinWorkId < 1 || conf.ServiceMaxWorkId > 63 {
			return errors.New("When WorkIdLength is 6, ServiceMinWorkId and ServiceMaxWorkId take values in the range of 1 to 63.")
		}
	case 7:
		if conf.ServiceMinWorkId < 1 || conf.ServiceMaxWorkId > 127 {
			return errors.New("When WorkIdLength is 7, ServiceMinWorkId and ServiceMaxWorkId take values in the range of 1 to 127.")
		}
	case 8:
		if conf.ServiceMinWorkId < 1 || conf.ServiceMaxWorkId > 255 {
			return errors.New("When WorkIdLength is 8, ServiceMinWorkId and ServiceMaxWorkId take values in the range of 1 to 255.")
		}
	case 9:
		if conf.ServiceMinWorkId < 1 || conf.ServiceMaxWorkId > 511 {
			return errors.New("When WorkIdLength is 9, ServiceMinWorkId and ServiceMaxWorkId take values in the range of 1 to 511.")
		}
	case 10:
		if conf.ServiceMinWorkId < 1 || conf.ServiceMaxWorkId > 1023 {
			return errors.New("When WorkIdLength is 10, ServiceMinWorkId and ServiceMaxWorkId take values in the range of 1 to 1023.")
		}
	}

	if time.Now().Unix() < conf.StartTimeStamp.Unix() {
		return consts.ErrMsgStartTimeStampError
	}

	return nil
}
