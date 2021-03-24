package file

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
)

//
func GetSize(f multipart.File)(int,error){
	content,err:=ioutil.ReadAll(f)
	return len(content),err
}

func GetExt(fileName string)string{
	return path.Ext(fileName)
}

func CheckNotExist(src string) bool{
	_,err:=os.Stat(src)

	return os.IsNotExist(err)  // 能够接受ErrNotExist、syscall的一些错误，它会返回一个布尔值，能够得知文件不存在或目录不存在
}

//
func CheckPermission(src string)bool{
	_,err:=os.Stat(src) // 返回文件信息结构描述文件

	return os.IsPermission(err)  // 能够接受ErrPermission、syscall的一些错误，它会返回一个布尔值，能够得知权限是否满足
}

//
func IsNotExistMKDir(src string)error{
	if notExist:=CheckNotExist(src); notExist==true{
		if err:=MkDir(src); err!=nil{
			return err
		}
	}

	return nil
}

func MkDir(src string)error{
	err:=os.MkdirAll(src,os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func Open(name string,flag int,perm os.FileMode)(*os.File,error){
	f,err:=os.OpenFile(name,flag,perm)
	if err != nil {
		return nil, err
	}

	return f,nil
}

// MustOpen maximize trying to open the file
func MustOpen(fileName, filePath string) (*os.File, error) {
	// 1、获取根路径名
	dir, err := os.Getwd() // 获取当前目录对应的根路径名  如：F:\golearn1\src\go-gin-example2
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}
	// 2、进行目的目录的拼接
	src := dir + "/" + filePath // 对目的路径进行拼接 F:\golearn1\src\go-gin-example2 +  \runtime\logs
	// 3、查看权限
	perm := CheckPermission(src) // 检查权限是否满足
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}
	// 4、查看目录是否存在
	err = IsNotExistMkDir(src) // 再检查这个目录是否存在
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	// 5、打开文件，返回文件
	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}

// IsNotExistMkDir create a directory if it does not exist
// 检查当前路径的文件夹是否存在，不存在的话则创建文件夹
func IsNotExistMkDir(src string) error {
	// 如果不存在的话创建这个文件夹
	notExist := CheckNotExist(src)
	if  notExist == true {
		err := MkDir(src)
		if err != nil {
			return err
		}
	}

	return nil
}