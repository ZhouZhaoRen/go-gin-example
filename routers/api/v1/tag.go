package v1

// 标签tag的专属controller层
import (
	"fmt"
	"go-gin-example/pkg/app"
	"go-gin-example/pkg/e"
	"go-gin-example/pkg/export"
	"go-gin-example/pkg/logging"
	"go-gin-example/pkg/setting"
	"go-gin-example/pkg/util"
	"go-gin-example/service/tag_service"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// GetTags 获取文章标签
func GetTags(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.Query("name")

	//maps := make(map[string]interface{})
	///data := make(map[string]interface{})

	if name != "" {
		//maps["name"] = name
	}
	//
	//var state int = -1
	state := -1

	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		//maps["state"] = state
	}

	//code := e.SUCCESS
	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageSize: setting.AppSetting.PageSize,
		PageNum:  util.GetPage(c),
	}
	tags, err := tagService.GetAll()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	count, err := tagService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": tags,
		"count": count,
	})

	// 返回json数据 有了上面的，这里变的多余了
	//c.JSON(http.StatusOK, gin.H{
	//	"code": code,
	//	"msg":  e.GetMsg(code),
	//	"data": tags,
	//})
}

// AddFormTag 定义一个结构体接收前端传过来的结构体
type AddFormTag struct {
	Name      string `form:"name" valid:"Required;MaxSize(100)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	State     int    `form:"state" valid:"Range(0,1)"`
}

// AddTag 新增文章标签
// @Summary 新增文章标签
// @Accept  json
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} string "{"code":200,"data":{},"msg":"ok"}" // 这里的string要改成go本身拥有的类型，用json是不可以的
// @Router /api/v1/tags [post]
// AddTag 新增文章标签
func AddTag(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form AddFormTag
	)
	//name:=c.Query("name")
	//state:=com.StrTo(c.DefaultQuery("state","0")).MustInt()
	//createdBy:=c.Query("created_by")
	//
	//// 通过beego的validation进行判断输入的数据是否符合要求
	//valid:=validation.Validation{}
	//valid.Required(name,"name").Message("名称不能为空")
	//valid.MaxSize(name,100,"name").Message("名称最长为100字符")
	//valid.Required(createdBy,"createdBy").Message("创建人不能为空")
	//valid.MaxSize(createdBy,100,"createdBy").Message("创建人最长为100字符")
	//valid.Range(state,0,1,"state").Message("状态只允许为0或1")
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	//code:=e.INVALID_PARAMS // 不可用的参数
	tagService := tag_service.Tag{
		Name:      form.Name,
		CreatedBy: form.CreatedBy,
		State:     form.State,
	}
	exists, err := tagService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if exists {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG, nil)
		return
	}

	err = tagService.Add()

	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

// EditTagForm 修改标签
type EditTagForm struct {
	ID         int    `form:"id" valid:"Required;Min(1)"`
	Name       string `form:"name" valid:"Required;MaxSize(100)"`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State      int    `form:"state" valid:"Range(0,1)"`
}

//EditTag  修改文章标签
func EditTag(c *gin.Context) {
	//id:=com.StrTo(c.Param("id")).MustInt() // 127.0.0.1:8000/api/v1/tags/1?id=1&state=0&name=zzr&modified_by=zzr  id和别的参数不一样
	//name:=c.Query("name")
	//modifiedBy:=c.Query("modified_by")
	//
	//valid:=validation.Validation{}
	//
	//state:=-1
	//if arg:=c.Query("state"); arg!=""{
	//	state=com.StrTo(arg).MustInt()
	//	valid.Range(state,0,1,"state").Message("状态只允许为0或1")
	//}
	////
	//valid.Required(id,"id").Message("id不能为空")
	//valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	//valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	//valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	var (
		appG = app.Gin{C: c}
		form EditTagForm
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	//
	tagService := tag_service.Tag{
		ID:         form.ID,
		Name:       form.Name,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}
	// 先判断标签是否存在
	exists, err := tagService.ExistById()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}
	err = tagService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}
	// 结束
	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

// DeleteTag 删除文章标签
func DeleteTag(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	tagService := tag_service.Tag{ID: id}
	exists, err := tagService.ExistById()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}
	err = tagService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// ExportTag 将所有的标签信息作为表格导出
func ExportTag(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.PostForm("name")
	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}
	fmt.Println("开始实例化service层")
	tagService := tag_service.Tag{
		Name:  name,
		State: state,
	}
	//
	fileName, err := tagService.Export()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXPORT_TAG_FAIL, nil)
		return
	}
	//
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"export_url":      export.GetExcelFullUrl(fileName),
		"export_save_url": export.GetExcelPath() + fileName,
	})
}

// ImportTag 导入标签
func ImportTag(c *gin.Context) {
	appG := app.Gin{C: c}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logging.Info(err)
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	//
	tagService := tag_service.Tag{}
	err = tagService.Import(file)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, e.ERROR_IMPORT_TAG_FAIL, nil)
		return
	}
	//
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
