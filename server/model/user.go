package model

import (
	"time"

	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
)

type Users struct {
	ID         int        `json:"id"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `gorm:"index"`
	Name       string     `json:"name"`
	Account    string     `json:"account"`
	Password   string
	InviteCode uint16 `json:"inviteCode"`
	Desc       string `json:"desc"`
	Avatar     string `json:"avatar"` // todo 暂时为url
}

func User() *gorm.DB {
	return mysql.GetInstance().Model(&Users{})
}
