package utils

import (
	"log"
	"os"
	"path"
)

func Exist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}
func MakeDir(filePath string) {
	if !Exist(filePath) {
		MkDirAll(filePath)
	}
}

func MkDirAll(path string) bool {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func Ext(fileName string, defaultExt string) string {
	t := path.Ext(fileName)
	if len(t) == 0 {
		return defaultExt
	}
	return t
}
