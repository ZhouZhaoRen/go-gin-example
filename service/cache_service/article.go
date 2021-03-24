package cache_service

import (
	"strconv"
	"strings"
	"go-gin-example/pkg/logging"
	"go-gin-example/pkg/e"
)

type Article struct {
	ID    int
	TagID int
	State int

	PageNum  int
	PageSize int
}

// 为缓存创建一个key
func (a *Article) GetArticleKey() string {
	return e.CACHE_ARTICLE + "_" + strconv.Itoa(a.ID) //ARTICLE_ID
}

func (a *Article) GetArticlesKey() string {
	keys := []string{
		e.CACHE_ARTICLE,
		"LIST",
	}

	if a.ID > 0 {
		keys = append(keys, strconv.Itoa(a.ID))
	}
	if a.TagID > 0 {
		keys = append(keys, strconv.Itoa(a.TagID))
	}
	if a.State >= 0 {
		keys = append(keys, strconv.Itoa(a.State))
	}
	if a.PageNum > 0 {
		keys = append(keys, strconv.Itoa(a.PageNum))
	}
	if a.PageSize > 0 {
		keys = append(keys, strconv.Itoa(a.PageSize))
	}
	s:=strings.Join(keys, "_") // 这里是为了测试返回的是怎样的key
	logging.Info(s)
	return s
}