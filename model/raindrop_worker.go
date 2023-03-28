package model

import (
	"github.com/treeyh/raindrop/consts"
	"time"
)

type RaindropWorker struct {
	Id int64 `json:"id"`

	Code string `json:"code"`

	TimeUnit consts.TimeUnit `json:"timeUnit"`

	HeartbeatTime time.Time `json:"heartbeatTime"`

	CreateTime time.Time `json:"createTime"`

	UpdateTime time.Time `json:"updateTime"`

	Version int64 `json:"version"`

	DelFlag int `json:"delFlag"`
}
