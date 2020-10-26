package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go-gin-example/models"
	"go-gin-example/pkg/e"
	"go-gin-example/pkg/setting"
	"go-gin-example/pkg/util"
	"log"
	"net/http"
)

// 获取多个文章标签
func GetTags(c *gin.Context) {
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}
	//
	var state int = -1

	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := e.SUCCESS
	data["lists"] = models.GetTags(util.GetPage(c), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	// 返回json数据
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// 新增文章标签
func AddTag(c *gin.Context) {
	name:=c.Query("name")
	state:=com.StrTo(c.DefaultQuery("state","0")).MustInt()
	createdBy:=c.Query("created_by")

	// 通过beego的validation进行判断输入的数据是否符合要求
	valid:=validation.Validation{}
	valid.Required(name,"name").Message("名称不能为空")
	valid.MaxSize(name,100,"name").Message("名称最长为100字符")
	valid.Required(createdBy,"createdBy").Message("创建人不能为空")
	valid.MaxSize(createdBy,100,"createdBy").Message("创建人最长为100字符")
	valid.Range(state,0,1,"state").Message("状态只允许为0或1")

	code:=e.INVALID_PARAMS // 不可用的参数

	if ! valid.HasErrors() {
		if !models.ExistTagByName(name){
			code=e.SUCCESS
			models.AddTag(name,state,createdBy)
		}else{
			code=e.ERROR_EXIST_TAG // 返回标签已存在的错误
		}
	}else{
		// 参数错误的话通过日志的形式打印
		log.Fatalf("参数出现错误：%v",valid.Errors)
	}

	//
	c.JSON(http.StatusOK,gin.H{
		"code":code,
		"msg":e.GetMsg(code),
		"data":make(map[string]string),
	})


}

// 修改文章标签
func EditTag(c *gin.Context) {
	id:=com.StrTo(c.Param("id")).MustInt() // 127.0.0.1:8000/api/v1/tags/1?id=1&state=0&name=zzr&modified_by=zzr  id和别的参数不一样
	name:=c.Query("name")
	modifiedBy:=c.Query("modified_by")

	valid:=validation.Validation{}

	state:=-1
	if arg:=c.Query("state"); arg!=""{
		state=com.StrTo(arg).MustInt()
		valid.Range(state,0,1,"state").Message("状态只允许为0或1")
	}
	//
	valid.Required(id,"id").Message("id不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	//
	code:=e.INVALID_PARAMS
	if !valid.HasErrors(){
		code=e.SUCCESS
		if models.ExistTagById(id){
			data:=make(map[string]interface{})
			data["modified_by"]=modifiedBy
			if name!=""{
				data["name"]=name
			}
			if state!=-1{
				data["state"]=state
			}
			models.UpdateTag(id,data)
		}else{
			code=e.ERROR_NOT_EXIST_TAG // 标签不存在
		}
	}else{
		// 参数错误的话通过日志的形式打印
		//log.Fatalf("参数出现错误：%v",valid.Errors)
	}
	//
	c.JSON(http.StatusOK,gin.H{
		"code":code,
		"msg":e.GetMsg(code),
		"data":make(map[string]string),
	})
}

// 删除文章标签
func DeleteTag(c *gin.Context) {
	id:=com.StrTo(c.Param("id")).MustInt()
	valid:=validation.Validation{}
	valid.Min(id,1,"id").Message("ID必须大于0")
	code:=e.INVALID_PARAMS

	if !valid.HasErrors() {
		code=e.SUCCESS
		if models.ExistTagById(id){
			models.DeleteTag(id)
		}else{
			code=e.ERROR_NOT_EXIST_TAG
		}
	}else{
		// 参数错误的话通过日志的形式打印
		log.Fatalf("参数出现错误：%v",valid.Errors)
	}
	//
	c.JSON(http.StatusOK,gin.H{
		"code":code,
		"msg":e.GetMsg(code),
		"data":make(map[string]string),
	})
}
