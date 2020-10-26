package v1

import (
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go-gin-example/models"
	"go-gin-example/pkg/e"
	"go-gin-example/pkg/setting"
	"go-gin-example/pkg/util"
	"net/http"
)
// 获取单个文章
func GetArticle(c *gin.Context){
	id:=com.StrTo(c.Param("id")).MustInt()
	valid:=validation.Validation{}
	valid.Min(id,1,"id").Message("id必须大于0")
	code:=e.INVALID_PARAMS
	// 实例化一个结构体接数据
	var data models.Article
	if valid.HasErrors(){
		fmt.Println("输入的id有问题：",valid.Errors)
	}else{
		if models.ExistArticleByID(id){
			data=models.GetArticle(id)
			code=e.SUCCESS
		}else{
			code=e.ERROR_NOT_EXIST_ARTICLE
		}
	}
	// 返回json数据
	c.JSON(http.StatusOK,gin.H{
		"code":code,
		"msg":e.GetMsg(code),
		"data":data,
	})
}

// 获取多个文章
func GetArticles(c *gin.Context){
	data:=make(map[string]interface{})
	maps:=make(map[string]interface{})
	valid:=validation.Validation{}
	var state int=-1
	if arg:=c.Query("state"); arg!=""{
		state=com.StrTo(arg).MustInt()
		maps["state"]=state
		valid.Range(state,0,1,"state").Message("状态只允许0或1")
	}
	//
	var tagId int =-1
	if arg:=c.Query("tag_id");arg!=""{
		tagId=com.StrTo(arg).MustInt()
		maps["tag_id"]=tagId
		valid.Min(tagId,1,"tag_id").Message("标签ID必须大于0")
	}
	//
	code:=e.INVALID_PARAMS
	if valid.HasErrors(){
		fmt.Println("输入数据出错：",valid.Errors)
	}else{
		data["lists"]=models.GetArticles(util.GetPage(c),setting.PageSize,maps)
		data["total"]=models.GetArticleTotal(maps)
		code=e.SUCCESS
	}
	// 返回一个json
	c.JSON(http.StatusOK,gin.H{
		"code":code,
		"msg":e.GetMsg(code),
		"data":data,
	})
}

// 新增文章
func AddArticle(c *gin.Context){
	tagId :=com.StrTo(c.Query("tag_id")).MustInt()
	title:=c.Query("title")
	desc:=c.Query("desc")
	content:=c.Query("content")
	createdBy:=c.Query("created_by")
	state:=com.StrTo(c.Query("state")).MustInt()

	// 对输入进行判断
	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code:=e.INVALID_PARAMS
	if !valid.HasErrors(){
		if models.ExistTagById(tagId){
			data:=make(map[string]interface{})
			data["tag_id"]=tagId
			data["title"]=title
			data["desc"]=desc
			data["content"]=content
			data["created_by"]=createdBy
			data["state"]=state
			// 进行插入操作
			models.AddArticle(data)

			code=e.SUCCESS
		}else{
			code=e.ERROR_NOT_EXIST_TAG
		}
	}else{
		fmt.Println("输入有问题：",valid.Errors)
	}

	// 向前端返回json数据
	c.JSON(http.StatusOK,gin.H{
		"code":code,
		"msg":e.GetMsg(code),
		"data":make(map[string]string),
	})
}

// 修改文章
func EditArticle(c *gin.Context){
	valid := validation.Validation{}
	id:=com.StrTo(c.Param("id")).MustInt()
	tagId :=com.StrTo(c.Query("tag_id")).MustInt()
	title:=c.Query("title")
	desc:=c.Query("desc")
	content:=c.Query("content")
	modifiedBy:=c.Query("modified_by")
	//state:=com.StrTo(c.Query("state")).MustInt()
	var state=-1
	if arg:=c.Query("state"); arg!=""{
		state=com.StrTo(arg).MustInt()
		valid.Range(state,0,1,"state").Message("状态只允许为0或1")
	}
	// 对输入进行判断
	valid.Min(id,1,"id").Message("id必须大于0")
	valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(modifiedBy, "modified_by").Message("创建人不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code:=e.INVALID_PARAMS
	//
	if valid.HasErrors(){
		fmt.Println("输入参数有问题：",valid.Errors)
	}else{
		if models.ExistArticleByID(id){
			data:=make(map[string]interface{})
			data["tag_id"]=tagId
			data["title"]=title
			data["desc"]=desc
			data["content"]=content
			data["modified_by"]=modifiedBy
			data["state"]=state
			// 开始修改
			models.EditArticle(id,data)
			code=e.SUCCESS
		}else {
			code=e.ERROR_NOT_EXIST_ARTICLE // 文章不存在
		}
	}
	//
	c.JSON(http.StatusOK,gin.H{
		"code":code,
		"msg":e.GetMsg(code),
		"data":make(map[string]string),
	})
}

// 删除文章
func DeleteArticle(c *gin.Context){
	id:=com.StrTo(c.Param("id")).MustInt()
	valid:=validation.Validation{}
	valid.Min(id,1,"id").Message("id必须大于0")
	code:=e.INVALID_PARAMS
	if valid.HasErrors(){
		fmt.Println("输入的id有问题：",valid.Errors)
	}else{
		if models.ExistArticleByID(id){
			models.DeleteArticle(id)
			code=e.SUCCESS
		}else{
			code=e.ERROR_NOT_EXIST_ARTICLE
		}
	}
	// 返回json数据
	c.JSON(http.StatusOK,gin.H{
		"code":code,
		"msg":e.GetMsg(code),
		"data":make(map[string]string),
	})

}