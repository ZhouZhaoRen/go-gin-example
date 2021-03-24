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

// 由于手动实现了callback函数，这里已经没有用了
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
func ExistTagByName(name string) (bool,error) {
	var tag Tag
	err:=db.Select("id").Where("name=?", name).First(&tag).Error
	// 由于不存在也算是错误，所以得考虑不存在的情况
	if err!=nil && err!=gorm.ErrRecordNotFound{
		return false,err
	}
	if tag.ID > 0 {
		return true,nil
	}
	return false,err
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
func AddTag(name string, state int, createBy string) error {
	err:=db.Create(&Tag{
		Name: name,
		//CreatedOn
		State:     state,
		CreatedBy: createBy,
	}).Error
	if err!=nil{
		return err
	}
	return nil
}

// 根据页码获取全部的标签
func GetTags(pageNum int, pageSize int, maps interface{}) ( []Tag, error) {
	var tags []Tag
	//err:=db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	err:=db.Where(maps).Find(&tags).Error
	if err != nil {
		return nil,err
	}
	return tags,nil
}

// 获取全部的标签总数
func GetTagTotal(maps interface{}) ( int, error) {
	var count int
	err:=db.Model(&Tag{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count,nil
}

// 根据id删除一个标签
func DeleteTag(id int) error {
	err:=db.Where("id=?", id).Delete(&Tag{}).Error
	if err!=nil{
		return err
	}
	return nil
}

// 根据id去修改标签
func UpdateTag(id int, data interface{}) bool {
	db.Model(&Tag{}).Where("id=?", id).Update(data)
	return true
}

// 删除deletedOn不为0的数据
func DeleteTags()bool{
	db.Unscoped().Where("deleted_on != ?",0).Delete(&Tag{})
	return true
}

