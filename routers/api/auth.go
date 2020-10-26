package api

import (
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"go-gin-example/models"
	"go-gin-example/pkg/e"
	"go-gin-example/pkg/logging"
	"go-gin-example/pkg/util"
	"log"
	"net/http"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)
	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	//
	if ok {
		// 判断用户是否存在
		isExist := models.CheckAuth(username, password)
		if isExist {
			// 用户存在，开始生成token
			token, err := util.GenerateToken(username, password)
			if err != nil {
				fmt.Println("生成token失败1：",err)
				code = e.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token
				code=e.SUCCESS
			}
		} else {
			logging.Info(e.GetMsg(code))
			code = e.ERROR_AUTH
		}
	} else { //
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
	}
	//
	c.JSON(http.StatusOK,gin.H{
		"code":code,
		"msg":e.GetMsg(code),
		"data":data, // 将token返回到前端
	})
}