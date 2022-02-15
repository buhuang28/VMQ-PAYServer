package route

import (
	"Pay-Server/controller"
	"github.com/gin-gonic/gin"
)

func WebStart() {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	oauthController := controller.OauthController{}
	router.POST("/oauth", Oauth)
	router.POST("/AddOauthCode", oauthController.AddOauthCode)
	router.POST("/ActiviteOuathCode", oauthController.ActiviteOuathCode)
	router.POST("/CheckOauthCode", oauthController.CheckOauthCode)
	router.POST("/CheckOauthCode2", oauthController.CheckOauthCode2)
	router.POST("/AuthOauthCode", oauthController.AuthOauthCode)
	router.POST("/RecoverTime", oauthController.RecoverTime)

	payController := controller.PayController{}
	router.POST("/CreateOrder", payController.CreateOrder)
	router.POST("/SelectOrder", payController.SelectOrder)
	router.GET("/PayCallBackOrder", payController.PayCallBack)
	router.GET("/WaitCall", payController.WaitCall)
	router.POST("/DeleteOrder", payController.DeleteOrder)

	router.Run(":8882")
}

type TempResult struct {
	Code int64 `json:"code"`
}

func Oauth(c *gin.Context) {
	c.JSON(0, TempResult{0})
	return
}
