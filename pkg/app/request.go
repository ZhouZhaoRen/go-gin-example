package app

import (
	"go-gin-example/pkg/logging"

	"github.com/astaxie/beego/validation"
)

// MarkErrors 的作用是循环将错误数组记录到日志当中
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Info(err.Key, err.Message)
	}

	return
}
