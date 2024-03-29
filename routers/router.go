package routers

import (
	"github.com/gin-gonic/gin"
	"go-gin-example/middleware/jwt"
	"go-gin-example/pkg/setting"
	"go-gin-example/routers/api"
	v1 "go-gin-example/routers/api/v1"
)

func InitRouter() *gin.Engine{
	r:=gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(setting.RunMode)
	// token验证
	r.GET("/auth",api.GetAuth)
	//r.GET("/test", func(c *gin.Context) {
	//	c.JSON(200,gin.H{
	//		"message":"test",
	//	})
	//})


	apiv1:=r.Group("/api/v1")
	apiv1.Use(jwt.JWT()) // 在这里添加中间件，每次执行都会先走这里进行token验证
	{
		// 对标签的操作
		// 获取标签列表
		apiv1.GET("/tags",v1.GetTags)
		// 新建标签
		apiv1.POST("/tags",v1.AddTag)
		// 更新指定标签
		apiv1.PUT("/tags/:id",v1.EditTag)
		// 删除指定标签
		apiv1.DELETE("/tags/:id",v1.DeleteTag)


		// 对文章的操作
		// 获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		// 获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		// 新建文章
		apiv1.POST("/articles", v1.AddArticle)
		// 更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		// 删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}
