/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type Group struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	TelegramId    int    `json:"telegram_id"`    //对应的每个机器人id
	UsernameGroup string `json:"username_group"` //群名称
	TypeGroup     string `json:"type_group"`     //类型
	IdGroup       int64  `json:"id_group"`       //群组的 id
	Created       int64  `json:"created"`        //创建时间
	Updated       int64  `json:"updated"`        //更新时间
	TelegramName  string `json:"telegram_name" gorm:"-"`
}

func CheckIsExistModelGroup(db *gorm.DB) {
	if db.HasTable(&Group{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Group{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Group{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}

func (g *Group) Add(db *gorm.DB) bool {
	err := db.Where("id_group=?", g.IdGroup).First(&Group{}).Error
	if err == nil {
		return false
	}
	g.Created = time.Now().Unix()
	g.Updated = time.Now().Unix()
	db.Save(&g)
	return true
}
