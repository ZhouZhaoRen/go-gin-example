package app

import (
	"go-gin-example/pkg/e"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// BindAndValid 的作用是将上下文中传进来的信息和表格进行对应
func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, e.ERROR
	}

	if !check {
		MarkErrors(valid.Errors) // 出现错误，将错误记录到日志中，错误是个数组，循环获取记录
		return http.StatusBadRequest, e.INVALID_PARAMS
	}
	//
	return http.StatusOK, e.SUCCESS
}
