/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangyi/TelegramSend/dao/sqlite"
	"github.com/wangyi/TelegramSend/model"
	"github.com/wangyi/TelegramSend/tools"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func SetTelegram(c *gin.Context) {
	action := c.Query("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		var total int = 0
		Db := sqlite.DB
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
		err = tt.Add(sqlite.DB)
		if err != nil {
			tools.Return101(c, err.Error())
			return
		}
		tools.Return200(c, "添加成功")
		return
	}
	if action == "UPDATE" {

	}

	//群发消息
	if action == "SEND" {
		kinds := c.PostForm("kinds")
		telegramId, _ := strconv.Atoi(c.PostForm("telegram_id"))
		text := c.PostForm("text")
		//群发
		tr := model.Telegram{}
		err := sqlite.DB.Where("id=?", telegramId).First(&tr).Error
		if err != nil {
			tools.Return101(c, "群发失败,机器人不存在")
			return
		}
		if ids, isE := c.GetPostForm("id"); isE == false {
			//只发文本内容
			if kinds == "text" {
				go model.MassTexting(tr.Token, telegramId, sqlite.DB, text)
			} else if kinds == "photo" {
				file, err := c.FormFile("file")
				if err != nil {
					tools.Return101(c, "fail")
					return
				}
				if file.Size > 499823 {
					tools.Return101(c, "Picture is too big")
					return
				}
				//判断是否是图片
				nameArray := strings.Split(file.Filename, ".")
				f, _ := file.Open()
				switch strings.ToUpper(nameArray[1]) {
				case "JPG", "JPEG":
					_, err = jpeg.Decode(f)
				case "PNG":
					_, err = png.Decode(f)
					fmt.Println(err)
				case "GIF":
					_, err = gif.Decode(f)
				default:
					tools.Return101(c, " Invalid file")
					return
				}
				if err != nil {
					fmt.Println(err)
					tools.Return101(c, " image is  error")
					return
				}
				nowStr := time.Now().Format("20060102150405")
				filepath := "./static/upload/" + nowStr + ".png"
				err = c.SaveUploadedFile(file, filepath)
				go model.MassPhotosTwo(tr.Token, telegramId, sqlite.DB, text, filepath)
			}
		} else {
			arrayString := strings.Split(ids, "@")
			if len(arrayString) > 0 {
				if kinds == "text" {
					go model.MassTextingIn(tr.Token, arrayString, sqlite.DB, text)
				} else if kinds == "photo" {
					file, err := c.FormFile("file")
					if err != nil {
						tools.Return101(c, "fail")
						return
					}
					if file.Size > 499823 {
						tools.Return101(c, "Picture is too big")
						return
					}
					//判断是否是图片
					nameArray := strings.Split(file.Filename, ".")
					f, _ := file.Open()
					switch strings.ToUpper(nameArray[1]) {
					case "JPG", "JPEG":
						_, err = jpeg.Decode(f)
					case "PNG":
						_, err = png.Decode(f)
						fmt.Println(err)
					case "GIF":
						_, err = gif.Decode(f)
					default:
						tools.Return101(c, " Invalid file")
						return
					}
					if err != nil {
						fmt.Println(err)
						tools.Return101(c, " image is  error")
						return
					}
					nowStr := time.Now().Format("20060102150405")
					filepath := "./static/upload/" + nowStr + ".png"
					err = c.SaveUploadedFile(file, filepath)
					go model.MassPhotosTwoIn(tr.Token, arrayString, sqlite.DB, text, filepath)
				}

			}

		}
		tools.Return200(c, "执行成功")
		return
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
		Db := sqlite.DB
		fish := make([]model.Group, 0)
		Db.Model(&model.Group{}).Count(&total)
		Db = Db.Model(&fish).Offset((page - 1) * limit).Limit(limit)
		if err := Db.Find(&fish).Error; err != nil {
			tools.JsonWrite(c, -101, nil, err.Error())
			return
		}

		for i, i2 := range fish {
			te := model.Telegram{}
			err := sqlite.DB.Where("id=?", i2.TelegramId).First(&te).Error
			if err == nil {
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
