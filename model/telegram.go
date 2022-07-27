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
	"github.com/wangyi/TelegramSend/dao/redis"
	eeor "github.com/wangyi/TelegramSend/error"
	tele "gopkg.in/telebot.v3"

	"log"
	"time"
)

type Telegram struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	RobotId string `json:"robot_id"` //机器人的 id
	Token   string `json:"token"`
	White   string `json:"white"`   //白名单
	Remark  string `json:"remark"`  //备注
	Status  int    `json:"status"`  //状态1 开启 2关闭
	Created int64  `json:"created"` //创建时间
	Updated int64  `json:"updated"` //更新时间
}

//添加一条数据
func (t *Telegram) Add(db *gorm.DB) error {
	//查询是否有重复的
	err := db.Where("robot_id=?", t.RobotId).First(&Telegram{}).Error
	if err == nil {
		return eeor.OtherError("不要重复添加")
	}
	t.Updated = time.Now().Unix()
	t.Created = time.Now().Unix()
	t.Status = 1
	err = db.Save(&t).Error
	if err != nil {
		return eeor.OtherError("对不起,添加失败,原因:" + err.Error())
	}
	return nil
}

//修改小飞机数据
func (t *Telegram) Update(db *gorm.DB) error {
	//查询是否有重复的
	err := db.Where("id=?", t.ID).First(&Telegram{}).Error
	if err == nil {
		return eeor.OtherError("您查询的数据不存在")
	}
	err = db.Model(Telegram{}).Where("id=?", t.ID).Update(&t).Error
	if err != nil {
		return eeor.OtherError("修改失败,原因:" + err.Error())
	}
	return nil
}

//机器人启动

func StartBot(token string, db *gorm.DB, TelegramId int, white string) {
	pref := tele.Settings{
		Token:     token,
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		ParseMode: tele.ModeHTML,
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	//机器人开始  并且权限校验
	b.Handle("/start", func(c tele.Context) error {
		user := c.Sender()
		if user.Username == white {
			return c.Send("请选择您要的操作", SetMenu())
		}
		return nil
	})
	//设置群发的内容
	//返回首页
	b.Handle(&btn1, func(c tele.Context) error {
		redis.Rdb.HSet("StartBot_"+token, "switch", "1") //1说明可以发送
		return c.Send("请输入您要群发的内容")
	})
	//对普通文字的处理
	b.Handle(tele.OnText, func(c tele.Context) error {

		b := tele.User{IsBot: true, ID: 5498054478}
		fmt.Println(c.Send(&b, "??????"))
		data := c.Update().Message.Chat
		g := Group{IdGroup: data.ID, UsernameGroup: data.Title, TypeGroup: string(data.Type), TelegramId: TelegramId}
		g.Add(db)
		return nil
	})
	//对图片的 处理
	b.Handle(tele.OnPhoto, func(c tele.Context) error {
		fmt.Println(c.Message().Photo.FileID)
		data := c.Update().Message.Chat
		g := Group{IdGroup: data.ID, UsernameGroup: data.Title, TypeGroup: string(data.Type), TelegramId: TelegramId}
		g.Add(db)
		return nil
	})
	//添加到群主 触发
	b.Handle(tele.OnAddedToGroup, func(context tele.Context) error {
		//机器人加入群的时候触发
		data := context.Update().Message.Chat
		g := Group{IdGroup: data.ID, UsernameGroup: data.Title, TypeGroup: string(data.Type), TelegramId: TelegramId}
		g.Add(db)
		return nil
	})
	b.Start()
}

var (
	menu = &tele.ReplyMarkup{ResizeKeyboard: true}
	btn1 = menu.Text("设置群发")
)

//设置菜单
func SetMenu() *tele.ReplyMarkup {
	var a []tele.Btn
	a = append(a, btn1)
	menu.Reply(
		menu.Split(4, a)...,
	)
	return menu
}

func CheckIsExistModelTelegram(db *gorm.DB) {
	if db.HasTable(&Telegram{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Telegram{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Telegram{}).Error
		if err == nil {
			fmt.Println("数据库已经存在了!")
		}
	}
}

func MassTexting(Token string, telId int, db *gorm.DB, text string) {
	pref := tele.Settings{
		Token:  Token, //hash 机器人地址
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		return
	}
	te := make([]Group, 0)
	db.Where("telegram_id=?", telId).Find(&te)
	for _, i2 := range te {
		a := tele.User{ID: i2.IdGroup}
		_, _ = b.Send(&a, text)
	}
	b.Stop()
}

func MassTextingIn(Token string, telId []string, db *gorm.DB, text string) {
	pref := tele.Settings{
		Token:  Token, //hash 机器人地址
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		return
	}
	te := make([]Group, 0)
	db.Find(&te, telId)
	for _, i2 := range te {
		a := tele.User{ID: i2.IdGroup}
		_, _ = b.Send(&a, text)
	}
	b.Stop()
}

func MassPhotos(Token string, telId int, db *gorm.DB, text string, filepath string) {
	pref := tele.Settings{
		Token:  Token, //hash 机器人地址
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		return
	}
	te := make([]Group, 0)
	db.Where("telegram_id=?", telId).Find(&te)
	for _, i2 := range te {

		a := tele.User{ID: i2.IdGroup}
		p := &tele.Photo{File: tele.FromURL(filepath), Caption: text}
		_, _ = b.Send(&a, p)
	}
	b.Stop()
}

func MassPhotosTwo(Token string, telId int, db *gorm.DB, text string, filepath string) {
	pref := tele.Settings{
		Token:  Token, //hash 机器人地址
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		return
	}
	te := make([]Group, 0)
	db.Where("telegram_id=?", telId).Find(&te)
	for _, i2 := range te {
		a := tele.User{ID: i2.IdGroup}
		p := &tele.Photo{File: tele.FromDisk(filepath), Caption: text}
		_, _ = b.Send(&a, p)
	}
	b.Stop()
}

func MassPhotosTwoIn(Token string, telId []string, db *gorm.DB, text string, filepath string) {
	pref := tele.Settings{
		Token:  Token, //hash 机器人地址
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		return
	}
	te := make([]Group, 0)
	db.Find(&te, telId)
	for _, i2 := range te {
		a := tele.User{ID: i2.IdGroup}
		p := &tele.Photo{File: tele.FromDisk(filepath), Caption: text}
		_, _ = b.Send(&a, p)
	}
	b.Stop()
}
