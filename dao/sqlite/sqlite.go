/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package sqlite

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/wangyi/TelegramSend/model"
)

//type Product struct {
//	gorm.Model
//	Title string
//	Code  string
//	Price uint
//}
var (
	DB  *gorm.DB
	err error
)

func Init() error {
	DB, err = gorm.Open("sqlite3", "telegram.db")
	if err != nil {
		fmt.Println("sqlite数据库链接失败", err)
		panic("failed to connect  sqlite database")
		return err
	}
	//设置连接池
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
	////////////////////////////////////////////////////////////////////////模型初始化
	model.CheckIsExistModelTelegram(DB)
	model.CheckIsExistModelGroup(DB)
	////////////////////////////////////////////////////////////////////////模型初始化
	return err

}

func Close() {
	DB.Close()
}
