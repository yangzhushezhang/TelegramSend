/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package process

import (
	"github.com/jinzhu/gorm"
	"github.com/wangyi/TelegramSend/model"
	"github.com/wangyi/TelegramSend/tools"
)

var ListBot []string

//机器的开始
func RobotMonitor(db *gorm.DB) {
	for true {
		t := make([]model.Telegram, 0)
		db.Where("status=?", 1).Find(&t)
		for _, i2 := range t {
			if tools.InArray(i2.RobotId, ListBot) == false {
				//启动机器人程序
				go model.StartBot(i2.Token, db, int(i2.ID),i2.White)
				ListBot = append(ListBot, i2.RobotId)
			}
		}
	}

}
