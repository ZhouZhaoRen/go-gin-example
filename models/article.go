package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag Tag `json:"tag"`
	Title string `json:"title"`
	Desc string `json:"desc"`
	Content string `json:"content"`
	CreatedBy string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State int `json:"state"`
}
//artilce的callback函数
func (article *Article) BeforeCreate(scope *gorm.Scope)error{
	scope.SetColumn("CreatedOn",time.Now().Unix())
	return nil
}
func (article *Article) BeforeUpdate(scope *gorm.Scope)error{
	scope.SetColumn("ModifiedOn",time.Now().Unix())
	return nil
}

func (article *Article) BeforeDelete(scope *gorm.Scope)error{
	scope.SetColumn("DeletedOn",time.Now().Unix())
	return nil
}

// 通过id判断文章是否存在
func ExistArticleByID(id int)bool{
	var article Article
	db.Select("id").Where("id=?",id).First(&article)
	if article.ID>0{
		return true
	}
	return false
}

// 统计符合条件的文章总数
func GetArticleTotal(maps interface{})(count int){
	db.Model(&Article{}).Where(maps).Count(&count)
	return
}

// 查到符合条件的文章，返回一个数组，文章包含具体的标签，所以得连表查询
func GetArticles(pageNum,pageSize int,maps interface{})(articles []Article){
	// 先查询出所有的文章，再根据文章的tagId找到对应的tag，再通过映射逻辑，将其填充到Article的Tag中
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles) // 预加载器
	return
}

// 通过id获取具体的文章
func GetArticle(id int)(article Article){
	db.Where("id=?",id).First(&article)
	db.Model(&article).Related(&article.Tag) // 进行关联查询
	return
}

// 通过id去修改文章
func EditArticle(id int,maps interface{})bool{
	db.Model(&Article{}).Where("id=?",id).Update(maps)
	return true
}

// 添加文章
func AddArticle(data map[string]interface{})bool{
	db.Create(&Article{
		TagID: data["tag_id"].(int), // 这个实际就是 Golang 中的类型断言，用于判断一个接口值的实际类型是否为某个类型，或一个非接口值的类型是否实现了某个接口类型
		Title: data["title"].(string),
		Desc: data["desc"].(string),
		Content: data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State: data["state"].(int),
	})

	return true
}

// 通过id删除对应的文章
func DeleteArticle(id int)bool{
	db.Where("id=?",id).Delete(&Article{})
	return true
}