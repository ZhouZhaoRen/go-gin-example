package tag_service

import (
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/EDDYCJY/go-gin-example/pkg/file"
	"github.com/tealeg/xlsx"
	"go-gin-example/models"
	"go-gin-example/pkg/export"
	"go-gin-example/pkg/gredis"
	"go-gin-example/pkg/logging"
	"go-gin-example/service/cache_service"
	"io"
	"strconv"
	"time"
)

// 之所以这样，是因为之后的crud可能会用到这几个
type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

func (t *Tag) ExistById() (bool, error) {
	return models.ExistArticleByID(t.ID)
}

func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

func (t *Tag) Edit() error {

	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	if t.State >= 0 {
		data["state"] = t.State
	}

	return models.EditArticle(t.ID, data)
}

func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.getMpas())
}

func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTags []models.Tag
	)

	cache := cache_service.Tag{
		State: t.State,

		PageSize: t.PageSize,
		PageNum:  t.PageNum,
	}
	// 如果存在则从Redis中读取，否则从数据库读取
	key := cache.GetTagsKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}
	//
	logging.Info("Redis中不存在当前数据")
	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMpas())
	if err != nil {
		return nil, err
	}
	gredis.Set(key, tags, 3600)
	return tags, nil
}

// 导出Excel表格
func (t *Tag) Export() (string, error) {
	tags,err:=t.GetAll()
	if err != nil {
		return "", err
	}
	//
	xlsFile:=xlsx.NewFile()
	sheet,err:=xlsFile.AddSheet("标签信息")
	if err != nil {
		return "",err
	}

	titles :=[]string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	row:=sheet.AddRow()

	var cell *xlsx.Cell
	for _,title:=range titles{
		cell=row.AddCell()
		cell.Value=title
	}

	for _,v:=range tags {
		values:=[]string{
			strconv.Itoa(v.ID),
			v.Name,
			v.CreatedBy,
			strconv.Itoa(v.CreatedOn),
			v.ModifiedBy,
			strconv.Itoa(v.ModifiedOn),
		}

		row=sheet.AddRow()
		for _,value:=range values {
			cell=row.AddCell()
			cell.Value=value
		}
	}
	//
	time:=strconv.Itoa(int(time.Now().Unix()))
	fileName:="tags-"+time+export.EXT

	dirFullPath :=export.GetExcelFullPath()
	err=file.IsNotExistMkDir(dirFullPath)
	if err != nil {
		return "",err
	}

	err=xlsFile.Save(dirFullPath+fileName)
	if err != nil {
		return "", err
	}

	return fileName,nil

}

func (t *Tag) Import(r io.Reader) error{
	xlsx,err:=excelize.OpenReader(r)
	if err != nil {
		return err
	}

	rows :=xlsx.GetRows("标签信息")
	// 遍历表格中的全部数据，然后一行一行的添加到数据库中
	for irow,row:=range rows {
		if irow>0{
			var data []string
			for _,cell:=range row {
				data=append(data,cell)
			}
			// 将表格的数据一行一行读出来，然后封装添加到数据库中
			models.AddTag(data[1],1,data[2])
		}
	}
	return nil
}

func (t *Tag) getMpas() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}
	return maps
}
