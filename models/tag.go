package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Tag struct {
	Model
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

// Tag的callback函数
func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}

//
func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())
	return nil
}

// 根据名字判断标签是否存在，返回true或者false
func ExistTagByName(name string) bool {
	var tag Tag
	db.Select("id").Where("name=?", name).First(&tag)
	if tag.ID > 0 {
		return true
	}
	return false
}

// 根据id判断标签是否存在
func ExistTagById(id int) bool {
	var tag Tag
	db.Where("id=?", id).First(&tag)
	if tag.ID > 0 {
		return true
	}
	return false
}

// 新增一个标签
func AddTag(name string, state int, createBy string) bool {
	db.Create(&Tag{
		Name: name,
		//CreatedOn
		State:     state,
		CreatedBy: createBy,
	})

	return true
}

// 根据页码获取全部的标签
func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)

	return
}

// 获取全部的标签总数
func GetTagTotal(maps interface{}) (count int) {
	db.Model(&Tag{}).Where(maps).Count(&count)

	return
}

// 根据id删除一个标签
func DeleteTag(id int) bool {
	db.Where("id=?", id).Delete(&Tag{})
	return true
}

// 根据id去修改标签
func UpdateTag(id int, data interface{}) bool {
	db.Model(&Tag{}).Where("id=?", id).Update(data)
	return true
}
