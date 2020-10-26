package main

import (
	"fmt"
	"go-gin-example/pkg/setting"
	"go-gin-example/routers"
	"net/http"
)

func main(){
	//router:=gin.Default() // 相当于创建了一个路由Handlers.可以后期绑定各类的路由规则和函数、中间件等
	//router.GET("/test", func(c *gin.Context) {
	//	c.JSON(200,gin.H{
	//		"message":"test",
	//	})
	//})
	router:=routers.InitRouter()

	s:=&http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort), // 监听的TCP地址
		Handler:        router, // http句柄，用于处理程序响应http请求
		ReadTimeout:    setting.ReadTimeout, //允许读取请求头的最大数时间
		WriteTimeout:   setting.WriteTimeout, // 允许写入的最大时间
		MaxHeaderBytes: 1 << 20, // 请求头的最大字节数
	}

	// 这里和r.Run()没有本质上的区别
	s.ListenAndServe() // 开始监听服务，监听 TCP 网络地址，Addr 和调用应用程序处理连接上的请求。
}