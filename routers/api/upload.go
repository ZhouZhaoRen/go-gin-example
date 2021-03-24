package api

import (
	"github.com/gin-gonic/gin"
	"go-gin-example/pkg/e"
	"go-gin-example/pkg/logging"
	"go-gin-example/pkg/upload"
	"net/http"
)

func UploadImage(c *gin.Context){
	code:=e.SUCCESS
	data:=make(map[string]interface{})
	// 获取上传的图片（返回提供的表单键的第一个文件）
	file,image,err:=c.Request.FormFile("image")
	// 若是出现错误，走这里
	if err != nil {
		logging.Warn(err)
		code=e.ERROR
		c.JSON(http.StatusOK,gin.H{
			"code":code,
			"msg":e.GetMsg(code),
			"data":data,
		})
	}
	// 正常接收到前端传过来的数据
	if image==nil{
		code=e.INVALID_PARAMS
		logging.Warn(e.GetMsg(code))
	}else{
		imageName:=upload.GetImageName(image.Filename)
		fullPath:=upload.GetImageFullPath()
		savePath:=upload.GetImagePath()

		src:=fullPath+imageName
		// 检查图片大小、检查图片后缀
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
			code=e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
			logging.Warn(e.GetMsg(code))
		}else{
			// 检查上传图片所需（权限，文件夹）
			err:=upload.CheckImage(fullPath)
			if err != nil {
				logging.Warn(err)
				code=e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
				logging.Warn(e.GetMsg(code))
				// 保存图片
			}else if err:=c.SaveUploadedFile(image,src);err !=nil {
				logging.Warn(err)
				code=e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
				logging.Warn(e.GetMsg(code))
			}else{
				data["image_url"]=upload.GetImageFullUrl(imageName)
				data["image_save_url"]=savePath+imageName
			}
		}
	}
	// 返回json数据
	c.JSON(http.StatusOK,gin.H{
		"code":code,
		"msg":e.GetMsg(code),
		"data":data,
	})

}
