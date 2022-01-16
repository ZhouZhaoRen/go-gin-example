package routers

import (
	_ "go-gin-example/docs" // 这个要添加才可以看到生成的文档
	"go-gin-example/middleware/jwt"
	"go-gin-example/pkg/export"
	"go-gin-example/pkg/qrcode"
	"go-gin-example/pkg/setting"
	"go-gin-example/pkg/upload"
	"go-gin-example/routers/api"
	v1 "go-gin-example/routers/api/v1"
	"net/http"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(setting.ServerSetting.RunMode)
	// 前端访问图片的路径
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	// 前端访问Excel表格的路径
	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	// 前端访问二维码以及海报
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))
	// token验证
	r.GET("/auth", api.GetAuth)
	// 添加swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 上传图片走这里
	r.POST("/upload", api.UploadImage)
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200,gin.H{
			"message":"test",
		})
	})

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT()) // 在这里添加中间件，每次执行都会先走这里进行token验证
	{
		// 对标签的操作
		// 获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		// 新建标签
		apiv1.POST("/tags", v1.AddTag)
		// 更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		// 删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)
		// 导出标签成Excel表格
		r.POST("/tags/export", v1.ExportTag)
		r.POST("/tags/import", v1.ImportTag)

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
		// 产生二维码
		apiv1.POST("/articles/poster/generate", v1.GenerateArticlePoster)
	}

	return r
}
