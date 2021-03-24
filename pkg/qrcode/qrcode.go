package qrcode

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"go-gin-example/pkg/file"
	"go-gin-example/pkg/setting"
	"go-gin-example/pkg/util"
	"image/jpeg"
)

type QrCode struct {
	URL    string
	Width  int
	Height int
	Ext    string
	Level  qr.ErrorCorrectionLevel
	Mode   qr.Encoding
}

const (
	EXT_JPG = ".jpg"
)

// 返回一个二维码的实例化结构体
func NewQrCode(url string, width, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *QrCode {
	return &QrCode{
		URL:    url,
		Width:  width,
		Height: height,
		Level:  level,
		Mode:   mode,
		Ext:    EXT_JPG,
	}
}

func GetQrCodePath() string {
	return setting.AppSetting.QrCodeSavePath
}

// 完整路径  runtime/qrcode
func GetQrCodeFullPath() string {
	return setting.AppSetting.RuntimeRootPath + setting.AppSetting.QrCodeSavePath
}

func GetQrCodeFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetQrCodePath() + name
}

// 对二维码指定的路径进行加密，防止泄露
func GetQrCodeFileName(value string) string {
	return util.EncodeMD5(value)
}

func (q *QrCode) GetQrCodeExt() string {
	return q.Ext
}

func (q *QrCode) CheckEncode(path string) bool {
	src := path + GetQrCodeFileName(q.URL) + q.GetQrCodeExt()
	if file.CheckNotExist(src) == true {
		return false
	}
	return true
}

func (q *QrCode) Encode(path string) (string, string, error) {
	// 获取二维码生成路径
	name := GetQrCodeFileName(q.URL) + q.GetQrCodeExt()  // 7387cc83f2af55579f4047a5fa239f83.jpg  作为文件名
	src := path + name  //src=/runtime/qrcode/7387cc83f2af55579f4047a5fa239f83+.jpg
	if file.CheckNotExist(src) == true { // 检查文件是否存在
		//  创建二维码
		code, err := qr.Encode(q.URL, q.Level, q.Mode)
		if err != nil {
			return "", "", err
		}

		// 缩放二维码到指定位置
		code, err = barcode.Scale(code, q.Width, q.Height)
		if err != nil {
			return "", "", err
		}

		f, err := file.MustOpen(name, path) // 传入文件名和文件路径
		if err != nil {
			return "", "", err
		}
		defer f.Close()
		// 新建存放二维码图片的文件
		err = jpeg.Encode(f, code, nil) // 将二维码写到指定的文件
		if err != nil {
			return "", "", err
		}
	}

	return name, path, nil
}
