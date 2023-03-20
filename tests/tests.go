package tests

import (
	"context"
	"github.com/treeyh/raindrop/config"
	"github.com/treeyh/raindrop/consts"
	"time"
)

func GetContext() context.Context {
	return context.Background()
}

func GetMySqlConfig() config.RainDropDbConfig {
	return config.RainDropDbConfig{
		DbType: "mysql",
		DbUrl:  "root:mysqlpw@(192.168.80.137:3306)/raindrop_db?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai",
	}
}

func GetConfig() config.RainDropConfig {
	return config.RainDropConfig{
		DbConfig:         GetMySqlConfig(),
		Logger:           nil,
		TimeUnit:         consts.TimeUnitSecond,
		StartTimeStamp:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local),
		TimeLength:       35,
		WorkIdLength:     5,
		ServiceMinWorkId: 10,
		ServiceMaxWorkId: 200,
	}
}
