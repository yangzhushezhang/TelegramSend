/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/TelegramSend/dao/mysql"
	"github.com/wangyi/TelegramSend/model"
	"github.com/wangyi/TelegramSend/tools"
	"net/http"
	"strconv"
)

func SetTelegram(c *gin.Context) {
	action := c.Query("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		var total int = 0
		Db := mysql.DB
		fish := make([]model.Telegram, 0)
		Db.Model(&model.Telegram{}).Count(&total)
		Db = Db.Model(&fish).Offset((page - 1) * limit).Limit(limit)
		if err := Db.Find(&fish).Error; err != nil {
			tools.JsonWrite(c, -101, nil, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": fish,
		})
		return

	}
	if action == "ADD" {
		var at AddTelegramVerification
		err := c.ShouldBind(&at)
		if err != nil {
			tools.Return101(c, err.Error())
			return
		}
		tt := model.Telegram{RobotId: at.RobotId, Remark: at.Remark, Token: at.Token, White: at.White}
		err = tt.Add(mysql.DB)
		if err != nil {
			tools.Return101(c, err.Error())
			return
		}
		tools.Return200(c, "添加成功")
		return
	}
	if action == "UPDATE" {

	}
	tools.Return101(c, "no action")
	return
}

func GetGroup(c *gin.Context) {
	action := c.Query("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		var total int = 0
		Db := mysql.DB
		fish := make([]model.Group, 0)
		Db.Model(&model.Group{}).Count(&total)
		Db = Db.Model(&fish).Offset((page - 1) * limit).Limit(limit)
		if err := Db.Find(&fish).Error; err != nil {
			tools.JsonWrite(c, -101, nil, err.Error())
			return
		}

		for i, i2 := range fish {
			te := model.Telegram{}
			err := mysql.DB.Where("id=?", i2.TelegramId).First(&te).Error
			if err != nil {
				fish[i].TelegramName = te.Remark
			}

		}

		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": fish,
		})
		return

	}
}
