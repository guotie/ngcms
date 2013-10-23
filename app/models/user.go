package models

import (
	"time"
)

type UserInfo struct {
	Uid      int64
	Descp    string
	Github   string
	Weibo    string
	Blog     string
	Location string
}

type User struct {
	Id          int64
	Username    string
	Password    string
	Signature   string
	Email       string
	CardID      string
	Createtm    time.Time
	createtm    int64
	Birthday    time.Time
	birthday    int64
	Protrait_xl string
	Protrait_l  string
	Protrait_m  string
	Protrait_s  string

	*UserInfo
}
