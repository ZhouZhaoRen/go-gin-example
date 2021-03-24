package main

import (
	"github.com/robfig/cron"
	"go-gin-example/models"
	"log"
	"time"
)
// 要开启定时任务的时候，这里要改为main方法
func main2() {
	log.Println("Starting ...")

	c := cron.New()
	// 往定时器中添加删除文章的任务
	c.AddFunc("* * * * * *", func() {
		log.Println("Run models.DeleteArticles...")
		models.DeleteArticles()
	})
	// 往定时器中添加删除标签的任务
	c.AddFunc("* * * * * *", func() {
		log.Println("Run models.DeleteTags...")
		models.DeleteTags()
	})
	// 启动定时器
	c.Run()

	// 为了阻塞主程序，下面两种方法都可以
	//select {
	//
	//}
	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}
