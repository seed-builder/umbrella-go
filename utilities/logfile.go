package utilities

import (
	"fmt"
	"time"
	"os"
)

type LogFile struct {

}

func (lf *LogFile) CurrentFile()(*os.File){
	fileName := fmt.Sprintf("log.%s.txt", time.Now().Format("2006_01_02"))
	cur, _ := os.Getwd()
	logDir := cur + "/log"
	if b, _ := lf.PathExists(logDir); !b {
		os.Mkdir(logDir, os.ModePerm)
	}
	file,_ := os.OpenFile(logDir + "/" + fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE,0666)
	return file
}

func (lf *LogFile) Write(p []byte) (n int, err error) {
	file := lf.CurrentFile()
	return file.Write(p)
}

func (lf *LogFile) PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}