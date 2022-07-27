/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	v1 "github.com/wangyi/TelegramSend/controller/telegram/v1"
	Weer "github.com/wangyi/TelegramSend/error"
)

func Setup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(Weer.ErrHandler())
	r.NoMethod(Weer.HandleNotFound)
	r.NoRoute(Weer.HandleNotFound)

	r.Static("/static", "./static")

	t := r.Group("/telegram/v1")
	{
		t.GET("setTelegram", v1.SetTelegram)
		t.POST("setTelegram", v1.SetTelegram)
		t.GET("setGroup", v1.GetGroup)
	}

	r.Run(fmt.Sprintf(":%d", viper.GetInt("app.port")))
	return r
}
