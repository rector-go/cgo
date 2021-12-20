package cgo

import (
	"cgo/constant"
	"cgo/utils"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type Controller struct {
	Data interface{}
}
type FileInfoTO struct {
	// 图片ID
	ID int64
	// 缩略图路径
	CompressPath string
	// 原图路径
	Path string
	// 原始文件名
	OriginalFileName string
	FileName         string
	FileSize         int64
}

func (p *Controller) GetFileNum(r *http.Request, keys ...string) int {
	m := r.MultipartForm
	if m == nil {
		return 0
	}
	if len(keys) == 0 {
		var num int
		for _, fileHeaders := range m.File {
			num += len(fileHeaders)
		}
		return num
	} else {
		var num int
		for _, value := range keys {
			num += len(m.File[value])
		}
		return num
	}
}

func (p *Controller) saveFile(filePath, relativePath string, fileHeader *multipart.FileHeader) *FileInfoTO {
	file, err := fileHeader.Open()
	if err != nil {
		log.Println(err)
		return nil
	}
	defer file.Close()
	name := utils.RandomUUID()
	fileType := utils.Ext(fileHeader.Filename, ".jpg")
	newName := name.String() + fileType
	dst, err := os.Create(filePath + newName)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer dst.Close()
	fileSize, err := io.Copy(dst, file)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &FileInfoTO{
		Path:             relativePath + newName,
		OriginalFileName: fileHeader.Filename,
		FileName:         newName,
		FileSize:         fileSize,
	}
}

func (p *Controller) SaveFiles(r *http.Request, relativePath string, keys ...string) []*FileInfoTO {
	r.ParseMultipartForm(32 << 20)
	m := r.MultipartForm
	if m == nil {
		log.Println("not multipartfrom !")
		return nil
	}
	fileInfos := make([]*FileInfoTO,0)

	filePath := constant.BASE_IMAGE_ADDRESS + relativePath
	utils.MakeDir(filePath)

	if len(keys) == 0 {
		for _,fileHeaders := range m.File { //遍历所有的所有的字段名(filename)获取FileHeaders
			for _,fileHeader := range fileHeaders{
				to := p.saveFile(filePath,relativePath,fileHeader)
				fileInfos = append(fileInfos,to)
			}
		}
	} else {
		for _,value := range keys {
			fileHeaders := m.File[value]//根据上传文件时指定的字段名(filename)获取FileHeaders
			for _,fileHeader := range fileHeaders{
				to := p.saveFile(filePath,relativePath,fileHeader)
				fileInfos = append(fileInfos,to)
			}
		}
	}

	return fileInfos
}
